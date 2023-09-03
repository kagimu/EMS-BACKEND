package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/hamsoftug/wholesomeEd/api/config"
	"github.com/hamsoftug/wholesomeEd/domain/auth/models"
	models2 "github.com/hamsoftug/wholesomeEd/domain/curriculum/strand/models"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func NewStrand(context *gin.Context) {

	type SubjectClusterInfo struct {
		Name          string `json:"name"`
		SubjectNumber string `json:"subject_number"`
	}

	var subjectClusterInfo SubjectClusterInfo
	err := context.BindJSON(&subjectClusterInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	name := subjectClusterInfo.Name
	subjectNumber := subjectClusterInfo.SubjectNumber

	userNumber := ""
	if len(context.Request.Header["User-Number"]) > 0 {
		userNumber = context.Request.Header["User-Number"][0]
	} else {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Not Authenticated!"})
		return
	}

	if name == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Name is required!"})
		return
	}

	if subjectNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Subject is required!"})
		return
	}

	var userData models.User
	//check if user is a valid user
	validateUser := config.DB.
		Where("number = ?", userNumber).
		First(&userData)
	if validateUser.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	var duplicateStrand models2.Strand
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("subject_number = ?", subjectNumber).
		First(&duplicateStrand).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Strand already exists"})
		return
	}

	number := uuid.NewV4().String()
	///save user data
	saveStrand := models2.Strand{
		Number:        number,
		Name:          name,
		SubjectNumber: subjectNumber,
		UserNumber:    userNumber,
	}
	config.DB.Save(&saveStrand)

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Subject cluster Added Successfully!", "data": saveStrand})
}

func FetchStrand(context *gin.Context) {

	subjectNumber := context.Query("subject_number")

	if subjectNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Subject is required!"})
		return
	}

	var strand []models2.Strand
	if err := config.DB.
		Where("subject_number = ?", subjectNumber).
		Order("created_at desc").
		Find(&strand).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No strand found!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": strand})
}
