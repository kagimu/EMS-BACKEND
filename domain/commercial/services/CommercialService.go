package services

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	commercialModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commercial/models"
	addonCatalog "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/addons"
	auditServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/audit/services"
	licensingModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// SeedingEnabled reports whether the illustrative installation-package
// catalogue may be seeded. It is opt-in via DEV_SEED_PACKAGES=true and is the
// single gate used by startup seeding and by the catalogue read self-heal.
func SeedingEnabled() bool {
	return strings.EqualFold(strings.TrimSpace(os.Getenv("DEV_SEED_PACKAGES")), "true")
}

var (
	ErrOrderNotFound       = errors.New("order not found")
	ErrOrderNotPayable     = errors.New("this order is not awaiting payment")
	ErrPaymentNotPaid      = errors.New("payment was not successful; no token issued")
	ErrTokenInvalid        = errors.New("activation token is not valid")
	ErrTokenWrongPurpose   = errors.New("activation token is not valid for this action")
	ErrTokenAlreadyUsed    = errors.New("activation token has already been used")
	ErrOrganisationMissing = errors.New("an organisation is required for this activation")
)

// newUUID returns a fresh random identifier.
func newUUID() string { return uuid.NewV4().String() }

// randomKey returns a 64-char hex key (fits the licence key column).
func randomKey() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

// defaultAddonsJSON returns the addon codes that are default-enabled in the EMS
// catalogue (excluding core). A new subscription licence is entitled to these,
// so addons such as Education are reachable after activation. It uses the
// existing catalogue — no parallel entitlement system.
func defaultAddonsJSON() string {
	addons := make([]string, 0)
	for _, def := range addonCatalog.Catalog {
		if def.Type == addonCatalog.ModuleTypeAddon && def.DefaultEnabled {
			addons = append(addons, def.Code)
		}
	}
	data, _ := json.Marshal(addons)
	return string(data)
}

func audit(action, entityType, entityID, actor, org string, newValues interface{}) {
	auditServices.RecordFact(auditServices.AuditRecord{
		OrganisationNumber: org,
		ActorUserNumber:    actor,
		Action:             action,
		EntityType:         entityType,
		EntityID:           entityID,
		NewValues:          newValues,
	})
}

// PackageInput is an administrator-supplied installation package definition.
type PackageInput struct {
	Code              string
	Name              string
	Description       string
	BasePriceUGX      int64
	IncludedDevices   int
	IncludedLocations int
	DevicePriceUGX    int64
	LocationPriceUGX  int64
	TransportBaseUGX  int64
	TransportPerKmUGX int64
	Methods           string
	ServicesJSON      string
	EstimatedDays     int
	Active            bool
}

// CreatePackage persists a new catalogue entry (authorised administrators only;
// the route enforces the permission).
func CreatePackage(actor string, in PackageInput) (*commercialModels.InstallationPackage, error) {
	pkg := commercialModels.InstallationPackage{
		Number:            newUUID(),
		Code:              strings.TrimSpace(in.Code),
		Name:              strings.TrimSpace(in.Name),
		Description:       in.Description,
		BasePriceUGX:      in.BasePriceUGX,
		IncludedDevices:   in.IncludedDevices,
		IncludedLocations: in.IncludedLocations,
		DevicePriceUGX:    in.DevicePriceUGX,
		LocationPriceUGX:  in.LocationPriceUGX,
		TransportBaseUGX:  in.TransportBaseUGX,
		TransportPerKmUGX: in.TransportPerKmUGX,
		Methods:           in.Methods,
		ServicesJSON:      in.ServicesJSON,
		EstimatedDays:     in.EstimatedDays,
		Active:            in.Active,
		CreatedByNumber:   actor,
	}
	if pkg.Name == "" {
		return nil, errors.New("package name is required")
	}
	if err := config.DB.Create(&pkg).Error; err != nil {
		return nil, err
	}
	audit("commercial.package.created", "installation_package", pkg.Number, actor, "", map[string]interface{}{
		"code": pkg.Code, "base_price_ugx": pkg.BasePriceUGX,
	})
	return &pkg, nil
}

