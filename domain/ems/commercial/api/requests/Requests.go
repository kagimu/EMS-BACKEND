package requests

import "time"

type CustomerProfileRequest struct {
	CustomerSince             *time.Time `json:"customer_since"`
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
	Reason                    string     `json:"reason"`
}
