package controllers

import (
	services2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/services"
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
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input- " + err.Error()})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	campusClass, err := services2.ResolveCampusClass(payload.ClassNumber, payload.CampusNumber, payload.CurriculumNumber, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to resolve campus class " + err.Error()})
		return
	}

	for _, val := range payload.TeacherIds {

		for _, stream := range payload.Streams {

			saveStreamTeacher := models.StreamTeacher{
				Number:               uuid.NewV4().String(),
				UserNumber:           val,
				StreamNumber:         stream,
				CampusClassNumber:    campusClass.Number,
				AddedBy:              user.Number,
				AssessmentAreaNumber: payload.AssessmentAreaNumber,
				SchoolNumber:         payload.SchoolNumber,
				CurriculumNumber:     payload.CurriculumNumber,
				CampusNumber:         payload.CampusNumber,
			}
			if err := services.AttachStaffToStream(&saveStreamTeacher); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to attach staff to stream"})
				return
			}
		}

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

	if c.Query("school_number") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Please select the school to view the list of teachers"})
		return
	}

	staff, err := services.GetStaffStreamsAndClasses(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get staff on stream"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": staff})
}
