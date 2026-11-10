package responses

import "time"

type CustomerProfileResponse struct {
	Number                    string     `json:"number"`
	OrganisationNumber        string     `json:"organisation_number"`
	CustomerSince             *time.Time `json:"customer_since,omitempty"`
	CommercialStatus          string     `json:"commercial_status"`
	CustomerCategory          string     `json:"customer_category"`
	RegisteredBusinessName    string     `json:"registered_business_name"`
	TaxIdentification         string     `json:"tax_identification"`
	BillingContactName        string     `json:"billing_contact_name"`
	BillingEmail              string     `json:"billing_email"`
	BillingPhone              string     `json:"billing_phone"`
	PhysicalAddress           string     `json:"physical_address"`
	Country                   string     `json:"country"`
	RelationshipManagerNumber string     `json:"relationship_manager_number"`
	InternalNotes             string     `json:"internal_notes"`
	CreatedByNumber           string     `json:"created_by_number"`
	UpdatedByNumber           string     `json:"updated_by_number"`
	CreatedAt                 time.Time  `json:"created_at"`
	UpdatedAt                 time.Time  `json:"updated_at"`
}

type CommercialEventResponse struct {
	ID                 uint      `json:"id"`
	Number             string    `json:"number"`
	EntityType         string    `json:"entity_type"`
	EntityNumber       string    `json:"entity_number"`
	OrganisationNumber string    `json:"organisation_number"`
	Action             string    `json:"action"`
	BeforeState        string    `json:"before_state,omitempty"`
	AfterState         string    `json:"after_state,omitempty"`
	Reason             string    `json:"reason,omitempty"`
	ActorUserNumber    string    `json:"actor_user_number"`
	CreatedAt          time.Time `json:"created_at"`
}
