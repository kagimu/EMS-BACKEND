package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	supportRequests "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/support/api/requests"
	supportModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/support/models"
	supportServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/support/services"
	"github.com/gin-gonic/gin"
)

func page(c *gin.Context) (int, int) {
	value, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	return value, perPage
}

func actor(c *gin.Context) (string, bool) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil || user == nil || user.Number == "" {
		c.JSON(http.StatusUnauthorized, commons.Feedback("Unauthorized"))
		return "", false
	}
	return user.Number, true
}

func writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, supportServices.ErrTicketNotFound):
		c.JSON(http.StatusNotFound, commons.Feedback("Support ticket not found"))
	case errors.Is(err, supportServices.ErrTicketAuthorisation):
		c.JSON(http.StatusForbidden, commons.Feedback("You cannot access this support ticket"))
	case errors.Is(err, supportServices.ErrTicketTransition):
		c.JSON(http.StatusConflict, commons.Feedback(err.Error()))
	case errors.Is(err, supportServices.ErrTicketValidation),
		errors.Is(err, supportServices.ErrTicketResolution),
		errors.Is(err, supportServices.ErrTicketMessage):
		c.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
	default:
		c.JSON(http.StatusInternalServerError, commons.Feedback("The support request could not be completed"))
	}
}

// OwnerOverview returns owner support counts.
func OwnerOverview(c *gin.Context) {
	overview, err := supportServices.Overview()
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Support overview fetched", "data": overview})
}

// OwnerTickets lists all tickets for the owner.
func OwnerTickets(c *gin.Context) {
	pageNumber, perPage := page(c)
	result, err := supportServices.ListTickets("", c.Query("organisation"), c.Query("search"), c.Query("status"), c.Query("category"), c.Query("priority"), c.Query("assigned_to"), pageNumber, perPage)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Support tickets fetched", "data": result})
}

// OwnerCreateTicket creates a ticket on behalf of a customer.
func OwnerCreateTicket(c *gin.Context) {
	actorNumber, ok := actor(c)
	if !ok {
		return
	}
	var request supportRequests.CreateTicketRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}
	detail, err := supportServices.CreateTicket(actorNumber, request.OrganisationNumber, request, true)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Support ticket created", "data": detail})
}

// OwnerTicket returns a ticket with internal notes and audit.
func OwnerTicket(c *gin.Context) {
	detail, err := supportServices.GetTicket(c.Param("number"), "", true)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Support ticket fetched", "data": detail})
}

// OwnerMessage adds an internal note or customer-visible reply.
func OwnerMessage(c *gin.Context) {
	actorNumber, ok := actor(c)
	if !ok {
		return
	}
	var request supportRequests.AddMessageRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}
	message, err := supportServices.AddMessage(actorNumber, c.Param("number"), "", supportModels.AuthorOwner, request.Visibility, request.Body)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Support message added", "data": message})
}

// OwnerUpdate changes status, priority, category, assignment or resolution.
func OwnerUpdate(c *gin.Context) {
	actorNumber, ok := actor(c)
	if !ok {
		return
	}
	var request supportRequests.UpdateTicketRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}
	detail, err := supportServices.UpdateTicket(actorNumber, c.Param("number"), request)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Support ticket updated", "data": detail})
}

// ClientTickets lists only the active organisation's tickets.
func ClientTickets(c *gin.Context) {
	organisationNumber, err := commons.GetUserEntityNumber(c)
	if err != nil || organisationNumber == "" {
		c.JSON(http.StatusConflict, commons.Feedback("No active organisation is selected"))
		return
	}
	pageNumber, perPage := page(c)
	result, err := supportServices.ListTickets(organisationNumber, "", c.Query("search"), c.Query("status"), c.Query("category"), c.Query("priority"), "", pageNumber, perPage)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Support tickets fetched", "data": result})
}

// ClientCreateTicket creates a ticket for the active organisation.
func ClientCreateTicket(c *gin.Context) {
	actorNumber, ok := actor(c)
	if !ok {
		return
	}
	organisationNumber, err := commons.GetUserEntityNumber(c)
	if err != nil || organisationNumber == "" {
		c.JSON(http.StatusConflict, commons.Feedback("No active organisation is selected"))
		return
	}
	var request supportRequests.CreateTicketRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}
	detail, err := supportServices.CreateTicket(actorNumber, organisationNumber, request, false)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Support ticket created", "data": detail})
}

// ClientTicket returns one ticket in the active organisation; internal notes
// and audit records are omitted.
func ClientTicket(c *gin.Context) {
	organisationNumber, err := commons.GetUserEntityNumber(c)
	if err != nil || organisationNumber == "" {
		c.JSON(http.StatusConflict, commons.Feedback("No active organisation is selected"))
		return
	}
	detail, err := supportServices.GetTicket(c.Param("number"), organisationNumber, false)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Support ticket fetched", "data": detail})
}

// ClientMessage adds a customer-visible reply.
func ClientMessage(c *gin.Context) {
	actorNumber, ok := actor(c)
	if !ok {
		return
	}
	organisationNumber, err := commons.GetUserEntityNumber(c)
	if err != nil || organisationNumber == "" {
		c.JSON(http.StatusConflict, commons.Feedback("No active organisation is selected"))
		return
	}
	var request supportRequests.AddMessageRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}
	message, err := supportServices.AddMessage(actorNumber, c.Param("number"), organisationNumber, supportModels.AuthorCustomer, supportModels.VisibilityCustomer, request.Body)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Reply added", "data": message})
}
