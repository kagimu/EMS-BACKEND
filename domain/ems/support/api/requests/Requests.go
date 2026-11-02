package requests

type CreateTicketRequest struct {
	OrganisationNumber string `json:"organisation_number"`
	RequestedByNumber  string `json:"requested_by_number"`
	Subject            string `json:"subject"`
	Description        string `json:"description"`
	Category           string `json:"category"`
	Priority           string `json:"priority"`
	LicenceNumber      string `json:"licence_number"`
	InstallationNumber string `json:"installation_number"`
}

type AddMessageRequest struct {
	Body       string `json:"body"`
	Visibility string `json:"visibility"`
}

type UpdateTicketRequest struct {
	Status            *string `json:"status"`
	Priority          *string `json:"priority"`
	Category          *string `json:"category"`
	AssignedToNumber  *string `json:"assigned_to_number"`
	AssignmentReason  string  `json:"assignment_reason"`
	ResolutionSummary *string `json:"resolution_summary"`
}

type EscalateRequest struct {
	Reason string `json:"reason"`
}
