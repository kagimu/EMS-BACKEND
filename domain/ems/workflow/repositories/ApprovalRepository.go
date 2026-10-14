package repositories

import (
	"errors"
	"strings"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/stakeholders/models"
	workflowModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/workflow/models"
	"gorm.io/gorm"
)

// The read paths accept an optional transaction so multi-record operations stay
// inside one unit of work; every other call uses the shared connection.
func db(txn ...*gorm.DB) *gorm.DB {
	if len(txn) > 0 && txn[0] != nil {
		return txn[0]
	}
	return config.DB
}

func CreateRequest(request *workflowModels.ApprovalRequest, txn ...*gorm.DB) error {
	return db(txn...).Create(request).Error
}

func SaveRequest(request *workflowModels.ApprovalRequest, txn ...*gorm.DB) error {
	return db(txn...).Save(request).Error
}

func CreateStep(step *workflowModels.ApprovalStep, txn ...*gorm.DB) error {
	return db(txn...).Create(step).Error
}

func SaveStep(step *workflowModels.ApprovalStep, txn ...*gorm.DB) error {
	return db(txn...).Save(step).Error
}

func CreateAssignment(assignment *workflowModels.ApprovalAssignment, txn ...*gorm.DB) error {
	return db(txn...).Create(assignment).Error
}

func SaveAssignment(assignment *workflowModels.ApprovalAssignment, txn ...*gorm.DB) error {
	return db(txn...).Save(assignment).Error
}

// CreateStepAnchor writes one snapshotted unit-aware authority anchor of a step
// (Foundation D, ratified D6).
func CreateStepAnchor(anchor *workflowModels.ApprovalStepAnchor, txn ...*gorm.DB) error {
	return db(txn...).Create(anchor).Error
}

// FindStepAnchorsByStepNumbers returns the normalised anchor rows of the supplied
// steps in ONE query, so a request list is never an N+1 read.
func FindStepAnchorsByStepNumbers(stepNumbers []string, txn ...*gorm.DB) ([]workflowModels.ApprovalStepAnchor, error) {
	if len(stepNumbers) == 0 {
		return []workflowModels.ApprovalStepAnchor{}, nil
	}

	var anchors []workflowModels.ApprovalStepAnchor
	if err := db(txn...).
		Where("step_number IN ?", stepNumbers).
		Order("step_number asc, id asc").
		Find(&anchors).Error; err != nil {
		return nil, err
	}
	return anchors, nil
}

// FindRequestByNumberAndOrganisation is always organisation-scoped, which
// prevents cross-organisation (IDOR) disclosure of another organisation's request.
func FindRequestByNumberAndOrganisation(number, organisationNumber string, txn ...*gorm.DB) (*workflowModels.ApprovalRequest, error) {
	var request workflowModels.ApprovalRequest
	if err := db(txn...).
		Where("number = ? AND organisation_number = ?", number, organisationNumber).
		First(&request).Error; err != nil {
		return nil, err
	}
	return &request, nil
}

// LockRequestByNumberAndOrganisation takes a row lock (SELECT ... FOR UPDATE) so
// concurrent decisions on one request are serialised.
func LockRequestByNumberAndOrganisation(number, organisationNumber string, txn ...*gorm.DB) (*workflowModels.ApprovalRequest, error) {
	var request workflowModels.ApprovalRequest
	if err := db(txn...).
		Clauses(lockForUpdate()).
		Where("number = ? AND organisation_number = ?", number, organisationNumber).
		First(&request).Error; err != nil {
		return nil, err
	}
	return &request, nil
}

