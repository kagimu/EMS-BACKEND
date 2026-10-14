package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	commercialModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commercial/models"
)

// Authoritative commercial pricing rules. The backend always recalculates from
// these; a client-supplied total is never trusted.

const (
	// SubscriptionDailyRateUGX is the single authoritative subscription rate.
	SubscriptionDailyRateUGX int64 = 1200
	// SubscriptionMinDays is the minimum permitted subscription duration.
	SubscriptionMinDays = 7
)

var (
	ErrSubscriptionDuration = errors.New("subscription duration must be a whole number of at least 7 days")
	ErrNegativeQuantity     = errors.New("quantities must not be negative")
	ErrUnknownPackage       = errors.New("installation package not found or inactive")
	ErrInvalidMethod        = errors.New("installation method is not offered by this package")
	ErrUnknownService       = errors.New("one or more selected services are not offered by this package")
	ErrCustomQuoteRequired  = errors.New("this package requires a custom quotation after assessment")
	ErrTransportManualQuote = errors.New("this travel distance requires a manual transport quotation")
	ErrUnknownLocation      = errors.New("the installation location is not recognised")
)

// SubscriptionAmount computes the authoritative subscription total for a whole
// number of days. Zero, negative and sub-minimum durations are rejected.
func SubscriptionAmount(days int) (int64, error) {
	if days < SubscriptionMinDays {
		return 0, ErrSubscriptionDuration
	}
	return int64(days) * SubscriptionDailyRateUGX, nil
}

// ServiceOption is one selectable additional service on a package.
type ServiceOption struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	PriceUGX int64  `json:"price_ugx"`
}

// QuoteInput is the client-supplied installation configuration. It contains no
// total — the client cannot influence the authoritative price.
type QuoteInput struct {
	Users        int      `json:"users"`
	Devices      int      `json:"devices"`
	Locations    int      `json:"locations"`
	Method       string   `json:"method"`
	Location     string   `json:"location"`
	TransportKm  int      `json:"transport_km"`
	ServiceCodes []string `json:"service_codes"`
}

// TransportBand is one per-round-trip travel band: distances up to MaxKm cost
// AmountUGX. Bands are ordered ascending by MaxKm.
type TransportBand struct {
	MaxKm     int   `json:"max_km"`
	AmountUGX int64 `json:"amount_ugx"`
}

// QuoteLine is one itemised quotation line.
type QuoteLine struct {
	Label     string `json:"label"`
	AmountUGX int64  `json:"amount_ugx"`
}

// Quote is an itemised, authoritative quotation.
type Quote struct {
	Currency      string      `json:"currency"`
	Lines         []QuoteLine `json:"lines"`
	TotalUGX      int64       `json:"total_ugx"`
	EstimatedDays int         `json:"estimated_days"`
}

// PackageServices decodes the package's additional-service catalogue.
func PackageServices(pkg commercialModels.InstallationPackage) []ServiceOption {
	if strings.TrimSpace(pkg.ServicesJSON) == "" {
		return nil
	}
	var options []ServiceOption
	if err := json.Unmarshal([]byte(pkg.ServicesJSON), &options); err != nil {
		return nil
	}
	return options
}

// PackageMethods returns the installation methods offered by a package.
func PackageMethods(pkg commercialModels.InstallationPackage) []string {
	if strings.TrimSpace(pkg.Methods) == "" {
		return nil
	}
	parts := strings.Split(pkg.Methods, ",")
	methods := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			methods = append(methods, trimmed)
		}
	}
	return methods
}

// PackageTransportBands decodes a package's fixed transport bands (if any).
func PackageTransportBands(pkg commercialModels.InstallationPackage) []TransportBand {
	if strings.TrimSpace(pkg.TransportBandsJSON) == "" {
		return nil
	}
	var bands []TransportBand
	if err := json.Unmarshal([]byte(pkg.TransportBandsJSON), &bands); err != nil {
		return nil
	}
	return bands
}

