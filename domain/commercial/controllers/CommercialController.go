package controllers

import (
	"errors"
	"net/http"
	"os"
	"strings"

	commercialServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commercial/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	licensingServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/services"
	"github.com/gin-gonic/gin"
)

// SimulatedPaymentsEnabled reports whether the test-payment routes may be
// registered/served. Dev-only; never enabled in production by a query param,
// browser flag or frontend value.
func SimulatedPaymentsEnabled() bool {
	return strings.EqualFold(strings.TrimSpace(os.Getenv("DEV_SIMULATED_PAYMENTS_ENABLED")), "true")
}

func fail(c *gin.Context, status int, err error) {
	c.JSON(status, commons.Feedback(err.Error()))
}

func commercialError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, commercialServices.ErrOrderNotFound):
		fail(c, http.StatusNotFound, err)
	case errors.Is(err, commercialServices.ErrPaymentNotPaid):
		fail(c, http.StatusPaymentRequired, err)
	case errors.Is(err, commercialServices.ErrTokenInvalid),
		errors.Is(err, commercialServices.ErrTokenWrongPurpose),
		errors.Is(err, commercialServices.ErrTokenAlreadyUsed):
		fail(c, http.StatusBadRequest, err)
	case errors.Is(err, commercialServices.ErrSubscriptionDuration),
		errors.Is(err, commercialServices.ErrNegativeQuantity),
		errors.Is(err, commercialServices.ErrUnknownService),
		errors.Is(err, commercialServices.ErrInvalidMethod),
		errors.Is(err, commercialServices.ErrCustomQuoteRequired),
		errors.Is(err, commercialServices.ErrUnknownLocation),
		errors.Is(err, commercialServices.ErrTransportManualQuote):
		fail(c, http.StatusBadRequest, err)
	default:
		fail(c, http.StatusInternalServerError, errors.New("the request could not be completed"))
	}
}

func currentOrganisation(c *gin.Context) (string, bool) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		fail(c, http.StatusUnauthorized, errors.New("Unauthorized"))
		return "", false
	}
	org, err := licensingServices.OrganisationNumberForUser(user.Number)
	if err != nil {
		fail(c, http.StatusConflict, errors.New("no active organisation for this account"))
		return "", false
	}
	return org, true
}

// ---- Installation packages ----

