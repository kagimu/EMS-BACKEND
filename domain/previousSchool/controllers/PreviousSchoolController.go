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

type Class struct {
	Name string `json:"name"`
	Year string `json:"year"`
}
type PreviousSchoolInfo struct {
	Name    string  `json:"name"`
	Classes []Class `json:"classes"`
}

func CreatePreviousSchools(c *gin.Context) {

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

	tx := config.DB.Begin()

	for _, item := range previousSchoolInfo {
		savePreviousSchool := previousSchoolModel.PreviousSchool{
			Number:            uuid.NewV4().String(),
			Name:              item.Name,
			ApplicationNumber: applicationNumber,
		}

		if err := tx.Create(&savePreviousSchool).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"debug":   err.Error(),
				"message": "failed to save previous school",
				"data":    savePreviousSchool,
			})
			return
		}

		// save previous school classes

		for _, class := range item.Classes {
			savePreviousSchoolClass := previousSchoolModel.PreviousSchoolClass{
				PreviousSchoolNumber: savePreviousSchool.Number,
				Class:                class.Name,
				Year:                 class.Year,
			}

			if err := tx.Create(&savePreviousSchoolClass).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "failed to save previous school class",
				})
			}
		}

		// set application status of the application to SUBMITTED

		application.Status = "SUBMITTED"
		if err := tx.Save(&application).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "failed to update application status",
			})
			return

		}

	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"message": "student previous schools saved successfully",
	})
}

func UpdatePreviousSchool(ctx *gin.Context) {

	var previousSchoolInfo []PreviousSchoolInfo

	if err := ctx.ShouldBindJSON(&previousSchoolInfo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Please format the correct request body"})
		return
	}
	applicationNumber := ctx.Param("application_number")
	var previousSchools []previousSchoolModel.PreviousSchool
	if err := config.DB.Where("application_number = ?", applicationNumber).Find(&previousSchools).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Previous schools not found"})
		return
	}

	tx := config.DB.Begin()

	for _, prevSchool := range previousSchools {

		var previousSchoolClasses previousSchoolModel.PreviousSchoolClass

		// remove previous school classes
		if err := tx.Where("previous_school_number = ?", prevSchool.Number).Delete(&previousSchoolClasses).Error; err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "failed to update previous school classes"})
			return
		}

		// remove previous schools

	}
	var previousSchool previousSchoolModel.PreviousSchool

	if err := tx.Where("application_number = ?", applicationNumber).Delete(&previousSchool).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "failed to update previous schools"})
		return

	}

	for _, item := range previousSchoolInfo {
		savePreviousSchool := previousSchoolModel.PreviousSchool{
			Number:            uuid.NewV4().String(),
			Name:              item.Name,
			ApplicationNumber: applicationNumber,
		}

		if err := tx.Create(&savePreviousSchool).Error; err != nil {

			tx.Rollback()
			ctx.JSON(http.StatusBadRequest, gin.H{
				"debug":   err.Error(),
				"message": "failed to save previous school",
			})
			return
		}

		// save previous school classes

		for _, class := range item.Classes {
			savePreviousSchoolClass := previousSchoolModel.PreviousSchoolClass{
				PreviousSchoolNumber: savePreviousSchool.Number,
				Class:                class.Name,
				Year:                 class.Year,
			}

			if err := tx.Create(&savePreviousSchoolClass).Error; err != nil {
				tx.Rollback()
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": "failed to save previous school class",
				})

				return
			}
		}

	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "student previous schools updated successfully"})
}
