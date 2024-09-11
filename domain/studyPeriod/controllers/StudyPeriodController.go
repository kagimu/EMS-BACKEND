package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/api/responses"
	academicYearModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/models"
	responses3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	levelModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/api/requests"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/api/responses"
	studyPeriodModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

var studyPeriodService = services.NewStudyPeriodService()

func NewStudyPeriod(context *gin.Context) {

	user, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var studyPeriodInfo requests.StudyPeriodRequest
	err = context.BindJSON(&studyPeriodInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed parsing request: " + err.Error()})
		return
	}

	err = studyPeriodInfo.Validate()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = studyPeriodService.NewStudyPeriod(studyPeriodInfo, user.Number)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Study period created successfully"})
}

func EditStudyPeriod(context *gin.Context) {

	user, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var studyPeriodInfo requests.UpdateStudyPeriodRequest
	err = context.BindJSON(&studyPeriodInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed parsing request: " + err.Error()})
		return
	}

	err = studyPeriodInfo.Validate()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = studyPeriodService.UpdateStudyPeriod(studyPeriodInfo, user.Number)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Study Period Edited Successfully!"})
}

func FetchStudyPeriodByAcademicYear(context *gin.Context) {
	academicYearNumber := context.Query("academic_year_number")

	var studyPeriodInfo []studyPeriodModel.StudyPeriod
	var studyPeriodInfo_ []responses2.StudyPeriodList
	if err := config.DB.
		Where("academic_year_number = ?", academicYearNumber).
		Order("code asc").
		Find(&studyPeriodInfo).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Academic year found"})
		return
	}

	// Populate the selected columns from the fetched data
	for _, studyPeriodList := range studyPeriodInfo {
		studyPeriodInfo_ = append(studyPeriodInfo_, responses2.StudyPeriodList{
			StudyPeriodNumber:    studyPeriodList.Number,
			Code:                 studyPeriodList.Code,
			StudyPeriodName:      "studyPeriodList.StudyPeriodName",
			StudyPeriodStartDate: (studyPeriodList.StudyPeriodStartDate).Format("2-Jan-2006 15:04:05"),
			StudyPeriodEndDate:   (studyPeriodList.StudyPeriodEndDate).Format("2-Jan-2006 15:04:05"),
			//AcademicYearNumber:   studyPeriodList.AcademicYearNumber,
			//AcademicYearName:     studyPeriodList.AcademicYear.AcademicYearName,
			UserNumber: studyPeriodList.UserNumber,
		})
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": studyPeriodInfo_})
}

func FetchStudyPeriodDetails(context *gin.Context) {

	studyPeriodNumber := context.Query("study_period_number")

	if studyPeriodNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Study Period is required!"})
		return
	}

	studyPeriodInfo, err := studyPeriodService.StudyPeriodDetails(studyPeriodNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": studyPeriodInfo})
}

func FetchStudyPeriod(context *gin.Context) {

	campusCurriculumNumber := context.Query("campus_curriculum_number")

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

		var academicYearInfo []academicYearModels.AcademicYear
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
				studyPeriodInfo_ = append(studyPeriodInfo_, responses2.StudyPeriodList{
					StudyPeriodNumber:    studyPeriodList.Number,
					Code:                 studyPeriodList.Code,
					StudyPeriodName:      "studyPeriodList.StudyPeriodName",
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

func NewStudyPeriodName(context *gin.Context) {

	user, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var studyPeriodNameRequest requests.StudyPeriodNameRequest
	err = context.BindJSON(&studyPeriodNameRequest)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Processing body request error: " + err.Error()})
		return
	}

	err = studyPeriodNameRequest.Validate()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = studyPeriodService.NewStudyPeriodName(studyPeriodNameRequest, user.Number)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Study period name created successfully"})
	return
}

func UpdateStudyPeriodName(context *gin.Context) {

	user, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var studyPeriodNameRequest requests.UpdateStudyPeriodNameRequest
	err = context.BindJSON(&studyPeriodNameRequest)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Processing body request error: " + err.Error()})
		return
	}

	err = studyPeriodNameRequest.Validate()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = studyPeriodService.UpdateStudyPeriodName(studyPeriodNameRequest, user.Number)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Study period name updated successfully"})
	return
}

func RemoveStudyPeriodName(context *gin.Context) {

	_, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	studyPeriodNameNumber := context.Param("study_period_name_number")

	if studyPeriodNameNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Choose study period name"})
		return
	}

	err = studyPeriodService.RemoveStudyPeriodName(studyPeriodNameNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Study period name removed successfully"})
}

func AllStudyPeriodName(context *gin.Context) {
	schoolNumber := context.Query("school_number")
	if schoolNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Choose school"})
		return
	}

	all, err := studyPeriodService.FindAllStudyPeriodNames(schoolNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": all})
}

func RemoveStudyPeriod(context *gin.Context) {
	studyPeriodNumber := context.Param("study_period_number")
	if studyPeriodNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Choose study period"})
		return
	}

	err := studyPeriodService.RemoveStudyPeriod(studyPeriodNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Study period removed successfully"})
}

func AllStudyPeriods(context *gin.Context) {
	academicYearNumber := context.Query("academic_year_number")
	calendarEventNumber := context.Query("calendar_event_number")
	studyPeriodNameNumber := context.Query("study_period_name_number")
	classNumber := context.Query("class_number")
	schoolNumber := context.Query("school_number")

	if schoolNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Choose a school"})
		return
	}

	if academicYearNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Choose the academic year"})
		return
	}

	all, err := studyPeriodService.FindAllStudyPeriods(academicYearNumber, calendarEventNumber, studyPeriodNameNumber, classNumber, schoolNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": all})
	return
}
