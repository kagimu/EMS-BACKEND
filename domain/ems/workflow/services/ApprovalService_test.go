package services

import (
	"strings"
	"testing"
	"unicode/utf8"

	authorityResponses "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/authority/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/workflow/models"
)

// TestSerializeParseAnchorsRoundTrip proves the snapshotted anchor list survives
// storage exactly, including the multi-anchor list Foundation D can return for a
// single level.
func TestSerializeParseAnchorsRoundTrip(t *testing.T) {
	original := []authorityResponses.AuthorityAnchorResponse{
		{Type: "role", Value: "role-1"},
		{Type: "stakeholder_type", Value: "owner"},
	}

	serialized := serializeAnchors(original)
	if serialized != "role:role-1|stakeholder_type:owner" {
		t.Fatalf("serializeAnchors = %q", serialized)
	}

	parsed := parseAnchors(serialized)
	if len(parsed) != 2 || parsed[0].Type != "role" || parsed[0].Value != "role-1" ||
		parsed[1].Type != "stakeholder_type" || parsed[1].Value != "owner" {
		t.Fatalf("parseAnchors = %#v", parsed)
	}
}

func TestSerializeAnchorsSkipsIncompleteFragments(t *testing.T) {
	serialized := serializeAnchors([]authorityResponses.AuthorityAnchorResponse{
		{Type: "role", Value: " "},
		{Type: "", Value: "x"},
		{Type: "role", Value: "ok"},
	})
	if serialized != "role:ok" {
		t.Fatalf("serializeAnchors = %q, want %q", serialized, "role:ok")
	}
}

func TestParseAnchorsIsDefensive(t *testing.T) {
	cases := map[string]int{
		"":                              0,
		"role":                          0,
		"role:":                         0,
		":value":                        0,
		"role:a":                        1,
		"role:a|stakeholder_type:owner": 2,
	}
	for input, want := range cases {
		if got := len(parseAnchors(input)); got != want {
			t.Fatalf("parseAnchors(%q) = %d anchors, want %d", input, got, want)
		}
	}
}

// TestRemoveUserImplementsSelfApprovalProhibition covers BD-E6: the requester is
// removed from the candidate set and no replacement is invented.
func TestRemoveUserImplementsSelfApprovalProhibition(t *testing.T) {
	candidates := []string{"a", "b", "c"}
	if got := removeUser(candidates, "b"); len(got) != 2 || contains(got, "b") {
		t.Fatalf("removeUser = %#v", got)
	}
	if got := removeUser([]string{"a"}, "a"); len(got) != 0 {
		t.Fatalf("a sole candidate who is the requester must leave an empty set, got %#v", got)
	}
}

func TestValidateRequestInput(t *testing.T) {
	amount := func(v float64) *float64 { return &v }

	cases := []struct {
		name           string
		module         string
		resourceType   string
		resourceNumber string
		decisionType   string
		amount         *float64
		wantErr        error
	}{
		{"valid", "finance", "purchase_order", "po-1", "purchase_request", amount(10), nil},
		{"valid without amount", "finance", "purchase_order", "po-1", "purchase_request", nil, nil},
		{"missing module", "", "purchase_order", "po-1", "purchase_request", nil, ErrModuleRequired},
		{"missing resource type", "finance", "", "po-1", "purchase_request", nil, ErrResourceTypeRequired},
		{"missing resource number", "finance", "purchase_order", "", "purchase_request", nil, ErrResourceNumberRequired},
		{"missing decision type", "finance", "purchase_order", "po-1", "", nil, ErrDecisionTypeRequired},
		{"negative amount", "finance", "purchase_order", "po-1", "purchase_request", amount(-1), ErrInvalidAmount},
		{"over long module", string(make([]byte, 51)), "purchase_order", "po-1", "purchase_request", nil, ErrInvalidFieldLength},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			err := validateRequestInput(testCase.module, testCase.resourceType, testCase.resourceNumber,
				testCase.decisionType, testCase.amount)
			if err != testCase.wantErr {
				t.Fatalf("err = %v, want %v", err, testCase.wantErr)
			}
		})
	}
}

