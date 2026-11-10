package models

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

const (
	CommercialStatusProspect  = "PROSPECT"
	CommercialStatusActive    = "ACTIVE"
	CommercialStatusSuspended = "SUSPENDED"
	CommercialStatusChurned   = "CHURNED"
)

const (
	CustomerCategoryBusiness   = "BUSINESS"
	CustomerCategoryGovernment = "GOVERNMENT"
	CustomerCategoryEducation  = "EDUCATION"
	CustomerCategoryNonProfit  = "NON_PROFIT"
	CustomerCategoryOther      = "OTHER"
)

var CommercialStatuses = []string{
	CommercialStatusProspect,
	CommercialStatusActive,
	CommercialStatusSuspended,
	CommercialStatusChurned,
}

var CustomerCategories = []string{
	CustomerCategoryBusiness,
	CustomerCategoryGovernment,
	CustomerCategoryEducation,
	CustomerCategoryNonProfit,
	CustomerCategoryOther,
}

// CustomerProfile is the commercial view of an existing EMS organisation. It
// never duplicates the organisation identity and does not affect licensing.
type CustomerProfile struct {
	gorm.Model
	Number                    string     `json:"number" gorm:"size:36;unique;not null"`
	OrganisationNumber        string     `json:"organisation_number" gorm:"size:36;uniqueIndex;not null"`
	CustomerSince             *time.Time `json:"customer_since,omitempty"`
	CommercialStatus          string     `json:"commercial_status" gorm:"size:20;not null;index"`
	CustomerCategory          string     `json:"customer_category" gorm:"size:20;not null;index"`
	RegisteredBusinessName    string     `json:"registered_business_name" gorm:"size:191"`
	TaxIdentification         string     `json:"tax_identification" gorm:"size:80"`
	BillingContactName        string     `json:"billing_contact_name" gorm:"size:191"`
	BillingEmail              string     `json:"billing_email" gorm:"size:191"`
	BillingPhone              string     `json:"billing_phone" gorm:"size:40"`
	PhysicalAddress           string     `json:"physical_address" gorm:"type:text"`
	Country                   string     `json:"country" gorm:"size:80"`
	RelationshipManagerNumber string     `json:"relationship_manager_number" gorm:"size:36;index"`
	InternalNotes             string     `json:"internal_notes" gorm:"type:text"`
	CreatedByNumber           string     `json:"created_by_number" gorm:"size:36"`
	UpdatedByNumber           string     `json:"updated_by_number" gorm:"size:36"`
}

func (CustomerProfile) TableName() string { return "ems_customer_profiles" }

// CommercialEvent is append-only commercial audit history. It is separate from
// licence events and cannot be edited through normal profile endpoints.
type CommercialEvent struct {
	ID                 uint      `json:"id" gorm:"primaryKey"`
	Number             string    `json:"number" gorm:"size:36;unique;not null"`
	EntityType         string    `json:"entity_type" gorm:"size:40;not null;index"`
	EntityNumber       string    `json:"entity_number" gorm:"size:36;not null;index"`
	OrganisationNumber string    `json:"organisation_number" gorm:"size:36;index"`
	Action             string    `json:"action" gorm:"size:60;not null;index"`
	BeforeState        string    `json:"before_state" gorm:"type:text"`
	AfterState         string    `json:"after_state" gorm:"type:text"`
	Reason             string    `json:"reason" gorm:"type:text"`
	ActorUserNumber    string    `json:"actor_user_number" gorm:"size:36;not null"`
	CreatedAt          time.Time `json:"created_at"`
}

func (CommercialEvent) TableName() string { return "ems_commercial_events" }

func NormaliseCommercialStatus(value string) (string, error) {
	value = strings.ToUpper(strings.TrimSpace(value))
	for _, candidate := range CommercialStatuses {
		if value == candidate {
			return value, nil
		}
	}
	return "", errors.New("invalid commercial status")
}

func NormaliseCustomerCategory(value string) (string, error) {
	value = strings.ToUpper(strings.TrimSpace(value))
	for _, candidate := range CustomerCategories {
		if value == candidate {
			return value, nil
		}
	}
	return "", errors.New("invalid customer category")
}
