package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/sourceOfSchoolInfo/models"
	"github.com/gin-gonic/gin"
)

func NewSourceOfSchoolInfo(c *gin.Context) {
	var sourceOfSchoolInfo models.SourceOfSchoolInfo
	var err error

	err = c.ShouldBindJSON(&sourceOfSchoolInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Internal server error"})
		return
	}

	err = sourceOfSchoolInfo.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Internal server error"})
		return
	}

	sourceOfSchoolInfo.Prepare()

	err = config.DB.Create(&sourceOfSchoolInfo).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusOK,
		"message": "source of school information created successfully",
		"data":    sourceOfSchoolInfo,
	})

}

func GetSourcesOfSchoolInfo(c *gin.Context) {
	var sourcesOfSchoolInfo []models.SourceOfSchoolInfo

	err := config.DB.Find(sourcesOfSchoolInfo).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to fetch sources of school info",
		})
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": sourcesOfSchoolInfo})
}
