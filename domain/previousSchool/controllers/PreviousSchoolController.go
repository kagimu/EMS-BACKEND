package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	applicationModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/models"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/previousSchool/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreatePreviousSchool(c *gin.Context) {
	type PreviousSchoolInfo struct {
		Name        string   `json:"name"`
		StartYear   string   `json:"start_year"`
		EndYear     string   `json:"end_year"`
		ClassNumber []string `json:"class_number"`
	}

	var previousSchoolInfo PreviousSchoolInfo

	if err := c.ShouldBindJSON(&previousSchoolInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Internal server errror"})
		return
	}

	classes := previousSchoolInfo.ClassNumber
	if len(classes) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "at least one class is required!"})
		return
	}

	userNumber := c.Request.Header["User-Number"][0]

	var user authModel.User
	if err := config.DB.Where("number = ?", userNumber).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user not found"})
		return
	}

	studentNumber := c.Param("student_number")
	applicationNumber := c.Param("application_number")

	var application applicationModel.Application
	if err := config.DB.
		Where("number = ?", applicationNumber).
		First(&application).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "application not found"})
		return
	}

	savePreviousSchool := models.PreviousSchool{
		Number:            uuid.NewV4().String(),
		Name:              previousSchoolInfo.Name,
		StartYear:         previousSchoolInfo.StartYear,
		EndYear:           previousSchoolInfo.EndYear,
		StudentNumber:     studentNumber,
		ApplicationNumber: applicationNumber,
		UserNumber:        userNumber,
	}

	if err := config.DB.Create(&savePreviousSchool).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     err.Error(),
			"message": "failed to save previous schools"})
		return
	}

	for _, class := range classes {
		saveClasses := models.PreviousSchoolClass{
			StudentNumber: studentNumber,
			Number:        uuid.NewV4().String(),
			ClassNumber:   class,
		}

		err := config.DB.Create(&saveClasses).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err":     err.Error(),
				"message": "failed   to save previous schools"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "previous school saved successfully",
		})

	}
}
