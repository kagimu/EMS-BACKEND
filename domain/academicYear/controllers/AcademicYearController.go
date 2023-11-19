package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/api/responses"
	academicYearModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	campusModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	eosModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strconv"
	"time"
)

func NewAcademicYear(context *gin.Context) {

	type AcademicYearInfo struct {
		AcademicYearName string   `json:"academic_year_name"`
		NoOfStudyPeriod  string   `json:"no_of_study_period"`
		AcademicYear     []string `json:"academic_year"`
		LevelNumber      string   `json:"level_number"`
		Code             string   `json:"code"`
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
	code := academicYearInfo.Code
	userNumber := context.Request.Header["User-Number"][0]

	if academicYearName == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Academic Year Name is required!"})
		return
	}

	if noOfStudyPeriod == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Number Of Study Periods is required!"})
		return
	}

	if len(academicYear) == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Academic Year is required!"})
		return
	}

	if levelNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Level Year is required!"})
		return
	}

	if code == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code is required!"})
		return
	}

	if code == "0" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code must be greater than zero!"})
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

	var duplicateCode academicYearModels.AcademicYear
	if duplicate := config.DB.
		Where("code = ?", code).
		First(&duplicateCode).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code already exists"})
		return
	}

	var duplicateAcademicYear academicYearModels.AcademicYear
	if duplicate := config.DB.
		Where("academic_year_name = ?", academicYearName).
		Where("level_number = ?", levelNumber).
		First(&duplicateAcademicYear).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Academic Year already exists"})
		return
	}

	layout := "2006-01-02T15:04:05.000Z"
	startDate, errorDate := time.Parse(layout, academicYear[0])

	if errorDate != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": errorDate.Error()})
		return
	}

	endDate, errorDate := time.Parse(layout, academicYear[1])

	if errorDate != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": errorDate.Error()})
		return
	}

	noOfStudyPeriodData, err := strconv.ParseInt(noOfStudyPeriod, 10, 32)
	if err != nil {
		// Handle the error if the conversion fails
		//fmt.Println("Error:", err)
		//return
	}
	config.DB.Begin()

	// create EOS calendar
	var generalCalendar = eosModels.GeneralCalendar{}

	generalCalendar.Number = uuid.NewV4().String()
	generalCalendar.EntityNumber = "TBD"
	generalCalendar.StartDate = startDate
	generalCalendar.EndDate = endDate
	generalCalendar.Name = academicYearName
	generalCalendar.UserNumber = userNumber
	generalCalendar.Description = "Academic Year"

	result := config.DB.Create(&generalCalendar)

	if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "failed to create calendar with " + result.Error.Error()})
		config.DB.Rollback()
		return

	}

	number := uuid.NewV4().String()
	///save user data
	saveAcademicYear := academicYearModels.AcademicYear{
		Number:                number,
		Code:                  codeNo,
		CalendarNumber:        generalCalendar.Number,
		AcademicYearName:      academicYearName,
		NoOfStudyPeriod:       int32(noOfStudyPeriodData),
		AcademicYearStartDate: startDate,
		AcademicYearEndDate:   endDate,
		UserNumber:            userNumber,
		LevelNumber:           levelNumber,
	}
	save := config.DB.Save(&saveAcademicYear)

	if save.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "failed to save academic year"})
		config.DB.Rollback()
		return
	}

	config.DB.Commit()

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Academic year Added Successfully!", "data": saveAcademicYear})
}

func FetchAcademicYearInfo(context *gin.Context) {

	levelNumber := context.Query("level")

	if levelNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Level is required!"})
		return
	}

	var academicYearInfo []academicYearModels.AcademicYear
	var academicYearInfo_ []responses.AcademicYearList
	if err := config.DB.
		//Select("name", "number").
		Where("level_number = ?", levelNumber).
		Order("code asc").
		Find(&academicYearInfo).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Academic year found"})
		return
	}

	// Populate the selected columns from the fetched data
	for _, academicList := range academicYearInfo {
		academicYearInfo_ = append(academicYearInfo_, responses.AcademicYearList{
			Number:                academicList.Number,
			Code:                  academicList.Code,
			AcademicYearName:      academicList.AcademicYearName,
			NoOfStudyPeriod:       academicList.NoOfStudyPeriod,
			AcademicYearStartDate: (academicList.AcademicYearStartDate).Format("2-Jan-2006 15:04:05"),
			AcademicYearEndDate:   (academicList.AcademicYearEndDate).Format("2-Jan-2006 15:04:05"),
			UserNumber:            academicList.UserNumber,
			LevelNumber:           academicList.LevelNumber,
		})
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": academicYearInfo_})

}

func FetchAcademicByClass(context *gin.Context) {

	classNumber := context.Query("class_number")

	if classNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Class number is required!"})
		return
	}

	var classes classModels.Class
	if err := config.DB.
		Where("number = ?", classNumber).
		First(&classes).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Class!"})
		return
	}

	var levelsInfo []campusModels.Level
	if err := config.DB.
		Where("name = ?", classes.LevelNumber).
		Find(&levelsInfo).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Level!"})
		return
	}

	var academicYearInfo_ []responses.AcademicYearListByClass

	for _, lvlList := range levelsInfo {

		var academicYearInfo []academicYearModels.AcademicYear
		if err := config.DB.
			//Select("name", "number").
			Where("level_number = ?", lvlList.Number).
			Preload("Level").
			Order("code asc").
			Find(&academicYearInfo).Error; err != nil {
			context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Academic year found"})
			return
		}

		// Populate the selected columns from the fetched data
		for _, academicList := range academicYearInfo {
			academicYearInfo_ = append(academicYearInfo_, responses.AcademicYearListByClass{
				Number:                academicList.Number,
				Code:                  academicList.Code,
				AcademicYearName:      academicList.AcademicYearName,
				NoOfStudyPeriod:       academicList.NoOfStudyPeriod,
				AcademicYearStartDate: (academicList.AcademicYearStartDate).Format("2-Jan-2006 15:04:05"),
				AcademicYearEndDate:   (academicList.AcademicYearEndDate).Format("2-Jan-2006 15:04:05"),
				UserNumber:            academicList.UserNumber,
				LevelNumber:           academicList.LevelNumber,
				LevelName:             academicList.Level.LearningLevelNumber,
			})
		}
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": academicYearInfo_})

}
