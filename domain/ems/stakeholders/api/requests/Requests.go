package requests

import "time"

// CreateStakeholderRequest is the create payload.
//
// The organisation is never supplied by the client; it comes from the
// authenticated active-organisation context. There is deliberately no field for
// creating a person/login — the referenced user must already exist.
type CreateStakeholderRequest struct {
	UserNumber      string     `json:"user_number"`
	StakeholderType string     `json:"stakeholder_type"`
	Title           string     `json:"title"`
	SharePercentage *float64   `json:"share_percentage"`
	Shares          *int64     `json:"shares"`
	ValidFrom       *time.Time `json:"valid_from"`
	ValidTo         *time.Time `json:"valid_to"`
}

// UpdateStakeholderRequest updates descriptive/lifecycle fields only.
//
// The person and the stakeholder type are immutable: changing the type means
// ending one relationship and creating another.
type UpdateStakeholderRequest struct {
	Title           string     `json:"title"`
	SharePercentage *float64   `json:"share_percentage"`
	Shares          *int64     `json:"shares"`
	ValidFrom       *time.Time `json:"valid_from"`
	ValidTo         *time.Time `json:"valid_to"`
}

// UpdateStakeholderStatusRequest is the status payload.
type UpdateStakeholderStatusRequest struct {
	Status string `json:"status"`
}

// CreateEmployeeRequest creates a new staff identity and its stakeholder record
// in one organisation-scoped operation. The server generates the user number.
type CreateEmployeeRequest struct {
	FirstName       string     `json:"first_name"`
	SecondName      string     `json:"second_name"`
	LastName        string     `json:"last_name"`
	Gender          string     `json:"gender"`
	Email           string     `json:"email"`
	Phone           string     `json:"phone"`
	StakeholderType string     `json:"stakeholder_type"`
	Title           string     `json:"title"`
	ValidFrom       *time.Time `json:"valid_from"`
	ValidTo         *time.Time `json:"valid_to"`
}
