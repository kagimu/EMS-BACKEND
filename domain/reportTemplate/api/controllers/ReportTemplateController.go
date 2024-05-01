package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
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