// ListPackages returns the active installation-package catalogue (public read).
func ListPackages(c *gin.Context) {
	packages, err := commercialServices.ListActivePackages()
	if err != nil {
		fail(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": packages})
}

type packageRequest struct {
	Code              string `json:"code"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	BasePriceUGX      int64  `json:"base_price_ugx"`
	IncludedDevices   int    `json:"included_devices"`
	IncludedLocations int    `json:"included_locations"`
	DevicePriceUGX    int64  `json:"device_price_ugx"`
	LocationPriceUGX  int64  `json:"location_price_ugx"`
	TransportBaseUGX  int64  `json:"transport_base_ugx"`
	TransportPerKmUGX int64  `json:"transport_per_km_ugx"`
	Methods           string `json:"methods"`
	ServicesJSON      string `json:"services_json"`
	EstimatedDays     int    `json:"estimated_days"`
	Active            bool   `json:"active"`
}

// CreatePackage persists a catalogue entry. The route is permission-gated
// (manage_ems_packages); administrators bypass permission checks.
func CreatePackage(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		fail(c, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	var req packageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, errors.New("invalid package payload"))
		return
	}
	pkg, err := commercialServices.CreatePackage(user.Number, commercialServices.PackageInput{
		Code: req.Code, Name: req.Name, Description: req.Description,
		BasePriceUGX: req.BasePriceUGX, IncludedDevices: req.IncludedDevices, IncludedLocations: req.IncludedLocations,
		DevicePriceUGX: req.DevicePriceUGX, LocationPriceUGX: req.LocationPriceUGX,
		TransportBaseUGX: req.TransportBaseUGX, TransportPerKmUGX: req.TransportPerKmUGX,
		Methods: req.Methods, ServicesJSON: req.ServicesJSON, EstimatedDays: req.EstimatedDays, Active: req.Active,
	})
	if err != nil {
		commercialError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": pkg})
}

// ---- Installation setup order ----

// ListLocations returns the self-hosted Ugandan location catalogue (name +
// approximate distance from Kampala city centre). Public, read-only; used by the
// installation form so transport can be resolved without any external map API.
func ListLocations(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": commercialServices.UgandaLocations()})
}

type setupOrderRequest struct {
	PackageNumber string   `json:"package_number"`
	ContactName   string   `json:"contact_name"`
	ContactEmail  string   `json:"contact_email"`
	ContactPhone  string   `json:"contact_phone"`
	Users         int      `json:"users"`
	Devices       int      `json:"devices"`
	Locations     int      `json:"locations"`
	Method        string   `json:"method"`
	Location      string   `json:"location"`
	TransportKm   int      `json:"transport_km"`
	ServiceCodes  []string `json:"service_codes"`
	Notes         string   `json:"notes"`
}

// CreateSetupOrder quotes and persists an installation order (public).
func CreateSetupOrder(c *gin.Context) {
	var req setupOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, errors.New("invalid configuration payload"))
		return
	}
	order, quote, err := commercialServices.CreateSetupOrder(commercialServices.SetupOrderInput{
		PackageNumber: req.PackageNumber, ContactName: req.ContactName, ContactEmail: req.ContactEmail,
		ContactPhone: req.ContactPhone, Users: req.Users, Devices: req.Devices, Locations: req.Locations,
		Method: req.Method, Location: req.Location, TransportKm: req.TransportKm,
		ServiceCodes: req.ServiceCodes, Notes: req.Notes,
	})
	if err != nil {
		commercialError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": gin.H{"order": order, "quote": quote}})
}

type scenarioRequest struct {
	Scenario string `json:"scenario"`
}

// SimulateSetupPayment records a test payment (dev-only). Only "success" issues
// a setup token, returned once.
func SimulateSetupPayment(c *gin.Context) {
	if !SimulatedPaymentsEnabled() {
		fail(c, http.StatusNotFound, errors.New("not found"))
		return
	}
	var req scenarioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, errors.New("invalid payment payload"))
		return
	}
	order, rawToken, err := commercialServices.SimulateSetupPayment(c.Param("number"), req.Scenario)
	if err != nil {
		if errors.Is(err, commercialServices.ErrPaymentNotPaid) {
			c.JSON(http.StatusPaymentRequired, gin.H{"message": err.Error(), "data": gin.H{"order": order, "token": ""}})
			return
		}
		commercialError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"order": order, "token": rawToken}})
}

type tokenRequest struct {
	Token string `json:"token"`
}

// ValidateSetupToken validates a raw setup token (public, one-shot read).
func ValidateSetupToken(c *gin.Context) {
	var req tokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, errors.New("invalid token payload"))
		return
	}
	order, err := commercialServices.ValidateSetupToken(req.Token)
	if err != nil {
		commercialError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": order})
}

// RedeemSetupToken consumes a setup token for the caller's organisation.
func RedeemSetupToken(c *gin.Context) {
	org, ok := currentOrganisation(c)
	if !ok {
		return
	}
	user, _ := commons.GetUserFromHeader(c)
	var req tokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, errors.New("invalid token payload"))
		return
	}
	order, err := commercialServices.RedeemSetupToken(req.Token, org, user.Number)
	if err != nil {
		commercialError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": order})
}

// ---- Subscription ----

type subscriptionDaysRequest struct {
	DurationDays int `json:"duration_days"`
}

// SubscriptionQuote returns the authoritative subscription amount (public).
func SubscriptionQuote(c *gin.Context) {
	var req subscriptionDaysRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, errors.New("invalid duration payload"))
		return
	}
	total, err := commercialServices.SubscriptionAmount(req.DurationDays)
	if err != nil {
		commercialError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"duration_days":  req.DurationDays,
		"daily_rate_ugx": commercialServices.SubscriptionDailyRateUGX,
		"total_ugx":      total,
		"currency":       "UGX",
		"min_days":       commercialServices.SubscriptionMinDays,
	}})
}

// CreateSubscriptionOrder quotes and persists a subscription order (auth).
func CreateSubscriptionOrder(c *gin.Context) {
	org, ok := currentOrganisation(c)
	if !ok {
		return
	}
	var req subscriptionDaysRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, errors.New("invalid duration payload"))
		return
	}
	order, err := commercialServices.CreateSubscriptionOrder(req.DurationDays, org)
	if err != nil {
		commercialError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": order})
}

// SimulateSubscriptionPayment records a test subscription payment (dev-only).
func SimulateSubscriptionPayment(c *gin.Context) {
	if !SimulatedPaymentsEnabled() {
		fail(c, http.StatusNotFound, errors.New("not found"))
		return
	}
	if _, ok := currentOrganisation(c); !ok {
		return
	}
	var req scenarioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, errors.New("invalid payment payload"))
		return
	}
	order, rawToken, err := commercialServices.SimulateSubscriptionPayment(c.Param("number"), req.Scenario)
	if err != nil {
		if errors.Is(err, commercialServices.ErrPaymentNotPaid) {
			c.JSON(http.StatusPaymentRequired, gin.H{"message": err.Error(), "data": gin.H{"order": order, "token": ""}})
			return
		}
		commercialError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"order": order, "token": rawToken}})
}

// CreateOnboardingSubscription creates a subscription order before sign-in,
// bound to a paid setup token (public). The successful payment will issue the
// account activation token.
func CreateOnboardingSubscription(c *gin.Context) {
	var req struct {
		SetupToken   string `json:"setup_token"`
		DurationDays int    `json:"duration_days"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, errors.New("invalid subscription payload"))
		return
	}
	order, err := commercialServices.CreateOnboardingSubscriptionOrder(req.SetupToken, req.DurationDays)
	if err != nil {
		commercialError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": order})
}

// SimulateOnboardingSubscriptionPayment records a test subscription payment
// (dev-only) and, only on success, issues the account activation token.
func SimulateOnboardingSubscriptionPayment(c *gin.Context) {
	if !SimulatedPaymentsEnabled() {
		fail(c, http.StatusNotFound, errors.New("not found"))
		return
	}
	var req scenarioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, errors.New("invalid payment payload"))
		return
	}
	order, rawToken, err := commercialServices.SimulateOnboardingSubscriptionPayment(c.Param("number"), req.Scenario)
	if err != nil {
		if errors.Is(err, commercialServices.ErrPaymentNotPaid) {
			c.JSON(http.StatusPaymentRequired, gin.H{"message": err.Error(), "data": gin.H{"order": order, "token": ""}})
			return
		}
		commercialError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"order": order, "token": rawToken}})
}

