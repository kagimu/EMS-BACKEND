package services

import (
	"errors"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	authorityResponses "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/authority/api/responses"
	authorityServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/authority/services"
	organisationServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/organisation/services"
	stakeholderRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/stakeholders/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/workflow/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/workflow/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/workflow/models"
	approvalRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/workflow/repositories"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Field limits mirror the column widths of the three EMS approval tables.
const (
	maxModuleLen         = 50
	maxResourceTypeLen   = 50
	maxResourceNumberLen = 36
	maxDecisionTypeLen   = 50
	// maxCommentLen is the shared limit for approval comments, rejection reasons
	// and return reasons. Ratified E2 decision D3 fixes it at 1,000 **Unicode
	// characters** (runes), not UTF-8 bytes, matching the widened
	// ems_approval_assignments.comment column (varchar(1000)).
	maxCommentLen = 1000

	// The request list is always bounded so a long approval history cannot be
	// loaded in one response. Callers may page with limit/offset.
	defaultListLimit = 50
	maxListLimit     = 200
)

var (
	ErrRequestNotFound        = errors.New("approval request not found")
	ErrModuleRequired         = errors.New("module is required")
	ErrResourceTypeRequired   = errors.New("resource_type is required")
	ErrResourceNumberRequired = errors.New("resource_number is required")
	ErrDecisionTypeRequired   = errors.New("decision_type is required")
	ErrInvalidFieldLength     = errors.New("a supplied field exceeds its maximum length")
	ErrInvalidAmount          = errors.New("amount must not be negative")
	ErrInvalidDeadline        = errors.New("deadline must be a future UTC timestamp")
	ErrInvalidRequestStatus   = errors.New("invalid request status")
	ErrInvalidPagination      = errors.New("limit and offset must be integers")
	ErrUnitNotFound           = errors.New("organisation unit was not found in this organisation")

	// Foundation D outcomes that do not yield an approval path.
	ErrDecisionNotConfigured    = errors.New("no decision authority is configured for this decision type")
	ErrNoApplicableAuthority    = errors.New("no decision authority applies to this request")
	ErrExceedsDeclaredAuthority = errors.New("the amount exceeds the declared decision authority")

	// Request/step action guards.
	ErrRequestNotActionable    = errors.New("approval request is not awaiting a decision")
	ErrStepNotActionable       = errors.New("approval step is not awaiting a decision")
	ErrNotEligible             = errors.New("the user is not an eligible approver for this step")
	ErrAlreadyDecided          = errors.New("this approval has already been decided")
	ErrRejectionReasonRequired = errors.New("a reason is required to reject an approval request")
	ErrReturnReasonRequired    = errors.New("a reason is required to return an approval request")
	ErrNotRequestOwner         = errors.New("only the requester may cancel this approval request")
	ErrRequestNotCancellable   = errors.New("approval request cannot be cancelled in its current status")

	// Foundation E2 resubmission linkage (ratified D5-D8). One generic error is
	// used for every failure so a create never becomes an existence oracle for
	// another organisation's or another user's request numbers.
	ErrPreviousRequestInvalid = errors.New("the referenced previous approval request is not resubmittable by this user")
)

// builtInDefaultDeadlineHours is the Foundation E4 (ratified E4-D4) built-in
// default approval deadline (7 days, 168 hours) used when an organisation has no
// explicit workflow deadline configuration. It is also the fallback for requests
// created before the E4-2 persistence invariant.
const builtInDefaultDeadlineHours = 168

// maxDeadlineHours bounds the configurable default at 365 days. A longer value is
// operationally indistinguishable from "no deadline", and 365 days is the
// largest horizon over which an advisory read-time deadline is meaningful.
const maxDeadlineHours = 8760

// reminderWindow is the Foundation E4 (ratified C4) reminder window. The derived
// "approaching deadline" condition is active when the request is PENDING, the
// effective deadline is still in the future, and no more than 24 hours remain.
// Reminder and overdue are mutually exclusive (overdue requires nowUTC > DeadlineAt).
const reminderWindow = 24 * time.Hour

// anchor is one snapshotted authority anchor.
//
// Foundation D (ratified D5): the unit is part of the anchor's identity. An empty
// unit means organisation-wide scope.
type anchor struct {
	Type  string
	Value string
	Unit  string
}

// isValidRequestStatus reports whether the value is a request state defined by
// Foundation E1/E2. No other state exists.
func isValidRequestStatus(status string) bool {
	switch status {
	case models.RequestStatusPending, models.RequestStatusApproved, models.RequestStatusRejected,
		models.RequestStatusCancelled, models.RequestStatusBlocked, models.RequestStatusReturned:
		return true
	default:
		return false
	}
}

// serializeAnchors stores the anchors of one level in a compact, parseable form.
//
// Foundation D returns a LIST of anchors per level; the eligible candidates for
// that level are the union across its anchors. The serialised value is the
// historical explanation of "why these people were asked" and is never rebuilt
// from current authority (BD-E3/BD-E14).
func serializeAnchors(anchors []authorityResponses.AuthorityAnchorResponse) string {
	parts := make([]string, 0, len(anchors))
	for _, item := range anchors {
		if strings.TrimSpace(item.Type) == "" || strings.TrimSpace(item.Value) == "" {
			continue
		}
		parts = append(parts, item.Type+":"+item.Value)
	}
	return strings.Join(parts, "|")
}

// parseAnchors is the inverse of serializeAnchors. Malformed fragments are
// skipped rather than guessed.
func parseAnchors(value string) []anchor {
	fragments := strings.Split(value, "|")
	result := make([]anchor, 0, len(fragments))

	for _, fragment := range fragments {
		index := strings.Index(fragment, ":")
		if index <= 0 || index == len(fragment)-1 {
			continue
		}
		result = append(result, anchor{Type: fragment[:index], Value: fragment[index+1:]})
	}
	return result
}

// loadStepAnchors returns the snapshotted anchors of the supplied steps, keyed by
// step number, in ONE query (no N+1).
//
// Foundation D (ratified D6 / D9 / D11): steps created after the amendment carry
// their anchors as normalised rows in ems_approval_step_anchors, with the unit as
// part of each anchor. Steps created before it have no such rows and fall back to
// the legacy serialised column, whose anchors carry NO unit and are therefore read
// as organisation-wide — so a legacy request is never retroactively narrowed.
func loadStepAnchors(steps []models.ApprovalStep) (map[string][]anchor, error) {
	result := make(map[string][]anchor, len(steps))
	if len(steps) == 0 {
		return result, nil
	}

	numbers := make([]string, 0, len(steps))
	for _, step := range steps {
		numbers = append(numbers, step.Number)
	}

	rows, err := approvalRepositories.FindStepAnchorsByStepNumbers(numbers)
	if err != nil {
		return nil, err
	}

	byStep := map[string][]anchor{}
	for _, row := range rows {
		byStep[row.StepNumber] = append(byStep[row.StepNumber], anchor{
			Type:  row.AnchorType,
			Value: row.AnchorValue,
			Unit:  row.OrganisationUnitNumber,
		})
	}

	for _, step := range steps {
		if anchors, ok := byStep[step.Number]; ok && len(anchors) > 0 {
			result[step.Number] = anchors
			continue
		}
		// Legacy step: read the unit-less serialised column, which is
		// organisation-wide by definition.
		result[step.Number] = parseAnchors(step.Anchors)
	}
	return result, nil
}

// stepAnchors returns the snapshotted anchors of one step. It is the reader used
// by the action-time eligibility re-check.
func stepAnchors(step models.ApprovalStep) ([]anchor, error) {
	byStep, err := loadStepAnchors([]models.ApprovalStep{step})
	if err != nil {
		return nil, err
	}
	return byStep[step.Number], nil
}

// removeUser drops one user number from a candidate set. It implements BD-E6
// (self-approval is prohibited): the requester is removed from every level. If
// that leaves no candidate the level is BLOCKED; no other candidate is invented.
func removeUser(candidates []string, userNumber string) []string {
	result := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		if candidate != userNumber {
			result = append(result, candidate)
		}
	}
	return result
}

