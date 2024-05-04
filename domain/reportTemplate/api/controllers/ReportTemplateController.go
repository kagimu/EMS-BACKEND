package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	newTemplateRequest "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportTemplate/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportTemplate/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportTemplate/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AllReportTemplates(context *gin.Context) {

	var templateResponses []responses.ReportTemplateResponse

	search := context.Query("search")

	schoolNumber := context.Query("school_number")

	templateResponses, errorResponse := services.FetchAllReportTemplates(schoolNumber, search)

	if errorResponse != nil {
		fmt.Println(errorResponse.Error())
		context.JSON(http.StatusInternalServerError, commons.Feedback("We could not get results, try again!"))
	}

	context.JSON(http.StatusOK, templateResponses)
}

func NewReportTemplate(context *gin.Context) {

	currentUser, err := commons.GetUserFromHeader(context)

	if err != nil {
		context.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var templateRequest newTemplateRequest.NewTemplateRequest

	fmt.Println("1: I reach here")
	fmt.Println(templateRequest)

	if noBodyError := context.BindJSON(&templateRequest); noBodyError != nil {
		context.JSON(http.StatusBadRequest, commons.Feedback("Set valid request body"))
		return
	}

	fmt.Println("2: I reach here two")
	fmt.Println(templateRequest)

	if validationError := templateRequest.Validate(); validationError != nil {
		context.JSON(http.StatusBadRequest, commons.Feedback(validationError.Error()))
		return
	}

	savingError, err := services.PostNewReportTemplate(templateRequest, currentUser.Number)

	if err != nil {
		context.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	if savingError != "" {
		context.JSON(http.StatusBadRequest, commons.Feedback(savingError))
		return
	}

	context.JSON(http.StatusOK, commons.Feedback("New template created successfully"))

}

func ClassReportTemplates(context *gin.Context) {

	var classReportTemplateResponses []responses.ClassReportTemplateResponse

	templateNumber := context.Query("template_number")
	schoolNumber := context.Query("school_number")

	classReportTemplateResponses, errorResponse := services.FetchAllClassReportTemplate(schoolNumber, templateNumber)

	if errorResponse != nil {
		fmt.Println(errorResponse.Error())
		context.JSON(http.StatusInternalServerError, commons.Feedback("We could not get results, try again!"))
	}

	context.JSON(http.StatusOK, classReportTemplateResponses)
}

func AddTemplateToClass(context *gin.Context) {

	currentUser, err := commons.GetUserFromHeader(context)

	if err != nil {
		context.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var templateRequest newTemplateRequest.ClassReportTemplateRequest

	if noBodyError := context.BindJSON(&templateRequest); noBodyError != nil {
		context.JSON(http.StatusBadRequest, commons.Feedback("Set valid request body"))
		return
	}

	if validationError := templateRequest.Validate(); validationError != nil {
		context.JSON(http.StatusBadRequest, commons.Feedback(validationError.Error()))
		return
	}

	savingError, validationErrors := services.AddTemplateToClass(templateRequest, currentUser.Number)

	if savingError != "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": savingError, "data": validationErrors})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Classes added to report template", "data": validationErrors})
}
