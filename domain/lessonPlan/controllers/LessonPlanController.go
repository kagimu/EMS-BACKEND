package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/lessonPlan/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/lessonPlan/api/responses"
	lessonPlanModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/lessonPlan/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strconv"
)

func NewLessonPlan(context *gin.Context) {

	var lessonPlanInfo requests.LessonPlanInfo

	err := context.BindJSON(&lessonPlanInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error! " + err.Error()})
		return
	}

	if err := lessonPlanInfo.Validate(); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	schemeSubTopicNumber := lessonPlanInfo.SchemeSubTopicNumber
	lessonPlanNumber := lessonPlanInfo.LessonPlanNumber
	lessonPlanConceptInfo := lessonPlanInfo.LessonPlanConceptInfo

	userNumber := context.Request.Header["User-Number"][0]

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
		Where("scheme_sub_topic_number = ?", schemeSubTopicNumber).
		Where("number != ?", lessonPlanNumber).
		First(&duplicateLessonPlan).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Lesson plan already exists"})
		return
	}

	// Begin the transaction
	tx := config.DB.Begin()

	number := uuid.NewV4().String()
	///save lesson plan data
	saveLessonPlan := lessonPlanModel.LessonPlan{
		Number:               number,
		SchemeSubTopicNumber: schemeSubTopicNumber,
		Status:               "draft",
		UserNumber:           userNumber,
	}

	if lessonPlanNumber == "" {
		if err := tx.Create(&saveLessonPlan).Error; err != nil {
			// If any error occurs, rollback the transaction
			tx.Rollback()
			// Handle or return the error as needed
			//context.JSON(http.StatusBadRequest, gin.H{"message": err})
			context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save Lesson plan"})
			return
		}
	}

	for _, concept := range lessonPlanConceptInfo {

		var duplicateLessonPlanConcept lessonPlanModel.LessonPlanConcepts
		if duplicate := tx.
			Where("concept_number = ?", concept.LessonPlanConceptNumber).
			Where("lesson_plan_number = ?", saveLessonPlan.Number).
			First(&duplicateLessonPlanConcept).Error; duplicate == nil {
			tx.Rollback()
			context.JSON(http.StatusBadRequest, gin.H{"message": "Duplicates Concepts not allowed"})
			return
		}

		saveLessonConcepts := lessonPlanModel.LessonPlanConcepts{
			Number:           uuid.NewV4().String(),
			ConceptNumber:    concept.LessonPlanConceptNumber,
			Code:             concept.Code,
			LessonPlanNumber: saveLessonPlan.Number,
			UserNumber:       userNumber,
		}

		//	adds new concepts when Concept number doesnt exists
		if concept.Number == "" {
			if err := tx.Create(&saveLessonConcepts).Error; err != nil {
				// If any error occurs, rollback the transaction
				tx.Rollback()
				// Handle or return the error as needed
				context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save Lesson Plan Concepts!"})
				return
			}
		}

		var referenceDataInfo lessonPlanModel.LessonPlanStep
		conceptNo := ""
		if concept.Number == "" {
			conceptNo = saveLessonConcepts.Number
		} else {
			conceptNo = concept.Number
		}

		tx.Unscoped().
			Where("lesson_plan_concept_number = ?", conceptNo).
			Delete(&referenceDataInfo)

		lessonPlanSteps := concept.StepList
		if len(lessonPlanSteps) > 0 {

			for _, step := range lessonPlanSteps {
				if step.Step != "" {
					saveLearningSteps := lessonPlanModel.LessonPlanStep{
						Number:                  uuid.NewV4().String(),
						Code:                    step.Code,
						Step:                    step.Step,
						LessonPlanConceptNumber: conceptNo,
					}

					if err := tx.Create(&saveLearningSteps).Error; err != nil {
						// If any error occurs, rollback the transaction
						tx.Rollback()
						// Handle or return the error as needed
						context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save steps"})
						return

					}
				}

			}

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

func AddLessonPlanReference(ctx *gin.Context) {
	var lessonPlanReferenceRequest requests.LessonPlanReferenceRequest

	err := ctx.Bind(&lessonPlanReferenceRequest)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, commons.Feedback(err.Error()))
		return
	}

	err = lessonPlanReferenceRequest.Validate()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	code, err := strconv.Atoi(lessonPlanReferenceRequest.Code)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))

		return
	}

	var lessonPlanReference = lessonPlanModel.LessonPlanReference{
		LessonPlanNumber:       lessonPlanReferenceRequest.LessonPlanNumber,
		LearningMaterialNumber: lessonPlanReferenceRequest.LearningMaterialNumber,
		Name:                   lessonPlanReferenceRequest.Name,
		Code:                   int64(code),
		UserNumber:             ctx.GetHeader("User-Number"),
		Number:                 uuid.NewV4().String(),
	}

	if config.DB.Create(&lessonPlanReference).Error != nil {
		ctx.JSON(http.StatusInternalServerError, commons.Feedback("Failed to add lesson plan reference"))
		return
	}

	ctx.JSON(http.StatusOK, commons.Feedback("Lesson plan reference added successfully!"))

}