// contains reports whether the slice holds the value.
func contains(values []string, value string) bool {
	for _, item := range values {
		if item == value {
			return true
		}
	}
	return false
}

// validateRequestInput validates the generic request fields E1 owns. It performs
// no module-specific business validation, and it never trusts a client-supplied
// organisation.
func validateRequestInput(module, resourceType, resourceNumber, decisionType string, amount *float64) error {
	if module == "" {
		return ErrModuleRequired
	}
	if resourceType == "" {
		return ErrResourceTypeRequired
	}
	if resourceNumber == "" {
		return ErrResourceNumberRequired
	}
	if decisionType == "" {
		return ErrDecisionTypeRequired
	}
	if len(module) > maxModuleLen || len(resourceType) > maxResourceTypeLen ||
		len(resourceNumber) > maxResourceNumberLen || len(decisionType) > maxDecisionTypeLen {
		return ErrInvalidFieldLength
	}
	if amount != nil && *amount < 0 {
		return ErrInvalidAmount
	}
	return nil
}

// mapAuthorityError keeps Foundation D's errors at the E boundary so callers do
// not depend on D directly. It changes nothing in D.
func mapAuthorityError(err error) error {
	switch {
	case errors.Is(err, authorityServices.ErrUnitNotFound):
		return ErrUnitNotFound
	case errors.Is(err, authorityServices.ErrInvalidDecisionType):
		return ErrDecisionTypeRequired
	case errors.Is(err, authorityServices.ErrInvalidAmount):
		return ErrInvalidAmount
	default:
		return err
	}
}

// userNumbersForAnchor returns the users currently implied by one snapshotted
// anchor. Snapshot-time candidate resolution and action-time eligibility
// re-checking both go through this single function so the two derivations can
// never drift apart.
//
// Foundation D (ratified D7 / D3 = S2 / D4): when the anchor carries a unit, the
// anchor holders are narrowed to people whose ems_user_assignments placement is
// that unit or a descendant of it, is ACTIVE and is currently inside its validity
// window. An EMPTY unit is organisation-wide and applies NO placement filter.
//
// Placement is only ever a NARROWING condition: authority still comes from the
// anchor, so being placed somewhere never grants approval authority by itself.
func userNumbersForAnchor(organisationNumber string, item anchor, now time.Time) ([]string, error) {
	var holders []string

	switch item.Type {
	case "role":
		numbers, err := approvalRepositories.FindUserNumbersByRole(organisationNumber, item.Value, now)
		if err != nil {
			return nil, err
		}
		holders = numbers
	case "stakeholder_type":
		numbers, err := approvalRepositories.FindUserNumbersByStakeholderType(organisationNumber, item.Value, now)
		if err != nil {
			return nil, err
		}
		holders = numbers
	default:
		// An anchor type E does not understand yields no candidate; it must never
		// be interpreted as "everyone".
		return nil, nil
	}

	if strings.TrimSpace(item.Unit) == "" {
		return holders, nil
	}

	placed, err := approvalRepositories.FindUserNumbersPlacedWithinUnitScope(organisationNumber, item.Unit, now)
	if err != nil {
		return nil, err
	}

	allowed := make(map[string]bool, len(placed))
	for _, number := range placed {
		allowed[number] = true
	}

	narrowed := make([]string, 0, len(holders))
	for _, number := range holders {
		if allowed[number] {
			narrowed = append(narrowed, number)
		}
	}
	return narrowed, nil
}

// toAnchors converts the resolver's anchor list into the workflow's internal
// unit-aware anchor form. Foundation D (ratified D5): the unit is part of the
// anchor's identity.
func toAnchors(items []authorityResponses.AuthorityAnchorResponse) []anchor {
	result := make([]anchor, 0, len(items))
	for _, item := range items {
		result = append(result, anchor{Type: item.Type, Value: item.Value, Unit: item.OrganisationUnitNumber})
	}
	return result
}

// resolveCandidates resolves the eligible approvers for one level: the union of
// the users implied by its anchors, restricted to people who can actually sign
// in (Foundation C2's login derivation) and ordered deterministically for a
// stable snapshot.
//
// The ordering is storage stability only: it never selects a primary approver.
// Every snapshot candidate is equally actionable (BD-E5).
func resolveCandidates(organisationNumber string, anchors []anchor) ([]string, error) {
	set := map[string]bool{}
	now := time.Now()

	for _, item := range anchors {
		numbers, err := userNumbersForAnchor(organisationNumber, item, now)
		if err != nil {
			return nil, err
		}
		for _, number := range numbers {
			set[number] = true
		}
	}

	candidates := make([]string, 0, len(set))
	for number := range set {
		candidates = append(candidates, number)
	}
	if len(candidates) == 0 {
		return candidates, nil
	}
	sort.Strings(candidates)

	authenticatable, err := stakeholderRepositories.FindAuthenticatableUserNumbers(candidates)
	if err != nil {
		return nil, err
	}

	eligible := make([]string, 0, len(candidates))
	for _, number := range candidates {
		if authenticatable[number] {
			eligible = append(eligible, number)
		}
	}
	return eligible, nil
}