// ListActivePackages returns the client-facing catalogue.
func ListActivePackages() ([]commercialModels.InstallationPackage, error) {
	var packages []commercialModels.InstallationPackage
	err := config.DB.Where("active = ?", true).Order("created_at asc").Find(&packages).Error
	return packages, err
}

func findPackage(number string) (*commercialModels.InstallationPackage, error) {
	var pkg commercialModels.InstallationPackage
	if err := config.DB.Where("number = ? AND active = ?", number, true).First(&pkg).Error; err != nil {
		return nil, ErrUnknownPackage
	}
	return &pkg, nil
}

// SetupOrderInput is the client's installation configuration (no total).
type SetupOrderInput struct {
	PackageNumber string
	ContactName   string
	ContactEmail  string
	ContactPhone  string
	Users         int
	Devices       int
	Locations     int
	Method        string
	Location      string
	TransportKm   int
	ServiceCodes  []string
	Notes         string
}

// CreateSetupOrder quotes and persists an installation order (AWAITING_PAYMENT).
// The authoritative total is calculated here and snapshotted.
func CreateSetupOrder(in SetupOrderInput) (*commercialModels.SetupOrder, Quote, error) {
	pkg, err := findPackage(in.PackageNumber)
	if err != nil {
		return nil, Quote{}, err
	}
	if err := ValidateMethod(*pkg, in.Method); err != nil {
		return nil, Quote{}, err
	}
	quote, err := QuoteInstallation(*pkg, QuoteInput{
		Users:        in.Users,
		Devices:      in.Devices,
		Locations:    in.Locations,
		Method:       in.Method,
		Location:     in.Location,
		TransportKm:  in.TransportKm,
		ServiceCodes: in.ServiceCodes,
	})
	if err != nil {
		return nil, Quote{}, err
	}

	configJSON, _ := json.Marshal(map[string]interface{}{
		"users": in.Users, "devices": in.Devices, "locations": in.Locations, "method": in.Method,
		"location": in.Location, "transport_km": in.TransportKm, "service_codes": in.ServiceCodes, "notes": in.Notes,
	})
	quoteJSON, _ := json.Marshal(quote)

	order := commercialModels.SetupOrder{
		Number:            newUUID(),
		PackageNumber:     pkg.Number,
		Status:            "AWAITING_PAYMENT",
		Currency:          "UGX",
		ContactName:       strings.TrimSpace(in.ContactName),
		ContactEmail:      strings.TrimSpace(in.ContactEmail),
		ContactPhone:      strings.TrimSpace(in.ContactPhone),
		ConfigurationJSON: string(configJSON),
		QuoteSnapshotJSON: string(quoteJSON),
		SetupTotalUGX:     quote.TotalUGX,
	}
	if err := config.DB.Create(&order).Error; err != nil {
		return nil, Quote{}, err
	}
	audit("commercial.setup_order.created", "setup_order", order.Number, "", "", map[string]interface{}{
		"total_ugx": order.SetupTotalUGX, "package": pkg.Number,
	})
	return &order, quote, nil
}