// FindAllByOrganisation lists requests of one organisation, optionally filtered.
//
// The result is always bounded: limit/offset come from the caller so an
// organisation with a long approval history cannot be loaded in one query.
func FindAllByOrganisation(organisationNumber, status, module, resourceType, resourceNumber string, limit, offset int) ([]workflowModels.ApprovalRequest, error) {
	var requests []workflowModels.ApprovalRequest

	query := config.DB.Where("organisation_number = ?", organisationNumber)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if module != "" {
		query = query.Where("module = ?", module)
	}
	if resourceType != "" {
		query = query.Where("resource_type = ?", resourceType)
	}
	if resourceNumber != "" {
		query = query.Where("resource_number = ?", resourceNumber)
	}

	if err := query.Order("created_at desc").Limit(limit).Offset(offset).Find(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}

// FindRequestsByNumbers returns the supplied requests in one query.
func FindRequestsByNumbers(requestNumbers []string, txn ...*gorm.DB) ([]workflowModels.ApprovalRequest, error) {
	if len(requestNumbers) == 0 {
		return []workflowModels.ApprovalRequest{}, nil
	}

	var approvals []workflowModels.ApprovalRequest
	if err := db(txn...).Where("number IN ?", requestNumbers).Find(&approvals).Error; err != nil {
		return nil, err
	}
	return approvals, nil
}

// FindStepsByNumbers returns the supplied steps in one query.
func FindStepsByNumbers(stepNumbers []string, txn ...*gorm.DB) ([]workflowModels.ApprovalStep, error) {
	if len(stepNumbers) == 0 {
		return []workflowModels.ApprovalStep{}, nil
	}

	var steps []workflowModels.ApprovalStep
	if err := db(txn...).Where("number IN ?", stepNumbers).Find(&steps).Error; err != nil {
		return nil, err
	}
	return steps, nil
}

// FindStepsByRequestNumbers returns the steps of the supplied requests in one
// query, ordered by level (never one query per request).
func FindStepsByRequestNumbers(requestNumbers []string, txn ...*gorm.DB) ([]workflowModels.ApprovalStep, error) {
	if len(requestNumbers) == 0 {
		return []workflowModels.ApprovalStep{}, nil
	}

	var steps []workflowModels.ApprovalStep
	if err := db(txn...).
		Where("request_number IN ?", requestNumbers).
		Order("request_number asc, approval_level asc").
		Find(&steps).Error; err != nil {
		return nil, err
	}
	return steps, nil
}

// FindAssignmentsByStepNumbers returns the candidate assignments of the supplied
// steps in one query.
func FindAssignmentsByStepNumbers(stepNumbers []string, txn ...*gorm.DB) ([]workflowModels.ApprovalAssignment, error) {
	if len(stepNumbers) == 0 {
		return []workflowModels.ApprovalAssignment{}, nil
	}

	var assignments []workflowModels.ApprovalAssignment
	if err := db(txn...).
		Where("step_number IN ?", stepNumbers).
		Order("step_number asc, created_at asc").
		Find(&assignments).Error; err != nil {
		return nil, err
	}
	return assignments, nil
}

// LockCurrentStep locks the step at the request's current level.
func LockCurrentStep(requestNumber string, approvalLevel int, organisationNumber string, txn ...*gorm.DB) (*workflowModels.ApprovalStep, error) {
	var step workflowModels.ApprovalStep
	if err := db(txn...).
		Clauses(lockForUpdate()).
		Where("request_number = ? AND approval_level = ? AND organisation_number = ?",
			requestNumber, approvalLevel, organisationNumber).
		First(&step).Error; err != nil {
		return nil, err
	}
	return &step, nil
}

// LockAssignment locks one candidate row and returns it only when it belongs to
// the given actor: a person who was never assigned to the step is not an
// approver for it, irrespective of permission.
func LockAssignment(stepNumber, userNumber, organisationNumber string, txn ...*gorm.DB) (*workflowModels.ApprovalAssignment, error) {
	var assignment workflowModels.ApprovalAssignment
	if err := db(txn...).
		Clauses(lockForUpdate()).
		Where("step_number = ? AND user_number = ? AND organisation_number = ?",
			stepNumber, userNumber, organisationNumber).
		First(&assignment).Error; err != nil {
		return nil, err
	}
	return &assignment, nil
}

// FindNextPendingStep returns the lowest step above the given level that is
// still awaiting a decision, or (nil, nil) when the request has no further level.
func FindNextPendingStep(requestNumber string, aboveLevel int, txn ...*gorm.DB) (*workflowModels.ApprovalStep, error) {
	var step workflowModels.ApprovalStep
	err := db(txn...).
		Where("request_number = ? AND approval_level > ? AND status = ?",
			requestNumber, aboveLevel, workflowModels.StepStatusPending).
		Order("approval_level asc").
		First(&step).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &step, nil
}

// CancelPendingSteps closes every step of a request that has not been decided.
//
// Both PENDING and BLOCKED steps move to CANCELLED: the approved blueprint's
// step state machine is `PENDING` / `BLOCKED` -> `CANCELLED` when the request is
// cancelled. Approved/rejected steps are left untouched as immutable history.
func CancelPendingSteps(requestNumber, organisationNumber string, txn ...*gorm.DB) error {
	return db(txn...).
		Model(&workflowModels.ApprovalStep{}).
		Where("request_number = ? AND organisation_number = ? AND status IN ?",
			requestNumber, organisationNumber,
			[]string{workflowModels.StepStatusPending, workflowModels.StepStatusBlocked}).
		Update("status", workflowModels.StepStatusCancelled).Error
}

// FindUserNumbersByRole returns the distinct users who currently hold a role in
// the organisation. Expired role assignments are excluded (BD-E4). The role is
// scoped to the organisation exactly as Foundation D scopes it.
func FindUserNumbersByRole(organisationNumber, roleNumber string, now time.Time) ([]string, error) {
	var rows []struct{ UserNumber string }
	err := config.DB.Table("user_roles AS ur").
		Select("ur.user_number AS user_number").
		Joins("JOIN roles AS r ON r.number = ur.role_number").
		Where("ur.deleted_at IS NULL AND r.deleted_at IS NULL").
		Where("ur.role_number = ?", roleNumber).
		Where("r.entity_number = ?", organisationNumber).
		Where("ur.expiry_date >= ?", now).
		Group("ur.user_number").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]string, 0, len(rows))
	for _, row := range rows {
		result = append(result, row.UserNumber)
	}
	return result, nil
}

