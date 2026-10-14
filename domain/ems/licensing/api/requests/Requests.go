package requests

import "time"

// IssueLicenceRequest is the control-plane issue/assign payload.
//
// The licence key is generated server-side and never supplied by the client.
// The grace period is fixed at 14 days by the ratified architecture (F-Q2) and
// is therefore not part of this payload.
type IssueLicenceRequest struct {
	OrganisationNumber string    `json:"organisation_number"`
	ValidFrom          time.Time `json:"valid_from"`
	ValidTo            time.Time `json:"valid_to"`
	MaxInstallations   int       `json:"max_installations"`
	LicensedAddons     []string  `json:"licensed_addons"`
}

// ConfigureLicenceRequest updates editable licence configuration. A nil field
// is left unchanged.
type ConfigureLicenceRequest struct {
	ValidFrom        *time.Time `json:"valid_from"`
	ValidTo          *time.Time `json:"valid_to"`
	MaxInstallations *int       `json:"max_installations"`
	LicensedAddons   *[]string  `json:"licensed_addons"`
}

// RenewLicenceRequest extends the validity of an existing licence. Renewal
// never creates a new licence identity, key or entitlements (F-Q6).
type RenewLicenceRequest struct {
	ValidTo time.Time `json:"valid_to"`
}

// ActivateInstallationRequest is the client activation payload. The licence
// key is validated against the registry; the installation identity is an
// opaque client-generated identifier that never grants authority.
//
// AddonCode is optional: when present the licence must entitle that add-on.
// ReplacesInstallationNumber is optional: it links this activation to a
// previously owner-released installation on the same licence (§13 history
// linkage only; it never transfers capacity or historical records).
type ActivateInstallationRequest struct {
	LicenceKey                 string `json:"licence_key"`
	InstallationIdentity       string `json:"installation_identity"`
	AddonCode                  string `json:"addon_code"`
	ReplacesInstallationNumber string `json:"replaces_installation_number"`
}

// ReleaseInstallationRequest is the control-plane release payload. Release is
// the only action that frees installation capacity (F-Q3A).
type ReleaseInstallationRequest struct {
	Reason string `json:"reason"`
}

// ReplaceLicenceRequest is the control-plane new-licence issuance payload. It
// supersedes an ACTIVE licence with a brand-new licence (new number, new opaque
// key, its own validity, entitlements and capacity). Installations are never
// inherited or reattached; the superseded licence's capacity-holding
// installations are explicitly released with the mandatory reason (§13, §24
// detail 1).
type ReplaceLicenceRequest struct {
	OrganisationNumber string    `json:"organisation_number"`
	ValidFrom          time.Time `json:"valid_from"`
	ValidTo            time.Time `json:"valid_to"`
	MaxInstallations   int       `json:"max_installations"`
	LicensedAddons     []string  `json:"licensed_addons"`
	Reason             string    `json:"reason"`
}
