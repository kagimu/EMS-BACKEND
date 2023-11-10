package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/strand/api/responses"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/strand/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func NewStrand(context *gin.Context) {

	type SubjectClusterInfo struct {
		Code          string `json:"code"`
		Name          string `json:"name"`
		SubjectNumber string `json:"subject_number"`
	}

	var subjectClusterInfo SubjectClusterInfo
	err := context.BindJSON(&subjectClusterInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	code := subjectClusterInfo.Code
	name := subjectClusterInfo.Name
	subjectNumber := subjectClusterInfo.SubjectNumber

	userNumber := context.Request.Header["User-Number"][0]

	if code == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code is required!"})
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

	codeNo, err := strconv.ParseInt(code, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
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
		Code:          codeNo,
	}
	config.DB.Save(&saveStrand)

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Subject cluster Added Successfully!", "data": saveStrand})
}

func FetchStrand(context *gin.Context) {

	subjectNumber := context.Query("subject_number")

	var strand []models2.Strand
	var strand_ []responses2.StrandDataList
	if err := config.DB.
		Scopes(subject(subjectNumber)).
		Preload("Subject.SubjectCluster").
		Preload("Subject").
		Order("code asc").
		Find(&strand).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No strand found!"})
		return
	}

	for _, list := range strand {
		strand_ = append(strand_, responses2.StrandDataList{
			Number:               list.Number,
			Name:                 list.Name,
			Code:                 list.Code,
			UserNumber:           list.UserNumber,
			SubjectNumber:        list.SubjectNumber,
			SubjectName:          list.Subject.Name,
			SubjectClusterNumber: list.Subject.SubjectCluster.Number,
			SubjectClusterName:   list.Subject.SubjectCluster.Name,
		})
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": strand_})
}

func subject(subjectNumber string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if subjectNumber != "" {
			return db.Where("subject_number = ?", subjectNumber)
		}
		return db
	}
}
