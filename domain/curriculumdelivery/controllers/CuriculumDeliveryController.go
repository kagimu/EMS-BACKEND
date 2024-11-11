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

	err = services.SaveCurriculumDeliverySession(ctx, request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, commons.Feedback("Failed to add session"))
	}

	ctx.JSON(200, commons.Feedback("Session added successfully"))
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