// SimulateSetupPayment records a test payment. Only a "success" scenario issues
// the setup token. The provider is simulated and this is never authoritative in
// production (routes are dev-gated).
func SimulateSetupPayment(orderNumber, scenario string) (*commercialModels.SetupOrder, string, error) {
	var order commercialModels.SetupOrder
	if err := config.DB.Where("number = ?", orderNumber).First(&order).Error; err != nil {
		return nil, "", ErrOrderNotFound
	}
	if order.Status != "AWAITING_PAYMENT" {
		return nil, "", ErrOrderNotPayable
	}

	status := "PENDING"
	switch scenario {
	case "success":
		status = "PAID"
	case "failed":
		status = "FAILED"
	case "pending":
		status = "PENDING"
	default:
		return nil, "", errors.New("unknown payment scenario")
	}

	payment := commercialModels.SimulatedPayment{
		Number:    newUUID(),
		Reference: order.Number,
		Kind:      "setup",
		AmountUGX: order.SetupTotalUGX,
		Scenario:  scenario,
		Status:    status,
		Provider:  "simulated",
	}
	if err := config.DB.Create(&payment).Error; err != nil {
		return nil, "", err
	}
	order.SimulatedPaymentNumber = payment.Number

	if status != "PAID" {
		config.DB.Save(&order)
		audit("commercial.setup_payment."+strings.ToLower(status), "setup_order", order.Number, "", "", map[string]interface{}{"status": status})
		return &order, "", ErrPaymentNotPaid
	}

	raw, err := GenerateRawToken()
	if err != nil {
		return nil, "", err
	}
	now := time.Now().UTC()
	token := commercialModels.ActivationToken{
		Number:           newUUID(),
		Purpose:          TokenPurposeSetup,
		TokenHash:        HashToken(raw),
		SetupOrderNumber: order.Number,
		Status:           TokenStatusIssued,
		IssuedAt:         &now,
	}
	if err := config.DB.Create(&token).Error; err != nil {
		return nil, "", err
	}
	order.Status = "TOKEN_ISSUED"
	if err := config.DB.Save(&order).Error; err != nil {
		return nil, "", err
	}
	audit("commercial.setup_token.issued", "setup_order", order.Number, "", "", map[string]interface{}{"token": token.Number})
	return &order, raw, nil
}

// ValidateSetupToken validates a raw setup token and returns its order. It is
// idempotent for an already-validated token and rejects redeemed/invalid ones.
func ValidateSetupToken(raw string) (*commercialModels.SetupOrder, error) {
	token, err := findToken(raw, TokenPurposeSetup)
	if err != nil {
		return nil, err
	}
	if token.Status == TokenStatusRedeemed {
		return nil, ErrTokenAlreadyUsed
	}
	if token.Status == TokenStatusIssued {
		now := time.Now().UTC()
		token.Status = TokenStatusValidated
		token.ValidatedAt = &now
		if err := config.DB.Save(token).Error; err != nil {
			return nil, err
		}
	}
	var order commercialModels.SetupOrder
	if err := config.DB.Where("number = ?", token.SetupOrderNumber).First(&order).Error; err != nil {
		return nil, ErrOrderNotFound
	}
	return &order, nil
}

// RedeemSetupToken consumes a setup token for the given organisation (one-shot).
func RedeemSetupToken(raw, organisationNumber, actor string) (*commercialModels.SetupOrder, error) {
	if strings.TrimSpace(organisationNumber) == "" {
		return nil, ErrOrganisationMissing
	}
	token, err := findToken(raw, TokenPurposeSetup)
	if err != nil {
		return nil, err
	}
	if token.Status == TokenStatusRedeemed {
		return nil, ErrTokenAlreadyUsed
	}
	var order commercialModels.SetupOrder
	if err := config.DB.Where("number = ?", token.SetupOrderNumber).First(&order).Error; err != nil {
		return nil, ErrOrderNotFound
	}

	now := time.Now().UTC()
	token.Status = TokenStatusRedeemed
	token.RedeemedAt = &now
	token.RedeemedOrganisation = organisationNumber
	if err := config.DB.Save(token).Error; err != nil {
		return nil, err
	}
	order.Status = "COMPLETED"
	order.OrganisationNumber = organisationNumber
	if err := config.DB.Save(&order).Error; err != nil {
		return nil, err
	}
	audit("commercial.setup_token.redeemed", "setup_order", order.Number, actor, organisationNumber, nil)
	return &order, nil
}

func findToken(raw, purpose string) (*commercialModels.ActivationToken, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, ErrTokenInvalid
	}
	var token commercialModels.ActivationToken
	if err := config.DB.Where("token_hash = ?", HashToken(raw)).First(&token).Error; err != nil {
		return nil, ErrTokenInvalid
	}
	if token.Purpose != purpose {
		return nil, ErrTokenWrongPurpose
	}
	return &token, nil
}

