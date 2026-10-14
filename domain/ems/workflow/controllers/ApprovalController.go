package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	auditHooks "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/audit/hooks"
	notificationDecision "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/notifications/decision"
	organisationServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/organisation/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/workflow/api/requests"
	workflowModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/workflow/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/workflow/services"
	"github.com/gin-gonic/gin"
)

// FetchApprovalRequests lists the organisation's approval requests.
func FetchApprovalRequests(c *gin.Context) {
	_, organisationNumber, ok := approvalOrganisationContext(c)
	if !ok {
		return
	}

	limit, err := parsePaginationParam(c.Query("limit"))
	if err != nil {
		writeServiceError(c, err)
		return
	}
	offset, err := parsePaginationParam(c.Query("offset"))
	if err != nil {
		writeServiceError(c, err)
		return
	}

	approvals, err := services.ListApprovalRequests(organisationNumber, c.Query("status"),
		c.Query("module"), c.Query("resource_type"), c.Query("resource_number"), limit, offset)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Approval requests fetched", "data": approvals})
}

// parsePaginationParam parses an optional non-negative page parameter. An absent
// value is 0, which the service maps to its default page size / first page.
func parsePaginationParam(raw string) (int, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return 0, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 0 {
		return 0, services.ErrInvalidPagination
	}
	return parsed, nil
}

// FetchApprovalRequest returns one approval request with its snapshotted path.
func FetchApprovalRequest(c *gin.Context) {
	_, organisationNumber, ok := approvalOrganisationContext(c)
	if !ok {
		return
	}

	approval, err := services.GetApprovalRequest(organisationNumber, c.Param("number"))
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Approval request fetched", "data": approval})
}

// CreateApprovalRequest creates a generic approval request. The originator is the
// authenticated user, who can never approve their own request.
func CreateApprovalRequest(c *gin.Context) {
	actorNumber, organisationNumber, ok := approvalOrganisationContext(c)
	if !ok {
		return
	}

	var request requests.CreateApprovalRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}

	approval, err := services.CreateApprovalRequest(actorNumber, organisationNumber, request)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	// Foundation G (G2): best-effort audit fact. Never changes the result.
	auditHooks.Record(c, organisationNumber, "approval.create", "ems_approval_request", approval.Number, approval)

	// Foundation E4 (E4-4, ratified E4-D9): an explicit per-request deadline
	// override is an authoritative deadline action and is recorded as its own
	// append-only audit fact. The configured default deadline is not a user
	// action and is therefore not recorded. Best-effort: it never changes the
	// business result.
	if request.DeadlineAt != nil && approval.DeadlineAt != nil {
		auditHooks.Record(c, organisationNumber, "approval.deadline_set", "ems_approval_request", approval.Number,
			map[string]interface{}{"deadline_at": approval.DeadlineAt})
	}

	// Foundation G (G4): notify the current-level approvers that a decision is
	// required. Best-effort and outside the business transaction (§19).
	notificationDecision.NotifyApprovalCreated(organisationNumber, approval.Number, approval.CurrentLevel)

	c.JSON(http.StatusCreated, gin.H{"message": "Approval request created", "data": approval})
}

// ApproveApprovalRequest records a valid approval for the current level.
func ApproveApprovalRequest(c *gin.Context) {
	actorNumber, organisationNumber, ok := approvalOrganisationContext(c)
	if !ok {
		return
	}

	var request requests.ApprovalDecisionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}

	approval, err := services.ApproveApprovalRequest(actorNumber, organisationNumber, c.Param("number"), request.Comment)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	// Foundation G (G2): best-effort audit fact. Never changes the result.
	auditHooks.Record(c, organisationNumber, "approval.approve", "ems_approval_request", approval.Number, approval)

	// Foundation G (G4): an advance notifies the new current level; a full
	// approval informs the requester (D-GQ1-1). Best-effort.
	switch approval.Status {
	case workflowModel.RequestStatusApproved:
		notificationDecision.NotifyApprovalApproved(organisationNumber, approval.Number, approval.RequestedByNumber)
	case workflowModel.RequestStatusPending:
		notificationDecision.NotifyApprovalAdvanced(organisationNumber, approval.Number, approval.CurrentLevel)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Approval recorded", "data": approval})
}

// RejectApprovalRequest rejects a request terminally. A reason is mandatory.
func RejectApprovalRequest(c *gin.Context) {
	actorNumber, organisationNumber, ok := approvalOrganisationContext(c)
	if !ok {
		return
	}

	var request requests.ApprovalDecisionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}

	approval, err := services.RejectApprovalRequest(actorNumber, organisationNumber, c.Param("number"), request.Comment)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	// Foundation G (G2): best-effort audit fact. Never changes the result.
	auditHooks.Record(c, organisationNumber, "approval.reject", "ems_approval_request", approval.Number, approval)

	// Foundation G (G4): inform the requester of the rejection (D-GQ1-1).
	notificationDecision.NotifyApprovalRejected(organisationNumber, approval.Number, approval.RequestedByNumber)

	c.JSON(http.StatusOK, gin.H{"message": "Approval request rejected", "data": approval})
}