// eligibilityResolver derives CURRENT candidate eligibility at READ TIME.
//
// Foundation E3 (ratified D2 + D3): the stalled condition and the inbox filter are
// derived when the data is read — never stored and never scheduled. The resolver
// memoises anchor -> eligible users so that a page of requests or a task list
// resolves each DISTINCT anchor once (E3 §21), with no N+1 pattern.
//
// It holds nothing beyond the call, is never persisted, and cannot become an
// authoritative state. The authoritative gate for acting remains the action-time
// check inside the locked decision transaction.
type eligibilityResolver struct {
	organisationNumber string
	now                time.Time
	anchorsByStep      map[string][]anchor
	usersByAnchor      map[string]map[string]bool
	authenticatable    map[string]bool
}

func newEligibilityResolver(organisationNumber string) *eligibilityResolver {
	return &eligibilityResolver{
		organisationNumber: organisationNumber,
		now:                time.Now(),
		usersByAnchor:      map[string]map[string]bool{},
	}
}

// withStepAnchors pre-seeds the snapshotted anchors of every step in the read, so
// a list never queries per step.
func (r *eligibilityResolver) withStepAnchors(anchorsByStep map[string][]anchor) *eligibilityResolver {
	r.anchorsByStep = anchorsByStep
	return r
}

// withAuthenticatable pre-seeds the login-eligibility set, so a list resolves it in
// ONE batched query rather than one query per candidate.
func (r *eligibilityResolver) withAuthenticatable(users map[string]bool) *eligibilityResolver {
	r.authenticatable = users
	return r
}

// stepAnchorItems returns a step's snapshotted anchors, using the pre-seeded set
// when available.
func (r *eligibilityResolver) stepAnchorItems(step models.ApprovalStep) ([]anchor, error) {
	if r.anchorsByStep != nil {
		if items, ok := r.anchorsByStep[step.Number]; ok {
			return items, nil
		}
	}
	return stepAnchors(step)
}

// anchorUsers returns the currently eligible users implied by one anchor, memoised
// per distinct anchor. An empty unit is organisation-wide (no placement filter).
func (r *eligibilityResolver) anchorUsers(item anchor) (map[string]bool, error) {
	if strings.TrimSpace(item.Type) == "" || strings.TrimSpace(item.Value) == "" {
		return map[string]bool{}, nil
	}

	key := item.Type + "|" + item.Value + "|" + item.Unit
	if cached, ok := r.usersByAnchor[key]; ok {
		return cached, nil
	}

	numbers, err := userNumbersForAnchor(r.organisationNumber, item, r.now)
	if err != nil {
		return nil, err
	}
	set := make(map[string]bool, len(numbers))
	for _, number := range numbers {
		set[number] = true
	}
	r.usersByAnchor[key] = set
	return set, nil
}

// authenticatableUser reports whether a user can currently sign in, using the
// pre-seeded set when available.
func (r *eligibilityResolver) authenticatableUser(userNumber string) (bool, error) {
	if r.authenticatable != nil {
		return r.authenticatable[userNumber], nil
	}
	users, err := stakeholderRepositories.FindAuthenticatableUserNumbers([]string{userNumber})
	if err != nil {
		return false, err
	}
	return users[userNumber], nil
}

// userEligibleForStep reports whether one user is currently eligible for a step,
// using the SAME anchor + placement derivation as action-time authorization. It
// never writes and never modifies the snapshot.
func (r *eligibilityResolver) userEligibleForStep(step models.ApprovalStep, userNumber string) (bool, error) {
	ok, err := r.authenticatableUser(userNumber)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}

	items, err := r.stepAnchorItems(step)
	if err != nil {
		return false, err
	}

	for _, item := range items {
		users, err := r.anchorUsers(item)
		if err != nil {
			return false, err
		}
		if users[userNumber] {
			return true, nil
		}
	}
	return false, nil
}

// stepHasEligibleCandidate reports whether ANY snapshotted candidate of the step is
// still currently eligible to act. It reads the FROZEN snapshot only: it never
// substitutes, reassigns or re-resolves the level.
func (r *eligibilityResolver) stepHasEligibleCandidate(step models.ApprovalStep, assignments []models.ApprovalAssignment) (bool, error) {
	if len(assignments) == 0 {
		return false, nil
	}

	items, err := r.stepAnchorItems(step)
	if err != nil {
		return false, err
	}
	if len(items) == 0 {
		return false, nil
	}

	for _, item := range items {
		users, err := r.anchorUsers(item)
		if err != nil {
			return false, err
		}
		for _, assignment := range assignments {
			if !users[assignment.UserNumber] {
				continue
			}
			ok, err := r.authenticatableUser(assignment.UserNumber)
			if err != nil {
				return false, err
			}
			if ok {
				return true, nil
			}
		}
	}
	return false, nil
}

// isCurrentlyEligible re-derives the actor's eligibility at action time (BD-E8).
// A user whose role was removed or expired, whose stakeholder relationship was
// ended, whose placement lapsed, or whose login was disabled is no longer an
// approver — without any silent reassignment or path rebuild.
func isCurrentlyEligible(organisationNumber string, step models.ApprovalStep, userNumber string) (bool, error) {
	return newEligibilityResolver(organisationNumber).userEligibleForStep(step, userNumber)
}

