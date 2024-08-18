package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	repositories3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/repositories"
	"net/http"
	"strconv"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	requests2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/api/responses"
	academicYearModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	responses3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/api/responses"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	eosModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	schemeOfWorkModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/schemeOfWork/models"
	levelModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/models"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/api/responses"
	studyPeriodModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func NewAcademicYear(context *gin.Context) {

	var academicYearInfo requests2.AcademicYearInfo
	err := context.BindJSON(&academicYearInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	academicYearInfo.Prepare()

	if err := academicYearInfo.Validate(); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	academicYearName := academicYearInfo.AcademicYearName
	noOfStudyPeriod := academicYearInfo.NoOfStudyPeriod
	academicYearStartDate := academicYearInfo.AcademicYearStartDate
	academicYearEndDate := academicYearInfo.AcademicYearEndDate
	levelNumber := academicYearInfo.LevelNumber
	code := academicYearInfo.Code
	userNumber := context.Request.Header["User-Number"][0]

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
		fmt.Println("Validation failed:", err)
		return
	}

	entityNumber, _ := commons.GetUserEntityNumber(context)
	if entityNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Set current entity"})
		return
	}

	tx := config.DB.Begin()

	// create EOS calendar
	var generalCalendar = eosModels.GeneralCalendar{
		Number:       uuid.NewV4().String(),
		EntityNumber: entityNumber,
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

func EditAcademicYear(context *gin.Context) {

	var academicYearInfo requests2.AcademicYearInfo
	err := context.BindJSON(&academicYearInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	academicYearInfo.Prepare()

	if err := academicYearInfo.Validate(); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if academicYearInfo.AcademicYearNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Academic year is required!"})
		return
	}

	academicYearName := academicYearInfo.AcademicYearName
	noOfStudyPeriod := academicYearInfo.NoOfStudyPeriod
	academicYearStartDate := academicYearInfo.AcademicYearStartDate
	academicYearEndDate := academicYearInfo.AcademicYearEndDate
	levelNumber := academicYearInfo.LevelNumber
	code := academicYearInfo.Code
	academicYearNumber := academicYearInfo.AcademicYearNumber
	userNumber := context.Request.Header["User-Number"][0]

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
		Where("number != ?", academicYearNumber).
		First(&duplicateCode).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code already exists"})
		return
	}

	var duplicateAcademicYear academicYearModel.AcademicYear
	if duplicate := config.DB.
		Where("academic_year_name = ?", academicYearName).
		Where("campus_level_number = ?", levelNumber).
		Where("number != ?", academicYearNumber).
		First(&duplicateAcademicYear).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Academic Year already exists"})
		return
	}

	format := "2006-01-02T15:04:05.000Z"
	startDate, errorDate := time.Parse(format, academicYearStartDate)

	if errorDate != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": errorDate.Error()})
		return
	}

	endDate, errorDate := time.Parse(format, academicYearEndDate)

	if errorDate != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": errorDate.Error()})
		return
	}

	noOfStudyPeriodData, err := strconv.ParseInt(noOfStudyPeriod, 10, 32)
	if err != nil {
		fmt.Println("Validation failed:", err)
		return
	}

	entityNumber, _ := commons.GetUserEntityNumber(context)
	if entityNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Set current entity"})
		return
	}

	tx := config.DB.Begin()

	//var calendarEvent eosModels.CalendarEvent
	var generalCalendar eosModels.GeneralCalendar
	updateCalendar := tx.Model(&generalCalendar).
		Where("number = ?", academicYearNumber).
		Updates(eosModels.GeneralCalendar{EntityNumber: entityNumber, StartDate: startDate, EndDate: endDate, Name: academicYearName, UserNumber: userNumber, Description: "Academic Year"})
	if updateCalendar.Error != nil {
		tx.Rollback()
		context.JSON(http.StatusNotFound, gin.H{"message": "Failed to update!"})
		return
	}

	// create EOS calendar
	/*	var generalCalendar = eosModels.GeneralCalendar{
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
		}*/

	//var calendarEvent eosModels.CalendarEvent
	var academicYearUpdate academicYearModel.AcademicYear
	updateAcademicYear := tx.Model(&academicYearUpdate).
		Where("number = ?", academicYearNumber).
		Updates(academicYearModel.AcademicYear{Code: codeNo, CalendarNumber: generalCalendar.Number, AcademicYearName: academicYearName, NoOfStudyPeriod: int32(noOfStudyPeriodData), AcademicYearStartDate: startDate, AcademicYearEndDate: endDate, UserNumber: userNumber, CampusLevelNumber: levelNumber})
	if updateAcademicYear.Error != nil {
		tx.Rollback()
		context.JSON(http.StatusNotFound, gin.H{"message": "Failed to update!"})
		return
	}

	///save user data
	/*saveAcademicYear := academicYearModel.AcademicYear{
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
	}*/

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create academic year!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Academic year Updated Successfully!", "data": academicYearUpdate})
}

