package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	lessonPlanModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/lessonPlan/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strconv"
)

func NewLessonPlan(context *gin.Context) {

	type Steps struct {
		Code int64  `json:"code"`
		Step string `json:"step"`
	}

	type LessonPlanInfo struct {
		Code                    string  `json:"code"`
		LessonPlanConceptNumber string  `json:"lesson_plan_concept_number"`
		SchemeSubTopicNumber    string  `json:"scheme_sub_topic_number"`
		StepList                []Steps `json:"steps"`
	}

	var lessonPlanInfo LessonPlanInfo
	err := context.BindJSON(&lessonPlanInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error! " + err.Error()})
		return
	}

	code := lessonPlanInfo.Code
	lessonPlanConceptNumber := lessonPlanInfo.LessonPlanConceptNumber
	schemeSubTopicNumber := lessonPlanInfo.SchemeSubTopicNumber
	stepList := lessonPlanInfo.StepList

	userNumber := context.Request.Header["User-Number"][0]
	fmt.Println(userNumber)
	if code == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code is required!"})
		return
	}

	if lessonPlanConceptNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Concept is required!"})
		return
	}

	if schemeSubTopicNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Scheme is required!"})
		return
	}

	if len(stepList) == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Steps are required!"})
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

	var duplicateLessonPlan lessonPlanModel.LessonPlan
	if duplicate := config.DB.
		Where("lesson_plan_concept_number = ?", lessonPlanConceptNumber).
		Where("scheme_sub_topic_number = ?", schemeSubTopicNumber).
		First(&duplicateLessonPlan).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Concept already exists"})
		return
	}

	codeNo, err := strconv.ParseInt(code, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Begin the transaction
	tx := config.DB.Begin()

	number := uuid.NewV4().String()
	///save lesson plan data
	saveLessonPlan := lessonPlanModel.LessonPlan{
		Number:                  number,
		LessonPlanConceptNumber: lessonPlanConceptNumber,
		SchemeSubTopicNumber:    schemeSubTopicNumber,
		Code:                    codeNo,
		UserNumber:              userNumber,
	}

	if err := tx.Create(&saveLessonPlan).Error; err != nil {
		// If any error occurs, rollback the transaction
		tx.Rollback()
		// Handle or return the error as needed
		//context.JSON(http.StatusBadRequest, gin.H{"message": err})
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create Lesson plan"})
		return
	}

	for _, step := range stepList {
		stepNo := uuid.NewV4().String()
		saveSchemeOfWorkSubtopic := lessonPlanModel.LessonPlanSteps{
			Number:           stepNo,
			LessonPlanNumber: saveLessonPlan.Number,
			Code:             step.Code,
			Step:             step.Step,
		}
		if err := tx.Create(&saveSchemeOfWorkSubtopic).Error; err != nil {
			// If any error occurs, rollback the transaction
			tx.Rollback()
			// Handle or return the error as needed
			context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create steps"})
			return
		}
	}

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create scheme"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Lesson Plan Added Successfully!", "data": lessonPlanInfo})
}
