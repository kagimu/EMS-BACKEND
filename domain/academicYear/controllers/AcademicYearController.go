package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/api/responses"
	academicYearModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	eosModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

func NewAcademicYear(context *gin.Context) {

	type AcademicYearInfo struct {
		AcademicYearName      string `json:"academic_year_name"`
		NoOfStudyPeriod       string `json:"no_of_study_period"`
		AcademicYearStartDate string `json:"academic_year_start_date"`
		AcademicYearEndDate   string `json:"academic_year_end_date"`
		LevelNumber           string `json:"level_number"`
		Code                  string `json:"code"`
	}

	var academicYearInfo AcademicYearInfo
	err := context.BindJSON(&academicYearInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	academicYearName := academicYearInfo.AcademicYearName
	noOfStudyPeriod := academicYearInfo.NoOfStudyPeriod
	academicYearStartDate := academicYearInfo.AcademicYearStartDate
	academicYearEndDate := academicYearInfo.AcademicYearEndDate
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

	if academicYearStartDate == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Start date is required!"})
		return
	}

	if academicYearEndDate == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "End date is required!"})
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

	var duplicateCode academicYearModel.AcademicYear
	if duplicate := config.DB.
		Where("code = ?", code).
		Where("campus_level_number = ?", levelNumber).
		First(&duplicateCode).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code already exists"})
		return
	}

	var duplicateAcademicYear academicYearModel.AcademicYear
	if duplicate := config.DB.
		Where("academic_year_name = ?", academicYearName).
		Where("campus_level_number = ?", levelNumber).
		First(&duplicateAcademicYear).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Academic Year already exists"})
		return
	}

	layout := "2006-01-02T15:04:05.000Z"
	startDate, errorDate := time.Parse(layout, academicYearStartDate)

	if errorDate != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": errorDate.Error()})
		return
	}

	endDate, errorDate := time.Parse(layout, academicYearEndDate)

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

	tx := config.DB.Begin()

	// create EOS calendar
	var generalCalendar = eosModels.GeneralCalendar{
		Number:       uuid.NewV4().String(),
		EntityNumber: "TBD",
		StartDate:    startDate,
		EndDate:      endDate,
		Name:         academicYearName,
		UserNumber:   userNumber,
		Description:  "Academic Year",
	}

	if err := tx.Create(&generalCalendar).Error; err != nil {
		// If any error occurs, rollback the transaction
		tx.Rollback()
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "failed to create calendar with " + err.Error()})
		return
	}

	///save user data
	saveAcademicYear := academicYearModel.AcademicYear{
		Number:                uuid.NewV4().String(),
		Code:                  codeNo,
		CalendarNumber:        generalCalendar.Number,
		AcademicYearName:      academicYearName,
		NoOfStudyPeriod:       int32(noOfStudyPeriodData),
		AcademicYearStartDate: startDate,
		AcademicYearEndDate:   endDate,
		UserNumber:            userNumber,
		CampusLevelNumber:     levelNumber,
	}

	if err := tx.Create(&saveAcademicYear).Error; err != nil {
		// If any error occurs, rollback the transaction
		tx.Rollback()
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save academic year"})
		return
	}

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create academic year!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Academic year Added Successfully!", "data": saveAcademicYear})
}

func FetchAcademicYearInfo(context *gin.Context) {

	levelNumber := context.Query("level")

	var academicYearInfo []academicYearModel.AcademicYear
	var academicYearInfo_ []responses.AcademicYearList
	if err := config.DB.
		Scopes(academicYearFilter(levelNumber)).
		Preload("CampusLevel.CurriculumLevel.Level").
		Preload("CampusLevel.CurriculumLevel").
		Preload("CampusLevel").
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
			LevelNumber:           academicList.CampusLevelNumber,
			LevelName:             academicList.CampusLevel.CurriculumLevel.Level.Name,
		})
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": academicYearInfo_})

}

func academicYearFilter(levelNumber string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if levelNumber != "" {
			return db.Where("campus_level_number = ?", levelNumber)
		}
		return db
	}
}

func FetchAcademicByClass(context *gin.Context) {

	classNumber := context.Query("class_number")

	if classNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Class number is required!"})
		return
	}

	var classes models2.CampusClass
	if err := config.DB.
		Where("number = ?", classNumber).
		First(&classes).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Class!"})
		return
	}

	var academicYearInfo []academicYearModel.AcademicYear
	var academicYearInfo_ []responses.AcademicYearList
	if err := config.DB.
		Scopes(academicYearFilter(classes.CampusLevelNumber)).
		Preload("CampusLevel.CurriculumLevel.Level").
		Preload("CampusLevel.CurriculumLevel").
		Preload("CampusLevel").
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
			LevelNumber:           academicList.CampusLevelNumber,
			LevelName:             academicList.CampusLevel.CurriculumLevel.Level.Name,
		})
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": academicYearInfo_})

}
