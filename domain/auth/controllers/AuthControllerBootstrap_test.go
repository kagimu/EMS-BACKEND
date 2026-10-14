package controllers

import "testing"

// TestResolveInitialUserType covers the first-user bootstrap rule (TEST1/TEST5
// logic). A fresh installation (zero users) bootstraps an administrator; every
// other state yields an ordinary user — it never grants Admin merely because no
// active admin happens to exist.
func TestResolveInitialUserType(t *testing.T) {
	cases := []struct {
		name  string
		count int64
		want  string
	}{
		{"fresh installation bootstraps an admin", 0, "Admin"},
		{"second registrant is an ordinary user", 1, "User"},
		{"many existing users stays ordinary", 10, "User"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := resolveInitialUserType(tc.count); got != tc.want {
				t.Fatalf("resolveInitialUserType(%d) = %q, want %q", tc.count, got, tc.want)
			}
		})
	}
}

// TestResolveInitialUserTypeNeverGrantsOnNonZeroCount is an explicit guard that
// a non-empty installation can never promote a registrant.
func TestResolveInitialUserTypeNeverGrantsOnNonZeroCount(t *testing.T) {
	for _, count := range []int64{1, 2, 100} {
		if got := resolveInitialUserType(count); got == "Admin" {
			t.Fatalf("resolveInitialUserType(%d) granted Admin, want User", count)
		}
	}
}