// TestNormalizePagination pins the request-list bounds: absent/zero means the
// default page, oversized limits are clamped, and a negative offset is reset.
func TestNormalizePagination(t *testing.T) {
	cases := []struct {
		limit, offset         int
		wantLimit, wantOffset int
	}{
		{0, 0, defaultListLimit, 0},
		{-5, -1, defaultListLimit, 0},
		{10, 20, 10, 20},
		{maxListLimit + 1, 0, maxListLimit, 0},
		{1, 0, 1, 0},
	}
	for _, testCase := range cases {
		limit, offset := normalizePagination(testCase.limit, testCase.offset)
		if limit != testCase.wantLimit || offset != testCase.wantOffset {
			t.Fatalf("normalizePagination(%d, %d) = (%d, %d), want (%d, %d)",
				testCase.limit, testCase.offset, limit, offset, testCase.wantLimit, testCase.wantOffset)
		}
	}
}

// TestRequestStatusVocabulary pins the E1+E2 state model: exactly six request
// states exist and no further workflow state was invented. Foundation E2 added
// RETURNED (ratified D4) and nothing else.
func TestRequestStatusVocabulary(t *testing.T) {
	valid := []string{
		models.RequestStatusPending, models.RequestStatusApproved, models.RequestStatusRejected,
		models.RequestStatusCancelled, models.RequestStatusBlocked, models.RequestStatusReturned,
	}
	for _, status := range valid {
		if !isValidRequestStatus(status) {
			t.Fatalf("%s must be a valid request status", status)
		}
	}
	for _, status := range []string{"", "WAITING", "APPROVED_PARTIALLY", "DELEGATED", "ESCALATED", "EXPIRED"} {
		if isValidRequestStatus(status) {
			t.Fatalf("%q must not be a valid request status in Foundation E1/E2", status)
		}
	}
}

// TestCommentWithinLimit pins the ratified D3 rule: the shared limit is 1,000
// *Unicode characters*, counted as runes rather than Go's byte-counting len().
func TestCommentWithinLimit(t *testing.T) {
	ascii := func(n int) string { return strings.Repeat("a", n) }
	multibyte := func(n int) string { return strings.Repeat("é", n) }

	cases := []struct {
		name    string
		comment string
		want    bool
	}{
		{"empty", "", true},
		{"ascii at the limit", ascii(1000), true},
		{"ascii over the limit", ascii(1001), false},
		// Each 'é' is one rune but two bytes, so a 1,000-rune string is 2,000
		// bytes: it must still be accepted, which len() would refuse.
		{"1000 multibyte runes (2000 bytes)", multibyte(1000), true},
		{"1001 multibyte runes", multibyte(1001), false},
	}
	for _, testCase := range cases {
		if got := commentWithinLimit(testCase.comment); got != testCase.want {
			t.Fatalf("%s: commentWithinLimit = %v, want %v (runes=%d bytes=%d)",
				testCase.name, got, testCase.want,
				utf8.RuneCountInString(testCase.comment), len(testCase.comment))
		}
	}
}

// TestParseAnchorsLegacyAnchorsAreOrganisationWide proves ratified D9/D11: a
// legacy serialised anchor carries no unit, so it is read as organisation-wide
// and a pre-amendment snapshot is therefore never retroactively narrowed.
func TestParseAnchorsLegacyAnchorsAreOrganisationWide(t *testing.T) {
	anchors := parseAnchors("role:role-fin|stakeholder_type:owner")
	if len(anchors) != 2 {
		t.Fatalf("expected 2 anchors, got %d: %+v", len(anchors), anchors)
	}
	for _, item := range anchors {
		if item.Unit != "" {
			t.Fatalf("a legacy anchor must carry no unit (organisation-wide), got %q", item.Unit)
		}
	}
}