// CreateSubscriptionOrder quotes and persists a subscription order.
func CreateSubscriptionOrder(days int, organisationNumber string) (*commercialModels.SubscriptionOrder, error) {
	total, err := SubscriptionAmount(days)
	if err != nil {
		return nil, err
	}
	order := commercialModels.SubscriptionOrder{
		Number:             newUUID(),
		Status:             "AWAITING_PAYMENT",
		Currency:           "UGX",
		DurationDays:       days,
		DailyRateUGX:       SubscriptionDailyRateUGX,
		TotalUGX:           total,
		OrganisationNumber: organisationNumber,
	}
	if err := config.DB.Create(&order).Error; err != nil {
		return nil, err
	}
	audit("commercial.subscription_order.created", "subscription_order", order.Number, "", organisationNumber, map[string]interface{}{
		"days": days, "total_ugx": total, "rate": SubscriptionDailyRateUGX,
	})
	return &order, nil
}

// CreateOnboardingSubscriptionOrder creates a subscription order BEFORE sign-in,
// bound to a paid setup token (no organisation yet). Used so the subscription
// payment can issue the account activation token.
func CreateOnboardingSubscriptionOrder(setupToken string, days int) (*commercialModels.SubscriptionOrder, error) {
	token, err := findToken(setupToken, TokenPurposeSetup)
	if err != nil {
		return nil, err
	}
	if token.Status == TokenStatusRedeemed {
		return nil, ErrTokenAlreadyUsed
	}
	total, err := SubscriptionAmount(days)
	if err != nil {
		return nil, err
	}
	order := commercialModels.SubscriptionOrder{
		Number:           newUUID(),
		Status:           "AWAITING_PAYMENT",
		Currency:         "UGX",
		DurationDays:     days,
		DailyRateUGX:     SubscriptionDailyRateUGX,
		TotalUGX:         total,
		SetupOrderNumber: token.SetupOrderNumber,
	}
	if err := config.DB.Create(&order).Error; err != nil {
		return nil, err
	}
	audit("commercial.onboarding_subscription.created", "subscription_order", order.Number, "", "", map[string]interface{}{
		"days": days, "total_ugx": total, "rate": SubscriptionDailyRateUGX, "setup_order": token.SetupOrderNumber,
	})
	return &order, nil
}

// SimulateOnboardingSubscriptionPayment records a test payment for an onboarding
// subscription and, only on success, issues the account activation token. It is
// the public counterpart of SimulateSubscriptionPayment (same underlying logic,
// no organisation required).
func SimulateOnboardingSubscriptionPayment(orderNumber, scenario string) (*commercialModels.SubscriptionOrder, string, error) {
	return SimulateSubscriptionPayment(orderNumber, scenario)
}

// ValidateSubscriptionToken validates a subscription (account) activation token
// and returns its order. Idempotent for an already-validated token; a redeemed
// token is rejected. It does NOT provision the licence (that happens later,
// against the organisation).
func ValidateSubscriptionToken(raw string) (*commercialModels.SubscriptionOrder, error) {
	token, err := findToken(raw, TokenPurposeSubscription)
	if err != nil {
		return nil, err
	}
	if token.Status == TokenStatusRedeemed {
		return nil, ErrTokenAlreadyUsed
	}
	if token.Status == TokenStatusIssued {
		now := time.Now().UTC()
		token.Status = TokenStatusValidated
		token.ValidatedAt = &now
		if err := config.DB.Save(token).Error; err != nil {
			return nil, err
		}
	}
	var order commercialModels.SubscriptionOrder
	if err := config.DB.Where("number = ?", token.SubscriptionOrderNumber).First(&order).Error; err != nil {
		return nil, ErrOrderNotFound
	}
	return &order, nil
}

