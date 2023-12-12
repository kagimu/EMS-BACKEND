package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
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
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
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
