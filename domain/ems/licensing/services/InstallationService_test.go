package services

import (
	"testing"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/models"
)

func TestValidInstallationIdentity(t *testing.T) {
	cases := []struct {
		name     string
		identity string
		want     bool
	}{
		{"uuid_like", "3f1d0a3a-9a4b-4c1e-9d2f-5b7c8e1a2b3c", true},
		{"short", "abc", false},
		{"empty", "", false},
		{"whitespace", "3f1d0a3a 9a4b", false},
		{"too_long", "x" + string(make([]byte, 128)), false},
	}
	for _, testCase := range cases {
		if got := validInstallationIdentity(testCase.identity); got != testCase.want {
			t.Errorf("%s: validInstallationIdentity = %v, want %v", testCase.name, got, testCase.want)
		}
	}
}

func TestInstallationCountsTowardCapacity(t *testing.T) {
	active := models.Installation{Status: models.InstallationStatusActive}
	if !installationCountsTowardCapacity(active) {
		t.Error("ACTIVE installation must consume capacity")
	}

	clientDeactivated := models.Installation{Status: models.InstallationStatusActive}
	now := time.Now()
	clientDeactivated.ClientDeactivatedAt = &now
	if !installationCountsTowardCapacity(clientDeactivated) {
		t.Error("client-deactivated installation must still consume capacity")
	}

	if installationCountsTowardCapacity(models.Installation{Status: models.InstallationStatusReleased}) {
		t.Error("RELEASED installation must not consume capacity")
	}
	if installationCountsTowardCapacity(models.Installation{Status: models.InstallationStatusRevoked}) {
		t.Error("REVOKED installation must not consume capacity")
	}
}

func TestActivationBlockedError(t *testing.T) {
	now := time.Date(2026, 3, 1, 12, 0, 0, 0, time.UTC)
	base := models.Licence{
		ValidFrom: now.Add(-24 * time.Hour),
		ValidTo:   now.Add(24 * time.Hour),
		GraceDays: 14,
		Status:    models.StatusActive,
	}

	if err := activationBlockedError(base, now); err != nil {
		t.Fatalf("ACTIVE licence must allow activation: %v", err)
	}

	inGrace := base
	inGrace.ValidTo = now.Add(-time.Hour)
	if err := activationBlockedError(inGrace, now); err != nil {
		t.Fatalf("IN_GRACE licence must allow activation: %v", err)
	}

	expired := base
	expired.ValidTo = now.Add(-20 * 24 * time.Hour)
	if err := activationBlockedError(expired, now); err != ErrLicenceExpired {
		t.Fatalf("expired licence: err = %v, want ErrLicenceExpired", err)
	}

	beforeWindow := base
	beforeWindow.ValidFrom = now.Add(24 * time.Hour)
	beforeWindow.ValidTo = now.Add(48 * time.Hour)
	if err := activationBlockedError(beforeWindow, now); err != ErrLicenceNotYetValid {
		t.Fatalf("future licence: err = %v, want ErrLicenceNotYetValid", err)
	}

	suspended := base
	suspended.Status = models.StatusSuspended
	if err := activationBlockedError(suspended, now); err != ErrLicenceNotActive {
		t.Fatalf("suspended licence: err = %v, want ErrLicenceNotActive", err)
	}

	revoked := base
	revoked.Status = models.StatusRevoked
	if err := activationBlockedError(revoked, now); err != ErrLicenceNotActive {
		t.Fatalf("revoked licence: err = %v, want ErrLicenceNotActive", err)
	}
}
