package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staffMembers/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staffMembers/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func generateJobCode() string {
	year := time.Now().Year()
	prefix := fmt.Sprintf("JOBCODE%d", year)

	var lastPosition models.JobPosition
	result := config.DB.Where("job_code LIKE ?", prefix+"%").Order("job_code desc").First(&lastPosition)

	if result.Error != nil {
		return fmt.Sprintf("%s001", prefix)
	}

	lastCode := lastPosition.JobCode
	if len(lastCode) == len(prefix)+3 {
		numStr := lastCode[len(prefix):]
		num, err := strconv.Atoi(numStr)
		if err == nil {
			return fmt.Sprintf("%s%03d", prefix, num+1)
		}
	}

	var count int64
	config.DB.Model(&models.JobPosition{}).Where("job_code LIKE ?", prefix+"%").Count(&count)
	return fmt.Sprintf("%s%03d", prefix, count+1)
}

func CreateJobPosition(ctx *gin.Context) {
	var request requests.JobPositionRequest
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

	jobPosition := models.JobPosition{
		Number:           uuid.NewV4().String(),
		PostName:         request.PostName,
		JobCode:          generateJobCode(),
		SchoolNumber:     request.SchoolNumber,
		DepartmentNumber: request.DepartmentNumber,
		UserNumber:       ctx.Request.Header.Get("User-Number"),
	}

	if err := config.DB.Create(&jobPosition).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Job position not saved."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Job position created successfully!", "data": jobPosition})
}

func FetchJobPositions(ctx *gin.Context) {
	var jobPositions []models.JobPosition
	
	schoolNumber := ctx.Query("school_number")
	departmentNumber := ctx.Query("department_number")
	
	query := config.DB.Model(&models.JobPosition{}).Preload("School").Preload("Department")
	
	if schoolNumber != "" {
		query = query.Where("school_number = ?", schoolNumber)
	}
	
	if departmentNumber != "" {
		query = query.Where("department_number = ?", departmentNumber)
	}

	if err := query.Find(&jobPositions).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch job positions"})
		return
	}

	var responseData []map[string]interface{}
	for _, jp := range jobPositions {
		responseData = append(responseData, map[string]interface{}{
			"number":            jp.Number,
			"post_name":         jp.PostName,
			"job_code":          jp.JobCode,
			"school_number":     jp.SchoolNumber,
			"department_number": jp.DepartmentNumber,
			"department_name":   jp.Department.Name,
		})
	}

	if len(responseData) == 0 {
		responseData = []map[string]interface{}{}
	}

	ctx.JSON(http.StatusOK, gin.H{"data": responseData})
}

func EditJobPosition(ctx *gin.Context) {
	var request requests.EditJobPositionRequest
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

	var jobPosition models.JobPosition
	if err := config.DB.Where("number = ?", request.Number).First(&jobPosition).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Job position not found"})
		return
	}

	jobPosition.PostName = request.PostName
	jobPosition.SchoolNumber = request.SchoolNumber
	jobPosition.DepartmentNumber = request.DepartmentNumber
	
	if err := config.DB.Save(&jobPosition).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Job position not updated."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Job position updated successfully!"})
}

func DeleteJobPosition(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Job position number is required"})
		return
	}

	var jobPosition models.JobPosition
	if err := config.DB.Where("number = ?", id).First(&jobPosition).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Job position not found"})
		return
	}

	var relatedVacanciesCount int64
	config.DB.Model(&models.JobVacancy{}).Where("job_position_number = ?", id).Count(&relatedVacanciesCount)
	if relatedVacanciesCount > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Cannot delete job position. It is already associated with one or more job vacancies."})
		return
	}

	if err := config.DB.Delete(&jobPosition).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Job position not deleted. Try again."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Job position deleted successfully!"})
}
