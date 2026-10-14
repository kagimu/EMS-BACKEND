package models

import (
	"time"

	"gorm.io/gorm"
)

// EMS Commercial Onboarding (test-payment workflow). All tables are additive;
// no existing licence/auth table is altered.
//
// NOTE: every field carries an explicit snake_case json tag. Without them Go
// serialises the exported Go names (PascalCase) and the frontend — which reads
// snake_case — would see undefined values.

// InstallationPackage is a configurable catalogue entry. No commercial prices
// are seeded as production defaults — authorised administrators create them.
//
// Fixed-price packages carry a base price and included allowances. A "custom"
// package (Custom=true) carries no fabricated price: it must be assessed and
// quoted manually, so the quotation flow refuses to auto-price it.
type InstallationPackage struct {
	gorm.Model
	Number            string `json:"number" gorm:"uniqueIndex;size:36"`
	Code              string `json:"code" gorm:"index;size:64"`
	Name              string `json:"name" gorm:"size:191"`
	Description       string `json:"description" gorm:"type:text"`
	BasePriceUGX      int64  `json:"base_price_ugx"`
	IncludedUsers     int    `json:"included_users"`
	IncludedDevices   int    `json:"included_devices"`
	IncludedLocations int    `json:"included_locations"`
	UserPriceUGX      int64  `json:"user_price_ugx"`
	DevicePriceUGX    int64  `json:"device_price_ugx"`
	LocationPriceUGX  int64  `json:"location_price_ugx"`
	TransportBaseUGX  int64  `json:"transport_base_ugx"`
	TransportPerKmUGX int64  `json:"transport_per_km_ugx"`
	// TransportStandardUGX is the standard transport fee for this package's
	// setup size (small setups 20,000; big setups 50,000). It is charged when
	// the installation location is outside Kampala city centre.
	TransportStandardUGX  int64  `json:"transport_standard_ugx"`
	TransportBandsJSON    string `json:"transport_bands_json" gorm:"type:text"`
	TrainingHours         int    `json:"training_hours"`
	Methods               string `json:"methods" gorm:"size:191"`
	ServicesJSON          string `json:"services_json" gorm:"type:text"`
	EstimatedDays         int    `json:"estimated_days"`
	DeploymentDescription string `json:"deployment_description" gorm:"type:text"`
	Custom                bool   `json:"custom" gorm:"default:false"`
	Active                bool   `json:"active" gorm:"default:true"`
	CreatedByNumber       string `json:"created_by_number" gorm:"size:36"`
}

// SetupOrder is one customer's installation quotation/order. The accepted
// quotation is snapshotted so later catalogue changes never alter it.
type SetupOrder struct {
	gorm.Model
	Number                 string `json:"number" gorm:"uniqueIndex;size:36"`
	PackageNumber          string `json:"package_number" gorm:"index;size:36"`
	Status                 string `json:"status" gorm:"index;size:32"`
	Currency               string `json:"currency" gorm:"size:8;default:UGX"`
	ContactName            string `json:"contact_name" gorm:"size:191"`
	ContactEmail           string `json:"contact_email" gorm:"size:191"`
	ContactPhone           string `json:"contact_phone" gorm:"size:64"`
	ConfigurationJSON      string `json:"configuration_json" gorm:"type:text"`
	QuoteSnapshotJSON      string `json:"quote_snapshot_json" gorm:"type:text"`
	SetupTotalUGX          int64  `json:"setup_total_ugx"`
	SimulatedPaymentNumber string `json:"simulated_payment_number" gorm:"size:36"`
	OrganisationNumber     string `json:"organisation_number" gorm:"index;size:36"`
}

// SubscriptionOrder is a software-subscription purchase (period-based). An
// onboarding subscription is created before sign-in and is bound to the paid
// setup order via SetupOrderNumber (OrganisationNumber stays empty until the
// account/organisation exists).
type SubscriptionOrder struct {
	gorm.Model
	Number                 string `json:"number" gorm:"uniqueIndex;size:36"`
	Status                 string `json:"status" gorm:"index;size:32"`
	Currency               string `json:"currency" gorm:"size:8;default:UGX"`
	DurationDays           int    `json:"duration_days"`
	DailyRateUGX           int64  `json:"daily_rate_ugx"`
	TotalUGX               int64  `json:"total_ugx"`
	OrganisationNumber     string `json:"organisation_number" gorm:"index;size:36"`
	SetupOrderNumber       string `json:"setup_order_number" gorm:"index;size:36"`
	SimulatedPaymentNumber string `json:"simulated_payment_number" gorm:"size:36"`
}

// ActivationToken is a one-shot activation token for either a setup order or a
// subscription order. Only the SHA-256 hash is stored; the raw token is never
// persisted or logged — and the hash is never serialised.
type ActivationToken struct {
	gorm.Model
	Number                  string     `json:"number" gorm:"uniqueIndex;size:36"`
	Purpose                 string     `json:"purpose" gorm:"index;size:16"`
	TokenHash               string     `json:"-" gorm:"uniqueIndex;size:64"`
	SetupOrderNumber        string     `json:"setup_order_number" gorm:"index;size:36"`
	SubscriptionOrderNumber string     `json:"subscription_order_number" gorm:"index;size:36"`
	Status                  string     `json:"status" gorm:"index;size:16"`
	IssuedAt                *time.Time `json:"issued_at"`
	ValidatedAt             *time.Time `json:"validated_at"`
	RedeemedAt              *time.Time `json:"redeemed_at"`
	RedeemedOrganisation    string     `json:"redeemed_organisation" gorm:"size:36"`
}

// SimulatedPayment records a test-provider transaction. Never real money.
type SimulatedPayment struct {
	gorm.Model
	Number    string `json:"number" gorm:"uniqueIndex;size:36"`
	Reference string `json:"reference" gorm:"index;size:36"`
	Kind      string `json:"kind" gorm:"size:16"`
	AmountUGX int64  `json:"amount_ugx"`
	Scenario  string `json:"scenario" gorm:"size:16"`
	Status    string `json:"status" gorm:"size:16"`
	Provider  string `json:"provider" gorm:"size:32;default:simulated"`
}

// TableName methods pin the commercial tables to the ems_ namespace.
func (InstallationPackage) TableName() string { return "ems_installation_packages" }
func (SetupOrder) TableName() string          { return "ems_setup_orders" }
func (SubscriptionOrder) TableName() string   { return "ems_subscription_orders" }
func (ActivationToken) TableName() string     { return "ems_activation_tokens" }
func (SimulatedPayment) TableName() string    { return "ems_simulated_payments" }
