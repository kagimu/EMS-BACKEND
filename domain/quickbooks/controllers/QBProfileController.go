package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/quickbooks/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

type ProfileRequest struct {
	Name             string `json:"name" binding:"required"`
	Username         string `json:"username" binding:"required"`
	Password         string `json:"password"`
	DepositToAccount string `json:"deposit_to_account"`
	IncomeAccount    string `json:"income_account"`
	IsDefault        bool   `json:"is_default"`
	IsActive         *bool  `json:"is_active"`
}

func GetProfiles(c *gin.Context) {
	var profiles []models.QBProfile
	if err := config.DB.Find(&profiles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch QuickBooks profiles"})
		return
	}

	type ProfileResponse struct {
		models.QBProfile
		HasPassword bool `json:"has_password"`
	}

	var response []ProfileResponse
	for _, p := range profiles {
		hasPass := p.PasswordHash != ""
		p.PasswordHash = "" // Do not expose password hash
		response = append(response, ProfileResponse{
			QBProfile:   p,
			HasPassword: hasPass,
		})
	}

	c.JSON(http.StatusOK, response)
}

func CreateProfile(c *gin.Context) {
	var req ProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if username is already taken
	var existing models.QBProfile
	if err := config.DB.Where("username = ?", req.Username).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A QuickBooks account with this username already exists"})
		return
	}

	depositTo := req.DepositToAccount
	if depositTo == "" {
		depositTo = "Undeposited Funds"
	}
	incomeAcc := req.IncomeAccount
	if incomeAcc == "" {
		incomeAcc = "Service Income"
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	// If this profile is default, unset default on others
	if req.IsDefault {
		config.DB.Model(&models.QBProfile{}).Where("1 = 1").Update("is_default", false)
	}

	fileID := fmt.Sprintf("{%s}", strings.ToUpper(uuid.NewV4().String()))
	ownerID := "{90A44FB7-33D9-4815-AC85-BC87A7E7D1EB}"

	profile := models.QBProfile{
		Number:           uuid.NewV4().String(),
		Name:             req.Name,
		Username:         req.Username,
		PasswordHash:     req.Password,
		DepositToAccount: depositTo,
		IncomeAccount:    incomeAcc,
		FileID:           fileID,
		OwnerID:          ownerID,
		IsDefault:        req.IsDefault,
		IsActive:         isActive,
		SchoolNumber:     "default",
	}

	if err := config.DB.Create(&profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create QuickBooks account profile"})
		return
	}

	profile.PasswordHash = ""
	c.JSON(http.StatusCreated, profile)
}

func UpdateProfile(c *gin.Context) {
	id := c.Param("id")
	var profile models.QBProfile

	if err := config.DB.Where("number = ?", id).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "QuickBooks account profile not found"})
		return
	}

	var req ProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// If username changed, ensure unique
	if req.Username != profile.Username {
		var existing models.QBProfile
		if err := config.DB.Where("username = ? AND number != ?", req.Username, id).First(&existing).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "A QuickBooks account with this username already exists"})
			return
		}
	}

	if req.IsDefault && !profile.IsDefault {
		config.DB.Model(&models.QBProfile{}).Where("number != ?", id).Update("is_default", false)
	}

	profile.Name = req.Name
	profile.Username = req.Username
	if req.Password != "" {
		profile.PasswordHash = req.Password
	}
	if req.DepositToAccount != "" {
		profile.DepositToAccount = req.DepositToAccount
	}
	if req.IncomeAccount != "" {
		profile.IncomeAccount = req.IncomeAccount
	}
	profile.IsDefault = req.IsDefault
	if req.IsActive != nil {
		profile.IsActive = *req.IsActive
	}

	if err := config.DB.Save(&profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update QuickBooks account profile"})
		return
	}

	profile.PasswordHash = ""
	c.JSON(http.StatusOK, profile)
}

func DeleteProfile(c *gin.Context) {
	id := c.Param("id")
	var profile models.QBProfile

	if err := config.DB.Where("number = ?", id).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "QuickBooks account profile not found"})
		return
	}

	if err := config.DB.Delete(&profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete QuickBooks account profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "QuickBooks account profile deleted successfully"})
}

func DownloadProfileQWC(c *gin.Context) {
	id := c.Param("id")
	var profile models.QBProfile

	if err := config.DB.Where("number = ?", id).First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "QuickBooks account profile not found"})
		return
	}

	proto := c.GetHeader("X-Forwarded-Proto")
	if proto == "" {
		proto = "https"
		if c.Request.TLS == nil && strings.Contains(c.Request.Host, "localhost") {
			proto = "http"
		}
	}
	baseURL := fmt.Sprintf("%s://%s", proto, c.Request.Host)

	fileID := profile.FileID
	if fileID == "" {
		fileID = fmt.Sprintf("{%s}", strings.ToUpper(uuid.NewV4().String()))
		profile.FileID = fileID
		config.DB.Save(&profile)
	}
	ownerID := profile.OwnerID
	if ownerID == "" {
		ownerID = "{90A44FB7-33D9-4815-AC85-BC87A7E7D1EB}"
	}

	appName := fmt.Sprintf("WholesomeEd - %s", profile.Name)

	qwc := fmt.Sprintf(`<?xml version="1.0"?>
<QBWCXML>
   <AppName>%s</AppName>
   <AppID></AppID>
   <AppURL>%s/api/quickbooks/soap</AppURL>
   <AppDescription>Sync to QuickBooks - %s</AppDescription>
   <AppSupport>%s/support</AppSupport>
   <UserName>%s</UserName>
   <OwnerID>%s</OwnerID>
   <FileID>%s</FileID>
   <QBType>QBFS</QBType>
   <Scheduler>
      <RunEveryNMinutes>5</RunEveryNMinutes>
   </Scheduler>
   <IsReadOnly>false</IsReadOnly>
</QBWCXML>`, escapeXML(appName), baseURL, escapeXML(profile.Name), baseURL, escapeXML(profile.Username), ownerID, fileID)

	cleanName := strings.ReplaceAll(profile.Name, " ", "_")
	filename := fmt.Sprintf("wholesome_ed_%s.qwc", cleanName)

	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Data(http.StatusOK, "application/x-qbwc", []byte(qwc))
}