// ReturnApprovalRequest implements Foundation E2 return-for-correction. An
// eligible approver of the current step sends the request back to the requester.
// A reason is mandatory and is limited to 1,000 Unicode characters.
func ReturnApprovalRequest(c *gin.Context) {
	actorNumber, organisationNumber, ok := approvalOrganisationContext(c)
	if !ok {
		return
	}

	var request requests.ApprovalDecisionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}

	approval, err := services.ReturnApprovalRequest(actorNumber, organisationNumber, c.Param("number"), request.Comment)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	// Foundation G (G2): best-effort audit fact. Never changes the result.
	auditHooks.Record(c, organisationNumber, "approval.return", "ems_approval_request", approval.Number, approval)

	// Foundation G (G4): the requester must correct and resubmit (D-GQ1-9).
	notificationDecision.NotifyApprovalReturned(organisationNumber, approval.Number, approval.RequestedByNumber)

	c.JSON(http.StatusOK, gin.H{"message": "Approval request returned", "data": approval})
}

// CancelApprovalRequest withdraws the caller's own request. Requester-only: a
// management permission never allows cancelling another user's request.
func CancelApprovalRequest(c *gin.Context) {
	actorNumber, organisationNumber, ok := approvalOrganisationContext(c)
	if !ok {
		return
	}

	approval, err := services.CancelApprovalRequest(actorNumber, organisationNumber, c.Param("number"))
	if err != nil {
		writeServiceError(c, err)
		return
	}

	// Foundation G (G2): best-effort audit fact. Never changes the result.
	auditHooks.Record(c, organisationNumber, "approval.cancel", "ems_approval_request", approval.Number, approval)

	c.JSON(http.StatusOK, gin.H{"message": "Approval request cancelled", "data": approval})
}

// FetchApprovalTasks returns the caller's actionable approvals.
func FetchApprovalTasks(c *gin.Context) {
	actorNumber, organisationNumber, ok := approvalOrganisationContext(c)
	if !ok {
		return
	}

	tasks, err := services.ListApprovalTasks(organisationNumber, actorNumber)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Approval tasks fetched", "data": tasks})
}

// approvalOrganisationContext resolves the authenticated actor and the
// authoritative organisation through the shared resolver. A client-supplied
// organisation number is never read.
func approvalOrganisationContext(c *gin.Context) (actorNumber string, organisationNumber string, ok bool) {
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
// lookups surface as 404 so another organisation's request is never disclosed.
func writeServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, organisationServices.ErrNoActiveOrganisation):
		c.JSON(http.StatusConflict, commons.Feedback(err.Error()))
	case errors.Is(err, services.ErrRequestNotFound):
		c.JSON(http.StatusNotFound, commons.Feedback(err.Error()))
	case errors.Is(err, services.ErrNotEligible),
		errors.Is(err, services.ErrNotRequestOwner):
		c.JSON(http.StatusForbidden, commons.Feedback(err.Error()))
	case errors.Is(err, services.ErrModuleRequired),
		errors.Is(err, services.ErrResourceTypeRequired),
		errors.Is(err, services.ErrResourceNumberRequired),
		errors.Is(err, services.ErrDecisionTypeRequired),
		errors.Is(err, services.ErrInvalidFieldLength),
		errors.Is(err, services.ErrInvalidAmount),
		errors.Is(err, services.ErrInvalidRequestStatus),
		errors.Is(err, services.ErrInvalidPagination),
		errors.Is(err, services.ErrUnitNotFound),
		errors.Is(err, services.ErrRejectionReasonRequired),
		errors.Is(err, services.ErrReturnReasonRequired),
		errors.Is(err, services.ErrPreviousRequestInvalid),
		errors.Is(err, services.ErrInvalidDeadline),
		errors.Is(err, services.ErrInvalidDeadlineHours):
		c.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
	case errors.Is(err, services.ErrDecisionNotConfigured),
		errors.Is(err, services.ErrNoApplicableAuthority),
		errors.Is(err, services.ErrExceedsDeclaredAuthority),
		errors.Is(err, services.ErrRequestNotActionable),
		errors.Is(err, services.ErrStepNotActionable),
		errors.Is(err, services.ErrAlreadyDecided),
		errors.Is(err, services.ErrRequestNotCancellable):
		c.JSON(http.StatusConflict, commons.Feedback(err.Error()))
	default:
		c.JSON(http.StatusInternalServerError, commons.Feedback("Something went wrong"))
	}
}
