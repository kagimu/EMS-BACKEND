package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/departments/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/departments/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateDepartment(ctx *gin.Context) {
	var request requests.DepartmentRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Please structure the response according to specifications!"})
		return
	}

	err = request.Validate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var existingCode models.Department
	if err := config.DB.Where("school_number = ? AND code = ?", request.SchoolNumber, request.Code).First(&existingCode).Error; err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Department with this code already exists in this school!"})
		return
	}

	var existingName models.Department
	if err := config.DB.Where("school_number = ? AND LOWER(name) = LOWER(?)", request.SchoolNumber, request.Name).First(&existingName).Error; err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Department with this name already exists in this school!"})
		return
	}

	department := models.Department{
		Number:       uuid.NewV4().String(),
		Name:         request.Name,
		Code:         request.Code,
		SchoolNumber: request.SchoolNumber,
		UserNumber:   ctx.Request.Header.Get("User-Number"),
	}

	if err := config.DB.Create(&department).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Department not saved."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Department created successfully!", "data": department})
}

func FetchAllDepartments(ctx *gin.Context) {
	var departments []models.Department
	
	// You might want to filter by school if school_number is provided in query params
	schoolNumber := ctx.Query("school_number")
	
	query := config.DB.Model(&models.Department{}).Preload("School")
	if schoolNumber != "" {
		query = query.Where("school_number = ?", schoolNumber)
	}

	if err := query.Find(&departments).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch departments"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": departments})
}

func EditDepartment(ctx *gin.Context) {
	var request requests.EditDepartmentRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Please structure the response according to specifications!"})
		return
	}

	err = request.Validate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var department models.Department
	if err := config.DB.Where("number = ?", request.Number).First(&department).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Department not found"})
		return
	}

	var existingCode models.Department
	if err := config.DB.Where("school_number = ? AND code = ? AND number != ?", request.SchoolNumber, request.Code, request.Number).First(&existingCode).Error; err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Department with this code already exists in this school!"})
		return
	}

	var existingName models.Department
	if err := config.DB.Where("school_number = ? AND LOWER(name) = LOWER(?) AND number != ?", request.SchoolNumber, request.Name, request.Number).First(&existingName).Error; err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Department with this name already exists in this school!"})
		return
	}

	department.Name = request.Name
	department.Code = request.Code
	department.SchoolNumber = request.SchoolNumber
	
	if err := config.DB.Save(&department).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Department not updated."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Department updated successfully!"})
}

func DeleteDepartment(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Department number is required"})
		return
	}

	var department models.Department
	if err := config.DB.Where("number = ?", id).First(&department).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Department not found"})
		return
	}

	if err := config.DB.Delete(&department).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Department not deleted. Try again."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Department deleted successfully!"})
}
