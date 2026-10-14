package decision

import (
	"log"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	stakeholderModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/stakeholders/models"
	workflowModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/workflow/models"
	approvalRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/workflow/repositories"
)

// permissionManageLicenceRegistry mirrors the explicit control-plane authority
// defined by Foundation F (licensing/middleware/RegistryAuthority.go). It is
// reproduced as a string here so the notification decision layer does not couple
// to the licensing middleware; it is never a generic Admin bypass.
const permissionManageLicenceRegistry = "manage_licence_registry"

// resolveRecipients resolves the user numbers for a rule's recipient kinds,
// constrained to the event's organisation. It returns only user numbers; it
// never invents a recipient and never widens to another organisation.
func resolveRecipients(kinds []recipientKind, event Event) []string {
	resolved := make([]string, 0, len(kinds))
	for _, kind := range kinds {
		switch kind {
		case requester:
			if event.RequestedByNumber != "" {
				resolved = append(resolved, event.RequestedByNumber)
			}
		case affectedUser:
			if event.UserNumber != "" {
				resolved = append(resolved, event.UserNumber)
			}
		case currentApprovers:
			resolved = append(resolved, currentLevelApprovers(event.OrganisationNumber, event.RequestNumber, event.Level)...)
		case organisationOwner:
			resolved = append(resolved, organisationOwners(event.OrganisationNumber)...)
		case licenceRegistryAuthority:
			resolved = append(resolved, licenceRegistryAuthorities()...)
		}
	}
	return resolved
}

// currentLevelApprovers returns the eligible (UNDECIDED) approvers of the
// request's current level. Eligibility is read from the authoritative approval
// assignment snapshot; the frozen snapshot is never modified (E2/E3).
func currentLevelApprovers(organisationNumber, requestNumber string, level int) []string {
	if organisationNumber == "" || requestNumber == "" || level <= 0 {
		return nil
	}
	var numbers []string
	err := config.DB.Table("ems_approval_assignments AS a").
		Joins("JOIN ems_approval_steps AS s ON s.number = a.step_number AND s.deleted_at IS NULL").
		Where("a.deleted_at IS NULL").
		Where("a.organisation_number = ? AND a.request_number = ? AND a.status = ?",
			organisationNumber, requestNumber, workflowModel.AssignmentStatusUndecided).
		Where("s.approval_level = ?", level).
		Distinct().Pluck("a.user_number", &numbers).Error
	if err != nil {
		log.Printf("notification: resolving current approvers for %s level %d failed: %v", requestNumber, level, err)
		return nil
	}
	return numbers
}

// organisationOwners returns the organisation's active owner stakeholder(s),
// using the same status + validity-window semantics as Foundation C/E2
// (ratified D10). If none exist the result is empty and the caller must FAIL
// CLOSED: no fallback recipient is invented (D-GQ1-4).
func organisationOwners(organisationNumber string) []string {
	if organisationNumber == "" {
		return nil
	}
	numbers, err := approvalRepositories.FindUserNumbersByStakeholderType(
		organisationNumber, stakeholderModel.StakeholderTypeOwner, time.Now())
	if err != nil {
		log.Printf("notification: resolving owner stakeholders for %s failed: %v", organisationNumber, err)
		return nil
	}
	return numbers
}

// licenceRegistryAuthorities returns the users holding the explicit
// manage_licence_registry permission. This is the D-GQ1-5 control-plane
// recipient rule; it is resolved from the existing user-role-permission model
// and never from the Admin user type.
func licenceRegistryAuthorities() []string {
	var numbers []string
	err := config.DB.Table("user_roles AS ur").
		Joins("JOIN role_permissions AS rp ON rp.role_number = ur.role_number AND rp.deleted_at IS NULL").
		Where("ur.deleted_at IS NULL AND ur.expiry_date >= ?", time.Now()).
		Where("rp.permission = ?", permissionManageLicenceRegistry).
		Distinct().Pluck("ur.user_number", &numbers).Error
	if err != nil {
		log.Printf("notification: resolving licensing-management recipients failed: %v", err)
		return nil
	}
	return numbers
}
