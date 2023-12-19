package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/schemeOfWork/api/responses"
	schemeOfWorkModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/schemeOfWork/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strconv"
)

func NewSchemeOfWork(context *gin.Context) {

	type SubtopicListInfo struct {
		Number        string `json:"number"`
		Name          string `json:"name"`
		Code          int64  `json:"code"`
		SortCode      int64  `json:"sort_code"`
		UserNumber    string `json:"user_number"`
		EstimatedTime int64  `json:"estimated_time"`
		TopicNumber   string `json:"topic_number"`
	}

	type SchemeOfWorkInfo struct {
		Name               string             `json:"name"`
		ClassNumber        string             `json:"class_number"`
		SubjectNumber      string             `json:"subject_number"`
		AcademicYearNumber string             `json:"academic_year_number"`
		SubtopicList       []SubtopicListInfo `json:"subtopicList"`
	}

	var schemeOfWorkInfo SchemeOfWorkInfo
	err := context.BindJSON(&schemeOfWorkInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error! " + err.Error()})
		return
	}

	name := schemeOfWorkInfo.Name
	classNumber := schemeOfWorkInfo.ClassNumber
	subjectNumber := schemeOfWorkInfo.SubjectNumber
	academicYearNumber := schemeOfWorkInfo.AcademicYearNumber
	subtopicList := schemeOfWorkInfo.SubtopicList
	userNumber := context.Request.Header["User-Number"][0]

	if name == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Name is required!"})
		return
	}

	if classNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Class is required!"})
		return
	}

	if subjectNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Subject is required!"})
		return
	}

	if academicYearNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Academic Year is required!"})
		return
	}

	if len(subtopicList) == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Sub topic list is required!"})
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

	var duplicateScheme schemeOfWorkModel.SchemeOfWork
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("class_number = ?", classNumber).
		Where("subject_number = ?", subjectNumber).
		Where("academic_year_number = ?", academicYearNumber).
		First(&duplicateScheme).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Scheme already exists"})
		return
	}

	// Begin the transaction
	tx := config.DB.Begin()

	number := uuid.NewV4().String()
	///save scheme of work data
	saveSchemeOfWork := schemeOfWorkModel.SchemeOfWork{
		Number:             number,
		Name:               name,
		Status:             "Pending",
		CampusClassNumber:  classNumber,
		SubjectNumber:      subjectNumber,
		AcademicYearNumber: academicYearNumber,
		UserNumber:         userNumber,
	}

	if err := tx.Create(&saveSchemeOfWork).Error; err != nil {
		// If any error occurs, rollback the transaction
		tx.Rollback()
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create scheme"})
		return
	}

	for _, schemeList := range subtopicList {

		var duplicateSchemeCode schemeOfWorkModel.SchemeOfWorkSubTopic
		if duplicate := tx.
			Where("code = ?", schemeList.SortCode).
			Where("scheme_of_work_number = ?", saveSchemeOfWork.Number).
			First(&duplicateSchemeCode).Error; duplicate == nil {
			tx.Rollback()
			context.JSON(http.StatusBadRequest, gin.H{"message": "Duplicate codes are not allowed!!!, Code: " + strconv.FormatInt(schemeList.SortCode, 10)})
			return
		}

		schemeNo := uuid.NewV4().String()
		saveSchemeOfWorkSubtopic := schemeOfWorkModel.SchemeOfWorkSubTopic{
			Number:             schemeNo,
			Code:               schemeList.SortCode,
			EstimatedTime:      schemeList.EstimatedTime,
			SchemeOfWorkNumber: saveSchemeOfWork.Number,
			SubTopicNumber:     schemeList.Number,
			UserNumber:         userNumber,
		}
		if err := tx.Create(&saveSchemeOfWorkSubtopic).Error; err != nil {
			// If any error occurs, rollback the transaction
			tx.Rollback()
			// Handle or return the error as needed
			context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create scheme"})
			return
		}
	}

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create scheme"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Scheme Added Successfully!", "data": saveSchemeOfWork})
}