// FindUserNumbersByStakeholderType returns the distinct users with an ACTIVE
// stakeholder relationship of the given type in the organisation (BD-E4).
//
// Foundation E2 (ratified D10): status alone is not sufficient. The stakeholder
// validity window must also contain now, with valid_to inclusive, matching the
// semantics established by Foundation C1 (which sets valid_to = now when a
// relationship is ended). A null valid_from/valid_to is unbounded on that side.
func FindUserNumbersByStakeholderType(organisationNumber, stakeholderType string, now time.Time) ([]string, error) {
	var rows []struct{ UserNumber string }
	err := config.DB.Model(&models.Stakeholder{}).
		Select("user_number AS user_number").
		Where("organisation_number = ? AND stakeholder_type = ? AND status = ?",
			organisationNumber, stakeholderType, constants.StatusActive).
		Where("(valid_from IS NULL OR valid_from <= ?)", now).
		Where("(valid_to IS NULL OR valid_to >= ?)", now).
		Group("user_number").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]string, 0, len(rows))
	for _, row := range rows {
		result = append(result, row.UserNumber)
	}
	return result, nil
}

// FindActionableAssignmentsForUser returns the caller's inbox: assignment rows
// that are still undecided on the current, pending step of a pending request of
// the organisation. One query, no N+1.
func FindActionableAssignmentsForUser(organisationNumber, userNumber string) ([]workflowModels.ApprovalAssignment, error) {
	var assignments []workflowModels.ApprovalAssignment

	err := config.DB.Table("ems_approval_assignments AS a").
		Select("a.*").
		Joins("JOIN ems_approval_steps AS s ON s.number = a.step_number AND s.deleted_at IS NULL").
		Joins("JOIN ems_approval_requests AS r ON r.number = a.request_number AND r.deleted_at IS NULL").
		Where("a.deleted_at IS NULL").
		Where("a.organisation_number = ? AND a.user_number = ?", organisationNumber, userNumber).
		Where("a.status = ?", workflowModels.AssignmentStatusUndecided).
		Where("s.status = ? AND s.approval_level = r.current_level", workflowModels.StepStatusPending).
		Where("r.status = ?", workflowModels.RequestStatusPending).
		Order("r.created_at asc").
		Scan(&assignments).Error
	if err != nil {
		return nil, err
	}
	return assignments, nil
}

