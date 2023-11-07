package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	topicModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/topic/models"
	subjectModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func NewSubject(context *gin.Context) {

	type SubjectInfo struct {
		Code                 string `json:"code"`
		Name                 string `json:"name"`
		SubjectClusterNumber string `json:"subject_cluster_number"`
	}

	var subjectInfo SubjectInfo
	err := context.BindJSON(&subjectInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	code := subjectInfo.Code
	name := subjectInfo.Name
	subjectClusterNumber := subjectInfo.SubjectClusterNumber

	userNumber := context.Request.Header["User-Number"][0]

	if code == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code is required!"})
		return
	}

	if name == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Name is required!"})
		return
	}

	if subjectClusterNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Level is required!"})
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

	var duplicateSubject subjectModel.Subject
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("level_number = ?", subjectClusterNumber).
		First(&duplicateSubject).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Subject already exists"})
		return
	}

	number := uuid.NewV4().String()
	///save user data
	saveSubject := subjectModel.Subject{
		Number:               number,
		Name:                 name,
		UserNumber:           userNumber,
		SubjectClusterNumber: subjectClusterNumber,
		Code:                 codeNo,
	}

	config.DB.Save(&saveSubject)
	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Campus Added Successfully!", "data": saveSubject})
}

func FetchSubject(context *gin.Context) {

	subjectClusterNumber := context.Query("subject_cluster_number")

	/*	if subjectClusterNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Subject cluster is required!"})
		return
	}*/

	var subjects []subjectModel.Subject
	if err := config.DB.
		Scopes(subjectCluster(subjectClusterNumber)).
		Order("code asc").
		Find(&subjects).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Subject!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": subjects})
}

func subjectCluster(subjectClusterNumber string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if subjectClusterNumber != "" {
			return db.Where("subject_cluster_number = ?", subjectClusterNumber)
		}
		return db
	}
}

type SubjectData struct {
	Number string `json:"number"`
	Name   string `json:"name"`
}

func FetchSubjectByClass(context *gin.Context) {

	classNumber := context.Query("class_number")

	var topic []topicModel.Topic
	if err := config.DB.
		Where("topics.class_number = ?", classNumber).
		Joins("join strands on strands.number = topics.strand_number").
		Preload("Strand.Subject").
		Preload("Strand").
		Order("created_at desc").
		Find(&topic).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No topic found!"})
		return
	}

	var subjectData []SubjectData
	for _, item := range topic {
		subjectData = append(subjectData, SubjectData{
			Number: item.Strand.Subject.Number,
			Name:   item.Strand.Subject.Name,
		})
	}

	subjectData = removeDuplicates(subjectData)

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": subjectData})
}

func removeDuplicates(data []SubjectData) []SubjectData {

	subjectMap := map[SubjectData]bool{}
	var result []SubjectData

	for _, subject := range data {
		if !subjectMap[subject] {
			subjectMap[subject] = true
			result = append(result, subject)
		}
	}
	return result
}

func FetchSubjectsByClassLevelNumber(ctx *gin.Context) {

	classLevelNumber := ctx.Param("class_level_number")

	var subjects []classModels.ClassLevelSubject
	config.DB.
		Where("class_level_number = ?", classLevelNumber).
		Preload("Subject").
		Find(&subjects)
	ctx.JSON(http.StatusOK, gin.H{"data": subjects})

}