func FetchLessonPlanReferences(ctx *gin.Context) {
	var lessonPlanReferences []responses.LessonPlanReferenceResponse
	var references []lessonPlanModel.LessonPlanReference

	if config.DB.Find(&references).Error != nil {
		ctx.JSON(http.StatusInternalServerError, commons.Feedback("Failed to fetch lesson plan references"))
		return
	}

	for _, reference := range references {
		lessonPlanReferences = append(lessonPlanReferences, responses.LessonPlanReferenceResponse{
			LessonPlanNumber:       reference.LessonPlanNumber,
			LearningMaterialNumber: reference.LearningMaterialNumber,
			Name:                   reference.Name,
			Code:                   reference.Code,
		})
	}

	ctx.JSON(http.StatusOK, lessonPlanReferences)
}

func FetchLessonPlan(context *gin.Context) {

	subTopicNumber := context.Query("sub_topic_number")

	if subTopicNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Subtopic is required!"})
		return
	}

	var lessonPlanInfo lessonPlanModel.LessonPlan
	if err := config.DB.
		Where("scheme_sub_topic_number = ?", subTopicNumber).
		//Preload("SchemeOfWork").
		Order("created_at desc").
		First(&lessonPlanInfo).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Lesson plan found"})
		return
	}

	var lessonPlanConceptsInfo []lessonPlanModel.LessonPlanConcepts
	var lessonPlanConceptsInfo_ []responses.LessonPlanConceptData

	if err := config.DB.
		Where("lesson_plan_number = ?", lessonPlanInfo.Number).
		Preload("Concept").
		Order("code asc").
		Find(&lessonPlanConceptsInfo).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Lesson plan concept found"})
		return
	}

	for _, lessonPlanConcept := range lessonPlanConceptsInfo {

		var lessonPlanStepsInfo []lessonPlanModel.LessonPlanStep
		var lessonPlanStepsInfo_ []responses.LessonPlanStepsData

		if err := config.DB.
			Where("lesson_plan_concept_number = ?", lessonPlanConcept.Number).
			Order("code asc").
			Find(&lessonPlanStepsInfo).Error; err != nil {
			context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Lesson plan Steps found"})
			return
		}

		for _, steps := range lessonPlanStepsInfo {
			lessonPlanStepsInfo_ = append(lessonPlanStepsInfo_, responses.LessonPlanStepsData{
				Number:                  steps.Number,
				Code:                    steps.Code,
				Step:                    steps.Step,
				LessonPlanConceptNumber: steps.LessonPlanConceptNumber,
			})
		}

		lessonPlanConceptsInfo_ = append(lessonPlanConceptsInfo_, responses.LessonPlanConceptData{
			Number:           lessonPlanConcept.Number,
			Code:             lessonPlanConcept.Code,
			ConceptNumber:    lessonPlanConcept.ConceptNumber,
			ConceptName:      lessonPlanConcept.Concept.Name,
			LessonPlanNumber: lessonPlanConcept.LessonPlanNumber,
			UserNumber:       lessonPlanConcept.UserNumber,
			LessonPlanSteps:  lessonPlanStepsInfo_,
		})
	}

	_lessonPlanResp := responses.LessonPlanData{
		Number:               lessonPlanInfo.Number,
		SchemeSubTopicNumber: lessonPlanInfo.SchemeSubTopicNumber,
		Status:               lessonPlanInfo.Status,
		UserNumber:           lessonPlanInfo.UserNumber,
		LessonPlanConcepts:   lessonPlanConceptsInfo_,
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _lessonPlanResp})

}
