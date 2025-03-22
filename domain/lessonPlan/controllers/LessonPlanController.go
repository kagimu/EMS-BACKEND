package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/lessonPlan/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/lessonPlan/api/responses"
	lessonPlanModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/lessonPlan/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/lessonPlan/services"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
)

var lessonPlanService = services.NewLessonPlanService()

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

	userNumber := context.Request.Header["User-Number"][0]

	err = lessonPlanService.NewLessonPlan(lessonPlanInfo, userNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
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

	err = lessonPlanService.AddLessonPlanReferences(lessonPlanReferenceRequest, ctx.GetHeader("User-Number"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
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
			LessonPlanNumber:          reference.LessonPlanNumber,
			LearningMaterialNumber:    reference.LearningMaterialNumber,
			Name:                      reference.Name,
			Code:                      reference.Code,
			LessonPlanReferenceNumber: reference.Number,
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

	var lessonPlanConceptsInfo []lessonPlanModel.LessonPlanConcept
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

	var lessonPlanReferenceInfo []lessonPlanModel.LessonPlanReference
	config.DB.Where("lesson_plan_number = ?", lessonPlanInfo.Number).Order("code asc").Find(&lessonPlanReferenceInfo)

	var lessonPlanReferenceInfo_ []responses.LessonPlanReferenceResponse

	for _, reference := range lessonPlanReferenceInfo {
		lessonPlanReferenceInfo_ = append(lessonPlanReferenceInfo_, responses.LessonPlanReferenceResponse{
			LessonPlanNumber:          reference.LessonPlanNumber,
			LearningMaterialNumber:    reference.LearningMaterialNumber,
			Name:                      reference.Name,
			Code:                      reference.Code,
			LessonPlanReferenceNumber: reference.Number,
		})
	}

	_lessonPlanResp := responses.LessonPlanData{
		Number:               lessonPlanInfo.Number,
		SchemeSubTopicNumber: lessonPlanInfo.SchemeSubTopicNumber,
		Status:               lessonPlanInfo.Status,
		UserNumber:           lessonPlanInfo.UserNumber,
		LessonPlanConcepts:   lessonPlanConceptsInfo_,
		References:           lessonPlanReferenceInfo_,
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _lessonPlanResp})

}

func ImportLessonPlan(context *gin.Context) {
	user, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusForbidden, gin.H{"message": "User not found"})
	}

	schemeSubTopicNumber := context.PostForm("scheme_sub_topic_number")
	if schemeSubTopicNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Choose the scheme subtopic for this lesson plan"})
		return
	}

	file, err := context.FormFile("file")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to get file: " + err.Error()})
		return
	}

	uploadedFile, err := file.Open()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to open file: " + err.Error()})
		return
	}
	defer func(uploadedFile multipart.File) {
		err = uploadedFile.Close()
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Closing file failed with error : " + err.Error()})
			return
		}
	}(uploadedFile)

	err = lessonPlanService.ImportPlan(uploadedFile, schemeSubTopicNumber, user.Number)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Concepts imported successfully"})
}

func ExportLessonPlanTemplate(context *gin.Context) {
	_, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusForbidden, gin.H{"message": "User not found"})
	}

	schemeSubTopicNumber := context.Param("scheme_sub_topic_number")
	if schemeSubTopicNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Choose the scheme subtopic for this lesson plan"})
		return
	}

	template, fileName, err := lessonPlanService.ExportLessonPlanTemplate(schemeSubTopicNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.Header("Content-Disposition", "attachment; filename="+fileName)
	context.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	if err := template.Write(context.Writer); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
}

func UpdateLessonPlanReference(context *gin.Context) {

	user, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusForbidden, gin.H{"message": "User not found"})
	}

	referenceNumber := context.Param("reference_number")
	if referenceNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Choose the reference number to update"})
		return
	}

	var referenceRequest requests.Material

	err = context.BindJSON(&referenceRequest)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = referenceRequest.Validate()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = lessonPlanService.UpdateLessonPlanReference(referenceRequest, referenceNumber, user.Number)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Reference updated successfully"})
}

func DeleteLessonPlanReference(context *gin.Context) {
	_, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusForbidden, gin.H{"message": "User not found"})
	}

	referenceNumber := context.Param("reference_number")
	if referenceNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Choose the reference number to delete"})
		return
	}

	err = lessonPlanService.DeleteLessonPlanReference(referenceNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Reference deleted successfully"})
}
