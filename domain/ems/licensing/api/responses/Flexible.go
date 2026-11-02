package responses

import "time"

type ScopeResponse struct {
	Number     string    `json:"number"`
	ScopeType  string    `json:"scope_type"`
	ScopeKey   string    `json:"scope_key,omitempty"`
	Quantity   int       `json:"quantity"`
	Used       int64     `json:"used"`
	Available  int64     `json:"available"`
	Status     string    `json:"status"`
	ValidFrom  time.Time `json:"valid_from"`
	ValidTo    time.Time `json:"valid_to"`
	ConfigJSON string    `json:"config_json,omitempty"`
}

type EntitlementResponse struct {
	Number     string    `json:"number"`
	Code       string    `json:"code"`
	Quantity   int       `json:"quantity"`
	Used       int64     `json:"used"`
	Status     string    `json:"status"`
	ValidFrom  time.Time `json:"valid_from"`
	ValidTo    time.Time `json:"valid_to"`
	ConfigJSON string    `json:"config_json,omitempty"`
}

type LicensedUnitResponse struct {
	Number             string     `json:"number"`
	LicenceNumber      string     `json:"licence_number"`
	ScopeNumber        string     `json:"scope_number"`
	UnitType           string     `json:"unit_type"`
	UnitIdentifier     string     `json:"unit_identifier"`
	OrganisationNumber string     `json:"organisation_number"`
	SiteNumber         string     `json:"site_number,omitempty"`
	Status             string     `json:"status"`
	RegisteredAt       time.Time  `json:"registered_at"`
	ReleasedAt         *time.Time `json:"released_at,omitempty"`
	ReleaseReason      string     `json:"release_reason,omitempty"`
}

type FlexibleLicenceResponse struct {
	LicenceNumber string                 `json:"licence_number"`
	Scopes        []ScopeResponse        `json:"scopes"`
	Entitlements  []EntitlementResponse  `json:"entitlements"`
	Units         []LicensedUnitResponse `json:"units"`
}
