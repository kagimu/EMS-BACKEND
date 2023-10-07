package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewAcademicYearController(context *gin.Context) {

	type AcademicYearInfo struct {
		AcademicYearName string `json:"academic_year_name"`
		NoOfStudyPeriod  string `json:"no_of_study_period"`
		AcademicYear     string `json:"academic_year"`
		LevelNumber      string `json:"level_number"`
	}

	var academicYearInfo AcademicYearInfo
	err := context.BindJSON(&academicYearInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	academicYearName := academicYearInfo.AcademicYearName
	noOfStudyPeriod := academicYearInfo.NoOfStudyPeriod
	academicYear := academicYearInfo.AcademicYear
	levelNumber := academicYearInfo.LevelNumber
	userNumber := context.Request.Header["User-Number"][0]

	/*if name == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Name is required!"})
		return
	}

	if subtopicNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Sub topic is required!"})
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

	var duplicateConcept models3.Concept
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("sub_topic_number = ?", subtopicNumber).
		First(&duplicateConcept).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Sub topic already exists"})
		return
	}

	number := uuid.NewV4().String()
	///save user data
	saveConcept := models3.Concept{
		Number:         number,
		Name:           name,
		SubTopicNumber: subtopicNumber,
		UserNumber:     userNumber,
	}
	config.DB.Save(&saveConcept)

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Concept Added Successfully!", "data": saveConcept})*/
}
