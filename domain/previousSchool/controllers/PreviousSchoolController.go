package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/repository"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/previousSchool/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

type Class struct {
	Class string `json:"class"`
	Year  string `json:"year"`
}
type PreviousSchoolInfo struct {
	Name    string  `json:"name"`
	Classes []Class `json:"classes"`
}

func CreatePreviousSchools(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Unauthorized"})
		return
	}

	var previousSchoolInfo []PreviousSchoolInfo

	if err := c.ShouldBindJSON(&previousSchoolInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Please format the correct request body"})
		return
	}

	applicationRepo := repository.NewApplicationRepositoryImpl(config.DB)

	applicationNumber := c.Param("application_number")

	application, err := applicationRepo.FindApplicationByNumber(applicationNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Application err" + err.Error()})
		return
	}

	if application.Number == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Application not found"})
		return
	}

	tx := config.DB.Begin()

	// delete previous classes
	var previousSchoolClasses models.PreviousSchoolClass
	if err := tx.Where("previous_school_number = ?", applicationNumber).Delete(&previousSchoolClasses).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to update previous school classes"})
		return
	}

	// delete previous schools
	var previousSchool models.PreviousSchool
	if err := tx.Where("application_number = ?", applicationNumber).Delete(&previousSchool).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to update previous schools"})
		return

	}

	for _, item := range previousSchoolInfo {
		savePreviousSchool := models.PreviousSchool{
			Number:            uuid.NewV4().String(),
			Name:              item.Name,
			ApplicationNumber: applicationNumber,
		}

		if err := tx.Create(&savePreviousSchool).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save previous school", "data": savePreviousSchool})
			return
		}

		// save previous school classes
		for _, class := range item.Classes {
			savePreviousSchoolClass := models.PreviousSchoolClass{
				Number:               uuid.NewV4().String(),
				PreviousSchoolNumber: savePreviousSchool.Number,
				Class:                class.Class,
				Year:                 class.Year,
			}

			if err := tx.Create(&savePreviousSchoolClass).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save previous school class", "err": err.Error()})
				return
			}
		}

	}
	/*if err := applicationRepo.UpdateApplicationStatus(applicationNumber, "SUBMITTED", tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update application status", "err": err.Error()})
		return
	}*/

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save previous schools"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Student previous schools saved successfully"})
}

func UpdatePreviousSchool(ctx *gin.Context) {

	var previousSchoolInfo []PreviousSchoolInfo

	if err := ctx.ShouldBindJSON(&previousSchoolInfo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Please format the correct request body"})
		return
	}
	applicationNumber := ctx.Param("application_number")
	var previousSchools []models.PreviousSchool
	if err := config.DB.Where("application_number = ?", applicationNumber).Find(&previousSchools).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Previous schools not found"})
		return
	}

	tx := config.DB.Begin()

	for _, prevSchool := range previousSchools {

		var previousSchoolClasses models.PreviousSchoolClass

		// remove previous school classes
		if err := tx.Where("previous_school_number = ?", prevSchool.Number).Delete(&previousSchoolClasses).Error; err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "failed to update previous school classes"})
			return
		}

		// remove previous schools

	}
	var previousSchool models.PreviousSchool

	if err := tx.Where("application_number = ?", applicationNumber).Delete(&previousSchool).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "failed to update previous schools"})
		return

	}

	for _, item := range previousSchoolInfo {
		savePreviousSchool := models.PreviousSchool{
			Number:            uuid.NewV4().String(),
			Name:              item.Name,
			ApplicationNumber: applicationNumber,
		}

		if err := tx.Create(&savePreviousSchool).Error; err != nil {

			tx.Rollback()
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save previous school"})
			return
		}

		// save previous school classes

		for _, class := range item.Classes {
			savePreviousSchoolClass := models.PreviousSchoolClass{
				PreviousSchoolNumber: savePreviousSchool.Number,
				Class:                class.Class,
				Year:                 class.Year,
			}

			if err := tx.Create(&savePreviousSchoolClass).Error; err != nil {
				tx.Rollback()
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": "Failed to save previous school class",
				})

				return
			}
		}

	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Student previous schools updated successfully"})
}
