package requests

import "time"

type FlexibleScopeInput struct {
	ScopeType  string    `json:"scope_type"`
	ScopeKey   string    `json:"scope_key"`
	Quantity   int       `json:"quantity"`
	Status     string    `json:"status"`
	ValidFrom  time.Time `json:"valid_from"`
	ValidTo    time.Time `json:"valid_to"`
	ConfigJSON string    `json:"config_json"`
}

type FlexibleEntitlementInput struct {
	Code       string    `json:"code"`
	Quantity   int       `json:"quantity"`
	Status     string    `json:"status"`
	ValidFrom  time.Time `json:"valid_from"`
	ValidTo    time.Time `json:"valid_to"`
	ConfigJSON string    `json:"config_json"`
}

type ConfigureFlexibleRequest struct {
	Scopes       []FlexibleScopeInput       `json:"scopes"`
	Entitlements []FlexibleEntitlementInput `json:"entitlements"`
	Reason       string                     `json:"reason"`
}

type RegisterUnitRequest struct {
	ScopeType      string `json:"scope_type"`
	UnitIdentifier string `json:"unit_identifier"`
	SiteNumber     string `json:"site_number"`
	ConfigJSON     string `json:"config_json"`
}

type ReleaseUnitRequest struct {
	Reason string `json:"reason"`
}
