package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/api/request"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func CreateStudent(c *gin.Context) {
	var payload request.StudentRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Internal server error!",
		})
		return
	}

	payload.Prepare()
	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to validate input",
		})
		return
	}

	userNumber := c.Request.Header["User-Number"][0]
	sourceOfSchoolInfoNumber := c.Param("source_of_school_info_number")

	var user authModels.User
	if err := config.DB.
		Where("number = ?", userNumber).
		First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	studentNumber := uuid.NewV4().String()
	year := strconv.Itoa(time.Now().Year())
	month := time.Now().Format("01")

	saveAPN := models.APN{
		Year:  year,
		Month: month,
	}

	if err := config.DB.Create(&saveAPN).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "failed to save APN",
			"data":    saveAPN,
		})
		return
	}

	savedID := saveAPN.ID + 1
	result := strings.Join([]string{saveAPN.Year, "01", fmt.Sprintf("%04s", strconv.FormatUint(savedID, 10))}, "")

	saveStudent := models.Student{
		Number:                    studentNumber,
		APN:                       result,
		Surname:                   payload.Surname,
		OtherName:                 payload.OtherName,
		Age:                       payload.Age,
		Gender:                    payload.Gender,
		Reason:                    payload.Reason,
		PreferredTypeOfEnrollment: payload.PreferredTypeOfEnrollment,
		SecondName:                payload.SecondName,
		UserNumber:                userNumber,
		SourceOfSchoolInfoNumber:  sourceOfSchoolInfoNumber,
	}

	if err := config.DB.Create(&saveStudent).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"Details": saveAPN.ID,
			"message": "failed to save student",
		})
		return
	}

	// Save previous schools
	var previousSchools []models.PreviousSchools
	var schoolNames []string

	for _, schoolName := range payload.PreviousSchools {
		schoolNames = append(schoolNames, schoolName)
	}

	pre := models.PreviousSchools{
		Number:        uuid.NewV4().String(),
		Name:          strings.Join(schoolNames, ""),
		StudentNumber: saveStudent.Number,
		UserNumber:    userNumber,
	}

	previousSchools = append(previousSchools, pre)

	if err := config.DB.Create(&previousSchools).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"details": err.Error(),
			"message": "failed to save previous school",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "student created successfully",
		"data":    saveStudent,
	})
}

func UpdateStudent(c *gin.Context) {
	student := models.Student{}
	studentAPN := c.Param("id")

	err := config.DB.First(&student, studentAPN).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "no student found",
		})
		return
	}

	if err := c.BindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "invalid request",
		})
		return
	}

	if err := config.DB.Save(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "failed to update student",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": student})

}

func DeleteStudent(c *gin.Context) {
	studentAPN := c.Param("id")
	student := models.Student{}

	err := config.DB.First(&student, studentAPN).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "no student found",
		})
		return
	}

	if err := config.DB.Delete(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "failed to delete student",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "student deleted successfully"})
}

func FetchStudentByAPN(c *gin.Context) {
	var student models.Student
	studentAPN := c.Param("apn")

	if studentAPN == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Internal server error!",
		})
		return
	}

	userNumber := c.Request.Header["User-Id"][0]

	var user authModels.User
	if err := config.DB.
		Where("number = ?", userNumber).
		First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	err := config.DB.First(&student, "apn = ?", studentAPN).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  http.StatusNotFound,
				"message": "no student found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": student})
}

func FetchAllStudents(c *gin.Context) {
	var students []models.Student

	err := config.DB.Find(&students).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "failed to fetch students",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": students})
}