// requestIsStalled derives, at READ TIME, whether a pending request's CURRENT level
// can still be advanced by its snapshotted candidates.
//
// Foundation E3 (ratified D2/D3): "stalled" means the request is PENDING, its
// current level still awaits a decision, and NONE of that level's snapshotted
// candidates is currently eligible to act. It is derived from current data on every
// read and is never stored.
//
// It is informational ONLY. It never substitutes a candidate, never re-resolves the
// path, never cancels, and never approves; and it deliberately considers ONLY the
// current level, so a future level cannot make an actionable request read as
// stalled. The action-time guard remains the sole authority for acting.
func requestIsStalled(resolver *eligibilityResolver, approval models.ApprovalRequest, steps []models.ApprovalStep, assignmentsByStep map[string][]models.ApprovalAssignment) (bool, error) {
	if approval.Status != models.RequestStatusPending {
		return false, nil
	}

	var current *models.ApprovalStep
	for index := range steps {
		if steps[index].ApprovalLevel == approval.CurrentLevel {
			current = &steps[index]
			break
		}
	}
	if current == nil {
		return false, nil
	}

	hasEligible, err := resolver.stepHasEligibleCandidate(*current, assignmentsByStep[current.Number])
	if err != nil {
		return false, err
	}
	return !hasEligible, nil
}

// currentStepCandidateUsers collects the distinct user numbers assigned to the
// CURRENT step of each PENDING request, so login eligibility can be resolved in ONE
// batched query for the whole read.
func currentStepCandidateUsers(approvals []models.ApprovalRequest, stepsByRequest map[string][]models.ApprovalStep, assignmentsByStep map[string][]models.ApprovalAssignment) []string {
	set := map[string]bool{}
	for _, approval := range approvals {
		if approval.Status != models.RequestStatusPending {
			continue
		}
		for _, step := range stepsByRequest[approval.Number] {
			if step.ApprovalLevel != approval.CurrentLevel {
				continue
			}
			for _, assignment := range assignmentsByStep[step.Number] {
				set[assignment.UserNumber] = true
			}
		}
	}

	result := make([]string, 0, len(set))
	for userNumber := range set {
		result = append(result, userNumber)
	}
	return result
}

// commentWithinLimit reports whether a comment fits the shared limit, counted in
// Unicode characters (runes) rather than Go's byte-counting len() (ratified E2
// decision D3). The database can therefore never receive a value that exceeds
// the varchar(1000) column: at most maxCommentLen runes is at most maxCommentLen
// characters.
func commentWithinLimit(comment string) bool {
	return utf8.RuneCountInString(comment) <= maxCommentLen
}

// normalizePagination applies the request-list bounds: a default page size, a
// hard maximum, and a non-negative offset.
func normalizePagination(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = defaultListLimit
	}
	if limit > maxListLimit {
		limit = maxListLimit
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}

// CreateApprovalRequest creates a request and its snapshotted path.
//
// It calls Foundation D's resolver exactly once, then stores the resulting path,
// anchors and resolved candidates. It never re-runs D for an existing request
// (BD-E14) and never fabricates authority when D declines to supply a path.
// effectiveDeadline returns the advisory deadline used for a request: the
// persisted deadline when present (an explicit override, or the default resolved
// and stored at creation), otherwise the built-in 7-day default applied to
// creation time for rows created before the E4-2 persistence invariant. Always
// returned in UTC (ratified E4-D5). A later configuration change never rewrites a
// persisted deadline (historical safety).
func effectiveDeadline(approval models.ApprovalRequest) time.Time {
	if approval.DeadlineAt != nil {
		return approval.DeadlineAt.UTC()
	}
	return approval.CreatedAt.UTC().Add(time.Duration(builtInDefaultDeadlineHours) * time.Hour)
}

// requestIsOverdue is the Foundation E4 (ratified E4-D1/D2/D6) read-time derived
// condition. It is true only for a PENDING request whose effective deadline has
// passed. It is never stored, never changes the request status, and never grants
// or removes the right to act; the action-time eligibility guard remains the
// sole authority.
func requestIsOverdue(approval models.ApprovalRequest, now time.Time) bool {
	if approval.Status != models.RequestStatusPending {
		return false
	}
	return now.UTC().After(effectiveDeadline(approval))
}

// requestIsApproachingDeadline is the Foundation E4 (ratified C1 Option B + C4)
// "approaching deadline" (reminder) condition. It is DERIVED AT READ TIME, is
// never stored, and creates NO notification:
//
//	status == PENDING && nowUTC <= DeadlineAt && DeadlineAt - nowUTC <= 24h
//
// It is mutually exclusive with overdue (overdue requires nowUTC > DeadlineAt),
// so an already-overdue request is never "approaching". Informational only — it
// never changes status, authority or eligibility.
func requestIsApproachingDeadline(approval models.ApprovalRequest, now time.Time) bool {
	if approval.Status != models.RequestStatusPending {
		return false
	}
	now = now.UTC()
	deadline := effectiveDeadline(approval)
	if now.After(deadline) {
		return false // Overdue; the reminder condition is not active once overdue.
	}
	return deadline.Sub(now) <= reminderWindow
}

