package responses

import (
	"time"

	commercialResponses "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/commercial/api/responses"
	licensingResponses "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/api/responses"
)

// OwnerAccessResponse confirms the authenticated actor is the configured
// software owner. It contains no secrets.
type OwnerAccessResponse struct {
	Owner      bool   `json:"owner"`
	UserNumber string `json:"user_number"`
}

// LicenceSummary is an owner-console licence row. It deliberately never
// carries the licence key.
type LicenceSummary struct {
	Number             string     `json:"number"`
	OrganisationNumber string     `json:"organisation_number"`
	OrganisationName   string     `json:"organisation_name"`
	ValidFrom          time.Time  `json:"valid_from"`
	ValidTo            time.Time  `json:"valid_to"`
	GraceDays          int        `json:"grace_days"`
	MaxInstallations   int        `json:"max_installations"`
	CapacityUsed       int64      `json:"capacity_used"`
	LicensedAddons     []string   `json:"licensed_addons"`
	Status             string     `json:"status"`
	EffectiveStatus    string     `json:"effective_status"`
	RenewalRequestedAt *time.Time `json:"renewal_requested_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
}

// InstallationSummary is an owner-console installation row. There is no
// operating-system, heartbeat or device-health data because the backend does
// not collect it.
type InstallationSummary struct {
	Number               string     `json:"number"`
	InstallationIdentity string     `json:"installation_identity"`
	OrganisationNumber   string     `json:"organisation_number"`
	OrganisationName     string     `json:"organisation_name"`
	LicenceNumber        string     `json:"licence_number"`
	Status               string     `json:"status"`
	ActivatedAt          time.Time  `json:"activated_at"`
	LastSeenAt           *time.Time `json:"last_seen_at,omitempty"`
	ReleasedAt           *time.Time `json:"released_at,omitempty"`
	ReleaseReason        string     `json:"release_reason,omitempty"`
}

// OrganisationSummary is a directory row. CustomerStatus is derived from
// existing commercial/licensing records; it is not a separately stored
// customer entity.
type OrganisationSummary struct {
	Number                  string     `json:"number"`
	Name                    string     `json:"name"`
	Email                   string     `json:"email"`
	Phone                   string     `json:"phone"`
	CreatedAt               time.Time  `json:"created_at"`
	CustomerStatus          string     `json:"customer_status"`
	LicenceCount            int64      `json:"licence_count"`
	ActiveLicenceCount      int64      `json:"active_licence_count"`
	InstallationCount       int64      `json:"installation_count"`
	OwnerNames              []string   `json:"owner_names"`
	LatestLicenceExpiry     *time.Time `json:"latest_licence_expiry,omitempty"`
	HasCommercialProfile    bool       `json:"has_commercial_profile"`
	CommercialStatus        string     `json:"commercial_status,omitempty"`
	CustomerCategory        string     `json:"customer_category,omitempty"`
	CommercialCustomerSince *time.Time `json:"commercial_customer_since,omitempty"`
}

// OverviewResponse is the aggregate owner dashboard payload. It contains counts
// and short safe lists only.
type OverviewResponse struct {
	Organisations           int64                                     `json:"organisations"`
	CommercialOrganisations int64                                     `json:"commercial_organisations"`
	LicensedOrganisations   int64                                     `json:"licensed_organisations"`
	ActiveLicences          int64                                     `json:"active_licences"`
	NearingExpiry           int64                                     `json:"nearing_expiry"`
	ExpiredLicences         int64                                     `json:"expired_licences"`
	SuspendedLicences       int64                                     `json:"suspended_licences"`
	RevokedLicences         int64                                     `json:"revoked_licences"`
	TotalInstallations      int64                                     `json:"total_installations"`
	ActiveInstallations     int64                                     `json:"active_installations"`
	PendingRenewals         int64                                     `json:"pending_renewals"`
	RecentEvents            []licensingResponses.LicenceEventResponse `json:"recent_events"`
	AttentionLicences       []LicenceSummary                          `json:"attention_licences"`
	RecentInstallations     []InstallationSummary                     `json:"recent_installations"`
}

// OrganisationDetail is the owner customer profile. It uses existing
// organisations, owners, licences, installations and licence events.
type OrganisationDetail struct {
	Organisation       OrganisationSummary                          `json:"organisation"`
	Licences           []licensingResponses.LicenceResponse         `json:"licences"`
	Installations      []licensingResponses.InstallationResponse    `json:"installations"`
	Events             []licensingResponses.LicenceEventResponse    `json:"events"`
	SetupOrders        int64                                        `json:"setup_orders"`
	SubscriptionOrders int64                                        `json:"subscription_orders"`
	CommercialProfile  *commercialResponses.CustomerProfileResponse `json:"commercial_profile"`
}

// LicenceDetail is an owner licence view plus its event timeline.
type LicenceDetail struct {
	Licence licensingResponses.LicenceResponse        `json:"licence"`
	Events  []licensingResponses.LicenceEventResponse `json:"events"`
}

// InstallationDetail is an owner installation view plus its event timeline.
type InstallationDetail struct {
	Installation licensingResponses.InstallationResponse   `json:"installation"`
	Licence      licensingResponses.LicenceResponse        `json:"licence"`
	Events       []licensingResponses.LicenceEventResponse `json:"events"`
}

// Page is the common paginated envelope.
type Page[T any] struct {
	Items   []T   `json:"items"`
	Total   int64 `json:"total"`
	Page    int   `json:"page"`
	PerPage int   `json:"per_page"`
}
