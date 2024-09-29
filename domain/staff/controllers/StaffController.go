package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staff/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staff/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staff/services"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func AttachStaffToStream(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var payload requests.StreamTeacherRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	saveStreamTeacher := models.StreamTeacher{
		Number:               uuid.NewV4().String(),
		UserNumber:           payload.UserNumber,
		StreamNumber:         payload.StreamNumber,
		CampusClassNumber:    payload.CampusClassNumber,
		AddedBy:              user.Number,
		AssessmentAreaNumber: payload.AssessmentAreaNumber,
		SchoolNumber:         payload.SchoolNumber,
		CurriculumNumber:     payload.CurriculumNumber,
		CampusNumber:         payload.CampusNumber,
	}

	if err := services.AttachStaffToStream(saveStreamTeacher); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to attach staff to stream"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Staff has been attached to stream"})
}

func DetachStaffFromStream(c *gin.Context) {

	staffNumber := c.Query("user_number")
	streamNumber := c.Query("stream_number")

	if err := services.DetachStaffFromStream(staffNumber, streamNumber); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to detach staff from stream"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Staff has been detached from stream"})
}

func GetStaffOnStream(c *gin.Context) {
	streamNumber := c.Param("stream_number")

	staff, err := services.GetStaffOnStream(streamNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get staff on stream"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": staff})
}
