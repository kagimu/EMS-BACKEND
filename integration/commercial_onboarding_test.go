package integration

import (
	"errors"
	"fmt"
	"testing"
	"time"

	commercialModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commercial/models"
	commercialServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commercial/services"
	licensingModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/models"
)

// TestCommercialOnboardingFlow exercises the end-to-end commercial workflow
// (packages → quotation → order → simulated setup payment → setup token →
// subscription → simulated subscription payment → subscription token →
// licence provisioning). It writes only to the dedicated onboarding database
// and skips otherwise (requireOnboardingDB).
func TestCommercialOnboardingFlow(t *testing.T) {
	db := requireOnboardingDB(t)

	suffix := time.Now().UnixNano() % 1000000000
	actor := fmt.Sprintf("commercial-actor-%d", suffix)
	org := fmt.Sprintf("commercial-org-%d", suffix)

	var createdPackages, createdTokens, createdPayments []string
	t.Cleanup(func() {
		db.Unscoped().Where("number IN ?", createdPackages).Delete(&commercialModels.InstallationPackage{})
		db.Unscoped().Where("number IN ?", createdTokens).Delete(&commercialModels.ActivationToken{})
		db.Unscoped().Where("number IN ?", createdPayments).Delete(&commercialModels.SimulatedPayment{})
		db.Unscoped().Where("organisation_number = ?", org).Delete(&commercialModels.SubscriptionOrder{})
		db.Unscoped().Where("organisation_number = ?", org).Delete(&commercialModels.SetupOrder{})
		db.Unscoped().Where("organisation_number = ?", org).Delete(&licensingModels.LicenceEvent{})
		db.Unscoped().Where("organisation_number = ?", org).Delete(&licensingModels.Licence{})
	})

	// --- Package catalogue ---
	pkg, err := commercialServices.CreatePackage(actor, commercialServices.PackageInput{
		Code: "std", Name: fmt.Sprintf("Standard %d", suffix), BasePriceUGX: 1000,
		IncludedDevices: 1, DevicePriceUGX: 500, Methods: "onsite,remote", EstimatedDays: 3, Active: true,
	})
	if err != nil {
		t.Fatalf("create package: %v", err)
	}
	createdPackages = append(createdPackages, pkg.Number)

	packages, err := commercialServices.ListActivePackages()
	if err != nil {
		t.Fatalf("list packages: %v", err)
	}
	found := false
	for _, p := range packages {
		if p.Number == pkg.Number {
			found = true
		}
	}
	if !found {
		t.Fatal("created package not returned by ListActivePackages")
	}

	// --- Installation order + authoritative quotation ---
	order, quote, err := commercialServices.CreateSetupOrder(commercialServices.SetupOrderInput{
		PackageNumber: pkg.Number, ContactName: "Buyer", Devices: 2, Method: "remote",
	})
	if err != nil {
		t.Fatalf("create setup order: %v", err)
	}
	if quote.TotalUGX != 1500 { // 1000 base + 1 extra device * 500
		t.Fatalf("quote total = %d, want 1500", quote.TotalUGX)
	}
	if order.SetupTotalUGX != 1500 || order.Status != "AWAITING_PAYMENT" {
		t.Fatalf("unexpected order: %+v", order)
	}

	// --- Failed payment issues no token ---
	if _, token, err := commercialServices.SimulateSetupPayment(order.Number, "failed"); !errors.Is(err, commercialServices.ErrPaymentNotPaid) || token != "" {
		t.Fatalf("failed payment: err=%v token=%q", err, token)
	}

	// --- Successful payment issues a setup token ---
	paid, setupRaw, err := commercialServices.SimulateSetupPayment(order.Number, "success")
	if err != nil || setupRaw == "" {
		t.Fatalf("successful setup payment: err=%v token=%q", err, setupRaw)
	}
	if paid.Status != "TOKEN_ISSUED" {
		t.Fatalf("order status = %q, want TOKEN_ISSUED", paid.Status)
	}

	// --- Validate setup token ---
	validated, err := commercialServices.ValidateSetupToken(setupRaw)
	if err != nil || validated.Number != order.Number {
		t.Fatalf("validate setup token: err=%v order=%v", err, validated)
	}

	// --- Subscription pricing ---
	if _, err := commercialServices.CreateSubscriptionOrder(6, org); !errors.Is(err, commercialServices.ErrSubscriptionDuration) {
		t.Fatalf("6-day subscription not rejected: %v", err)
	}
	subOrder, err := commercialServices.CreateSubscriptionOrder(7, org)
	if err != nil {
		t.Fatalf("create subscription order: %v", err)
	}
	if subOrder.TotalUGX != 8400 || subOrder.DailyRateUGX != 1200 {
		t.Fatalf("subscription total = %d rate = %d, want 8400/1200", subOrder.TotalUGX, subOrder.DailyRateUGX)
	}

	_, subRaw, err := commercialServices.SimulateSubscriptionPayment(subOrder.Number, "success")
	if err != nil || subRaw == "" {
		t.Fatalf("subscription payment: err=%v token=%q", err, subRaw)
	}

	// --- Token purpose separation ---
	if _, err := commercialServices.ValidateSetupToken(subRaw); !errors.Is(err, commercialServices.ErrTokenWrongPurpose) {
		t.Fatalf("subscription token validated as setup: %v", err)
	}
	if _, _, err := commercialServices.RedeemSubscriptionToken(setupRaw, org, actor); !errors.Is(err, commercialServices.ErrTokenWrongPurpose) {
		t.Fatalf("setup token activated a subscription: %v", err)
	}

	// --- Subscription activation provisions a licence ---
	activated, licence, err := commercialServices.RedeemSubscriptionToken(subRaw, org, actor)
	if err != nil {
		t.Fatalf("redeem subscription token: %v", err)
	}
	if activated.Status != "ACTIVE" || licence == nil || licence.Status != licensingModels.StatusActive {
		t.Fatalf("unexpected activation: order=%v licence=%v", activated, licence)
	}
	var first licensingModels.Licence
	if err := db.Where("number = ?", licence.Number).First(&first).Error; err != nil {
		t.Fatalf("reload licence: %v", err)
	}
	firstValidTo := first.ValidTo

	// --- Replaying the subscription token must not extend expiry again ---
	if _, _, err := commercialServices.RedeemSubscriptionToken(subRaw, org, actor); !errors.Is(err, commercialServices.ErrTokenAlreadyUsed) {
		t.Fatalf("replay not rejected: %v", err)
	}
	var after licensingModels.Licence
	if err := db.Where("number = ?", licence.Number).First(&after).Error; err != nil {
		t.Fatalf("reload licence: %v", err)
	}
	if !after.ValidTo.Equal(firstValidTo) {
		t.Fatalf("replay changed expiry: %v -> %v", firstValidTo, after.ValidTo)
	}

	// --- Setup token redemption (one-shot) ---
	if _, err := commercialServices.RedeemSetupToken(setupRaw, org, actor); err != nil {
		t.Fatalf("redeem setup token: %v", err)
	}
	if _, err := commercialServices.RedeemSetupToken(setupRaw, org, actor); !errors.Is(err, commercialServices.ErrTokenAlreadyUsed) {
		t.Fatalf("setup token replay not rejected: %v", err)
	}

	// --- Tokens are stored hashed, never raw ---
	var stored commercialModels.ActivationToken
	if err := db.Where("token_hash = ?", commercialServices.HashToken(setupRaw)).First(&stored).Error; err != nil {
		t.Fatalf("token not stored by hash: %v", err)
	}
	if stored.TokenHash == setupRaw {
		t.Fatal("raw token stored in plaintext")
	}
}
