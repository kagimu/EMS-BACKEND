package controllers

import (
	"errors"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	auditHooks "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/audit/hooks"
	notificationDecision "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/notifications/decision"
	organisationServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/organisation/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/stakeholders/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/stakeholders/services"
	"github.com/gin-gonic/gin"
)

// FetchStakeholders lists stakeholders for the authenticated organisation.
func FetchStakeholders(c *gin.Context) {
	_, organisationNumber, ok := stakeholderOrganisationContext(c)
	if !ok {
		return
	}

	stakeholders, err := services.ListStakeholders(organisationNumber, c.Query("stakeholder_type"), c.Query("status"))
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stakeholders fetched", "data": stakeholders})
}

// FetchStakeholder returns one stakeholder from the authenticated organisation.
func FetchStakeholder(c *gin.Context) {
	_, organisationNumber, ok := stakeholderOrganisationContext(c)
	if !ok {
		return
	}

	stakeholder, err := services.GetStakeholder(organisationNumber, c.Param("number"))
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stakeholder fetched", "data": stakeholder})
}

// CreateStakeholder creates a relationship for an existing person.
func CreateStakeholder(c *gin.Context) {
	actorNumber, organisationNumber, ok := stakeholderOrganisationContext(c)
	if !ok {
		return
	}

	var request requests.CreateStakeholderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}

	stakeholder, err := services.CreateStakeholder(actorNumber, organisationNumber, request)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	// Foundation G (G2): best-effort audit fact. Never changes the result.
	auditHooks.Record(c, organisationNumber, "stakeholder.create", "ems_stakeholder", stakeholder.Number, stakeholder)

	// Foundation G (G4): inform the affected stakeholder only (D-GQ1-11).
	notificationDecision.NotifyStakeholderChanged(organisationNumber, stakeholder.Number, stakeholder.UserNumber)

	c.JSON(http.StatusCreated, gin.H{"message": "Stakeholder created", "data": stakeholder})
}

// UpdateStakeholder updates descriptive/lifecycle fields.
func UpdateStakeholder(c *gin.Context) {
	_, organisationNumber, ok := stakeholderOrganisationContext(c)
	if !ok {
		return
	}

	var request requests.UpdateStakeholderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}

	stakeholder, err := services.UpdateStakeholder(organisationNumber, c.Param("number"), request)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	// Foundation G (G2): best-effort audit fact. Never changes the result.
	auditHooks.Record(c, organisationNumber, "stakeholder.update", "ems_stakeholder", stakeholder.Number, stakeholder)

	// Foundation G (G4): inform the affected stakeholder only (D-GQ1-11).
	notificationDecision.NotifyStakeholderChanged(organisationNumber, stakeholder.Number, stakeholder.UserNumber)

	c.JSON(http.StatusOK, gin.H{"message": "Stakeholder updated", "data": stakeholder})
}

// UpdateStakeholderStatus applies the ACTIVE / INACTIVE / ENDED lifecycle.
func UpdateStakeholderStatus(c *gin.Context) {
	_, organisationNumber, ok := stakeholderOrganisationContext(c)
	if !ok {
		return
	}

	var request requests.UpdateStakeholderStatusRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}

	stakeholder, err := services.SetStakeholderStatus(organisationNumber, c.Param("number"), request.Status)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	// Foundation G (G2): best-effort audit fact. Never changes the result.
	auditHooks.Record(c, organisationNumber, "stakeholder.status", "ems_stakeholder", stakeholder.Number, stakeholder)

	// Foundation G (G4): inform the affected stakeholder only (D-GQ1-11).
	notificationDecision.NotifyStakeholderChanged(organisationNumber, stakeholder.Number, stakeholder.UserNumber)

	c.JSON(http.StatusOK, gin.H{"message": "Stakeholder status updated", "data": stakeholder})
}

// stakeholderOrganisationContext resolves the authenticated actor and the
// authoritative organisation through the shared resolver.
func stakeholderOrganisationContext(c *gin.Context) (actorNumber string, organisationNumber string, ok bool) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, commons.Feedback("Unauthorized"))
		return "", "", false
	}

	organisationNumber, err = services.CurrentOrganisationNumber(user.Number)
	if err != nil {
		writeServiceError(c, err)
		return "", "", false
	}

	return user.Number, organisationNumber, true
}

// writeServiceError maps service errors to HTTP codes. Cross-organisation
// lookups surface as 404 so another organisation's data is never disclosed.
func writeServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, organisationServices.ErrNoActiveOrganisation):
		c.JSON(http.StatusConflict, commons.Feedback(err.Error()))
	case errors.Is(err, services.ErrStakeholderNotFound):
		c.JSON(http.StatusNotFound, commons.Feedback(err.Error()))
	case errors.Is(err, services.ErrDuplicateActiveStakeholder),
		errors.Is(err, services.ErrLastActiveOwner):
		c.JSON(http.StatusConflict, commons.Feedback(err.Error()))
	case errors.Is(err, services.ErrInvalidStakeholderType),
		errors.Is(err, services.ErrInvalidStatus),
		errors.Is(err, services.ErrInvalidSharePercentage),
		errors.Is(err, services.ErrInvalidDateRange),
		errors.Is(err, services.ErrUserRequired),
		errors.Is(err, services.ErrUserNotFound):
		c.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
	default:
		c.JSON(http.StatusInternalServerError, commons.Feedback("Something went wrong"))
	}
}
