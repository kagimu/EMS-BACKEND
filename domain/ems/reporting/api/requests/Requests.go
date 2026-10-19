package requests

import "time"

// CreateReportingLineRequest creates a reporting line between two existing
// assignments of the current organisation. The organisation is never supplied by
// the client; it is derived from the authenticated active-organisation context.
type CreateReportingLineRequest struct {
	SubordinateAssignmentNumber string     `json:"subordinate_assignment_number"`
	ManagerAssignmentNumber     string     `json:"manager_assignment_number"`
	ValidFrom                   *time.Time `json:"valid_from"`
	ValidTo                     *time.Time `json:"valid_to"`
}

// UpdateReportingLineRequest updates the manager and lifecycle fields. The
// subordinate is immutable: moving a reporting line to a different subordinate
// means ending this record and creating a new one, preserving history.
type UpdateReportingLineRequest struct {
	ManagerAssignmentNumber string     `json:"manager_assignment_number"`
	ValidFrom               *time.Time `json:"valid_from"`
	ValidTo                 *time.Time `json:"valid_to"`
}

// UpdateReportingLineStatusRequest is the status payload.
type UpdateReportingLineStatusRequest struct {
	Status string `json:"status"`
}
