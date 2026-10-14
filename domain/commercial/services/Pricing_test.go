package services

import (
	"errors"
	"testing"

	commercialModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commercial/models"
)

func TestSubscriptionAmount(t *testing.T) {
	cases := []struct {
		days int
		want int64
		ok   bool
	}{
		{7, 8400, true},
		{14, 16800, true},
		{28, 33600, true},
		{45, 54000, true},
		{6, 0, false},
		{0, 0, false},
		{-1, 0, false},
	}
	for _, tc := range cases {
		got, err := SubscriptionAmount(tc.days)
		if tc.ok {
			if err != nil || got != tc.want {
				t.Fatalf("SubscriptionAmount(%d) = (%d, %v), want (%d, nil)", tc.days, got, err, tc.want)
			}
		} else if !errors.Is(err, ErrSubscriptionDuration) {
			t.Fatalf("SubscriptionAmount(%d) = (%d, %v), want ErrSubscriptionDuration", tc.days, got, err)
		}
	}
}

func samplePackage() commercialModels.InstallationPackage {
	return commercialModels.InstallationPackage{
		Number:            "pkg-1",
		Name:              "Standard",
		BasePriceUGX:      1000,
		IncludedDevices:   2,
		IncludedLocations: 1,
		DevicePriceUGX:    500,
		LocationPriceUGX:  300,
		TransportBaseUGX:  100,
		TransportPerKmUGX: 50,
		Methods:           "onsite,remote",
		ServicesJSON:      `[{"code":"training","name":"Training","price_ugx":2000}]`,
		EstimatedDays:     3,
		Active:            true,
	}
}

func TestQuoteInstallation(t *testing.T) {
	quote, err := QuoteInstallation(samplePackage(), QuoteInput{
		Devices: 4, Locations: 2, Method: "onsite", TransportKm: 10, ServiceCodes: []string{"training"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 1000 base + (4-2)*500 devices + (2-1)*300 locations + (100+10*50) transport + 2000 service
	const want = 1000 + 1000 + 300 + 600 + 2000
	if quote.TotalUGX != want {
		t.Fatalf("total = %d, want %d", quote.TotalUGX, want)
	}
	if quote.EstimatedDays != 3 || quote.Currency != "UGX" {
		t.Fatalf("unexpected quote metadata: %+v", quote)
	}
}

func TestQuoteInstallationRejectsInvalid(t *testing.T) {
	if _, err := QuoteInstallation(samplePackage(), QuoteInput{Devices: -1}); !errors.Is(err, ErrNegativeQuantity) {
		t.Fatalf("negative devices not rejected: %v", err)
	}
	if _, err := QuoteInstallation(samplePackage(), QuoteInput{ServiceCodes: []string{"nope"}}); !errors.Is(err, ErrUnknownService) {
		t.Fatalf("unknown service not rejected: %v", err)
	}
}

func TestValidateMethod(t *testing.T) {
	pkg := samplePackage()
	if err := ValidateMethod(pkg, "remote"); err != nil {
		t.Fatalf("remote should be valid: %v", err)
	}
	if err := ValidateMethod(pkg, ""); err != nil {
		t.Fatalf("empty method should be valid: %v", err)
	}
	if err := ValidateMethod(pkg, "hybrid"); !errors.Is(err, ErrInvalidMethod) {
		t.Fatalf("hybrid should be rejected: %v", err)
	}
}
