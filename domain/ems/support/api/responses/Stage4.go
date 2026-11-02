package responses

import "time"

type StaffSummary struct {
	UserNumber      string `json:"user_number"`
	Role            string `json:"role"`
	Status          string `json:"status"`
	OpenTickets     int64  `json:"open_tickets"`
	AssignedTickets int64  `json:"assigned_tickets"`
	HighPriority    int64  `json:"high_priority"`
	Overdue         int64  `json:"overdue"`
}

type WorkloadSummary struct {
	Staff               []StaffSummary `json:"staff"`
	Unassigned          int64          `json:"unassigned"`
	OverloadedThreshold int64          `json:"overloaded_threshold"`
}

type SLATargetResponse struct {
	Priority             string `json:"priority"`
	FirstResponseMinutes int    `json:"first_response_minutes"`
	ResolutionMinutes    int    `json:"resolution_minutes"`
}

type AttachmentResponse struct {
	Number       string    `json:"number"`
	TicketNumber string    `json:"ticket_number"`
	OriginalName string    `json:"original_name"`
	ContentType  string    `json:"content_type"`
	SizeBytes    int64     `json:"size_bytes"`
	Visibility   string    `json:"visibility"`
	CreatedAt    time.Time `json:"created_at"`
}

type SupportReport struct {
	TicketsReceived             int64            `json:"tickets_received"`
	TicketsResolved             int64            `json:"tickets_resolved"`
	OpenTickets                 int64            `json:"open_tickets"`
	OverdueTickets              int64            `json:"overdue_tickets"`
	AverageFirstResponseMinutes float64          `json:"average_first_response_minutes"`
	AverageResolutionMinutes    float64          `json:"average_resolution_minutes"`
	ByPriority                  map[string]int64 `json:"by_priority"`
	ByCategory                  map[string]int64 `json:"by_category"`
	ByAgent                     map[string]int64 `json:"by_agent"`
}
