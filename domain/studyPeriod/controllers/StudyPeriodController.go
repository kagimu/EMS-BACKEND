package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	academicYearModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	eosModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/api/responses"
	studyPeriodModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

func NewStudyPeriod(context *gin.Context) {

	type StudyPeriodInfo struct {
		Code                 string   `json:"code"`
		StudyPeriodName      string   `json:"study_period_name"`
		AcademicYearNumber   string   `json:"academic_year_number"`
		StudyPeriodDateRange []string `json:"study_period_date_range"`
	}

	var studyPeriodInfo StudyPeriodInfo
	err := context.BindJSON(&studyPeriodInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	code := studyPeriodInfo.Code
	studyPeriodName := studyPeriodInfo.StudyPeriodName
	academicYearNumber := studyPeriodInfo.AcademicYearNumber
	studyPeriodDateRange := studyPeriodInfo.StudyPeriodDateRange
	userNumber := context.Request.Header["User-Number"][0]

	if academicYearNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Academic Year Name is required!"})
		return
	}

	if studyPeriodName == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Study Periods is required!"})
		return
	}

	if len(studyPeriodDateRange) == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Study period range is required!"})
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

	var duplicateCode studyPeriodModel.StudyPeriod
	if duplicate := config.DB.
		Where("code = ?", code).
		First(&duplicateCode).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code already exists"})
		return
	}

	var duplicateStudyPeriod studyPeriodModel.StudyPeriod
	if duplicate := config.DB.
		Where("study_period_name = ?", studyPeriodName).
		Where("academic_year_number = ?", academicYearNumber).
		First(&duplicateStudyPeriod).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Study Period already exists"})
		return
	}

	layout := "2006-01-02T15:04:05.000Z"
	startDate, errorDate := time.Parse(layout, studyPeriodDateRange[0])

	if errorDate != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": errorDate.Error()})
		return
	}

	endDate, errorDate := time.Parse(layout, studyPeriodDateRange[1])

	if errorDate != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": errorDate.Error()})
		return
	}

	//validate academic year scope validation (Boundary)
	//academicYearNumber

	var academicYearData academicYearModels.AcademicYear
	//check if user is a valid user
	validateAcademicYr := config.DB.
		Where("number = ?", academicYearNumber).
		First(&academicYearData)
	if validateAcademicYr.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Academic Year not found"})
		return
	}

	startDateAcademicYr := academicYearData.AcademicYearStartDate
	endDateAcademicYr := academicYearData.AcademicYearEndDate

	/*fmt.Println(startDateAcademicYr)
	fmt.Println("<==>")
	fmt.Println(startDate)*/

	//validate start date
	startDate = startDate.Add(3 * time.Hour) // Converting UTC to EAT for comparison
	if startDateAcademicYr.After(startDate) {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Start date on Study Period is out of bounds"})
		return
	}

	//validate end date
	endDate = endDate.Add(3 * time.Hour) // Converting UTC to EAT for comparison
	if endDateAcademicYr.Before(endDate) {
		context.JSON(http.StatusBadRequest, gin.H{"message": "End date on Study Period is out of bounds"})
		return
	}

	tx := config.DB.Begin()

	// create calendar event

	var calendarEvent = eosModels.CalendarEvent{}

	// Format the time as "15:04" (24-hour format)
	startTime := startDate.Format("15:04")
	endTime := endDate.Format("15:04")

	calendarEvent.Number = uuid.NewV4().String()
	calendarEvent.StartDate = startDate
	calendarEvent.EndDate = endDate
	calendarEvent.StartTime = startTime
	calendarEvent.EndTime = endTime
	calendarEvent.CalendarNumber = academicYearData.CalendarNumber
	calendarEvent.Name = studyPeriodName

	if err := tx.Create(&calendarEvent).Error; err != nil {

		tx.Rollback()

		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create calendar event"})
		return
	}

	number := uuid.NewV4().String()
	///save user data
	saveStudyPeriod := studyPeriodModel.StudyPeriod{
		Number:               number,
		Code:                 codeNo,
		StudyPeriodName:      studyPeriodName,
		StudyPeriodStartDate: startDate,
		StudyPeriodEndDate:   endDate,
		AcademicYearNumber:   academicYearNumber,
		UserNumber:           userNumber,
		CalendarEventNumber:  calendarEvent.Number,
	}
	//config.DB.Save(&saveStudyPeriod)

	if err := tx.Create(&saveStudyPeriod).Error; err != nil {
		// If any error occurs, rollback the transaction
		tx.Rollback()
		// Handle or return the error as needed
		//context.JSON(http.StatusBadRequest, gin.H{"message": err})
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create Study Period"})
		return
	}

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create Study Period"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Study Period Successfully!", "data": saveStudyPeriod})
}

func FetchStudyPeriod(context *gin.Context) {

	levelNumber := context.Query("level")
	academicYearNumber := context.Query("academic_year_number")

	if levelNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Level is required!"})
		return
	}

	var studyPeriodInfo []studyPeriodModel.StudyPeriod
	var studyPeriodInfo_ []responses2.StudyPeriodList
	if err := config.DB.
		Scopes(academicYearFilter(academicYearNumber)).
		Where("academic_years.campus_level_number = ?", levelNumber).
		Joins("join academic_years on academic_years.number = study_periods.academic_year_number").
		Preload("AcademicYear").
		Order("code asc").
		Find(&studyPeriodInfo).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Academic year found"})
		return
	}

	// Populate the selected columns from the fetched data
	for _, academicList := range studyPeriodInfo {
		studyPeriodInfo_ = append(studyPeriodInfo_, responses2.StudyPeriodList{
			Number:               academicList.Number,
			Code:                 academicList.Code,
			StudyPeriodName:      academicList.StudyPeriodName,
			StudyPeriodStartDate: (academicList.StudyPeriodStartDate).Format("2-Jan-2006 15:04:05"),
			StudyPeriodEndDate:   (academicList.StudyPeriodEndDate).Format("2-Jan-2006 15:04:05"),
			AcademicYearNumber:   academicList.AcademicYearNumber,
			AcademicYearName:     academicList.AcademicYear.AcademicYearName,
			UserNumber:           academicList.UserNumber,
		})
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": studyPeriodInfo_})

}

func academicYearFilter(academicYearNumber string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if academicYearNumber != "" {
			return db.Where("academic_year_number = ?", academicYearNumber)
		}
		return db
	}
}