// FindUserNumbersPlacedWithinUnitScope implements Foundation D's ratified
// candidate-resolution narrowing: D7 (source = placements), D3 = S2 (placement
// semantics) and D4 (any ACTIVE, in-window placement counts).
//
// A person qualifies when they hold an ems_user_assignments row that is
//   - in the same organisation,
//   - ACTIVE,
//   - currently inside its validity window (valid_from / valid_to, inclusive), and
//   - placed at the scope unit OR AT ANY DESCENDANT OF IT.
//
// A person placed ABOVE the scope unit — or in an unrelated branch — does NOT
// qualify: this is S2, and S1 is not implemented.
//
// An EMPTY scope unit is organisation-wide and is the caller's responsibility to
// skip before calling this function; it returns no users here.
//
// This deliberately does NOT reuse FindActiveByOrgUserUnitType, which filters
// status only and ignores valid_from/valid_to.
func FindUserNumbersPlacedWithinUnitScope(organisationNumber, scopeUnitNumber string, now time.Time) ([]string, error) {
	scope := strings.TrimSpace(scopeUnitNumber)
	if scope == "" {
		return []string{}, nil
	}

	units, err := UnitSubtree(organisationNumber, scope)
	if err != nil {
		return nil, err
	}
	if len(units) == 0 {
		return []string{}, nil
	}

	var rows []struct{ UserNumber string }
	err = config.DB.Table("ems_user_assignments").
		Select("user_number AS user_number").
		Where("organisation_number = ?", organisationNumber).
		Where("organisation_unit_number IN ?", units).
		Where("status = ?", constants.StatusActive).
		Where("(valid_from IS NULL OR valid_from <= ?)", now).
		Where("(valid_to IS NULL OR valid_to >= ?)", now).
		Group("user_number").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]string, 0, len(rows))
	for _, row := range rows {
		result = append(result, row.UserNumber)
	}
	return result, nil
}

// UnitSubtree returns the supplied unit together with every descendant unit in the
// same organisation, using the existing Foundation B hierarchy
// (ems_organisation_units.parent_number). It returns an empty slice when the unit
// does not exist in the organisation. Reads the unit set once — no per-level
// queries — mirroring Foundation D's ancestorSet.
func UnitSubtree(organisationNumber, rootUnitNumber string) ([]string, error) {
	var units []struct {
		Number       string
		ParentNumber string
	}
	if err := config.DB.Table("ems_organisation_units").
		Select("number AS number, parent_number AS parent_number").
		Where("organisation_number = ?", organisationNumber).
		Scan(&units).Error; err != nil {
		return nil, err
	}

	children := map[string][]string{}
	present := map[string]bool{}
	for _, unit := range units {
		present[unit.Number] = true
		if unit.ParentNumber != "" {
			children[unit.ParentNumber] = append(children[unit.ParentNumber], unit.Number)
		}
	}
	if !present[rootUnitNumber] {
		return []string{}, nil
	}

	subtree := []string{}
	seen := map[string]bool{}
	queue := []string{rootUnitNumber}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if seen[current] {
			continue
		}
		seen[current] = true
		subtree = append(subtree, current)
		queue = append(queue, children[current]...)
	}
	return subtree, nil
}
