package models

import "testing"

func TestPaymentStatusTransitions(t *testing.T) {
	allowed := [][2]string{
		{StatusInitiated, StatusPending},
		{StatusInitiated, StatusUnknown},
		{StatusPending, StatusUnknown},
		{StatusPending, StatusSuccessful},
		{StatusPending, StatusFailed},
		{StatusUnknown, StatusSuccessful},
		{StatusUnknown, StatusFailed},
		{StatusUnknown, StatusExpired},
	}
	for _, pair := range allowed {
		if !CanTransition(pair[0], pair[1]) {
			t.Fatalf("expected transition %s -> %s", pair[0], pair[1])
		}
	}
	rejected := [][2]string{
		{StatusSuccessful, StatusFailed},
		{StatusSuccessful, StatusPending},
		{StatusSuccessful, StatusUnknown},
		{StatusFailed, StatusSuccessful},
		{StatusCancelled, StatusSuccessful},
		{StatusExpired, StatusUnknown},
	}
	for _, pair := range rejected {
		if CanTransition(pair[0], pair[1]) {
			t.Fatalf("unexpected transition %s -> %s", pair[0], pair[1])
		}
	}
}

func TestUnknownIsUnresolved(t *testing.T) {
	if !IsUnresolved(StatusUnknown) {
		t.Fatal("UNKNOWN must remain unresolved")
	}
	if IsFinal(StatusUnknown) {
		t.Fatal("UNKNOWN must not be final")
	}
}
