package services

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/models"
)

func sampleLicence() models.Licence {
	return models.Licence{
		Number:             "lic-1",
		LicenceKey:         "opaque-key",
		OrganisationNumber: "org-1",
		ValidFrom:          time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		ValidTo:            time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
		GraceDays:          models.DefaultGraceDays,
		MaxInstallations:   3,
		LicensedAddons:     `["education"]`,
		Status:             models.StatusActive,
	}
}

func TestEffectiveStatus(t *testing.T) {
	licence := sampleLicence()

	cases := []struct {
		name string
		now  time.Time
		want string
	}{
		{"before_valid_from", time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC), models.EffectiveStatusActive},
		{"inside_valid_window", time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC), models.EffectiveStatusActive},
		{"exactly_at_valid_to", licence.ValidTo, models.EffectiveStatusInGrace},
		{"immediately_after_valid_to", licence.ValidTo.Add(time.Second), models.EffectiveStatusInGrace},
		{"throughout_grace", licence.ValidTo.Add(13 * 24 * time.Hour), models.EffectiveStatusInGrace},
		{"last_second_of_grace", licence.ValidTo.Add(14*24*time.Hour - time.Second), models.EffectiveStatusInGrace},
		{"exactly_at_grace_end", licence.ValidTo.Add(14 * 24 * time.Hour), models.EffectiveStatusExpired},
		{"after_grace", licence.ValidTo.Add(15 * 24 * time.Hour), models.EffectiveStatusExpired},
	}

	for _, testCase := range cases {
		if got := EffectiveStatus(licence, testCase.now); got != testCase.want {
			t.Errorf("%s: EffectiveStatus = %q, want %q", testCase.name, got, testCase.want)
		}
	}

	suspended := licence
	suspended.Status = models.StatusSuspended
	if got := EffectiveStatus(suspended, licence.ValidFrom.Add(time.Hour)); got != models.EffectiveStatusSuspended {
		t.Errorf("suspended: EffectiveStatus = %q, want %q", got, models.EffectiveStatusSuspended)
	}

	revoked := licence
	revoked.Status = models.StatusRevoked
	if got := EffectiveStatus(revoked, licence.ValidFrom.Add(time.Hour)); got != models.EffectiveStatusRevoked {
		t.Errorf("revoked: EffectiveStatus = %q, want %q", got, models.EffectiveStatusRevoked)
	}
}

func TestLicencePermitsAccess(t *testing.T) {
	licence := sampleLicence()

	cases := []struct {
		name string
		now  time.Time
		want bool
	}{
		{"before_valid_from", licence.ValidFrom.Add(-time.Hour), false},
		{"inside_valid_window", licence.ValidFrom.Add(time.Hour), true},
		{"exactly_at_valid_to", licence.ValidTo, true},
		{"inside_grace", licence.ValidTo.Add(13 * 24 * time.Hour), true},
		{"at_grace_end", licence.ValidTo.Add(14 * 24 * time.Hour), false},
	}

	for _, testCase := range cases {
		if got := LicencePermitsAccess(licence, testCase.now); got != testCase.want {
			t.Errorf("%s: LicencePermitsAccess = %v, want %v", testCase.name, got, testCase.want)
		}
	}

	suspended := licence
	suspended.Status = models.StatusSuspended
	if LicencePermitsAccess(suspended, licence.ValidFrom.Add(time.Hour)) {
		t.Error("suspended licence must not permit access")
	}

	revoked := licence
	revoked.Status = models.StatusRevoked
	if LicencePermitsAccess(revoked, licence.ValidFrom.Add(time.Hour)) {
		t.Error("revoked licence must not permit access")
	}
}

func TestNormaliseAddons(t *testing.T) {
	normalised, err := NormaliseAddons([]string{" education ", "education", "core", ""})
	if err != nil {
		t.Fatalf("NormaliseAddons returned error: %v", err)
	}
	if len(normalised) != 2 || normalised[0] != "education" || normalised[1] != "core" {
		t.Fatalf("NormaliseAddons = %v, want [education core]", normalised)
	}

	if _, err := NormaliseAddons([]string{"finance"}); !errors.Is(err, ErrUnknownAddon) {
		t.Fatalf("unknown add-on: err = %v, want ErrUnknownAddon", err)
	}
}

func TestLicenceContainsAddon(t *testing.T) {
	licence := sampleLicence()
	if !LicenceContainsAddon(licence, "education") {
		t.Error("licence must contain education")
	}
	if LicenceContainsAddon(licence, "core") {
		t.Error("licence must not contain an add-on it does not list")
	}
}

func TestGenerateLicenceKey(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 200; i++ {
		key, err := generateLicenceKey()
		if err != nil {
			t.Fatalf("generateLicenceKey returned error: %v", err)
		}
		if key == "" {
			t.Fatal("generated key must not be empty")
		}
		if len(key) != 43 {
			t.Fatalf("generated key length = %d, want 43", len(key))
		}
		if strings.Contains(key, "org-1") {
			t.Fatal("generated key must not be derived from organisation information")
		}
		for _, char := range key {
			if !strings.ContainsRune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_", char) {
				t.Fatalf("generated key contains unexpected character %q", char)
			}
		}
		if seen[key] {
			t.Fatal("generated key collision")
		}
		seen[key] = true
	}
}