// SimulateSubscriptionPayment records a test payment and, only on success,
// issues the subscription token.
func SimulateSubscriptionPayment(orderNumber, scenario string) (*commercialModels.SubscriptionOrder, string, error) {
	var order commercialModels.SubscriptionOrder
	if err := config.DB.Where("number = ?", orderNumber).First(&order).Error; err != nil {
		return nil, "", ErrOrderNotFound
	}
	if order.Status != "AWAITING_PAYMENT" {
		return nil, "", ErrOrderNotPayable
	}

	status := "PENDING"
	switch scenario {
	case "success":
		status = "PAID"
	case "failed":
		status = "FAILED"
	case "pending":
		status = "PENDING"
	default:
		return nil, "", errors.New("unknown payment scenario")
	}

	payment := commercialModels.SimulatedPayment{
		Number:    newUUID(),
		Reference: order.Number,
		Kind:      "subscription",
		AmountUGX: order.TotalUGX,
		Scenario:  scenario,
		Status:    status,
		Provider:  "simulated",
	}
	if err := config.DB.Create(&payment).Error; err != nil {
		return nil, "", err
	}
	order.SimulatedPaymentNumber = payment.Number

	if status != "PAID" {
		config.DB.Save(&order)
		return &order, "", ErrPaymentNotPaid
	}

	raw, err := GenerateRawToken()
	if err != nil {
		return nil, "", err
	}
	now := time.Now().UTC()
	token := commercialModels.ActivationToken{
		Number:                  newUUID(),
		Purpose:                 TokenPurposeSubscription,
		TokenHash:               HashToken(raw),
		SubscriptionOrderNumber: order.Number,
		Status:                  TokenStatusIssued,
		IssuedAt:                &now,
	}
	if err := config.DB.Create(&token).Error; err != nil {
		return nil, "", err
	}
	order.Status = "TOKEN_ISSUED"
	if err := config.DB.Save(&order).Error; err != nil {
		return nil, "", err
	}
	audit("commercial.subscription_token.issued", "subscription_order", order.Number, "", order.OrganisationNumber, nil)
	return &order, raw, nil
}

// RedeemSubscriptionToken consumes a subscription token (one-shot) and, on a
// paid order, provisions/extend the organisation's existing EMS licence by the
// purchased period using the existing licensing tables. Provisioning is
// idempotent because the token becomes REDEEMED.
func RedeemSubscriptionToken(raw, organisationNumber, actor string) (*commercialModels.SubscriptionOrder, *licensingModels.Licence, error) {
	if strings.TrimSpace(organisationNumber) == "" {
		return nil, nil, ErrOrganisationMissing
	}
	token, err := findToken(raw, TokenPurposeSubscription)
	if err != nil {
		return nil, nil, err
	}
	if token.Status == TokenStatusRedeemed {
		return nil, nil, ErrTokenAlreadyUsed
	}
	var order commercialModels.SubscriptionOrder
	if err := config.DB.Where("number = ?", token.SubscriptionOrderNumber).First(&order).Error; err != nil {
		return nil, nil, ErrOrderNotFound
	}
	if order.Status != "TOKEN_ISSUED" && order.Status != "PAID" {
		return nil, nil, ErrPaymentNotPaid
	}

	now := time.Now().UTC()
	token.Status = TokenStatusRedeemed
	token.RedeemedAt = &now
	token.RedeemedOrganisation = organisationNumber
	if err := config.DB.Save(token).Error; err != nil {
		return nil, nil, err
	}

	licence, err := provisionSubscriptionLicence(organisationNumber, order.DurationDays, actor)
	if err != nil {
		return nil, nil, err
	}

	order.Status = "ACTIVE"
	order.OrganisationNumber = organisationNumber
	if err := config.DB.Save(&order).Error; err != nil {
		return nil, nil, err
	}
	audit("commercial.subscription_token.redeemed", "subscription_order", order.Number, actor, organisationNumber, map[string]interface{}{
		"licence": licence.Number, "days": order.DurationDays,
	})
	return &order, licence, nil
}

// SubscriptionOrdersForOrganisation lists an organisation's subscription orders.
func SubscriptionOrdersForOrganisation(organisationNumber string) ([]commercialModels.SubscriptionOrder, error) {
	var orders []commercialModels.SubscriptionOrder
	err := config.DB.Where("organisation_number = ?", organisationNumber).Order("created_at desc").Find(&orders).Error
	return orders, err
}

