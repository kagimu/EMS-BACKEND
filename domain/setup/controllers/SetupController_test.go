package controllers

import "testing"

// TestNextInstallationAction covers the single-next-step rule for every
// installation state (audit states A–F).
func TestNextInstallationAction(t *testing.T) {
	cases := []struct {
		name            string
		users, admin    bool
		owner, org, sub bool
		want            string
	}{
		{"fresh install", false, false, false, false, false, NextActionCommercialSetup},
		{"user exists, admin not active, no owner", true, false, false, false, false, NextActionResumeActivation},
		{"admin active, no owner", true, true, false, false, false, NextActionCompleteSetup},
		{"owner exists, no organisation", true, true, true, false, false, NextActionCompleteSetup},
		{"owner + organisation, no subscription", true, true, true, true, false, NextActionActivateSubscription},
		{"fully ready", true, true, true, true, true, NextActionReady},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := nextInstallationAction(tc.users, tc.admin, tc.owner, tc.org, tc.sub)
			if got != tc.want {
				t.Fatalf("nextInstallationAction(%v,%v,%v,%v,%v) = %q, want %q",
					tc.users, tc.admin, tc.owner, tc.org, tc.sub, got, tc.want)
			}
		})
	}
}

// TestLicenceContainsAddon verifies the addon-entitlement read reuses the
// existing licensed_addons representation and never guesses.
func TestLicenceContainsAddon(t *testing.T) {
	cases := []struct {
		json string
		code string
		want bool
	}{
		{`["education"]`, "education", true},
		{`["Education"]`, "education", true},
		{`["finance","education"]`, "education", true},
		{`["finance"]`, "education", false},
		{`[]`, "education", false},
		{``, "education", false},
		{`not-json`, "education", false},
	}
	for _, tc := range cases {
		if got := licenceContainsAddon(tc.json, tc.code); got != tc.want {
			t.Fatalf("licenceContainsAddon(%q,%q) = %v, want %v", tc.json, tc.code, got, tc.want)
		}
	}
}