// transportAmount returns the authoritative transport charge for a round-trip
// distance. Fixed bands take precedence when configured; a distance beyond the
// last band requires a manual quotation.
func transportAmount(pkg commercialModels.InstallationPackage, km int) (int64, bool, error) {
	if km <= 0 {
		return 0, false, nil
	}
	bands := PackageTransportBands(pkg)
	if len(bands) > 0 {
		for _, band := range bands {
			if km <= band.MaxKm {
				return band.AmountUGX, band.AmountUGX > 0, nil
			}
		}
		return 0, false, ErrTransportManualQuote
	}
	transport := pkg.TransportBaseUGX + int64(km)*pkg.TransportPerKmUGX
	return transport, transport > 0, nil
}

// QuoteInstallation computes the authoritative installation quotation for a
// package given the client's configuration. Included allowances are never
// charged for a second time, and a custom package is never auto-priced.
func QuoteInstallation(pkg commercialModels.InstallationPackage, in QuoteInput) (Quote, error) {
	if pkg.Custom {
		return Quote{}, ErrCustomQuoteRequired
	}
	if in.Users < 0 || in.Devices < 0 || in.Locations < 0 || in.TransportKm < 0 {
		return Quote{}, ErrNegativeQuantity
	}

	lines := []QuoteLine{{Label: pkg.Name + " base package", AmountUGX: pkg.BasePriceUGX}}

	extraUsers := in.Users - pkg.IncludedUsers
	if extraUsers > 0 {
		lines = append(lines, QuoteLine{
			Label:     "Additional users",
			AmountUGX: int64(extraUsers) * pkg.UserPriceUGX,
		})
	}

	extraDevices := in.Devices - pkg.IncludedDevices
	if extraDevices > 0 {
		lines = append(lines, QuoteLine{
			Label:     "Additional client computers",
			AmountUGX: int64(extraDevices) * pkg.DevicePriceUGX,
		})
	}

	extraLocations := in.Locations - pkg.IncludedLocations
	if extraLocations > 0 {
		lines = append(lines, QuoteLine{
			Label:     "Additional locations",
			AmountUGX: int64(extraLocations) * pkg.LocationPriceUGX,
		})
	}

	method := strings.TrimSpace(in.Method)
	if method == "onsite" {
		location := strings.TrimSpace(in.Location)
		if location != "" {
			km, ok := LocationDistanceKm(location)
			if !ok {
				return Quote{}, ErrUnknownLocation
			}
			// Kampala city centre is local — no transport charge.
			if km > 0 {
				fee := pkg.TransportStandardUGX
				if fee <= 0 {
					// Fall back to any configured bands / base+per-km rate.
					fallback, _, err := transportAmount(pkg, km)
					if err != nil {
						return Quote{}, err
					}
					fee = fallback
				}
				if fee > 0 {
					lines = append(lines, QuoteLine{
						Label:     fmt.Sprintf("Transport (Kampala City Centre → %s, ~%d km)", location, km),
						AmountUGX: fee,
					})
				}
			}
		} else {
			amount, charge, err := transportAmount(pkg, in.TransportKm)
			if err != nil {
				return Quote{}, err
			}
			if charge {
				lines = append(lines, QuoteLine{Label: "Transport / travel", AmountUGX: amount})
			}
		}
	}

	if len(in.ServiceCodes) > 0 {
		available := map[string]ServiceOption{}
		for _, option := range PackageServices(pkg) {
			available[option.Code] = option
		}
		for _, code := range in.ServiceCodes {
			option, ok := available[code]
			if !ok {
				return Quote{}, ErrUnknownService
			}
			lines = append(lines, QuoteLine{Label: option.Name, AmountUGX: option.PriceUGX})
		}
	}

	var total int64
	for _, line := range lines {
		total += line.AmountUGX
	}

	return Quote{
		Currency:      "UGX",
		Lines:         lines,
		TotalUGX:      total,
		EstimatedDays: pkg.EstimatedDays,
	}, nil
}

// ValidateMethod ensures the requested method (if any) is offered by the
// package. An empty method is accepted (defaults to the first offered).
func ValidateMethod(pkg commercialModels.InstallationPackage, method string) error {
	if strings.TrimSpace(method) == "" {
		return nil
	}
	for _, offered := range PackageMethods(pkg) {
		if strings.EqualFold(offered, method) {
			return nil
		}
	}
	return ErrInvalidMethod
}
