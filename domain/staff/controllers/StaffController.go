package controllers

import (
	"errors"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staff/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staff/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func CreateStaff(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.StaffRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	saveStaff := models.Staff{
		Number:                     uuid.NewV4().String(),
		Surname:                    payload.Surname,
		StaffNumber:                payload.StaffNumber,
		Firstname:                  payload.FirstName,
		OtherName:                  payload.OtherName,
		Gender:                     payload.Gender,
		MainTelephoneNumber:        payload.MainTelephoneNumber,
		AlternativeTelephoneNumber: payload.AlternativeTelephoneNumber,
		MainEmailAddress:           payload.MainEmailAddress,
		AlternativeEmailAddress:    payload.AlternativeEmailAddress,
		JobTitleNumber:             payload.JobTitleNumber,
		CampusNumber:               payload.CampusNumber,
		SalaryScaleNumber:          payload.SalaryScaleNumber,
		UserNumber:                 user.Number,
	}

	if err := config.DB.Create(&saveStaff).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save staff member"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Staff member created successfully"})
}

func FetchStaff(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var staffMembers []models.Staff
	if err := config.DB.Find(&staffMembers).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch staff members"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": staffMembers})
}

func FetchMemberByStaffNumber(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var staff models.Staff

	staffNumber := c.Query("staff_number")

	if staffNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Staff number is missing"})
		return
	}

	if err := config.DB.First(&staff, "staff_number = ?", staffNumber).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "No staff member found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": staff})
}

func EditStaff(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.StaffRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	var staff models.Staff
	staffNumber := c.Query("staff_number")
	if err := config.DB.First(&staff, "staff_number = ?", staffNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "No staff found with the specified staff number"})
		return
	}

	saveStaff := models.Staff{
		Surname:                    payload.Surname,
		Firstname:                  payload.FirstName,
		OtherName:                  payload.OtherName,
		Gender:                     payload.Gender,
		MainTelephoneNumber:        payload.MainTelephoneNumber,
		AlternativeTelephoneNumber: payload.AlternativeTelephoneNumber,
		MainEmailAddress:           payload.MainEmailAddress,
		AlternativeEmailAddress:    payload.AlternativeEmailAddress,
		JobTitleNumber:             payload.JobTitleNumber,
		CampusNumber:               payload.CampusNumber,
		SalaryScaleNumber:          payload.SalaryScaleNumber,
	}

	if err := config.DB.Model(&staff).Where("staff_number = ?", staffNumber).Updates(saveStaff).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update staff member info"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": staff})
}

func DeleteStaff(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	staffNumber := c.Param("staff_number")
	var staff models.Staff

	if err := config.DB.First(&staff, staffNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "No staff member found"})
		return
	}

	if err := config.DB.Delete(&staff).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete staff"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Staff deleted successfully"})
}