// provisionSubscriptionLicence creates the organisation's ACTIVE EMS licence or
// extends it by the purchased period. It writes only the existing
// ems_licences / ems_licence_events tables (no parallel entitlement system) and
// never creates a second ACTIVE licence for the same organisation.
func provisionSubscriptionLicence(organisationNumber string, days int, actor string) (*licensingModels.Licence, error) {
	if days < SubscriptionMinDays {
		return nil, ErrSubscriptionDuration
	}
	period := time.Duration(days) * 24 * time.Hour
	now := time.Now().UTC()

	var licence licensingModels.Licence
	err := config.DB.Where("organisation_number = ? AND status = ?", organisationNumber, licensingModels.StatusActive).First(&licence).Error

	if err == nil {
		// Extend from the later of now and the current expiry.
		base := licence.ValidTo
		if base.Before(now) {
			base = now
		}
		before := map[string]interface{}{"valid_to": licence.ValidTo}
		licence.ValidTo = base.Add(period)
		if saveErr := config.DB.Save(&licence).Error; saveErr != nil {
			return nil, saveErr
		}
		recordLicenceEvent(licensingModels.LicenceEventRenewed, licence.Number, organisationNumber, actor, before,
			map[string]interface{}{"valid_to": licence.ValidTo})
		return &licence, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	key, keyErr := randomKey()
	if keyErr != nil {
		return nil, keyErr
	}
	org := organisationNumber
	created := licensingModels.Licence{
		Number:             newUUID(),
		LicenceKey:         key,
		OrganisationNumber: organisationNumber,
		ValidFrom:          now,
		ValidTo:            now.Add(period),
		GraceDays:          licensingModels.DefaultGraceDays,
		MaxInstallations:   1,
		LicensedAddons:     defaultAddonsJSON(),
		Status:             licensingModels.StatusActive,
		ActiveOrgNumber:    &org,
		IssuedBy:           actor,
	}
	if createErr := config.DB.Create(&created).Error; createErr != nil {
		return nil, createErr
	}
	recordLicenceEvent(licensingModels.LicenceEventIssued, created.Number, organisationNumber, actor, nil,
		map[string]interface{}{"valid_to": created.ValidTo, "source": "subscription"})
	return &created, nil
}

func recordLicenceEvent(eventType licensingModels.LicenceEventType, licenceNumber, organisationNumber, actor string, before, after interface{}) {
	beforeJSON, _ := json.Marshal(before)
	afterJSON, _ := json.Marshal(after)
	actorNumber := actor
	if strings.TrimSpace(actorNumber) == "" {
		actorNumber = "system"
	}
	// Non-fatal: licence facts must never roll back the business action.
	config.DB.Create(&licensingModels.LicenceEvent{
		EventType:          eventType,
		LicenceNumber:      licenceNumber,
		OrganisationNumber: organisationNumber,
		ActorUserNumber:    actorNumber,
		BeforeState:        string(beforeJSON),
		AfterState:         string(afterJSON),
	})
}

// testTransportBands is the illustrative per-round-trip travel pricing used for
// testing (>100 km requires a manual quotation; see PackageTransportBands).
const testTransportBandsJSON = `[{"max_km":10,"amount_ugx":20000},{"max_km":25,"amount_ugx":40000},{"max_km":50,"amount_ugx":80000},{"max_km":100,"amount_ugx":150000}]`

// testServicesJSON is the illustrative add-on service catalogue.
const testServicesJSON = `[{"code":"remote-training","name":"Additional remote training (per hour)","price_ugx":50000},{"code":"onsite-engineer-day","name":"Additional on-site engineer-day","price_ugx":150000}]`

type testPackage struct {
	Code                  string
	Name                  string
	BasePriceUGX          int64
	IncludedUsers         int
	IncludedDevices       int
	IncludedLocations     int
	TrainingHours         int
	EstimatedDays         int
	TransportStandardUGX  int64
	DeploymentDescription string
	Custom                bool
}

// SeedTestPackages inserts the illustrative installation-package catalogue into
// the connected database. It is idempotent (matched by code — never duplicated)
// and MUST only be called in the dedicated disposable onboarding environment
// (gated by DEV_SEED_PACKAGES). The prices are ILLUSTRATIVE TEST PRICES, not an
// approved commercial price list. Standard transport: small setups 20,000;
// big setups 50,000 (charged only outside Kampala city centre).
func SeedTestPackages() (int, error) {
	userPrice := int64(5000)
	devicePrice := int64(15000)
	locationPrice := int64(75000)

	definitions := []testPackage{
		{"ems-starter", "EMS Starter", 250000, 5, 5, 1, 1, 1, 20000, "Single-computer installation for a small business.", false},
		{"ems-business", "EMS Business", 450000, 15, 15, 1, 2, 2, 20000, "Client-owned server with local-network access (1–2 working days).", false},
		{"ems-professional", "EMS Professional", 750000, 30, 30, 2, 4, 3, 50000, "Client-owned server with local-network access (2–3 working days). Multi-location deployment must be confirmed as technically feasible.", false},
		{"ems-enterprise", "EMS Enterprise", 1500000, 75, 75, 3, 8, 5, 50000, "Client-owned server with local-network access (3–5 working days). Confirm infrastructure requirements before promising the deployment.", false},
		{"ems-multi-branch", "EMS Multi-Branch", 2500000, 100, 100, 5, 12, 7, 50000, "Assessment required (5–7 working days). Multi-branch synchronisation or connectivity is not supported until verified.", false},
		{"ems-custom", "EMS Custom Deployment", 0, 0, 0, 0, 0, 0, 0, "Custom deployment — scope, allowances, training and duration are agreed after assessment (custom quotation).", true},
	}

	created := 0
	for _, def := range definitions {
		var existing commercialModels.InstallationPackage
		err := config.DB.Where("code = ?", def.Code).First(&existing).Error

		if err == nil {
			// Converge the illustrative row (dev-only) so re-runs apply the
			// current test catalogue, including the standard transport fee.
			existing.Name = def.Name
			existing.Description = def.DeploymentDescription
			existing.BasePriceUGX = def.BasePriceUGX
			existing.IncludedUsers = def.IncludedUsers
			existing.IncludedDevices = def.IncludedDevices
			existing.IncludedLocations = def.IncludedLocations
			existing.TrainingHours = def.TrainingHours
			existing.EstimatedDays = def.EstimatedDays
			existing.TransportStandardUGX = def.TransportStandardUGX
			existing.DeploymentDescription = def.DeploymentDescription
			existing.Custom = def.Custom
			existing.Active = true
			if err := config.DB.Save(&existing).Error; err != nil {
				return created, err
			}
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return created, err
		}

		pkg := commercialModels.InstallationPackage{
			Number:                newUUID(),
			Code:                  def.Code,
			Name:                  def.Name,
			Description:           def.DeploymentDescription,
			BasePriceUGX:          def.BasePriceUGX,
			IncludedUsers:         def.IncludedUsers,
			IncludedDevices:       def.IncludedDevices,
			IncludedLocations:     def.IncludedLocations,
			UserPriceUGX:          userPrice,
			DevicePriceUGX:        devicePrice,
			LocationPriceUGX:      locationPrice,
			TransportStandardUGX:  def.TransportStandardUGX,
			TransportBandsJSON:    testTransportBandsJSON,
			TrainingHours:         def.TrainingHours,
			Methods:               "onsite,remote",
			ServicesJSON:          testServicesJSON,
			EstimatedDays:         def.EstimatedDays,
			DeploymentDescription: def.DeploymentDescription,
			Custom:                def.Custom,
			Active:                true,
			CreatedByNumber:       "dev-seed",
		}
		if err := config.DB.Create(&pkg).Error; err != nil {
			return created, err
		}
		created++
	}
	return created, nil
}