// ValidateSubscriptionToken validates the account activation token (public;
// idempotent read). It does not provision the licence.
func ValidateSubscriptionToken(c *gin.Context) {
	var req tokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, errors.New("invalid token payload"))
		return
	}
	order, err := commercialServices.ValidateSubscriptionToken(req.Token)
	if err != nil {
		commercialError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": order})
}

// ActivateSubscription redeems a subscription token and provisions/extend the
// organisation's existing EMS licence (auth).
func ActivateSubscription(c *gin.Context) {
	org, ok := currentOrganisation(c)
	if !ok {
		return
	}
	user, _ := commons.GetUserFromHeader(c)
	var req tokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, errors.New("invalid token payload"))
		return
	}
	order, licence, err := commercialServices.RedeemSubscriptionToken(req.Token, org, user.Number)
	if err != nil {
		commercialError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"order": order, "licence": licence}})
}

// SubscriptionStatus lists the organisation's subscription orders and its
// current licence status (auth).
func SubscriptionStatus(c *gin.Context) {
	org, ok := currentOrganisation(c)
	if !ok {
		return
	}
	user, _ := commons.GetUserFromHeader(c)
	orders, err := commercialServices.SubscriptionOrdersForOrganisation(org)
	if err != nil {
		fail(c, http.StatusInternalServerError, err)
		return
	}
	licence, licenceErr := licensingServices.ClientLicenceStatus(user.Number)
	if licenceErr != nil {
		licence = nil
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"orders": orders, "licence": licence}})
}
