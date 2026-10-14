package responses

import "time"

// LicenceResponse is the control-plane licence view. It never carries the
// opaque licence key: ordinary reads must not expose it.
type LicenceResponse struct {
	Number             string     `json:"number"`
	OrganisationNumber string     `json:"organisation_number"`
	ValidFrom          time.Time  `json:"valid_from"`
	ValidTo            time.Time  `json:"valid_to"`
	GraceDays          int        `json:"grace_days"`
	MaxInstallations   int        `json:"max_installations"`
	LicensedAddons     []string   `json:"licensed_addons"`
	Status             string     `json:"status"`
	EffectiveStatus    string     `json:"effective_status"`
	ReplacedByNumber   string     `json:"replaced_by_number,omitempty"`
	IssuedBy           string     `json:"issued_by"`
	LastActedBy        string     `json:"last_acted_by,omitempty"`
	RenewalRequestedBy string     `json:"renewal_requested_by,omitempty"`
	RenewalRequestedAt *time.Time `json:"renewal_requested_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// IssuedLicenceResponse is returned only by the issuance operation and is the
// only representation that carries the newly generated licence key, which the
// control plane delivers to the client.
type IssuedLicenceResponse struct {
	LicenceResponse
	LicenceKey string `json:"licence_key"`
}

// LicenceStatusResponse is the client-plane view of the organisation's active
// licence. It never carries the licence key.
type LicenceStatusResponse struct {
	Number             string     `json:"number"`
	EffectiveStatus    string     `json:"effective_status"`
	ValidFrom          time.Time  `json:"valid_from"`
	ValidTo            time.Time  `json:"valid_to"`
	GraceDays          int        `json:"grace_days"`
	LicensedAddons     []string   `json:"licensed_addons"`
	MaxInstallations   int        `json:"max_installations"`
	RenewalRequestedAt *time.Time `json:"renewal_requested_at,omitempty"`
}

// InstallationResponse is the installation view used by both planes. It never
// carries the licence key.
type InstallationResponse struct {
	Number               string     `json:"number"`
	InstallationIdentity string     `json:"installation_identity"`
	OrganisationNumber   string     `json:"organisation_number"`
	LicenceNumber        string     `json:"licence_number"`
	Status               string     `json:"status"`
	ActivatedAt          time.Time  `json:"activated_at"`
	ActivatedBy          string     `json:"activated_by,omitempty"`
	LastSeenAt           *time.Time `json:"last_seen_at,omitempty"`
	ClientDeactivatedAt  *time.Time `json:"client_deactivated_at,omitempty"`
	ClientDeactivatedBy  string     `json:"client_deactivated_by,omitempty"`
	ReleasedBy           string     `json:"released_by,omitempty"`
	ReleasedAt           *time.Time `json:"released_at,omitempty"`
	ReleaseReason        string     `json:"release_reason,omitempty"`
	ReplacedByNumber     string     `json:"replaced_by_number,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

// ActivationResponse is returned when an installation activates or
// reactivates. It never carries the licence key.
type ActivationResponse struct {
	InstallationNumber   string     `json:"installation_number"`
	InstallationIdentity string     `json:"installation_identity"`
	OrganisationNumber   string     `json:"organisation_number"`
	LicenceNumber        string     `json:"licence_number"`
	Status               string     `json:"status"`
	LicenceStatus        string     `json:"licence_status"`
	EffectiveStatus      string     `json:"effective_status"`
	ActivatedAt          time.Time  `json:"activated_at"`
	LastSeenAt           *time.Time `json:"last_seen_at,omitempty"`
	LicensedAddons       []string   `json:"licensed_addons"`
	MaxInstallations     int        `json:"max_installations"`
	CapacityUsed         int        `json:"capacity_used"`
}

// LicenceEventResponse is the control-plane view of an append-only licensing
// fact. It never carries the licence key: facts reference the licence by
// number only.
type LicenceEventResponse struct {
	Number               uint      `json:"number"`
	EventType            string    `json:"event_type"`
	LicenceNumber        string    `json:"licence_number"`
	OrganisationNumber   string    `json:"organisation_number"`
	ActorUserNumber      string    `json:"actor_user_number"`
	InstallationNumber   string    `json:"installation_number,omitempty"`
	InstallationIdentity string    `json:"installation_identity,omitempty"`
	Reason               string    `json:"reason,omitempty"`
	ClientInitiated      bool      `json:"client_initiated"`
	CreatedAt            time.Time `json:"created_at"`
}