func CreateApprovalRequest(actorNumber, organisationNumber string, request requests.CreateApprovalRequest) (*responses.ApprovalRequestResponse, error) {
	module := strings.TrimSpace(request.Module)
	resourceType := strings.TrimSpace(request.ResourceType)
	resourceNumber := strings.TrimSpace(request.ResourceNumber)
	decisionType := strings.TrimSpace(request.DecisionType)
	unitNumber := strings.TrimSpace(request.OrganisationUnitNumber)

	if err := validateRequestInput(module, resourceType, resourceNumber, decisionType, request.Amount); err != nil {
		return nil, err
	}

	// Foundation E2 resubmission linkage (ratified D5-D8). The linkage is
	// validated server-side BEFORE anything is written or resolved, and a single
	// generic error is returned for every failure so this endpoint can never be
	// used as an existence oracle for another organisation's or another user's
	// request numbers.
	previousRequestNumber := strings.TrimSpace(request.PreviousRequestNumber)
	if previousRequestNumber != "" {
		if err := validateResubmissionLink(previousRequestNumber, organisationNumber, actorNumber); err != nil {
			return nil, err
		}
	}

	// Foundation E4 (ratified E4-D4/D5): resolve the effective advisory deadline
	// at creation and PERSIST it. Persisting the resolved value means a later
	// change to the configured default can never rewrite this request's deadline
	// (historical safety). An explicit override takes precedence; otherwise the
	// organisation's configured default (built-in 7 days when unset) applies.
	now := time.Now().UTC()
	deadline := now.Add(time.Duration(deadlineDefaultHours(organisationNumber)) * time.Hour)
	if request.DeadlineAt != nil {
		deadline = request.DeadlineAt.UTC()
		if !deadline.After(now) {
			return nil, ErrInvalidDeadline
		}
	}
	request.DeadlineAt = &deadline

	amount := 0.0
	if request.Amount != nil {
		amount = *request.Amount
	}

	resolution, err := authorityServices.ResolveAuthority(organisationNumber, decisionType, unitNumber, amount)
	if err != nil {
		return nil, mapAuthorityError(err)
	}
	// Cases 1 and 3: nothing governed, or nothing applicable to this scope.
	// E creates no request and never reinterprets D's answer.
	if !resolution.Configured {
		return nil, ErrDecisionNotConfigured
	}
	if len(resolution.RequiredPath) == 0 {
		return nil, ErrNoApplicableAuthority
	}
	// Case 4: the amount is above every declared limit for this scope. E does
	// not pretend authority exists and does not escalate on its own.
	if resolution.ExceedsDeclaredAuthority {
		return nil, ErrExceedsDeclaredAuthority
	}

	type plannedStep struct {
		level      int
		anchors    []authorityResponses.AuthorityAnchorResponse
		candidates []string
	}
	plans := make([]plannedStep, 0, len(resolution.RequiredPath))
	blocked := false

	for _, level := range resolution.RequiredPath {
		candidates, err := resolveCandidates(organisationNumber, toAnchors(level.Anchors))
		if err != nil {
			return nil, err
		}
		// BD-E6: the requester never approves their own request.
		candidates = removeUser(candidates, actorNumber)
		if len(candidates) == 0 {
			blocked = true
		}

		plans = append(plans, plannedStep{
			level:      level.ApprovalLevel,
			anchors:    level.Anchors,
			candidates: candidates,
		})
	}

	if len(plans) == 0 {
		return nil, ErrNoApplicableAuthority
	}

	requestNumber := uuid.NewV4().String()
	approval := models.ApprovalRequest{
		Number:                 requestNumber,
		OrganisationNumber:     organisationNumber,
		OrganisationUnitNumber: unitNumber,
		Module:                 module,
		ResourceType:           resourceType,
		ResourceNumber:         resourceNumber,
		DecisionType:           decisionType,
		Amount:                 request.Amount,
		CurrentLevel:           plans[0].level,
		RequestedByNumber:      actorNumber,
		AddedByNumber:          actorNumber,
		PreviousRequestNumber:  previousRequestNumber,
		DeadlineAt:             request.DeadlineAt,
	}
	if blocked {
		approval.Status = models.RequestStatusBlocked
	} else {
		approval.Status = models.RequestStatusPending
	}

	// The request, its steps and its assignments are created in one transaction so
	// a partially consistent path can never be observed.
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		if err := approvalRepositories.CreateRequest(&approval, tx); err != nil {
			return err
		}

		for _, plan := range plans {
			step := models.ApprovalStep{
				Number:             uuid.NewV4().String(),
				RequestNumber:      requestNumber,
				OrganisationNumber: organisationNumber,
				ApprovalLevel:      plan.level,
				Anchors:            serializeAnchors(plan.anchors),
				Status:             models.StepStatusPending,
			}
			// BD-E7: a level with no eligible approver is recorded as BLOCKED.
			if len(plan.candidates) == 0 {
				step.Status = models.StepStatusBlocked
			}
			if err := approvalRepositories.CreateStep(&step, tx); err != nil {
				return err
			}

			// Foundation D (ratified D6): snapshot the unit-aware anchors as
			// normalised rows, keeping same-anchor/different-unit authorities
			// distinct. The legacy serialised column is still written (unit-less) so
			// it stays populated and pre-amendment readers keep working; the
			// normalised rows are the authoritative anchor record for this step.
			for _, item := range plan.anchors {
				if strings.TrimSpace(item.Type) == "" || strings.TrimSpace(item.Value) == "" {
					continue
				}
				stepAnchor := models.ApprovalStepAnchor{
					Number:                 uuid.NewV4().String(),
					StepNumber:             step.Number,
					OrganisationNumber:     organisationNumber,
					AnchorType:             item.Type,
					AnchorValue:            item.Value,
					OrganisationUnitNumber: item.OrganisationUnitNumber,
				}
				if err := approvalRepositories.CreateStepAnchor(&stepAnchor, tx); err != nil {
					return err
				}
			}

			for _, candidate := range plan.candidates {
				assignment := models.ApprovalAssignment{
					Number:             uuid.NewV4().String(),
					StepNumber:         step.Number,
					RequestNumber:      requestNumber,
					OrganisationNumber: organisationNumber,
					UserNumber:         candidate,
					Status:             models.AssignmentStatusUndecided,
				}
				if err := approvalRepositories.CreateAssignment(&assignment, tx); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return GetApprovalRequest(organisationNumber, requestNumber)
}

// validateResubmissionLink enforces the ratified E2 linkage rules (D5-D8). The
// referenced request must exist WITHIN the acting organisation, must be owned by
// the acting user (only the original requester may resubmit), and must be in
// RETURNED — the only linkable state, because BD-E9 keeps rejection terminal.
//
// Every failure returns the same error on purpose: distinguishing "not found"
// from "belongs to someone else" would leak the existence of requests the caller
// may not see.
func validateResubmissionLink(previousRequestNumber, organisationNumber, actorNumber string) error {
	previous, err := approvalRepositories.FindRequestByNumberAndOrganisation(previousRequestNumber, organisationNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPreviousRequestInvalid
		}
		return err
	}

	if previous.RequestedByNumber != actorNumber {
		return ErrPreviousRequestInvalid
	}
	if previous.Status != models.RequestStatusReturned {
		return ErrPreviousRequestInvalid
	}
	return nil
}

// ListApprovalRequests lists the organisation's requests, optionally filtered by
// status, module, resource type and resource number. The page is always bounded.
func ListApprovalRequests(organisationNumber, status, module, resourceType, resourceNumber string, limit, offset int) ([]responses.ApprovalRequestResponse, error) {
	if status != "" && !isValidRequestStatus(status) {
		return nil, ErrInvalidRequestStatus
	}

	limit, offset = normalizePagination(limit, offset)

	approvals, err := approvalRepositories.FindAllByOrganisation(organisationNumber, status,
		strings.TrimSpace(module), strings.TrimSpace(resourceType), strings.TrimSpace(resourceNumber), limit, offset)
	if err != nil {
		return nil, err
	}

	numbers := make([]string, 0, len(approvals))
	for _, approval := range approvals {
		numbers = append(numbers, approval.Number)
	}

	steps, err := approvalRepositories.FindStepsByRequestNumbers(numbers)
	if err != nil {
		return nil, err
	}
	stepNumbers := make([]string, 0, len(steps))
	for _, step := range steps {
		stepNumbers = append(stepNumbers, step.Number)
	}
	assignments, err := approvalRepositories.FindAssignmentsByStepNumbers(stepNumbers)
	if err != nil {
		return nil, err
	}

	stepsByRequest := map[string][]models.ApprovalStep{}
	for _, step := range steps {
		stepsByRequest[step.RequestNumber] = append(stepsByRequest[step.RequestNumber], step)
	}
	assignmentsByStep := map[string][]models.ApprovalAssignment{}
	for _, assignment := range assignments {
		assignmentsByStep[assignment.StepNumber] = append(assignmentsByStep[assignment.StepNumber], assignment)
	}

	anchorsByStep, err := loadStepAnchors(steps)
	if err != nil {
		return nil, err
	}

	// Foundation E3 (ratified D2/D3): derive the stalled condition at READ TIME.
	// One shared resolver memoises anchor lookups and one batched query resolves
	// login eligibility, so the page has no N+1 pattern (E3 §21).
	resolver := newEligibilityResolver(organisationNumber).withStepAnchors(anchorsByStep)
	authenticatable, err := stakeholderRepositories.FindAuthenticatableUserNumbers(
		currentStepCandidateUsers(approvals, stepsByRequest, assignmentsByStep))
	if err != nil {
		return nil, err
	}
	resolver.withAuthenticatable(authenticatable)

	// Foundation E4 (C1 Option B + C4): one UTC evaluation instant for the whole
	// page. The derived reminder/overdue conditions are computed in memory from
	// the already-loaded requests — no per-row query and no notification write.
	now := time.Now().UTC()

	result := make([]responses.ApprovalRequestResponse, 0, len(approvals))
	for _, approval := range approvals {
		stalled, err := requestIsStalled(resolver, approval, stepsByRequest[approval.Number], assignmentsByStep)
		if err != nil {
			return nil, err
		}
		result = append(result, toResponse(approval, stepsByRequest[approval.Number], assignmentsByStep, anchorsByStep, stalled, now))
	}
	return result, nil
}

// GetApprovalRequest returns one request of the organisation with its full,
// snapshotted path, including the reason each assignee was chosen.
func GetApprovalRequest(organisationNumber, number string) (*responses.ApprovalRequestResponse, error) {
	approval, err := approvalRepositories.FindRequestByNumberAndOrganisation(number, organisationNumber)
	if err != nil {
		return nil, ErrRequestNotFound
	}

	steps, err := approvalRepositories.FindStepsByRequestNumbers([]string{approval.Number})
	if err != nil {
		return nil, err
	}
	stepNumbers := make([]string, 0, len(steps))
	for _, step := range steps {
		stepNumbers = append(stepNumbers, step.Number)
	}
	assignments, err := approvalRepositories.FindAssignmentsByStepNumbers(stepNumbers)
	if err != nil {
		return nil, err
	}

	assignmentsByStep := map[string][]models.ApprovalAssignment{}
	for _, assignment := range assignments {
		assignmentsByStep[assignment.StepNumber] = append(assignmentsByStep[assignment.StepNumber], assignment)
	}

	anchorsByStep, err := loadStepAnchors(steps)
	if err != nil {
		return nil, err
	}

	// Foundation E3 (ratified D2/D3): the stalled condition is derived here, at read
	// time, and never stored.
	resolver := newEligibilityResolver(organisationNumber).withStepAnchors(anchorsByStep)
	authenticatable, err := stakeholderRepositories.FindAuthenticatableUserNumbers(
		currentStepCandidateUsers([]models.ApprovalRequest{*approval}, map[string][]models.ApprovalStep{approval.Number: steps}, assignmentsByStep))
	if err != nil {
		return nil, err
	}
	resolver.withAuthenticatable(authenticatable)

	stalled, err := requestIsStalled(resolver, *approval, steps, assignmentsByStep)
	if err != nil {
		return nil, err
	}

	// Foundation E4 (C1 Option B + C4): the derived reminder/overdue conditions are
	// computed in memory from the loaded request using one UTC instant — no
	// notification write and no extra query.
	response := toResponse(*approval, steps, assignmentsByStep, anchorsByStep, stalled, time.Now().UTC())
	return &response, nil
}

// ApproveApprovalRequest records one valid approval for the request's current
// level.
//
// Under BD-E1 (within a level: ANY) the first accepted decision completes the
// level; the other candidates of that level can no longer act. Under BD-E2 the
// levels are strictly sequential, so the next level only becomes actionable when
// this one completes. The whole operation is one transaction and the request row
// is locked, so two simultaneous approvals can never both complete one level.
func ApproveApprovalRequest(actorNumber, organisationNumber, number, comment string) (*responses.ApprovalRequestResponse, error) {
	comment = strings.TrimSpace(comment)
	if !commentWithinLimit(comment) {
		return nil, ErrInvalidFieldLength
	}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		return recordDecision(tx, actorNumber, organisationNumber, number, models.AssignmentStatusApproved, comment)
	})
	if err != nil {
		return nil, err
	}

	return GetApprovalRequest(organisationNumber, number)
}

