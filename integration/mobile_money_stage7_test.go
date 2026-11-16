package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/routes"
	commercialModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commercial/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	paymentModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/payments/models"
	paymentProviders "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/payments/providers"
)

type fakeMobileMoneyProvider struct {
	mu                    sync.Mutex
	status                string
	providerStatus        string
	providerTransactionID string
	providerAmount        int64
	initiateDelay         time.Duration
	initiateErr           error
	statusErr             error
}

func (f *fakeMobileMoneyProvider) Name() string        { return paymentProviders.MTNMoMo }
func (f *fakeMobileMoneyProvider) DisplayName() string { return "Test MTN MoMo" }
func (f *fakeMobileMoneyProvider) Enabled() bool       { return true }
func (f *fakeMobileMoneyProvider) Environment() string { return "test" }

func (f *fakeMobileMoneyProvider) Initiate(_ context.Context, request paymentProviders.Request) (paymentProviders.Result, error) {
	f.mu.Lock()
	delay := f.initiateDelay
	initiateErr := f.initiateErr
	f.mu.Unlock()
	if delay > 0 {
		time.Sleep(delay)
	}
	if initiateErr != nil {
		return paymentProviders.Result{}, initiateErr
	}
	return paymentProviders.Result{
		Status:            "PENDING",
		ProviderStatus:    "PENDING",
		ProviderReference: request.PaymentNumber,
	}, nil
}

func (f *fakeMobileMoneyProvider) Status(_ context.Context, providerReference string) (paymentProviders.Result, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.statusErr != nil {
		return paymentProviders.Result{}, f.statusErr
	}
	return paymentProviders.Result{
		Status:                f.status,
		ProviderStatus:        f.providerStatus,
		ProviderReference:     providerReference,
		ProviderTransactionID: f.providerTransactionID,
		ProviderAmount:        f.providerAmount,
	}, nil
}

func (f *fakeMobileMoneyProvider) VerifyCallback(ctx context.Context, providerReference string) (paymentProviders.Result, error) {
	return f.Status(ctx, providerReference)
}

func (f *fakeMobileMoneyProvider) CallbackReference(body []byte) (string, string, error) {
	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", "", err
	}
	reference, _ := payload["referenceId"].(string)
	return reference, "", nil
}

func (f *fakeMobileMoneyProvider) setStatus(status string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.status = status
	f.providerStatus = status
	if status == "SUCCESSFUL" {
		f.providerTransactionID = "fake-financial-transaction"
	}
}

