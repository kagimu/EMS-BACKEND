package controllers

import (
	"errors"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	auditHooks "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/audit/hooks"
	organisationServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/organisation/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/reporting/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/reporting/services"
	"github.com/gin-gonic/gin"
)

// FetchReportingLines lists reporting lines for the authenticated organisation.
func FetchReportingLines(c *gin.Context) {
	_, organisationNumber, ok := reportingOrganisationContext(c)
	if !ok {
		return
	}

	lines, err := services.ListReportingLines(
		organisationNumber,
		c.Query("subordinate_assignment_number"),
		c.Query("manager_assignment_number"),
		c.Query("status"),
	)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reporting lines fetched", "data": lines})
}

// FetchReportingLine returns one reporting line from the authenticated organisation.
func FetchReportingLine(c *gin.Context) {
	_, organisationNumber, ok := reportingOrganisationContext(c)
	if !ok {
		return
	}

	line, err := services.GetReportingLine(organisationNumber, c.Param("number"))
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reporting line fetched", "data": line})
}

// CreateReportingLine creates a line between two existing assignments.
func CreateReportingLine(c *gin.Context) {
	actorNumber, organisationNumber, ok := reportingOrganisationContext(c)
	if !ok {
		return
	}

	var request requests.CreateReportingLineRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}

	line, err := services.CreateReportingLine(actorNumber, organisationNumber, request)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	auditHooks.Record(c, organisationNumber, "reporting.create", "ems_reporting_line", line.Number, line)

	c.JSON(http.StatusCreated, gin.H{"message": "Reporting line created", "data": line})
}

// UpdateReportingLine updates the manager and lifecycle fields.
func UpdateReportingLine(c *gin.Context) {
	_, organisationNumber, ok := reportingOrganisationContext(c)
	if !ok {
		return
	}

	var request requests.UpdateReportingLineRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}

	line, err := services.UpdateReportingLine(organisationNumber, c.Param("number"), request)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	auditHooks.Record(c, organisationNumber, "reporting.update", "ems_reporting_line", line.Number, line)

	c.JSON(http.StatusOK, gin.H{"message": "Reporting line updated", "data": line})
}

// UpdateReportingLineStatus applies the ACTIVE / INACTIVE / ENDED lifecycle.
func UpdateReportingLineStatus(c *gin.Context) {
	_, organisationNumber, ok := reportingOrganisationContext(c)
	if !ok {
		return
	}

	var request requests.UpdateReportingLineStatusRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}

	line, err := services.SetReportingLineStatus(organisationNumber, c.Param("number"), request.Status)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	auditHooks.Record(c, organisationNumber, "reporting.status", "ems_reporting_line", line.Number, line)

	c.JSON(http.StatusOK, gin.H{"message": "Reporting line status updated", "data": line})
}

// reportingOrganisationContext resolves the authenticated actor and the
// authoritative organisation through the shared resolver.
func reportingOrganisationContext(c *gin.Context) (actorNumber string, organisationNumber string, ok bool) {
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

// writeServiceError maps reporting errors to HTTP codes. Cross-organisation
// assignment lookups surface as 404 so another organisation's data is never
// disclosed.
func writeServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, organisationServices.ErrNoActiveOrganisation):
		c.JSON(http.StatusConflict, commons.Feedback(err.Error()))
	case errors.Is(err, services.ErrReportingLineNotFound),
		errors.Is(err, services.ErrAssignmentNotFound):
		c.JSON(http.StatusNotFound, commons.Feedback(err.Error()))
	case errors.Is(err, services.ErrDuplicateActiveReportingLine),
		errors.Is(err, services.ErrReportingCycle),
		errors.Is(err, services.ErrReportingChangeBusy):
		c.JSON(http.StatusConflict, commons.Feedback(err.Error()))
	case errors.Is(err, services.ErrSubordinateRequired),
		errors.Is(err, services.ErrManagerRequired),
		errors.Is(err, services.ErrSelfReporting),
		errors.Is(err, services.ErrAssignmentNotActive),
		errors.Is(err, services.ErrInvalidStatus),
		errors.Is(err, services.ErrInvalidDateRange),
		errors.Is(err, services.ErrReactivationRequiresDates):
		c.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
	default:
		c.JSON(http.StatusInternalServerError, commons.Feedback("Something went wrong"))
	}
}
