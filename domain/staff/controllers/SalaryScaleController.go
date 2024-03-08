package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staff/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staff/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateSalaryScale(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.SalaryScaleRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process payload"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	saveSalaryScale := models.SalaryScale{
		Number:        uuid.NewV4().String(),
		Name:          payload.Name,
		Description:   payload.Description,
		MinimumAmount: payload.MinimumAmount,
		MaximumAmount: payload.MaximumAmount,
		UserNumber:    user.Number,
	}

	if err := config.DB.Create(&saveSalaryScale).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save salary scale"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Salary scale created successfully"})
}

func FetchSalaryScales(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var salaryScales []models.SalaryScale
	if err := config.DB.Find(&salaryScales).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch salary scales"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": salaryScales})
}

func EditSalaryScale(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.UpdateSalaryScaleRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process payload"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var salaryScale models.SalaryScale
	salaryScaleNumber := c.Query("salary_scale_number")
	if err := config.DB.Where("number = ?", salaryScaleNumber).First(&salaryScale).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Salary scale not found"})
		return
	}

	updateSalaryScale := models.SalaryScale{
		Name:          payload.Name,
		Description:   payload.Description,
		MinimumAmount: payload.MinimumAmount,
		MaximumAmount: payload.MaximumAmount,
	}

	if err := config.DB.Model(&updateSalaryScale).Where("number = ?", salaryScaleNumber).Updates(updateSalaryScale).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update salry scale"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": salaryScale})
}