func UpdateSchemeOfWork(context *gin.Context) {

	type SubtopicListInfo struct {
		Code               int    `json:"code"`
		SortCode           int    `json:"sort_code"`
		EstimatedTime      int    `json:"estimated_time"`
		Name               string `json:"name"`
		Number             string `json:"number"`
		UserNumber         string `json:"user_number"`
		SchemeOfWorkNumber string `json:"scheme_of_work_number"`
		SubTopicNumber     string `json:"sub_topic_number"`
	}

	type SchemeOfWorkInfo struct {
		Name          string             `json:"name"`
		ClassNumber   string             `json:"class_number"`
		SubjectNumber string             `json:"subject_number"`
		SchemeNumber  string             `json:"scheme_number"`
		SubtopicList  []SubtopicListInfo `json:"subtopicList"`
	}

	var schemeOfWorkInfo SchemeOfWorkInfo
	err := context.BindJSON(&schemeOfWorkInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error! " + err.Error()})
		return
	}

	name := schemeOfWorkInfo.Name
	schemeNumber := schemeOfWorkInfo.SchemeNumber
	classNumber := schemeOfWorkInfo.ClassNumber
	subjectNumber := schemeOfWorkInfo.SubjectNumber
	subtopicList := schemeOfWorkInfo.SubtopicList
	userNumber := context.Request.Header["User-Number"][0]

	if name == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Name is required!"})
		return
	}

	if schemeNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Scheme number is required!"})
		return
	}

	if classNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Class is required!"})
		return
	}

	if subjectNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Subject is required!"})
		return
	}

	if len(subtopicList) == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Sub topic list is required!"})
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

	var duplicateScheme schemeOfWorkModel.SchemeOfWork
	if duplicate := config.DB.
		Where("number != ?", schemeNumber).
		Where("name = ?", name).
		Where("class_number = ?", classNumber).
		Where("subject_number = ?", subjectNumber).
		First(&duplicateScheme).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Scheme already exists"})
		return
	}

	tx := config.DB.Begin()
	var schemeUpdate schemeOfWorkModel.SchemeOfWork
	updateScheme := tx.Model(&schemeUpdate).
		Where("number = ?", schemeNumber).
		Updates(schemeOfWorkModel.SchemeOfWork{Name: name})
	if updateScheme.Error != nil {
		tx.Rollback()
		context.JSON(http.StatusNotFound, gin.H{"message": "Failed to update!"})
		return
	}

	for _, schemeList := range subtopicList {
		codeNo, err := strconv.ParseInt(strconv.Itoa(schemeList.SortCode), 10, 64)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		var duplicateSchemeCode schemeOfWorkModel.SchemeOfWorkSubTopic
		if duplicate := tx.
			Where("code = ?", codeNo).
			Where("scheme_of_work_number = ?", schemeNumber).
			Where("number != ?", schemeList.Number).
			First(&duplicateSchemeCode).Error; duplicate == nil {
			tx.Rollback()
			context.JSON(http.StatusBadRequest, gin.H{"message": "Duplicate codes are not allowed!!!, Code: " + strconv.FormatInt(codeNo, 10)})
			return
		}

		var schemeSubtopicUpdate schemeOfWorkModel.SchemeOfWorkSubTopic
		updateScheme := tx.Model(&schemeSubtopicUpdate).
			Where("scheme_of_work_number = ?", schemeNumber).
			Where("number = ?", schemeList.Number).
			Updates(schemeOfWorkModel.SchemeOfWorkSubTopic{Code: codeNo})
		if updateScheme.Error != nil {
			tx.Rollback()
			context.JSON(http.StatusNotFound, gin.H{"message": "Failed to update!"})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		context.JSON(http.StatusNotFound, gin.H{"message": "Failed to commit!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": subtopicList})
}

func FetchSchemeOfWork(context *gin.Context) {

	academicYearNumber := context.Query("academic_year_number")
	classNumber := context.Query("class_number")
	subjectNumber := context.Query("subject_number")

	if classNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Class is required!"})
		return
	}

	if subjectNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Subject is required!"})
		return
	}

	var schemeOfWrkInfo schemeOfWorkModel.SchemeOfWork
	if err := config.DB.
		Where("academic_year_number = ?", academicYearNumber).
		Where("campus_class_number = ?", classNumber).
		Where("subject_number = ?", subjectNumber).
		Preload("AcademicYear").
		Preload("CampusClass.CurriculumClass.Class").
		Preload("CampusClass.CurriculumClass").
		Preload("CampusClass").
		Preload("Subject").
		Order("created_at desc").
		First(&schemeOfWrkInfo).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Scheme found"})
		return
	}

	var schemeOfWrkSubTopicInfo []schemeOfWorkModel.SchemeOfWorkSubTopic
	var schemeOfWrkSubTopicInfo_ []responses2.SchemeOfWorkSubtopicInfo
	if err := config.DB.
		Where("scheme_of_work_number = ?", schemeOfWrkInfo.Number).
		Preload("SubTopic.Topic.Strand").
		Preload("SubTopic.Topic").
		Preload("SubTopic").
		Order("code asc").
		Find(&schemeOfWrkSubTopicInfo).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Scheme found"})
		return
	}

	for _, subTopicList := range schemeOfWrkSubTopicInfo {

		//fetch lesson plan
		/*	var lessonPlanInfo []models2.LessonPlan
			var lessonPlanInfo_ []responses.LessonPlanResponseInfo
			if err := config.DB.
				Where("scheme_sub_topic_number = ?", subTopicList.Number).
				//Preload("Concept").
				Order("created_at desc").
				Find(&lessonPlanInfo).Error; err != nil {
				context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Scheme found"})
				return
			}*/

		/*for _, lessonPlanList := range lessonPlanInfo {

			//fetch lesson plan steps
			var lessonPlanStepsInfo []models2.LessonPlanSteps
			var lessonPlanStepsInfo_ []responses.LessonPlanStepsInfo
			//var lessonPlanInfo_ []responses.LessonPlanResponseInfo
			if err := config.DB.
				Where("lesson_plan_number = ?", lessonPlanList.Number).
				Order("code asc").
				Find(&lessonPlanStepsInfo).Error; err != nil {
				context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Scheme found"})
				return
			}

			for _, lessonPlanSteps := range lessonPlanStepsInfo {
				lessonPlanStepsInfo_ = append(lessonPlanStepsInfo_, responses.LessonPlanStepsInfo{
					Number:           lessonPlanSteps.Number,
					Code:             lessonPlanSteps.Code,
					Step:             lessonPlanSteps.Step,
					LessonPlanNumber: lessonPlanSteps.LessonPlanNumber,
				})
			}

			lessonPlanInfo_ = append(lessonPlanInfo_, responses.LessonPlanResponseInfo{
				Number: lessonPlanList.Number,
				//LessonPlanConceptNumber: lessonPlanList.LessonPlanConceptNumber,
				SchemeSubTopicNumber: lessonPlanList.SchemeSubTopicNumber,
				//ConceptName:             lessonPlanList.Concept.Name,
				//Code:                    lessonPlanList.Code,
				UserNumber:      lessonPlanList.UserNumber,
				LessonPlanSteps: lessonPlanStepsInfo_,
			})
		}*/

		schemeOfWrkSubTopicInfo_ = append(schemeOfWrkSubTopicInfo_, responses2.SchemeOfWorkSubtopicInfo{
			Number:             subTopicList.Number,
			Code:               subTopicList.Code,
			SortCode:           subTopicList.Code,
			Name:               subTopicList.SubTopic.Name,
			TopicName:          subTopicList.SubTopic.Topic.Name,
			StrandName:         subTopicList.SubTopic.Topic.Strand.Name,
			EstimatedTime:      subTopicList.EstimatedTime,
			SchemeOfWorkNumber: subTopicList.SchemeOfWorkNumber,
			SubTopicNumber:     subTopicList.SubTopicNumber,
			UserNumber:         subTopicList.UserNumber,
			//LessonPlan:         lessonPlanInfo_,
		})
	}

	schemeOfWrk_ := responses2.SchemeOfWorkResponseInfo{
		Number:                   schemeOfWrkInfo.Number,
		Name:                     schemeOfWrkInfo.Name,
		Status:                   schemeOfWrkInfo.Status,
		ClassNumber:              schemeOfWrkInfo.CampusClassNumber,
		ClassName:                schemeOfWrkInfo.CampusClass.CurriculumClass.Class.Name,
		SubjectNumber:            schemeOfWrkInfo.SubjectNumber,
		SubjectName:              schemeOfWrkInfo.Subject.Name,
		UserNumber:               schemeOfWrkInfo.UserNumber,
		AcademicYearNumber:       schemeOfWrkInfo.AcademicYearNumber,
		AcademicYearName:         schemeOfWrkInfo.AcademicYear.AcademicYearName,
		SchemeOfWorkSubtopicInfo: schemeOfWrkSubTopicInfo_,
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": schemeOfWrk_})

}
