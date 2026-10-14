package controllers

import (
	"errors"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	auditHooks "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/audit/hooks"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/organisation/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/organisation/services"
	"github.com/gin-gonic/gin"
)

// FetchOrganisationUnits lists units for the authenticated organisation.
func FetchOrganisationUnits(c *gin.Context) {
	_, organisationNumber, ok := organisationContext(c)
	if !ok {
		return
	}

	units, err := services.ListUnits(organisationNumber, c.Query("parent"), c.Query("status"))
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Organisation units fetched", "data": units})
}

// FetchOrganisationUnit returns a single unit from the authenticated organisation.
func FetchOrganisationUnit(c *gin.Context) {
	_, organisationNumber, ok := organisationContext(c)
	if !ok {
		return
	}

	unit, err := services.GetUnit(organisationNumber, c.Param("number"))
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Organisation unit fetched", "data": unit})
}

// CreateOrganisationUnit creates a unit in the authenticated organisation.
func CreateOrganisationUnit(c *gin.Context) {
	actorNumber, organisationNumber, ok := organisationContext(c)
	if !ok {
		return
	}

	var request requests.CreateOrganisationUnitRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}

	unit, err := services.CreateUnit(actorNumber, organisationNumber, request)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	// Foundation G (G2): best-effort audit fact. Never changes the result.
	auditHooks.Record(c, organisationNumber, "organisation_unit.create", "ems_organisation_unit", unit.Number, unit)

	c.JSON(http.StatusCreated, gin.H{"message": "Organisation unit created", "data": unit})
}

// UpdateOrganisationUnit updates editable fields of a unit.
func UpdateOrganisationUnit(c *gin.Context) {
	_, organisationNumber, ok := organisationContext(c)
	if !ok {
		return
	}

	var request requests.UpdateOrganisationUnitRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}

	unit, err := services.UpdateUnit(organisationNumber, c.Param("number"), request)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	// Foundation G (G2): best-effort audit fact. Never changes the result.
	auditHooks.Record(c, organisationNumber, "organisation_unit.update", "ems_organisation_unit", unit.Number, unit)

	c.JSON(http.StatusOK, gin.H{"message": "Organisation unit updated", "data": unit})
}

// UpdateOrganisationUnitStatus activates or deactivates a unit.
func UpdateOrganisationUnitStatus(c *gin.Context) {
	_, organisationNumber, ok := organisationContext(c)
	if !ok {
		return
	}

	var request requests.UpdateOrganisationUnitStatusRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}

	unit, err := services.SetUnitStatus(organisationNumber, c.Param("number"), request.Status)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	// Foundation G (G2): best-effort audit fact. Never changes the result.
	auditHooks.Record(c, organisationNumber, "organisation_unit.status", "ems_organisation_unit", unit.Number, unit)

	c.JSON(http.StatusOK, gin.H{"message": "Organisation unit status updated", "data": unit})
}

// DeleteOrganisationUnit soft-deletes a unit that has no children.
func DeleteOrganisationUnit(c *gin.Context) {
	_, organisationNumber, ok := organisationContext(c)
	if !ok {
		return
	}

	if err := services.DeleteUnit(organisationNumber, c.Param("number")); err != nil {
		writeServiceError(c, err)
		return
	}

	// Foundation G (G2): best-effort audit fact. Never changes the result.
	auditHooks.Record(c, organisationNumber, "organisation_unit.delete", "ems_organisation_unit", c.Param("number"), nil)

	c.JSON(http.StatusOK, commons.Feedback("Organisation unit deleted"))
}

// organisationContext resolves the authenticated actor and the authoritative
// organisation. It writes the error response itself and returns ok=false when
// the request must not proceed.
func organisationContext(c *gin.Context) (actorNumber string, organisationNumber string, ok bool) {
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

// writeServiceError maps service errors to HTTP status codes. Cross-organisation
// lookups surface as 404 so the API never discloses another organisation's data.
func writeServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrNoActiveOrganisation):
		c.JSON(http.StatusConflict, commons.Feedback(err.Error()))
	case errors.Is(err, services.ErrUnitNotFound):
		c.JSON(http.StatusNotFound, commons.Feedback(err.Error()))
	case errors.Is(err, services.ErrHasChildren),
		errors.Is(err, services.ErrHasActiveAssignments),
		errors.Is(err, services.ErrDuplicateCode):
		c.JSON(http.StatusConflict, commons.Feedback(err.Error()))
	case errors.Is(err, services.ErrParentNotFound),
		errors.Is(err, services.ErrParentInactive),
		errors.Is(err, services.ErrSelfParent),
		errors.Is(err, services.ErrCircularHierarchy),
		errors.Is(err, services.ErrInvalidStatus),
		errors.Is(err, services.ErrNameRequired),
		errors.Is(err, services.ErrUnitTypeRequired):
		c.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
	default:
		c.JSON(http.StatusInternalServerError, commons.Feedback("Something went wrong"))
	}
}
