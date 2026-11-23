package integration

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	commercialModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commercial/models"
	paymentModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/payments/models"
	paymentProviders "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/payments/providers"
	paymentServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/payments/services"
	uuid "github.com/satori/go.uuid"
)

func TestStage8ReconciliationRecoversMissingCallback(t *testing.T) {
	fx := setupStage0BFixture(t)
	t.Setenv("PAYMENTS_RECONCILIATION_MAX_ATTEMPTS", "3")

	fake := &fakeMobileMoneyProvider{status: "PENDING", providerStatus: "PENDING"}
	restore := paymentProviders.SetForTest(paymentProviders.MTNMoMo, fake)
	t.Cleanup(restore)

	prefix := fmt.Sprintf("S8-%d", time.Now().UnixNano()%10000000000)
	order := commercialModels.SetupOrder{
		Number:        prefix + "-setup",
		Status:        "AWAITING_PAYMENT",
		Currency:      "UGX",
		SetupTotalUGX: 7500,
	}
	if err := fx.db.Create(&order).Error; err != nil {
		t.Fatalf("seed order: %v", err)
	}
	providerReference := uuid.NewV4().String()
	lastChecked := time.Now().UTC().Add(-10 * time.Minute)
	expires := time.Now().UTC().Add(20 * time.Minute)
	payment := paymentModels.MobileMoneyPayment{
		Number:                prefix + "-payment",
		Provider:              paymentProviders.MTNMoMo,
		Purpose:               paymentModels.PurposeSetup,
		Reference:             order.Number,
		CommercialOrderNumber: order.Number,
		PayerPhone:            "256771234567",
		AmountUGX:             order.SetupTotalUGX,
		Currency:              "UGX",
		Status:                paymentModels.StatusPending,
		StatusVersion:         1,
		ProviderReference:     providerReference,
		LastCheckedAt:         &lastChecked,
		ExpiresAt:             &expires,
		FulfilmentStatus:      paymentModels.FulfilmentNotRequired,
	}
	if err := fx.db.Create(&payment).Error; err != nil {
		t.Fatalf("seed payment: %v", err)
	}
	t.Cleanup(func() {
		fx.db.Unscoped().Where("transaction_number = ?", payment.Number).Delete(&paymentModels.PaymentReconciliation{})
		fx.db.Unscoped().Where("transaction_number = ?", payment.Number).Delete(&paymentModels.PaymentEvent{})
		fx.db.Unscoped().Where("number = ?", payment.Number).Delete(&paymentModels.MobileMoneyPayment{})
		fx.db.Unscoped().Where("number = ?", order.Number).Delete(&commercialModels.SetupOrder{})
	})

	fake.setStatus("SUCCESSFUL")
	if _, err := paymentServices.ReconcileOnce(context.Background()); err != nil {
		t.Fatalf("reconcile: %v", err)
	}

	var refreshed paymentModels.MobileMoneyPayment
	if err := fx.db.Where("number = ?", payment.Number).First(&refreshed).Error; err != nil {
		t.Fatalf("load payment: %v", err)
	}
	if refreshed.Status != paymentModels.StatusSuccessful {
		t.Fatalf("reconciliation did not recover payment: %s", refreshed.Status)
	}
	var reconciliations []paymentModels.PaymentReconciliation
	fx.db.Where("transaction_number = ?", payment.Number).Find(&reconciliations)
	if len(reconciliations) == 0 {
		t.Fatal("reconciliation record was not written")
	}
	if reconciliations[0].Result != paymentModels.ReconciliationMatched {
		t.Fatalf("unexpected reconciliation result: %s", reconciliations[0].Result)
	}
}

func TestStage8AmbiguousTimeoutBecomesUnknownAndBlocksSecondCharge(t *testing.T) {
	fx := setupStage0BFixture(t)
	fake := &fakeMobileMoneyProvider{
		status: "PENDING", providerStatus: "PENDING",
		initiateErr: &paymentProviders.ProviderError{Operation: "test", Ambiguous: true, Err: errors.New("timeout")},
	}
	restore := paymentProviders.SetForTest(paymentProviders.MTNMoMo, fake)
	t.Cleanup(restore)

	prefix := fmt.Sprintf("S8U-%d", time.Now().UnixNano()%10000000000)
	order := commercialModels.SetupOrder{
		Number: prefix + "-setup", Status: "AWAITING_PAYMENT", Currency: "UGX", SetupTotalUGX: 3000,
	}
	if err := fx.db.Create(&order).Error; err != nil {
		t.Fatalf("seed order: %v", err)
	}
	t.Cleanup(func() {
		fx.db.Unscoped().Where("reference = ? AND purpose = ?", order.Number, paymentModels.PurposeSetup).Delete(&paymentModels.MobileMoneyPayment{})
		fx.db.Unscoped().Where("number = ?", order.Number).Delete(&commercialModels.SetupOrder{})
	})

	payment, err := paymentServices.Initiate(paymentServices.InitiateInput{
		Purpose: paymentModels.PurposeSetup, Reference: order.Number,
		Provider: paymentProviders.MTNMoMo, Phone: "0771234567", IdempotencyKey: "s8-unknown-1",
	})
	if err != nil {
		t.Fatalf("initiate ambiguous payment: %v", err)
	}
	if payment.Status != paymentModels.StatusUnknown {
		t.Fatalf("ambiguous timeout became %s, want UNKNOWN", payment.Status)
	}
	if _, err := paymentServices.Initiate(paymentServices.InitiateInput{
		Purpose: paymentModels.PurposeSetup, Reference: order.Number,
		Provider: paymentProviders.MTNMoMo, Phone: "0779999999", IdempotencyKey: "s8-unknown-2",
	}); !errors.Is(err, paymentServices.ErrPaymentInProgress) {
		t.Fatalf("second charge was not blocked: %v", err)
	}
	var count int64
	fx.db.Model(&paymentModels.MobileMoneyPayment{}).Where("reference = ?", order.Number).Count(&count)
	if count != 1 {
		t.Fatalf("expected one payment attempt, got %d", count)
	}
}

