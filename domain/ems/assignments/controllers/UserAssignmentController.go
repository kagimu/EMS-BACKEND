package controllers

import (
	"errors"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/assignments/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/assignments/services"
	auditHooks "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/audit/hooks"
	notificationDecision "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/notifications/decision"
	organisationServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/organisation/services"
	"github.com/gin-gonic/gin"
)

// FetchUserAssignments lists assignments for the authenticated organisation.
func FetchUserAssignments(c *gin.Context) {
	_, organisationNumber, ok := assignmentOrganisationContext(c)
	if !ok {
		return
	}

	assignments, err := services.ListAssignments(
		organisationNumber,
		c.Query("organisation_unit_number"),
		c.Query("user_number"),
		c.Query("status"),
	)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User assignments fetched", "data": assignments})
}

// FetchUserAssignment returns one assignment from the authenticated organisation.
func FetchUserAssignment(c *gin.Context) {
	_, organisationNumber, ok := assignmentOrganisationContext(c)
	if !ok {
		return
	}

	assignment, err := services.GetAssignment(organisationNumber, c.Param("number"))
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User assignment fetched", "data": assignment})
}

// CreateUserAssignment creates an assignment for an existing person at an
// existing active unit of the authenticated organisation.
func CreateUserAssignment(c *gin.Context) {
	actorNumber, organisationNumber, ok := assignmentOrganisationContext(c)
	if !ok {
		return
	}

	var request requests.CreateUserAssignmentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}

	assignment, err := services.CreateAssignment(actorNumber, organisationNumber, request)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	// Foundation G (G2): best-effort audit fact. Never changes the result.
	auditHooks.Record(c, organisationNumber, "assignment.create", "ems_user_assignment", assignment.Number, assignment)

	// Foundation G (G4): a new assignment is an actionable responsibility for
	// the assignee (D-GQ1-13).
	notificationDecision.NotifyAssignmentCreated(organisationNumber, assignment.Number, assignment.UserNumber)

	c.JSON(http.StatusCreated, gin.H{"message": "User assignment created", "data": assignment})
}

// UpdateUserAssignment updates placement and lifecycle fields.
func UpdateUserAssignment(c *gin.Context) {
	_, organisationNumber, ok := assignmentOrganisationContext(c)
	if !ok {
		return
	}

	var request requests.UpdateUserAssignmentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}

	assignment, err := services.UpdateAssignment(organisationNumber, c.Param("number"), request)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	// Foundation G (G2): best-effort audit fact. Never changes the result.
	auditHooks.Record(c, organisationNumber, "assignment.update", "ems_user_assignment", assignment.Number, assignment)

	// Foundation G (G4): inform the assignee (D-GQ1-12).
	notificationDecision.NotifyAssignmentChanged(organisationNumber, assignment.Number, assignment.UserNumber)

	c.JSON(http.StatusOK, gin.H{"message": "User assignment updated", "data": assignment})
}

// UpdateUserAssignmentStatus applies the ACTIVE / INACTIVE / ENDED lifecycle.
func UpdateUserAssignmentStatus(c *gin.Context) {
	_, organisationNumber, ok := assignmentOrganisationContext(c)
	if !ok {
		return
	}

	var request requests.UpdateUserAssignmentStatusRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}

	assignment, err := services.SetAssignmentStatus(organisationNumber, c.Param("number"), request.Status)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	// Foundation G (G2): best-effort audit fact. Never changes the result.
	auditHooks.Record(c, organisationNumber, "assignment.status", "ems_user_assignment", assignment.Number, assignment)

	// Foundation G (G4): inform the assignee (D-GQ1-12).
	notificationDecision.NotifyAssignmentChanged(organisationNumber, assignment.Number, assignment.UserNumber)

	c.JSON(http.StatusOK, gin.H{"message": "User assignment status updated", "data": assignment})
}

// assignmentOrganisationContext resolves the authenticated actor and the
// authoritative organisation through the shared resolver.
func assignmentOrganisationContext(c *gin.Context) (actorNumber string, organisationNumber string, ok bool) {
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
	case errors.Is(err, services.ErrAssignmentNotFound):
		c.JSON(http.StatusNotFound, commons.Feedback(err.Error()))
	case errors.Is(err, services.ErrDuplicateActiveAssignment):
		c.JSON(http.StatusConflict, commons.Feedback(err.Error()))
	case errors.Is(err, services.ErrInvalidAssignmentType),
		errors.Is(err, services.ErrInvalidStatus),
		errors.Is(err, services.ErrInvalidDateRange),
		errors.Is(err, services.ErrUserRequired),
		errors.Is(err, services.ErrUserNotFound),
		errors.Is(err, services.ErrUnitRequired),
		errors.Is(err, services.ErrUnitNotFound),
		errors.Is(err, services.ErrUnitNotActive):
		c.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
	default:
		c.JSON(http.StatusInternalServerError, commons.Feedback("Something went wrong"))
	}
}
