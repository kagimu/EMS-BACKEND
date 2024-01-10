package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/sourceOfSchoolInfo/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func NewSourceOfSchoolInfo(c *gin.Context) {
	type SourceOfSchoolInfoRequest struct {
		Name string `json:"name"`
	}

	var sourceOfSchoolInfo SourceOfSchoolInfoRequest

	err := c.ShouldBindJSON(&sourceOfSchoolInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Internal server error"})
		return
	}

	if sourceOfSchoolInfo.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "name is required"})
		return
	}

	userNumber := c.Request.Header["User-Number"][0]

	var user authModels.User
	if err := config.DB.
		Where("number = ?", userNumber).
		First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	tx := config.DB.Begin()
	saveSourceOfSchoolInfo := models.SourceOfSchoolInfo{
		Number:     uuid.NewV4().String(),
		Name:       sourceOfSchoolInfo.Name,
		UserNumber: userNumber,
	}

	if err := tx.Create(&saveSourceOfSchoolInfo).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to save source of school info"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to save school of school info"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusOK,
		"message": "source of school information created successfully",
		"data":    saveSourceOfSchoolInfo,
	})

}

func GetSourcesOfSchoolInfo(c *gin.Context) {
	var sourcesOfSchoolInfo []models.SourceOfSchoolInfo

	err := config.DB.Find(&sourcesOfSchoolInfo).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to fetch sources of school info",
		})
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": sourcesOfSchoolInfo})
}

func EditSourceOfSchoolInfo(c *gin.Context) {
	type SourceOfSchoolInfo struct {
		Name string `json:"name"`
	}

	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	var sourceOfSchoolInfo SourceOfSchoolInfo
	if err := c.BindJSON(&sourceOfSchoolInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	name := sourceOfSchoolInfo.Name
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Name is required!"})
		return
	}

	var sourceOfSchoolInfoModel models.SourceOfSchoolInfo
	sourceOfSchoolInfoNumber := c.Param("source_of_school_info_number")

	if err := config.DB.
		Where("number = ?", sourceOfSchoolInfoNumber).
		First(&sourceOfSchoolInfoModel).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     err.Error(),
			"message": "Source of school info not found"})
		return
	}

	var duplicateSourceOfSchoolInfoModel models.SourceOfSchoolInfo
	userNumber := user.Number
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("user_number = ?", userNumber).
		First(&duplicateSourceOfSchoolInfoModel).Error; duplicate == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "source of school info already exists"})
		return
	}

	tx := config.DB.Begin()

	if err := tx.Model(&sourceOfSchoolInfoModel).Where("number = ?", sourceOfSchoolInfoNumber).Update("name", sourceOfSchoolInfo.Name).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save source of school", "details": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save school!!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "source of school info edited successfully", "data": sourceOfSchoolInfoModel})
}

func DeleteSourceOfSchoolInfo(c *gin.Context) {
	userNumber := c.Request.Header["User-Number"][0]
	var user authModels.User

	if err := config.DB.
		Where("number = ?", userNumber).
		First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     err.Error(),
			"message": "User not found"})
		return
	}

	var sourceOfSchoolInfoModel models.SourceOfSchoolInfo
	sourceOfSchoolInfoNumber := c.Param("source_of_school_info_number")

	if err := config.DB.
		Where("number = ?", sourceOfSchoolInfoNumber).
		First(&sourceOfSchoolInfoModel).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     err.Error(),
			"message": "source of school info not found"})
		return
	}

	tx := config.DB.Begin()

	if err := tx.Delete(&sourceOfSchoolInfoModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete source of school info",
			"details": err.Error(),
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save school!!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Source of school info deleted successfully"})

}