func TestStage8AmountMismatchDoesNotMarkPaid(t *testing.T) {
	fx := setupStage0BFixture(t)
	fake := &fakeMobileMoneyProvider{
		status: "SUCCESSFUL", providerStatus: "SUCCESSFUL", providerAmount: 1,
	}
	restore := paymentProviders.SetForTest(paymentProviders.MTNMoMo, fake)
	t.Cleanup(restore)

	prefix := fmt.Sprintf("S8M-%d", time.Now().UnixNano()%10000000000)
	order := commercialModels.SetupOrder{
		Number: prefix + "-setup", Status: "AWAITING_PAYMENT", Currency: "UGX", SetupTotalUGX: 5000,
	}
	if err := fx.db.Create(&order).Error; err != nil {
		t.Fatalf("seed order: %v", err)
	}
	expires := time.Now().UTC().Add(20 * time.Minute)
	payment := paymentModels.MobileMoneyPayment{
		Number: prefix + "-payment", Provider: paymentProviders.MTNMoMo,
		Purpose: paymentModels.PurposeSetup, Reference: order.Number,
		PayerPhone: "256771234567", AmountUGX: order.SetupTotalUGX, Currency: "UGX",
		Status: paymentModels.StatusPending, StatusVersion: 1,
		ProviderReference: uuid.NewV4().String(), ExpiresAt: &expires,
		FulfilmentStatus: paymentModels.FulfilmentNotRequired,
	}
	if err := fx.db.Create(&payment).Error; err != nil {
		t.Fatalf("seed payment: %v", err)
	}
	t.Cleanup(func() {
		fx.db.Unscoped().Where("transaction_number = ?", payment.Number).Delete(&paymentModels.PaymentEvent{})
		fx.db.Unscoped().Where("number = ?", payment.Number).Delete(&paymentModels.MobileMoneyPayment{})
		fx.db.Unscoped().Where("number = ?", order.Number).Delete(&commercialModels.SetupOrder{})
	})

	refreshed, err := paymentServices.Refresh(payment.Number)
	if err != nil {
		t.Fatalf("refresh: %v", err)
	}
	if refreshed.Status != paymentModels.StatusUnknown {
		t.Fatalf("amount mismatch became %s, want UNKNOWN", refreshed.Status)
	}
	if refreshed.PaidAt != nil {
		t.Fatal("amount mismatch must not set PaidAt")
	}
}

func TestStage8ConcurrentInitiationCreatesOnePayment(t *testing.T) {
	fx := setupStage0BFixture(t)
	fake := &fakeMobileMoneyProvider{
		status: "PENDING", providerStatus: "PENDING", initiateDelay: 150 * time.Millisecond,
	}
	restore := paymentProviders.SetForTest(paymentProviders.MTNMoMo, fake)
	t.Cleanup(restore)

	prefix := fmt.Sprintf("S8C-%d", time.Now().UnixNano()%10000000000)
	order := commercialModels.SetupOrder{
		Number: prefix + "-setup", Status: "AWAITING_PAYMENT", Currency: "UGX", SetupTotalUGX: 2500,
	}
	if err := fx.db.Create(&order).Error; err != nil {
		t.Fatalf("seed order: %v", err)
	}
	t.Cleanup(func() {
		fx.db.Unscoped().Where("reference = ? AND purpose = ?", order.Number, paymentModels.PurposeSetup).Delete(&paymentModels.MobileMoneyPayment{})
		fx.db.Unscoped().Where("number = ?", order.Number).Delete(&commercialModels.SetupOrder{})
	})

	numbers := make([]string, 2)
	errs := make([]error, 2)
	var wg sync.WaitGroup
	for index := 0; index < 2; index++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			payment, err := paymentServices.Initiate(paymentServices.InitiateInput{
				Purpose: paymentModels.PurposeSetup, Reference: order.Number,
				Provider: paymentProviders.MTNMoMo, Phone: "0771234567",
				IdempotencyKey: fmt.Sprintf("s8-concurrent-%d", index),
			})
			if err != nil {
				errs[index] = err
				return
			}
			numbers[index] = payment.Number
		}(index)
	}
	wg.Wait()
	for index, err := range errs {
		if err != nil {
			t.Fatalf("concurrent initiation %d: %v", index, err)
		}
	}
	if numbers[0] == "" || numbers[0] != numbers[1] {
		t.Fatalf("concurrent initiation returned different payments: %v", numbers)
	}
	var count int64
	fx.db.Model(&paymentModels.MobileMoneyPayment{}).Where("reference = ?", order.Number).Count(&count)
	if count != 1 {
		t.Fatalf("expected one payment row, got %d", count)
	}
}