func TestStage7MobileMoneyPaymentAndOwnerTracking(t *testing.T) {
	fx := setupStage0BFixture(t)
	t.Setenv(commons.SoftwareOwnerUserNumbersEnv, s0bAdminA)
	t.Setenv("PAYMENTS_MAX_AMOUNT_UGX", "1000000")

	fake := &fakeMobileMoneyProvider{status: "PENDING", providerStatus: "PENDING"}
	restore := paymentProviders.SetForTest(paymentProviders.MTNMoMo, fake)
	t.Cleanup(restore)

	prefix := fmt.Sprintf("S7-%d", time.Now().UnixNano()%10000000000)
	order := commercialModels.SetupOrder{
		Number:        prefix + "-setup",
		Status:        "AWAITING_PAYMENT",
		Currency:      "UGX",
		SetupTotalUGX: 5000,
		ContactPhone:  "0771234567",
	}
	if err := fx.db.Create(&order).Error; err != nil {
		t.Fatalf("seed setup order: %v", err)
	}
	t.Cleanup(func() {
		fx.db.Unscoped().Where("setup_order_number = ?", order.Number).Delete(&commercialModels.ActivationToken{})
		fx.db.Unscoped().Where("number = ?", order.Number).Delete(&commercialModels.SetupOrder{})
		fx.db.Unscoped().Where("reference = ? AND purpose = ?", order.Number, paymentModels.PurposeSetup).Delete(&paymentModels.MobileMoneyPayment{})
	})

	server := httptest.NewServer(routes.SetupRoutes())
	t.Cleanup(server.Close)
	client := server.Client()

	status, body := stage1Request(t, client, server.URL, http.MethodPost,
		"/setup/orders/"+order.Number+"/pay", "", map[string]string{
			"provider": paymentProviders.MTNMoMo,
			"phone":    "0771234567",
		})
	if status != http.StatusCreated {
		t.Fatalf("initiate payment = %d (%s)", status, string(body))
	}
	var initiated struct {
		Data struct {
			Number         string `json:"number"`
			Status         string `json:"status"`
			PayerPhoneMask string `json:"payer_phone_masked"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &initiated); err != nil {
		t.Fatalf("decode initiated payment: %v", err)
	}
	if initiated.Data.Number == "" || initiated.Data.Status != "PENDING" {
		t.Fatalf("unexpected initiated payment: %s", string(body))
	}
	if initiated.Data.PayerPhoneMask == "256771234567" || initiated.Data.PayerPhoneMask == "" {
		t.Fatalf("payer phone was not masked: %s", string(body))
	}

	fake.setStatus("SUCCESSFUL")
	status, body = stage1Request(t, client, server.URL, http.MethodGet,
		"/setup/payments/"+initiated.Data.Number, "", nil)
	if status != http.StatusOK {
		t.Fatalf("payment status = %d (%s)", status, string(body))
	}
	var refreshed struct {
		Data struct {
			Status string `json:"status"`
		} `json:"data"`
	}
	_ = json.Unmarshal(body, &refreshed)
	if refreshed.Data.Status != "SUCCESSFUL" {
		t.Fatalf("payment was not confirmed: %s", string(body))
	}

	status, body = stage1Request(t, client, server.URL, http.MethodPost,
		"/setup/payments/"+initiated.Data.Number+"/claim", "", map[string]string{})
	if status != http.StatusOK {
		t.Fatalf("claim payment = %d (%s)", status, string(body))
	}
	var claimed struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	_ = json.Unmarshal(body, &claimed)
	if claimed.Data.Token == "" {
		t.Fatalf("claim did not return a setup token: %s", string(body))
	}
	status, body = stage1Request(t, client, server.URL, http.MethodPost,
		"/setup/payments/"+initiated.Data.Number+"/claim", "", map[string]string{})
	if status != http.StatusConflict {
		t.Fatalf("duplicate fulfilment claim = %d (%s), want 409", status, string(body))
	}

	var fulfilled paymentModels.MobileMoneyPayment
	if err := fx.db.Where("number = ?", initiated.Data.Number).First(&fulfilled).Error; err != nil {
		t.Fatalf("load fulfilled payment: %v", err)
	}
	if fulfilled.FulfilmentStatus != paymentModels.FulfilmentCompleted {
		t.Fatalf("fulfilment status = %q, want COMPLETED", fulfilled.FulfilmentStatus)
	}

	status, body = stage1Request(t, client, server.URL, http.MethodGet,
		"/main/ems/owner/payments", fx.tokenAdminA, nil)
	if status != http.StatusOK {
		t.Fatalf("owner payments = %d (%s)", status, string(body))
	}
	if status, _ := stage1Request(t, client, server.URL, http.MethodGet,
		"/main/ems/owner/payments", fx.tokenAdminNo, nil); status != http.StatusForbidden {
		t.Fatalf("non-owner owner-payments access = %d, want 403", status)
	}
	var listed struct {
		Data struct {
			Items   []map[string]interface{} `json:"items"`
			Summary struct {
				SuccessfulCount int64 `json:"successful_count"`
			} `json:"summary"`
		} `json:"data"`
	}
	_ = json.Unmarshal(body, &listed)
	if len(listed.Data.Items) == 0 || listed.Data.Summary.SuccessfulCount == 0 {
		t.Fatalf("owner payment list did not include the successful payment: %s", string(body))
	}

	status, body = stage1Request(t, client, server.URL, http.MethodPost,
		"/main/ems/owner/payments/"+initiated.Data.Number+"/verify", fx.tokenAdminA, map[string]string{})
	if status != http.StatusOK {
		t.Fatalf("owner verify = %d (%s)", status, string(body))
	}
	var verified struct {
		Data struct {
			Status   string `json:"status"`
			Verified string `json:"verified_at"`
		} `json:"data"`
	}
	_ = json.Unmarshal(body, &verified)
	if verified.Data.Status != "SUCCESSFUL" || verified.Data.Verified == "" {
		t.Fatalf("owner verification did not record the check: %s", string(body))
	}

	var stored paymentModels.MobileMoneyPayment
	if err := fx.db.Where("number = ?", initiated.Data.Number).First(&stored).Error; err != nil {
		t.Fatalf("load payment for callback: %v", err)
	}
	for _, method := range []string{http.MethodPost, http.MethodPut} {
		status, body = stage1Request(t, client, server.URL, method,
			"/payments/callbacks/mtn", "", map[string]string{"referenceId": stored.ProviderReference, "status": "SUCCESSFUL"})
		if status != http.StatusOK {
			t.Fatalf("mtn %s callback = %d (%s)", method, status, string(body))
		}
	}
}
