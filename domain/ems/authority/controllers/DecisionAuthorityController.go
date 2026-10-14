package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	auditHooks "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/audit/hooks"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/authority/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/authority/services"
	organisationServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/organisation/services"
	"github.com/gin-gonic/gin"
)

// FetchDecisionAuthorities lists authority definitions for the authenticated organisation.
func FetchDecisionAuthorities(c *gin.Context) {
	_, organisationNumber, ok := authorityOrganisationContext(c)
	if !ok {
		return
	}

	authorities, err := services.ListDecisionAuthorities(organisationNumber, c.Query("decision_type"), c.Query("status"))
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Decision authorities fetched", "data": authorities})
}

// FetchDecisionAuthority returns one authority definition.
func FetchDecisionAuthority(c *gin.Context) {
	_, organisationNumber, ok := authorityOrganisationContext(c)
	if !ok {
		return
	}

	authority, err := services.GetDecisionAuthority(organisationNumber, c.Param("number"))
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Decision authority fetched", "data": authority})
}

// CreateDecisionAuthority creates an authority definition.
func CreateDecisionAuthority(c *gin.Context) {
	actorNumber, organisationNumber, ok := authorityOrganisationContext(c)
	if !ok {
		return
	}

	var request requests.CreateDecisionAuthorityRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}

	authority, err := services.CreateDecisionAuthority(actorNumber, organisationNumber, request)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	// Foundation G (G2): best-effort audit fact. Never changes the result.
	auditHooks.Record(c, organisationNumber, "authority.create", "ems_decision_authority", authority.Number, authority)

	c.JSON(http.StatusCreated, gin.H{"message": "Decision authority created", "data": authority})
}

// UpdateDecisionAuthority updates an authority definition.
func UpdateDecisionAuthority(c *gin.Context) {
	_, organisationNumber, ok := authorityOrganisationContext(c)
	if !ok {
		return
	}

	var request requests.UpdateDecisionAuthorityRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}

	authority, err := services.UpdateDecisionAuthority(organisationNumber, c.Param("number"), request)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	// Foundation G (G2): best-effort audit fact. Never changes the result.
	auditHooks.Record(c, organisationNumber, "authority.update", "ems_decision_authority", authority.Number, authority)

	c.JSON(http.StatusOK, gin.H{"message": "Decision authority updated", "data": authority})
}

// UpdateDecisionAuthorityStatus applies the ACTIVE / INACTIVE / ENDED lifecycle.
func UpdateDecisionAuthorityStatus(c *gin.Context) {
	_, organisationNumber, ok := authorityOrganisationContext(c)
	if !ok {
		return
	}

	var request requests.UpdateDecisionAuthorityStatusRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}

	authority, err := services.SetDecisionAuthorityStatus(organisationNumber, c.Param("number"), request.Status)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	// Foundation G (G2): best-effort audit fact. Never changes the result.
	auditHooks.Record(c, organisationNumber, "authority.status", "ems_decision_authority", authority.Number, authority)

	c.JSON(http.StatusOK, gin.H{"message": "Decision authority status updated", "data": authority})
}

// ResolveAuthority returns the required approval path for a decision type and
// amount. Read-only: it describes who is required; it creates nothing.
func ResolveAuthority(c *gin.Context) {
	_, organisationNumber, ok := authorityOrganisationContext(c)
	if !ok {
		return
	}

	decisionType := c.Query("decision_type")
	if strings.TrimSpace(decisionType) == "" {
		c.JSON(http.StatusBadRequest, commons.Feedback("decision_type is required"))
		return
	}

	amount, err := strconv.ParseFloat(strings.TrimSpace(c.Query("amount")), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("amount is required and must be a number"))
		return
	}

	resolution, err := services.ResolveAuthority(organisationNumber, decisionType, c.Query("organisation_unit"), amount)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Authority resolved", "data": resolution})
}

// authorityOrganisationContext resolves the authenticated actor and the
// authoritative organisation through the shared resolver.
func authorityOrganisationContext(c *gin.Context) (actorNumber string, organisationNumber string, ok bool) {
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
	case errors.Is(err, services.ErrAuthorityNotFound):
		c.JSON(http.StatusNotFound, commons.Feedback(err.Error()))
	case errors.Is(err, services.ErrInvalidDecisionType),
		errors.Is(err, services.ErrInvalidApprovalLevel),
		errors.Is(err, services.ErrInvalidAmount),
		errors.Is(err, services.ErrInvalidAmountRange),
		errors.Is(err, services.ErrInvalidAnchor),
		errors.Is(err, services.ErrInvalidStakeholderType),
		errors.Is(err, services.ErrRoleNotFound),
		errors.Is(err, services.ErrInvalidStatus),
		errors.Is(err, services.ErrInvalidDateRange),
		errors.Is(err, services.ErrUnitNotFound),
		errors.Is(err, services.ErrUnitNotActive):
		c.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
	default:
		c.JSON(http.StatusInternalServerError, commons.Feedback("Something went wrong"))
	}
}