// RejectApprovalRequest rejects the request terminally (BD-E9). A reason is
// mandatory, the decision record is immutable, and the request can never
// progress afterwards.
func RejectApprovalRequest(actorNumber, organisationNumber, number, comment string) (*responses.ApprovalRequestResponse, error) {
	comment = strings.TrimSpace(comment)
	if comment == "" {
		return nil, ErrRejectionReasonRequired
	}
	if !commentWithinLimit(comment) {
		return nil, ErrInvalidFieldLength
	}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		return recordDecision(tx, actorNumber, organisationNumber, number, models.AssignmentStatusRejected, comment)
	})
	if err != nil {
		return nil, err
	}

	return GetApprovalRequest(organisationNumber, number)
}

// ReturnApprovalRequest implements Foundation E2 return-for-correction (ratified
// D1/D2/D3/D4).
//
// An eligible approver of the request's CURRENT step sends the request back to
// its requester. The request ends its active life in RETURNED; it can never be
// approved, rejected or cancelled afterwards, and its recorded decisions and
// snapshotted path stay immutable. The only way forward is resubmission as a new
// request (ReturnApprovalRequest does not create one).
//
// The reason is mandatory and limited to maxCommentLen Unicode characters. The
// write goes through the same guarded transaction as approve/reject, so the
// request-row lock, the current-step guard, the assignment guard and the
// action-time eligibility re-check are all reused unchanged.
func ReturnApprovalRequest(actorNumber, organisationNumber, number, comment string) (*responses.ApprovalRequestResponse, error) {
	comment = strings.TrimSpace(comment)
	if comment == "" {
		return nil, ErrReturnReasonRequired
	}
	if !commentWithinLimit(comment) {
		return nil, ErrInvalidFieldLength
	}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		return recordDecision(tx, actorNumber, organisationNumber, number, models.AssignmentStatusReturned, comment)
	})
	if err != nil {
		return nil, err
	}

	return GetApprovalRequest(organisationNumber, number)
}

