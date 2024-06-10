package controllers

import (
	"fmt"
	activityReportingRequests "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/activityreporting/api/requests"
	activityReportingServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/activityreporting/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SaveIndicator(context *gin.Context) {

	user, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var payload activityReportingRequests.IndicatorRequest
	if err = context.ShouldBindJSON(&payload); err != nil {
		fmt.Println(err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "Check input body"})
		return
	}

	if err = payload.Validate(); err != nil {
		context.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	if err = activityReportingServices.SaveIndicator(payload, user.Number); err != nil {
		context.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	context.JSON(http.StatusOK, commons.Feedback("Indicator added successfully"))

}

func GetIndicators(context *gin.Context) {
	_, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	entityNumber := context.Query("entity_number")
	if entityNumber == "" {
		context.JSON(http.StatusBadRequest, commons.Feedback("provide entity number"))
		return
	}

	name := context.Query("name")
	status := context.Query("status")

	response, err := activityReportingServices.GetIndicators(entityNumber, name, status)
	if err != nil {
		context.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": response})
}
