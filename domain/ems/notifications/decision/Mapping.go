package decision

// Notification event types and the ratified G-Q1 recipient rules.
//
// Foundation G Increment 4 implements the Owner-ratified G-Q1 event-to-recipient
// map. This file is the single, fixed (non-configurable) map. It contains no
// per-organisation configuration and no preference handling (G-Q1/G-Q3).
//
// Event type strings deliberately reuse the authoritative vocabulary:
//   - approval event names follow the Foundation E action names recorded by G2;
//   - licence event names are the exact Foundation F LicenceEventType values.
//
// Nothing here changes Foundation A–F semantics; a decision only decides
// whether a notification is created and for whom.

// Foundation E1/E2 approval events.
const (
	EventApprovalAssigned = "approval.create"
	EventApprovalAdvanced = "approval.advance"
	EventApprovalApproved = "approval.approve"
	EventApprovalRejected = "approval.reject"
	EventApprovalReturned = "approval.return"
)

// Foundation C events (one notification type per family; the affected user is
// the only recipient).
const (
	EventStakeholderChanged = "stakeholder.changed"
	EventAssignmentAssigned = "assignment.create"
	EventAssignmentChanged  = "assignment.update"
)

// Foundation F licensing events. These values are identical to the
// LicenceEventType constants so the licensing fact source can be passed through
// without translation, and no licensing fact is duplicated.
const (
	EventLicenceIssued           = "LICENCE_ISSUED"
	EventLicenceConfigured       = "LICENCE_CONFIGURED"
	EventLicenceRenewed          = "LICENCE_RENEWED"
	EventLicenceReplaced         = "LICENCE_REPLACED"
	EventLicenceSuspended        = "LICENCE_SUSPENDED"
	EventLicenceRevoked          = "LICENCE_REVOKED"
	EventInstallationReleased    = "INSTALLATION_RELEASED"
	EventClientRenewalRequested  = "CLIENT_RENEWAL_REQUESTED"
	EventInstallationActivated   = "INSTALLATION_ACTIVATED"
	EventInstallationReactivated = "INSTALLATION_REACTIVATED"
	EventInstallationDeactivated = "INSTALLATION_DEACTIVATED"
)

// recipientKind identifies an existing EMS recipient relationship. No new
// recipient concept is introduced.
type recipientKind int

const (
	// currentApprovers: the eligible (UNDECIDED) approvers of the request's
	// current level, from the authoritative approval assignment snapshot.
	currentApprovers recipientKind = iota
	// requester: the approval request's RequestedByNumber.
	requester
	// affectedUser: the subject of the C relationship (stakeholder/assignee).
	affectedUser
	// organisationOwner: the organisation's active owner stakeholder(s).
	organisationOwner
	// licenceRegistryAuthority: users holding the explicit
	// manage_licence_registry permission. Never a generic Admin.
	licenceRegistryAuthority
)

// recipientRule is one fixed G-Q1 mapping. actionRequired selects an
// action-required notification over an informational one (G-Q1/D-GQ1-13/16).
type recipientRule struct {
	kinds          []recipientKind
	actionRequired bool
	title          string
	body           string
}

// rules is the complete ratified G-Q1 map. An event absent from this map is
// audit-only / non-notifying (D-GQ1-17) and produces no notification.
var rules = map[string]recipientRule{
	// D-GQ1-1 approval events.
	EventApprovalAssigned: {[]recipientKind{currentApprovers}, true, "Approval required", "An approval request is awaiting your decision."},
	EventApprovalAdvanced: {[]recipientKind{currentApprovers}, true, "Approval required", "An approval request has advanced to your level and is awaiting your decision."},
	EventApprovalApproved: {[]recipientKind{requester}, false, "Approval completed", "Your approval request was fully approved."},
	EventApprovalRejected: {[]recipientKind{requester}, false, "Approval rejected", "Your approval request was rejected."},
	// D-GQ1-9: return is action-required for the requester (correct + resubmit).
	EventApprovalReturned: {[]recipientKind{requester}, true, "Approval returned", "Your approval request was returned for correction. Please revise and resubmit."},

	// D-GQ1-2/11: the affected stakeholder only, informational.
	EventStakeholderChanged: {[]recipientKind{affectedUser}, false, "Stakeholder record updated", "Your stakeholder record in the organisation was changed."},
	// D-GQ1-13: a new assignment is a genuine new responsibility, so it is
	// action-required; subsequent updates are informational (D-GQ1-3/12).
	EventAssignmentAssigned: {[]recipientKind{affectedUser}, true, "New assignment", "You have been given a new assignment. Please review it."},
	EventAssignmentChanged:  {[]recipientKind{affectedUser}, false, "Assignment updated", "Your assignment in the organisation was changed."},

	// D-GQ1-4/14: client-facing licence lifecycle, informational, to the
	// organisation's active owner.
	EventLicenceIssued:        {[]recipientKind{organisationOwner}, false, "Licence issued", "A licence was issued to your organisation."},
	EventLicenceConfigured:    {[]recipientKind{organisationOwner}, false, "Licence updated", "Your organisation's licence configuration was updated."},
	EventLicenceRenewed:       {[]recipientKind{organisationOwner}, false, "Licence renewed", "Your organisation's licence was renewed."},
	EventLicenceReplaced:      {[]recipientKind{organisationOwner}, false, "Licence replaced", "Your organisation's licence was replaced."},
	EventInstallationReleased: {[]recipientKind{organisationOwner}, false, "Installation released", "An installation for your organisation was released by the licensing team."},

	// D-GQ1-6/16: suspension/revocation is action-required for the client
	// owner, and the control-plane licensing recipient is also informed.
	EventLicenceSuspended: {[]recipientKind{organisationOwner, licenceRegistryAuthority}, true, "Licence suspended", "Your organisation's licence was suspended. Action is required."},
	EventLicenceRevoked:   {[]recipientKind{organisationOwner, licenceRegistryAuthority}, true, "Licence revoked", "Your organisation's licence was revoked. Action is required."},

	// D-GQ1-5: control-plane licensing recipient, never a generic Admin.
	EventClientRenewalRequested:  {[]recipientKind{licenceRegistryAuthority}, true, "Licence renewal requested", "A client has requested licence renewal. Action is required."},
	EventInstallationActivated:   {[]recipientKind{licenceRegistryAuthority}, false, "Installation activated", "A client activated an installation."},
	EventInstallationReactivated: {[]recipientKind{licenceRegistryAuthority}, false, "Installation reactivated", "A client reactivated an installation."},
	EventInstallationDeactivated: {[]recipientKind{licenceRegistryAuthority}, false, "Installation deactivated", "A client deactivated an installation."},
}
