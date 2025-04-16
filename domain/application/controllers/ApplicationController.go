package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/repository"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/pagination"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/parent/repositories"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/repositories"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
)

func FetchApplication(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User not found",
		})
		return
	}

	application, err := services.GetApplicationById(c.Param("application_number"))
	if err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Application not found"+err.Error()))
		return
	}

	c.JSON(http.StatusOK, application)
}

func FetchApplications(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User not found",
		})
		return
	}

	apps := make([]models.Application, 0)
	var paginationResponse pagination.PaginatedResponse

	var applications []models.Application

	//check for user account type and return applications accordingly
	if user.Type == "User" {

		if err := config.DB.Model(&models.Application{}).
			Scopes(
				filters.ByStatus(c.Query("status")),
				filters.BySchoolNumber(c.Query("school_number")),
				filters.ByCampusNumber(c.Query("campus_number")),
				pagination.Paginate(c, &paginationResponse),
			).
			Where("user_number = ?", user.Number).
			Preload("Student").Preload("Parent").
			Preload("PreviousSchools").
			Find(&applications).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "please contact support to fetch applications",
			})
			return
		}

	} else {
		schoolNumber := c.Query("school_number")
		if user.Type != constants.TypeClient {
			if schoolNumber == "" {
				c.JSON(http.StatusBadRequest, gin.H{"message": "Please select the school you want to view applications for"})
				return
			}
		}

		if err := config.DB.Model(&models.Application{}).
			Preload("Student").Preload("Parent").
			Preload("Student.CampusClass").
			Preload("Student.PreferredCampus").
			Preload("PreviousSchools").
			Joins("JOIN students ON students.application_number = applications.number").
			Scopes(
				filters.ByStatus(c.Query("applications.status")),
				filters.ByAnyField("students.school_number", schoolNumber),
				pagination.Paginate(c, &paginationResponse),
			).
			Find(&applications).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "application fetching failed",
			})
			return
		}

	}

	if applications != nil {
		apps = applications
	}
	paginationResponse.Content = apps

	c.JSON(http.StatusOK, paginationResponse)
}

func EditApplication(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User not found",
		})
		return
	}

	type ApplicationInfo struct {
		Status                string `json:"status"`
		ApplicationFeePayment bool   `json:"application_fee_payment"`
		InterviewMode         string `json:"interview_mode"`
	}

	var payload ApplicationInfo
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	if payload.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Status is required",
		})
		return
	}

	var application models.Application
	var applicationNumber = c.Param("application_number")

	if err := config.DB.First(&application, "number = ?", applicationNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "No application found with the specified number"})
		return
	}

	updateApplication := models.Application{
		Status:                payload.Status,
		ApplicationFeePayment: payload.ApplicationFeePayment,
		InterviewMode:         payload.InterviewMode,
	}

	if err := config.DB.Model(&application).Where("number = ?", applicationNumber).Updates(updateApplication).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update application"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    application,
		"message": "Application updated successfully",
	})
}

func FetchApplicationByNumber(c *gin.Context, applicationNumber string) (*models.Application, error) {

	var application models.Application
	if err := config.DB.Where("number = ?", applicationNumber).Preload("Parent").First(&application).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Application not found"})
		return nil, err
	}

	return &application, nil
}

func SubmitApplication(context *gin.Context) {

	var applicationNumber = context.Param("application_number")
	var application models.Application
	if err := config.DB.Where("number = ?", applicationNumber).First(&application).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Application not found"})
		return

	}

	application.Status = constants.StatusSubmitted

	if err := config.DB.Save(&application).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save application"})
		return
	}

	context.JSON(http.StatusOK, commons.Feedback("Application submitted successfully"))
}
func DeleteApplication(ctx *gin.Context) {
	txn := config.DB.Begin()
	applicationRepo := repository.NewApplicationRepositoryImpl(txn)
	parentRepo := repositories.NewParentRepositoryImpl(txn)
	studentRepo := repositories2.NewStudentRepositoryImpl(txn)

	application, err := applicationRepo.FindApplicationByNumber(ctx.Param("application_number"))
	if err != nil {
		txn.Rollback()
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Application not found"})
		return
	}

	parent, err := parentRepo.GetParentByApplicationNumber(application.Number)

	if err != nil {
		txn.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete parent"})
		return
	}

	if parent.Number != "" {
		if err := parentRepo.Delete(parent.Number); err != nil {
			txn.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete parent information"})
			return
		}
	}

	student, err := studentRepo.GetStudentByApplicationNumber(application.Number)
	if err != nil {
		txn.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete student information"})
		return
	}
	if student.Number != "" {
		if err := studentRepo.Delete(student.Number, txn); err != nil {
			txn.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete student information"})
			return
		}
	}

	if err := applicationRepo.DeleteDraftApplication(application.Number, txn); err != nil {
		txn.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete application"})
		return
	}

	txn.Commit()

	ctx.JSON(http.StatusOK, commons.Feedback("Application deleted successfully"))

}
