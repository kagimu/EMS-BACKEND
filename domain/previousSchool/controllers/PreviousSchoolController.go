package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	applicationModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/models"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	previousSchoolModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/previousSchool/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreatePreviousSchools(c *gin.Context) {
	type Class struct {
		Name string `json:"name"`
		Year string `json:"year"`
	}
	type PreviousSchoolInfo struct {
		Name    string  `json:"name"`
		Classes []Class `json:"classes"`
	}

	var previousSchoolInfo []PreviousSchoolInfo

	if err := c.ShouldBindJSON(&previousSchoolInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Please format the correct request body"})
		return
	}

	userNumber := c.Request.Header["User-Number"][0]

	var user authModel.User
	if err := config.DB.Where("number = ?", userNumber).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user not found"})
		return
	}

	applicationNumber := c.Param("application_number")

	var application applicationModel.Application
	if err := config.DB.
		Where("number = ?", applicationNumber).
		First(&application).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "application not found"})
		return
	}

	for _, item := range previousSchoolInfo {
		savePreviousSchool := previousSchoolModel.PreviousSchool{
			Number:            uuid.NewV4().String(),
			Name:              item.Name,
			ApplicationNumber: applicationNumber,
		}

		if err := config.DB.Create(&savePreviousSchool).Error; err != nil {

			c.JSON(http.StatusBadRequest, gin.H{
				"debug":   err.Error(),
				"message": "failed to save previous school",
				"data":    savePreviousSchool,
			})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "student previous schools saved successfully",
	})
}