// recordDecision is the shared approve/reject/return transaction. It locks the
// request, verifies it is actionable, verifies the actor holds an assignment for
// the current step, re-checks eligibility right now, and then writes the
// decision, completes the step, and advances, returns or completes the request.
func recordDecision(tx *gorm.DB, actorNumber, organisationNumber, number, decision, comment string) error {
	approval, err := approvalRepositories.LockRequestByNumberAndOrganisation(number, organisationNumber, tx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRequestNotFound
		}
		return err
	}
	if approval.Status != models.RequestStatusPending {
		return ErrRequestNotActionable
	}

	// The step the transaction locked decides which level this action targets.
	// After a concurrent completion the request either moves on or is finished,
	// so the second attempt cannot complete the same level twice.
	step, err := approvalRepositories.LockCurrentStep(approval.Number, approval.CurrentLevel, organisationNumber, tx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrStepNotActionable
		}
		return err
	}
	if step.Status != models.StepStatusPending {
		return ErrStepNotActionable
	}

	assignment, err := approvalRepositories.LockAssignment(step.Number, actorNumber, organisationNumber, tx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotEligible
		}
		return err
	}
	if assignment.Status != models.AssignmentStatusUndecided {
		return ErrAlreadyDecided
	}

	eligible, err := isCurrentlyEligible(organisationNumber, *step, actorNumber)
	if err != nil {
		return err
	}
	if !eligible {
		return ErrNotEligible
	}

	now := time.Now()
	assignment.Status = decision
	assignment.Comment = comment
	assignment.DecidedAt = &now
	if err := approvalRepositories.SaveAssignment(assignment, tx); err != nil {
		return err
	}

	// The step takes its status from the step vocabulary, not from the
	// assignment vocabulary: the constants are distinct even though the
	// approved/rejected strings currently coincide.
	stepStatus := models.StepStatusApproved
	switch decision {
	case models.AssignmentStatusRejected:
		stepStatus = models.StepStatusRejected
	case models.AssignmentStatusReturned:
		stepStatus = models.StepStatusReturned
	}
	step.Status = stepStatus
	if err := approvalRepositories.SaveStep(step, tx); err != nil {
		return err
	}

	if decision == models.AssignmentStatusRejected {
		// BD-E9: rejection is terminal. Later levels stay untouched and can never
		// become actionable because the request is no longer PENDING.
		approval.Status = models.RequestStatusRejected
		return approvalRepositories.SaveRequest(approval, tx)
	}

	if decision == models.AssignmentStatusReturned {
		// E2 D1/D2/D4: return goes to the requester and ends the request's active
		// life. Later levels stay PENDING but are never actionable because the
		// request is no longer PENDING. Prior decisions are untouched.
		approval.Status = models.RequestStatusReturned
		return approvalRepositories.SaveRequest(approval, tx)
	}

	next, err := approvalRepositories.FindNextPendingStep(approval.Number, step.ApprovalLevel, tx)
	if err != nil {
		return err
	}
	if next == nil {
		approval.Status = models.RequestStatusApproved
	} else {
		approval.CurrentLevel = next.ApprovalLevel
	}
	return approvalRepositories.SaveRequest(approval, tx)
}

// CancelApprovalRequest lets the requester withdraw their own request. It is
// available only while the request can still be acted on (PENDING or BLOCKED),
// and it never lets one user cancel another user's request, even with the
// management permission required by the route.
func CancelApprovalRequest(actorNumber, organisationNumber, number string) (*responses.ApprovalRequestResponse, error) {
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		approval, err := approvalRepositories.LockRequestByNumberAndOrganisation(number, organisationNumber, tx)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrRequestNotFound
			}
			return err
		}

		if approval.RequestedByNumber != actorNumber {
			return ErrNotRequestOwner
		}
		if approval.Status != models.RequestStatusPending && approval.Status != models.RequestStatusBlocked {
			return ErrRequestNotCancellable
		}

		approval.Status = models.RequestStatusCancelled
		if err := approvalRepositories.CancelPendingSteps(approval.Number, organisationNumber, tx); err != nil {
			return err
		}
		return approvalRepositories.SaveRequest(approval, tx)
	})
	if err != nil {
		return nil, err
	}

	return GetApprovalRequest(organisationNumber, number)
}

