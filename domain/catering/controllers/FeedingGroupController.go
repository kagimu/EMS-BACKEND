package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/services"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
)

func CreateFeedingGroup(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var payload requests.FeedingGroupRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = services.CreateFeedingGroup(payload, user.Number)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Create Feeding Group error: " + err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Feeding group has been created"})
}

func FetchFeedingGroups(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	schoolNumber := c.Query("school_number")
	if schoolNumber == "" {
		c.JSON(http.StatusBadRequest, commons.Feedback("Provide school number"))
		return
	}

	status := c.Query("status")

	feedingGroups, err := services.GetFeedingGroupsBySchoolNumber(schoolNumber, status)

	if err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Getting feeding groups error: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": feedingGroups})
}

func UpdateFeedingGroup(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var payload requests.UpdateFeedingGroupRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Check the provided body request."})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = services.UpdateFeedingGroup(payload, user.Number)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Updating feeding group error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Feeding group has been updated"})
}

func DeleteFeedingGroup(context *gin.Context) {

	_, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	feedingGroupNumber := context.Query("feeding_group_number")

	err = services.DeleteFeedingGroup(feedingGroupNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Removing feeding group error: " + err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Feeding group has been removed"})
}

func ActivateFeedingGroupOnCalendars(context *gin.Context) {

	user, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var payload requests.FeedingGroupCalendarRequest
	if err := context.ShouldBindJSON(&payload); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Check the provided body request."})
		return
	}

	if err = payload.Validate(); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = services.SaveFeedingGroupCalendar(payload, user.Number)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, commons.Feedback("Groups activated on calendars"))
}

func FetchFeedingGroupOnCalendars(context *gin.Context) {
	schoolNumber := context.Query("school_number")
	if schoolNumber == "" {
		context.JSON(http.StatusBadRequest, commons.Feedback("Provide school number"))
		return
	}

	entityNumber, _ := commons.GetUserEntityNumber(context)
	if entityNumber == "" {
		context.JSON(http.StatusBadRequest, commons.Feedback("Set current entity"))
		return
	}

	groupCalendarEventResponse, err := services.GetFeedingGroupCalendarEvent(schoolNumber, entityNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": groupCalendarEventResponse})
}
