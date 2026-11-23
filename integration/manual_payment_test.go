package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/routes"
	commercialModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commercial/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	paymentModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/payments/models"
)

func TestManualPaymentProofApprovalAndClaim(t *testing.T) {
	fx := setupStage0BFixture(t)
	t.Setenv(commons.SoftwareOwnerUserNumbersEnv, s0bAdminA)
	t.Setenv("MANUAL_PAYMENT_ENABLED", "true")
	t.Setenv("MANUAL_PAYMENT_PROOF_DIR", t.TempDir())
	t.Setenv("MANUAL_PAYMENT_MAX_PROOF_MB", "2")

	prefix := fmt.Sprintf("MP-%d", time.Now().UnixNano()%10000000000)
	order := commercialModels.SetupOrder{
		Number: prefix + "-setup", Status: "AWAITING_PAYMENT", Currency: "UGX", SetupTotalUGX: 125000,
		ContactName: "Manual Customer", ContactEmail: "manual@example.com", ContactPhone: "256771234567",
	}
	if err := fx.db.Create(&order).Error; err != nil {
		t.Fatalf("seed order: %v", err)
	}
	var submission paymentModels.ManualPaymentSubmission
	t.Cleanup(func() {
		if submission.Number != "" {
			fx.db.Unscoped().Where("transaction_number = ?", submission.Number).Delete(&paymentModels.PaymentEvent{})
			fx.db.Unscoped().Where("number = ?", submission.Number).Delete(&paymentModels.ManualPaymentSubmission{})
		}
		fx.db.Unscoped().Where("setup_order_number = ?", order.Number).Delete(&commercialModels.ActivationToken{})
		fx.db.Unscoped().Where("number = ?", order.Number).Delete(&commercialModels.SetupOrder{})
	})

	server := httptest.NewServer(routes.SetupRoutes())
	t.Cleanup(server.Close)
	client := server.Client()

	var uploadPayload bytes.Buffer
	writer := multipart.NewWriter(&uploadPayload)
	_ = writer.WriteField("method", paymentModels.ManualMethodBankTransfer)
	_ = writer.WriteField("requester_name", "Manual Customer")
	_ = writer.WriteField("payer_name", "Manual Customer")
	_ = writer.WriteField("payer_phone", "256771234567")
	_ = writer.WriteField("payer_reference", "BANK-REF-1")
	_ = writer.WriteField("note", "Paid at branch")
	part, err := writer.CreateFormFile("proof", "proof.pdf")
	if err != nil {
		t.Fatalf("create upload part: %v", err)
	}
	_, _ = part.Write([]byte("%PDF-1.4\nmanual payment proof\n"))
	_ = writer.Close()

	request, err := http.NewRequest(http.MethodPost, server.URL+"/setup/orders/"+order.Number+"/manual-payment", &uploadPayload)
	if err != nil {
		t.Fatalf("build upload request: %v", err)
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())
	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("upload proof: %v", err)
	}
	uploadBody, _ := io.ReadAll(response.Body)
	response.Body.Close()
	if response.StatusCode != http.StatusCreated {
		t.Fatalf("upload proof = %d (%s)", response.StatusCode, string(uploadBody))
	}
	var uploaded struct {
		Data struct {
			Number    string `json:"number"`
			Status    string `json:"status"`
			Requester string `json:"requester_name"`
		} `json:"data"`
	}
	if err := json.Unmarshal(uploadBody, &uploaded); err != nil {
		t.Fatalf("decode upload: %v", err)
	}
	if uploaded.Data.Number == "" || uploaded.Data.Status != paymentModels.ManualStatusSubmitted || uploaded.Data.Requester != "Manual Customer" {
		t.Fatalf("unexpected upload response: %s", string(uploadBody))
	}
	if err := fx.db.Where("number = ?", uploaded.Data.Number).First(&submission).Error; err != nil {
		t.Fatalf("load submission: %v", err)
	}

	status, responseBody := stage1Request(t, client, server.URL, http.MethodGet,
		"/setup/manual-payments/"+submission.Number, "", nil)
	if status != http.StatusOK || !bytes.Contains(responseBody, []byte(`"SUBMITTED"`)) {
		t.Fatalf("public manual status = %d (%s)", status, string(responseBody))
	}

	status, responseBody = stage1Request(t, client, server.URL, http.MethodGet,
		"/main/ems/owner/manual-payments", fx.tokenAdminA, nil)
	if status != http.StatusOK || !bytes.Contains(responseBody, []byte(submission.Number)) {
		t.Fatalf("owner manual list = %d (%s)", status, string(responseBody))
	}
	status, responseBody = stage1Request(t, client, server.URL, http.MethodGet,
		"/main/ems/owner/manual-payments/"+submission.Number+"/proof", fx.tokenAdminA, nil)
	if status != http.StatusOK || len(responseBody) == 0 {
		t.Fatalf("owner proof download = %d (%s)", status, string(responseBody))
	}

	status, responseBody = stage1Request(t, client, server.URL, http.MethodPost,
		"/main/ems/owner/manual-payments/"+submission.Number+"/approve", fx.tokenAdminA, map[string]string{"note": "Verified against bank statement"})
	if status != http.StatusOK {
		t.Fatalf("approve manual payment = %d (%s)", status, string(responseBody))
	}

	status, responseBody = stage1Request(t, client, server.URL, http.MethodGet,
		"/setup/manual-payments/"+submission.Number, "", nil)
	if status != http.StatusOK || !bytes.Contains(responseBody, []byte(`"APPROVED"`)) {
		t.Fatalf("approved manual status = %d (%s)", status, string(responseBody))
	}

	status, responseBody = stage1Request(t, client, server.URL, http.MethodPost,
		"/setup/manual-payments/"+submission.Number+"/claim", "", map[string]string{})
	if status != http.StatusOK {
		t.Fatalf("claim manual payment = %d (%s)", status, string(responseBody))
	}
	var claimed struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	_ = json.Unmarshal(responseBody, &claimed)
	if claimed.Data.Token == "" {
		t.Fatalf("manual claim did not return a token: %s", string(responseBody))
	}

	var refreshed commercialModels.SetupOrder
	if err := fx.db.Where("number = ?", order.Number).First(&refreshed).Error; err != nil {
		t.Fatalf("load order: %v", err)
	}
	if refreshed.Status != "TOKEN_ISSUED" {
		t.Fatalf("order status = %s, want TOKEN_ISSUED", refreshed.Status)
	}

	status, responseBody = stage1Request(t, client, server.URL, http.MethodGet,
		"/setup/manual-payments/"+submission.Number, "", nil)
	if status != http.StatusOK || !bytes.Contains(responseBody, []byte(`"APPROVED"`)) {
		t.Fatalf("approved manual status = %d (%s)", status, string(responseBody))
	}

	status, responseBody = stage1Request(t, client, server.URL, http.MethodPost,
		"/setup/manual-payments/"+submission.Number+"/claim", "", map[string]string{})
	if status != http.StatusConflict {
		t.Fatalf("duplicate manual claim = %d (%s), want 409", status, string(responseBody))
	}
}