// ListApprovalTasks returns the caller's actionable inbox: undecided assignments
// on the current step of pending requests. Only the caller's own tasks are
// returned, so no one can discover who else was asked to approve.
func ListApprovalTasks(organisationNumber, userNumber string) ([]responses.ApprovalTaskResponse, error) {
	assignments, err := approvalRepositories.FindActionableAssignmentsForUser(organisationNumber, userNumber)
	if err != nil {
		return nil, err
	}
	if len(assignments) == 0 {
		return []responses.ApprovalTaskResponse{}, nil
	}

	requestNumbers := make([]string, 0, len(assignments))
	stepNumbers := make([]string, 0, len(assignments))
	for _, assignment := range assignments {
		requestNumbers = append(requestNumbers, assignment.RequestNumber)
		stepNumbers = append(stepNumbers, assignment.StepNumber)
	}

	approvals, err := approvalRepositories.FindRequestsByNumbers(requestNumbers)
	if err != nil {
		return nil, err
	}
	steps, err := approvalRepositories.FindStepsByNumbers(stepNumbers)
	if err != nil {
		return nil, err
	}

	requestsByNumber := map[string]models.ApprovalRequest{}
	for _, approval := range approvals {
		requestsByNumber[approval.Number] = approval
	}
	stepsByNumber := map[string]models.ApprovalStep{}
	for _, step := range steps {
		stepsByNumber[step.Number] = step
	}

	// Foundation E3 (ratified D3): the inbox is a READ-TIME view that returns only
	// assignments whose actor is CURRENTLY eligible. It narrows what is shown; it
	// never re-resolves or rewrites the snapshot, never substitutes another person,
	// and the action-time guard remains the sole authority for acting. Anchor
	// lookups are memoised per distinct anchor, so this adds no query per row
	// (E3 §21).
	anchorsByStep, err := loadStepAnchors(steps)
	if err != nil {
		return nil, err
	}
	authenticatable, err := stakeholderRepositories.FindAuthenticatableUserNumbers([]string{userNumber})
	if err != nil {
		return nil, err
	}
	resolver := newEligibilityResolver(organisationNumber).
		withStepAnchors(anchorsByStep).
		withAuthenticatable(authenticatable)

	// Foundation E4 (C1 Option B + C4): one UTC instant for the whole inbox; the
	// derived reminder/overdue conditions are computed in memory from the
	// already-loaded requests (no per-row query, no notification write).
	now := time.Now().UTC()

	result := make([]responses.ApprovalTaskResponse, 0, len(assignments))
	for _, assignment := range assignments {
		approval, ok := requestsByNumber[assignment.RequestNumber]
		if !ok {
			continue
		}
		step, ok := stepsByNumber[assignment.StepNumber]
		if !ok {
			continue
		}

		eligible, err := resolver.userEligibleForStep(step, userNumber)
		if err != nil {
			return nil, err
		}
		if !eligible {
			continue
		}

		deadline := effectiveDeadline(approval)
		result = append(result, responses.ApprovalTaskResponse{
			RequestNumber:          approval.Number,
			StepNumber:             step.Number,
			AssignmentNumber:       assignment.Number,
			ApprovalLevel:          step.ApprovalLevel,
			Module:                 approval.Module,
			ResourceType:           approval.ResourceType,
			ResourceNumber:         approval.ResourceNumber,
			DecisionType:           approval.DecisionType,
			Amount:                 approval.Amount,
			OrganisationUnitNumber: approval.OrganisationUnitNumber,
			RequestStatus:          approval.Status,
			DeadlineAt:             &deadline,
			Reminder:               requestIsApproachingDeadline(approval, now),
			Overdue:                requestIsOverdue(approval, now),
		})
	}
	return result, nil
}

// CurrentOrganisationNumber is re-exported so callers depend on one resolver.
var CurrentOrganisationNumber = organisationServices.CurrentOrganisationNumber

func toResponse(approval models.ApprovalRequest, steps []models.ApprovalStep, assignmentsByStep map[string][]models.ApprovalAssignment, anchorsByStep map[string][]anchor, stalled bool, now time.Time) responses.ApprovalRequestResponse {
	// Foundation E4 (ratified E4-D1/D5/D6, C1 Option B, C4): the effective advisory
	// deadline and the derived reminder/overdue conditions are computed on every
	// read, in memory from the supplied UTC instant, and are never stored.
	now = now.UTC()
	deadline := effectiveDeadline(approval)
	overdue := requestIsOverdue(approval, now)
	reminder := requestIsApproachingDeadline(approval, now)

	stepResponses := make([]responses.ApprovalStepResponse, 0, len(steps))

	for _, step := range steps {
		// A step is actionable only when the request is pending, the step itself
		// is pending, and the step is the request's current level. BD-E2 keeps
		// later levels visible but not actionable.
		actionable := approval.Status == models.RequestStatusPending &&
			step.Status == models.StepStatusPending &&
			step.ApprovalLevel == approval.CurrentLevel

		// Foundation D (ratified D5/D6/D9): unit-aware anchors come from the
		// normalised rows for steps created after the amendment, and from the
		// unit-less legacy column (organisation-wide) for earlier steps.
		anchorItems := anchorsByStep[step.Number]
		anchorResponses := make([]responses.ApprovalAnchorResponse, 0, len(anchorItems))
		for _, item := range anchorItems {
			anchorResponses = append(anchorResponses, responses.ApprovalAnchorResponse{
				Type:                   item.Type,
				Value:                  item.Value,
				OrganisationUnitNumber: item.Unit,
			})
		}

		stepAssignments := assignmentsByStep[step.Number]
		assignmentResponses := make([]responses.ApprovalAssignmentResponse, 0, len(stepAssignments))
		for _, assignment := range stepAssignments {
			assignmentResponses = append(assignmentResponses, responses.ApprovalAssignmentResponse{
				Number:     assignment.Number,
				UserNumber: assignment.UserNumber,
				Status:     assignment.Status,
				Comment:    assignment.Comment,
				DecidedAt:  assignment.DecidedAt,
			})
		}

		stepResponses = append(stepResponses, responses.ApprovalStepResponse{
			Number:        step.Number,
			ApprovalLevel: step.ApprovalLevel,
			Status:        step.Status,
			Actionable:    actionable,
			Anchors:       anchorResponses,
			Assignments:   assignmentResponses,
		})
	}

	return responses.ApprovalRequestResponse{
		Number:                 approval.Number,
		OrganisationNumber:     approval.OrganisationNumber,
		OrganisationUnitNumber: approval.OrganisationUnitNumber,
		Module:                 approval.Module,
		ResourceType:           approval.ResourceType,
		ResourceNumber:         approval.ResourceNumber,
		DecisionType:           approval.DecisionType,
		Amount:                 approval.Amount,
		Status:                 approval.Status,
		CurrentLevel:           approval.CurrentLevel,
		Stalled:                stalled,
		DeadlineAt:             &deadline,
		Overdue:                overdue,
		Reminder:               reminder,
		RequestedByNumber:      approval.RequestedByNumber,
		PreviousRequestNumber:  approval.PreviousRequestNumber,
		Steps:                  stepResponses,
		CreatedAt:              approval.CreatedAt,
		UpdatedAt:              approval.UpdatedAt,
	}
}
