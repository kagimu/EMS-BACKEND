package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/pagination"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculumdelivery/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculumdelivery/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddSession(ctx *gin.Context) {

	var request requests.CurriculumDeliveryRequest
	err := ctx.ShouldBindJSON(&request)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, commons.Feedback("Invalid request data "+err.Error()))
		return
	}

	err = request.Validate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	err = services.SaveCurriculumDeliverySession(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, commons.Feedback("Failed to add session"+err.Error()))
		return
	}

	ctx.JSON(200, commons.Feedback("Sessions added successfully"))
}

func GetSessions(ctx *gin.Context) {
	paginatedResponse := pagination.PaginatedResponse{}
	sessions, err := services.GetCurriculumDeliveryByStreamNumberAndStepNumber(ctx, &paginatedResponse)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, commons.Feedback("Failed to get sessions"))
	}
	paginatedResponse.Content = sessions

	ctx.JSON(200, paginatedResponse)
}

func GetStudentsBySession(context *gin.Context) {
	paginatedResponse := pagination.PaginatedResponse{}
	students, err := services.GetStudentsBySession(context, &paginatedResponse)
	if err != nil {
		context.JSON(http.StatusInternalServerError, commons.Feedback("Failed to get students"))
	}
	paginatedResponse.Content = students

	context.JSON(200, paginatedResponse)
}

func UpdateSession(ctx *gin.Context) {
	var request requests.UpdateCurriculumDeliveryRequest
	err := ctx.ShouldBindJSON(&request)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, commons.Feedback("Invalid request data "+err.Error()))
		return
	}

	err = request.Validate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	err = services.UpdateCurriculumDeliverySession(request)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, commons.Feedback("Failed to update session"))
		return
	}

	ctx.JSON(200, commons.Feedback("Sessions updated successfully"))
}

func GetCurriculumDeliveryProgressByStream(context *gin.Context) {

	stream := context.Param("stream_number")
	if stream == "" {
		context.JSON(http.StatusBadRequest, commons.Feedback("Please provide stream number"))
		return
	}
	progress, err := services.GetCurriculumDeliveryProgressByStream(stream)
	if err != nil {
		context.JSON(http.StatusInternalServerError, commons.Feedback("Failed to get curriculum delivery progress"+err.Error()))
		return
	}

	context.JSON(200, progress)

}