func FetchAcademicYearData(context *gin.Context) {
	academicYearNumber := context.Query("academic_year_number")

	if academicYearNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Academic Year is required!"})
		return
	}

	var academicYearInfo academicYearModel.AcademicYear

	if err := config.DB.
		Where("number = ?", academicYearNumber).
		Preload("CampusLevel.CurriculumLevel.Level").
		Preload("CampusLevel.CurriculumLevel").
		Preload("CampusLevel").
		Order("code asc").
		First(&academicYearInfo).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Academic year found"})
		return
	}

	// Populate the selected columns from the fetched data
	academicYearInfo_ := responses.AcademicYearList{
		Number:                academicYearInfo.Number,
		Code:                  academicYearInfo.Code,
		AcademicYearName:      academicYearInfo.AcademicYearName,
		NoOfStudyPeriod:       academicYearInfo.NoOfStudyPeriod,
		AcademicYearStartDate: (academicYearInfo.AcademicYearStartDate).Format("2-Jan-2006 15:04:05"),
		AcademicYearEndDate:   (academicYearInfo.AcademicYearEndDate).Format("2-Jan-2006 15:04:05"),
		UserNumber:            academicYearInfo.UserNumber,
		LevelNumber:           academicYearInfo.CampusLevelNumber,
		LevelName:             academicYearInfo.CampusLevel.CurriculumLevel.Level.Name,
	}
	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": academicYearInfo_})
}
func FetchAcademicYearInfo(context *gin.Context) {

	campusCurriculumNumber := context.Query("campus_curriculum_number")
	//campusLevelNumber := context.Query("campus_level_number")

	if campusCurriculumNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Curriculum is required!"})
		return
	}

	var campusLevel []levelModels.CampusLevel
	var campusLevel_ []responses3.CampusLevelDataList
	if err := config.DB.
		Where("campus_curriculum_number = ?", campusCurriculumNumber).
		Preload("CurriculumLevel.Level").
		Preload("CurriculumLevel").
		Find(&campusLevel).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No curriculum found"})
		return
	}

	for _, items := range campusLevel {

		var academicYearInfo []academicYearModel.AcademicYear
		var academicYearInfo_ []responses.AcademicYearList
		if err := config.DB.
			Where("campus_level_number = ?", items.Number).
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

		campusLevel_ = append(campusLevel_, responses3.CampusLevelDataList{
			Number:                items.Number,
			CurriculumLevelNumber: items.CurriculumLevelNumber,
			Code:                  items.Code,
			Name:                  items.CurriculumLevel.Level.Name,
			AcademicYear:          academicYearInfo_,
		})
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": campusLevel_})
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

func DeleteAcademicYear(context *gin.Context) {

	var AcademicYearInfo requests2.AcademicYearDelRequest
	err := context.BindJSON(&AcademicYearInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	academicYearNumber := AcademicYearInfo.Number

	var studyPeriodChecker studyPeriodModel.StudyPeriod
	if err := config.DB.
		Where("academic_year_number = ?", academicYearNumber).
		First(&studyPeriodChecker).Error; err == nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "You can not delete this Academic year, it is used on Study Period"})
		return
	}

	var schemeOfWorkChecker schemeOfWorkModel.SchemeOfWork
	if err := config.DB.
		Where("academic_year_number = ?", academicYearNumber).
		First(&schemeOfWorkChecker).Error; err == nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "You can not delete this Academic year, it is used on Scheme of work"})
		return
	}

	var academicYearData academicYearModel.AcademicYear
	academicYearDataInfo := config.DB.
		Where("number = ?", academicYearNumber).
		First(&academicYearData)
	if academicYearDataInfo.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Academic Year not found!"})
		return
	}

	tx := config.DB.Begin()

	//delete Study Period
	if err := tx.
		Unscoped().
		Where("number = ?", academicYearNumber).
		Delete(&academicYearData).Error; err != nil {
		// If any error occurs, rollback the transaction
		tx.Rollback()
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to delete Academic Year"})
		return
	}

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to delete academic year"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Academic year deleted successfully!"})

}

func AcademicYearsByLevel(context *gin.Context) {

	campusLevelNumber := context.Query("campus_level_number")

	if campusLevelNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Level is required!"})
		return
	}

	var academicYearInfo []academicYearModel.AcademicYear
	var academicYearInfo_ []responses.AcademicYearList
	if err := config.DB.
		Where("campus_level_number = ?", campusLevelNumber).
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

func AcademicYearsByCampusClass(context *gin.Context) {
	campusClassNumber := context.Query("campus_class_number")

	if campusClassNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Class is required!"})
		return
	}

	var campusClassInfo models2.CampusClass
	if err := config.DB.
		Where("number = ?", campusClassNumber).
		First(&campusClassInfo).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Class found!"})
		return
	}

	var academicYearInfo []academicYearModel.AcademicYear
	var academicYearInfo_ []responses.AcademicYearList
	if err := config.DB.
		Where("campus_level_number = ?", campusClassInfo.CampusLevelNumber).
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

		var studyPeriodInfo []studyPeriodModel.StudyPeriod
		var studyPeriodInfo_ []responses2.StudyPeriodList
		if err := config.DB.
			Where("academic_year_number = ?", academicList.Number).
			Order("code asc").
			Find(&studyPeriodInfo).Error; err != nil {
			context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Academic year found"})
			return
		}

		// Populate the selected columns from the fetched data
		for _, studyPeriodList := range studyPeriodInfo {

			studyPeriodName, _ := repositories3.GetStudyPeriodNameByNumber(studyPeriodList.StudyPeriodNameNumber)

			studyPeriodInfo_ = append(studyPeriodInfo_, responses2.StudyPeriodList{
				Number:               studyPeriodList.Number,
				Code:                 studyPeriodList.Code,
				StudyPeriodName:      studyPeriodName.Name,
				StudyPeriodStartDate: (studyPeriodList.StudyPeriodStartDate).Format("2-Jan-2006 15:04:05"),
				StudyPeriodEndDate:   (studyPeriodList.StudyPeriodEndDate).Format("2-Jan-2006 15:04:05"),
				//AcademicYearNumber:   studyPeriodList.AcademicYearNumber,
				//AcademicYearName:     studyPeriodList.AcademicYear.AcademicYearName,
				UserNumber: studyPeriodList.UserNumber,
			})
		}

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
			StudyPeriod:           studyPeriodInfo_,
		})
	}
	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": academicYearInfo_})
}
