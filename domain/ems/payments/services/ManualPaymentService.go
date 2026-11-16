package services

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	commercialModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commercial/models"
	auditServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/audit/services"
	paymentFulfilment "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/payments/fulfilment"
	paymentModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/payments/models"
	paymentNotifier "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/payments/notifier"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrManualNotFound     = errors.New("manual payment submission not found")
	ErrManualNotPayable   = errors.New("this order is not awaiting payment")
	ErrManualDuplicate    = errors.New("a manual payment proof is already awaiting review for this order")
	ErrManualInvalid      = errors.New("invalid manual payment submission")
	ErrManualProof        = errors.New("a valid proof of payment file is required")
	ErrManualRejectReason = errors.New("a rejection reason is required")
	ErrManualNotApproved  = errors.New("this manual payment has not been approved")
	ErrManualRejected     = errors.New("this manual payment was rejected")
	ErrManualAlready      = errors.New("this manual payment has already been claimed")
)

type ManualPaymentOptions struct {
	Enabled      bool                `json:"enabled"`
	Configured   bool                `json:"configured"`
	Instructions string              `json:"instructions"`
	Methods      []ManualPaymentInfo `json:"methods"`
}

type ManualPaymentInfo struct {
	Method      string `json:"method"`
	Label       string `json:"label"`
	Name        string `json:"name"`
	Number      string `json:"number"`
	BankName    string `json:"bank_name"`
	AccountName string `json:"account_name"`
	Branch      string `json:"branch"`
}

type ManualUploadInput struct {
	OrderNumber        string
	Purpose            string
	OrganisationNumber string
	AllowOnboarding    bool
	Method             string
	PayerName          string
	PayerPhone         string
	PayerReference     string
	CustomerNote       string
	CustomerName       string
	CustomerEmail      string
	CustomerPhone      string
	RequesterName      string
	File               *multipart.FileHeader
}

type ManualPaymentView struct {
	Number             string     `json:"number"`
	OrderNumber        string     `json:"order_number"`
	Purpose            string     `json:"purpose"`
	OrganisationNumber string     `json:"organisation_number"`
	CustomerName       string     `json:"customer_name"`
	CustomerEmail      string     `json:"customer_email"`
	CustomerPhone      string     `json:"customer_phone"`
	RequesterName      string     `json:"requester_name"`
	Method             string     `json:"method"`
	AmountUGX          int64      `json:"amount_ugx"`
	Currency           string     `json:"currency"`
	PayerName          string     `json:"payer_name"`
	PayerPhone         string     `json:"payer_phone"`
	PayerReference     string     `json:"payer_reference"`
	CustomerNote       string     `json:"customer_note"`
	Status             string     `json:"status"`
	ProofName          string     `json:"proof_name"`
	ProofMime          string     `json:"proof_mime"`
	ProofSize          int64      `json:"proof_size"`
	HasProof           bool       `json:"has_proof"`
	SubmittedAt        time.Time  `json:"submitted_at"`
	ReviewedAt         *time.Time `json:"reviewed_at"`
	ReviewedByNumber   string     `json:"reviewed_by_number"`
	ReviewNote         string     `json:"review_note"`
	OrderStatus        string     `json:"order_status"`
	TokenClaimedAt     *time.Time `json:"token_claimed_at"`
}

type ManualPaymentDetail struct {
	ManualPaymentView
	ProofPath             string `json:"-"`
	SuggestedEmailSubject string `json:"suggested_email_subject"`
	SuggestedEmailBody    string `json:"suggested_email_body"`
}

type ManualListFilter struct {
	Status  string
	Search  string
	Page    int
	PerPage int
}

func manualEnabled() bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv("MANUAL_PAYMENT_ENABLED")))
	return value == "" || value == "true" || value == "1" || value == "yes"
}

func manualProofDir() string {
	value := strings.TrimSpace(os.Getenv("MANUAL_PAYMENT_PROOF_DIR"))
	if value == "" {
		value = "./storage/manual-payment-proofs"
	}
	return value
}

func manualMaxProofBytes() int64 {
	value := strings.TrimSpace(os.Getenv("MANUAL_PAYMENT_MAX_PROOF_MB"))
	mb, err := strconv.Atoi(value)
	if err != nil || mb < 1 {
		mb = 10
	}
	return int64(mb) << 20
}

func ManualOptions() ManualPaymentOptions {
	mtnNumber := strings.TrimSpace(os.Getenv("MANUAL_PAYMENT_MTN_NUMBER"))
	bankAccount := strings.TrimSpace(os.Getenv("MANUAL_PAYMENT_BANK_ACCOUNT_NUMBER"))
	options := ManualPaymentOptions{
		Enabled:      manualEnabled(),
		Configured:   mtnNumber != "" || bankAccount != "",
		Instructions: strings.TrimSpace(os.Getenv("MANUAL_PAYMENT_INSTRUCTIONS")),
		Methods: []ManualPaymentInfo{
			{
				Method: paymentModels.ManualMethodMTNMoMo, Label: "MTN Mobile Money",
				Name:   strings.TrimSpace(os.Getenv("MANUAL_PAYMENT_MTN_NAME")),
				Number: mtnNumber,
			},
			{
				Method: paymentModels.ManualMethodBankTransfer, Label: "Bank transfer",
				BankName:    strings.TrimSpace(os.Getenv("MANUAL_PAYMENT_BANK_NAME")),
				AccountName: strings.TrimSpace(os.Getenv("MANUAL_PAYMENT_BANK_ACCOUNT_NAME")),
				Number:      bankAccount,
				Branch:      strings.TrimSpace(os.Getenv("MANUAL_PAYMENT_BANK_BRANCH")),
			},
		},
	}
	return options
}

type manualOrderFacts struct {
	OrderNumber        string
	Purpose            string
	OrganisationNumber string
	AmountUGX          int64
	Currency           string
	CustomerName       string
	CustomerEmail      string
	CustomerPhone      string
	Status             string
}

func manualFacts(tx *gorm.DB, purpose, orderNumber, organisationNumber string, allowOnboarding bool) (manualOrderFacts, error) {
	switch purpose {
	case paymentModels.PurposeSetup:
		var order commercialModels.SetupOrder
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("number = ?", orderNumber).First(&order).Error; err != nil {
			return manualOrderFacts{}, ErrManualNotFound
		}
		if order.Status != "AWAITING_PAYMENT" {
			return manualOrderFacts{}, ErrManualNotPayable
		}
		return manualOrderFacts{
			OrderNumber: order.Number, Purpose: purpose, AmountUGX: order.SetupTotalUGX,
			Currency: firstNonEmpty(order.Currency, "UGX"), CustomerName: order.ContactName,
			CustomerEmail: order.ContactEmail, CustomerPhone: order.ContactPhone, Status: order.Status,
		}, nil
	case paymentModels.PurposeSubscription:
		var order commercialModels.SubscriptionOrder
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("number = ?", orderNumber).First(&order).Error; err != nil {
			return manualOrderFacts{}, ErrManualNotFound
		}
		if order.Status != "AWAITING_PAYMENT" {
			return manualOrderFacts{}, ErrManualNotPayable
		}
		if order.OrganisationNumber != "" && order.OrganisationNumber != organisationNumber {
			return manualOrderFacts{}, ErrManualNotFound
		}
		if order.OrganisationNumber == "" && !allowOnboarding {
			return manualOrderFacts{}, ErrManualNotFound
		}
		return manualOrderFacts{
			OrderNumber: order.Number, Purpose: purpose, OrganisationNumber: order.OrganisationNumber,
			AmountUGX: order.TotalUGX, Currency: firstNonEmpty(order.Currency, "UGX"), Status: order.Status,
		}, nil
	default:
		return manualOrderFacts{}, ErrManualNotFound
	}
}

func proofExtension(contentType string) (string, bool) {
	switch strings.ToLower(strings.Split(contentType, ";")[0]) {
	case "application/pdf":
		return ".pdf", true
	case "image/png":
		return ".png", true
	case "image/jpeg":
		return ".jpg", true
	default:
		return "", false
	}
}

func saveManualProof(file *multipart.FileHeader) (path, originalName, mime string, size int64, err error) {
	if file == nil {
		return "", "", "", 0, ErrManualProof
	}
	if file.Size <= 0 || file.Size > manualMaxProofBytes() {
		return "", "", "", 0, ErrManualProof
	}
	source, err := file.Open()
	if err != nil {
		return "", "", "", 0, ErrManualProof
	}
	defer source.Close()
	head := make([]byte, 512)
	read, err := source.Read(head)
	if err != nil && err != io.EOF {
		return "", "", "", 0, ErrManualProof
	}
	mime = http.DetectContentType(head[:read])
	extension, ok := proofExtension(mime)
	if !ok {
		return "", "", "", 0, ErrManualProof
	}
	if _, err := source.Seek(0, io.SeekStart); err != nil {
		return "", "", "", 0, ErrManualProof
	}
	dir := manualProofDir()
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return "", "", "", 0, err
	}
	name := uuid.NewV4().String() + extension
	destination := filepath.Join(dir, name)
	out, err := os.OpenFile(destination, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0o640)
	if err != nil {
		return "", "", "", 0, err
	}
	written, copyErr := io.Copy(out, io.LimitReader(source, manualMaxProofBytes()+1))
	closeErr := out.Close()
	if copyErr != nil || closeErr != nil || written > manualMaxProofBytes() {
		_ = os.Remove(destination)
		return "", "", "", 0, ErrManualProof
	}
	originalName = filepath.Base(file.Filename)
	originalName = strings.Map(func(r rune) rune {
		if r == '\r' || r == '\n' || r == '"' {
			return -1
		}
		return r
	}, originalName)
	if originalName == "." || originalName == "" {
		originalName = "payment-proof" + extension
	}
	return destination, originalName, mime, written, nil
}

func UploadManualPayment(input ManualUploadInput) (*paymentModels.ManualPaymentSubmission, error) {
	if !manualEnabled() {
		return nil, ErrManualInvalid
	}
	method := strings.ToUpper(strings.TrimSpace(input.Method))
	if method != paymentModels.ManualMethodMTNMoMo && method != paymentModels.ManualMethodBankTransfer {
		return nil, ErrManualInvalid
	}
	requesterName := strings.TrimSpace(input.RequesterName)
	if requesterName == "" {
		return nil, ErrManualInvalid
	}
	var facts manualOrderFacts
	var submission paymentModels.ManualPaymentSubmission
	var reuse bool
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		var err error
		facts, err = manualFacts(tx, input.Purpose, input.OrderNumber, input.OrganisationNumber, input.AllowOnboarding)
		if err != nil {
			return err
		}
		var existing paymentModels.ManualPaymentSubmission
		if err := tx.Where("order_number = ? AND status = ?", facts.OrderNumber, paymentModels.ManualStatusSubmitted).
			Order("created_at desc").First(&existing).Error; err == nil {
			submission = existing
			reuse = true
			return nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		customerName := firstNonEmpty(input.CustomerName, facts.CustomerName)
		customerEmail := firstNonEmpty(input.CustomerEmail, facts.CustomerEmail)
		customerPhone := firstNonEmpty(input.CustomerPhone, facts.CustomerPhone, input.PayerPhone)
		submission = paymentModels.ManualPaymentSubmission{
			Number: uuid.NewV4().String(), OrderNumber: facts.OrderNumber, Purpose: facts.Purpose,
			OrganisationNumber: facts.OrganisationNumber, CustomerName: customerName,
			CustomerEmail: customerEmail, CustomerPhone: customerPhone, RequesterName: requesterName, Method: method,
			AmountUGX: facts.AmountUGX, Currency: facts.Currency, PayerName: input.PayerName,
			PayerPhone: input.PayerPhone, PayerReference: input.PayerReference, CustomerNote: input.CustomerNote,
			Status: paymentModels.ManualStatusSubmitted, SubmittedAt: time.Now().UTC(),
		}
		if err := tx.Create(&submission).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if reuse {
		return &submission, nil
	}
	proofPath, proofName, proofMime, proofSize, err := saveManualProof(input.File)
	if err != nil {
		config.DB.Unscoped().Delete(&submission)
		return nil, err
	}
	submission.ProofPath = proofPath
	submission.ProofName = proofName
	submission.ProofMime = proofMime
	submission.ProofSize = proofSize
	if err := config.DB.Save(&submission).Error; err != nil {
		_ = os.Remove(proofPath)
		config.DB.Unscoped().Delete(&submission)
		return nil, err
	}
	auditServices.RecordFact(auditServices.AuditRecord{
		OrganisationNumber: submission.OrganisationNumber, Action: "manual_payment.submitted",
		EntityType: "manual_payment", EntityID: submission.Number,
		NewValues: map[string]interface{}{"order": submission.OrderNumber, "method": submission.Method, "amount": submission.AmountUGX},
	})
	paymentNotifier.NotifySubmitted(&submission)
	return &submission, nil
}

func manualOrderStatus(submission paymentModels.ManualPaymentSubmission) string {
	switch submission.Purpose {
	case paymentModels.PurposeSetup:
		var order commercialModels.SetupOrder
		if err := config.DB.Where("number = ?", submission.OrderNumber).First(&order).Error; err == nil {
			return order.Status
		}
	case paymentModels.PurposeSubscription:
		var order commercialModels.SubscriptionOrder
		if err := config.DB.Where("number = ?", submission.OrderNumber).First(&order).Error; err == nil {
			return order.Status
		}
	}
	return ""
}

func toManualView(submission paymentModels.ManualPaymentSubmission) ManualPaymentView {
	return ManualPaymentView{
		Number: submission.Number, OrderNumber: submission.OrderNumber, Purpose: submission.Purpose,
		OrganisationNumber: submission.OrganisationNumber,
		CustomerName:       firstNonEmpty(submission.CustomerName, submission.PayerName),
		RequesterName:      firstNonEmpty(submission.RequesterName, submission.PayerName, submission.CustomerName),
		CustomerEmail:      submission.CustomerEmail, CustomerPhone: submission.CustomerPhone,
		Method: submission.Method, AmountUGX: submission.AmountUGX, Currency: submission.Currency,
		PayerName: submission.PayerName, PayerPhone: submission.PayerPhone,
		PayerReference: submission.PayerReference, CustomerNote: submission.CustomerNote,
		Status: submission.Status, ProofName: submission.ProofName, ProofMime: submission.ProofMime,
		ProofSize: submission.ProofSize, HasProof: submission.ProofPath != "",
		SubmittedAt: submission.SubmittedAt, ReviewedAt: submission.ReviewedAt,
		ReviewedByNumber: submission.ReviewedByNumber, ReviewNote: submission.ReviewNote,
		OrderStatus: manualOrderStatus(submission), TokenClaimedAt: submission.TokenClaimedAt,
	}
}

func ManualView(submission *paymentModels.ManualPaymentSubmission) ManualPaymentView {
	if submission == nil {
		return ManualPaymentView{}
	}
	return toManualView(*submission)
}

func GetManualPayment(number, organisationNumber string) (*ManualPaymentView, error) {
	var submission paymentModels.ManualPaymentSubmission
	if err := config.DB.Where("number = ?", number).First(&submission).Error; err != nil {
		return nil, ErrManualNotFound
	}
	if organisationNumber != "" && submission.OrganisationNumber != organisationNumber {
		return nil, ErrManualNotFound
	}
	if organisationNumber == "" && submission.OrganisationNumber != "" {
		return nil, ErrManualNotFound
	}
	view := toManualView(submission)
	return &view, nil
}

func ListManualPayments(filter ManualListFilter) ([]ManualPaymentView, int64, error) {
	query := config.DB.Model(&paymentModels.ManualPaymentSubmission{})
	if strings.TrimSpace(filter.Status) != "" {
		query = query.Where("status = ?", strings.ToUpper(strings.TrimSpace(filter.Status)))
	}
	if search := strings.TrimSpace(filter.Search); search != "" {
		like := "%" + search + "%"
		query = query.Where("number LIKE ? OR order_number LIKE ? OR customer_name LIKE ? OR customer_email LIKE ? OR payer_name LIKE ? OR payer_phone LIKE ? OR payer_reference LIKE ?", like, like, like, like, like, like, like)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	page := filter.Page
	if page < 1 {
		page = 1
	}
	perPage := filter.PerPage
	if perPage < 1 {
		perPage = 25
	}
	var rows []paymentModels.ManualPaymentSubmission
	if err := query.Order("submitted_at desc").Limit(perPage).Offset((page - 1) * perPage).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	views := make([]ManualPaymentView, 0, len(rows))
	for _, row := range rows {
		views = append(views, toManualView(row))
	}
	return views, total, nil
}

func GetManualPaymentDetail(number string) (*ManualPaymentDetail, error) {
	var submission paymentModels.ManualPaymentSubmission
	if err := config.DB.Where("number = ?", number).First(&submission).Error; err != nil {
		return nil, ErrManualNotFound
	}
	subject, body := paymentNotifier.ApprovalEmail(&submission)
	return &ManualPaymentDetail{
		ManualPaymentView:     toManualView(submission),
		ProofPath:             submission.ProofPath,
		SuggestedEmailSubject: subject,
		SuggestedEmailBody:    body,
	}, nil
}

func ApproveManualPayment(number, actor, note, emailSubject, emailBody, emailTo string) (*ManualPaymentView, error) {
	var submission paymentModels.ManualPaymentSubmission
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("number = ?", number).First(&submission).Error; err != nil {
			return ErrManualNotFound
		}
		if submission.Status == paymentModels.ManualStatusApproved {
			return nil
		}
		if submission.Status != paymentModels.ManualStatusSubmitted {
			return ErrManualInvalid
		}
		now := time.Now().UTC()
		if err := markManualOrderPaid(tx, submission); err != nil {
			return err
		}
		submission.Status = paymentModels.ManualStatusApproved
		submission.ReviewedAt = &now
		submission.ReviewedByNumber = actor
		submission.ReviewNote = strings.TrimSpace(note)
		return tx.Save(&submission).Error
	})
	if err != nil {
		return nil, err
	}
	view := toManualView(submission)
	auditServices.RecordFact(auditServices.AuditRecord{
		OrganisationNumber: submission.OrganisationNumber, ActorUserNumber: actor,
		Action: "manual_payment.approved", EntityType: "manual_payment", EntityID: submission.Number,
		NewValues: map[string]interface{}{"order": submission.OrderNumber, "amount": submission.AmountUGX, "note": submission.ReviewNote},
	})
	paymentNotifier.NotifyReviewed(&submission, true, emailSubject, emailBody, emailTo)
	return &view, nil
}

func RejectManualPayment(number, actor, note string) (*ManualPaymentView, error) {
	if strings.TrimSpace(note) == "" {
		return nil, ErrManualRejectReason
	}
	var submission paymentModels.ManualPaymentSubmission
	if err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("number = ?", number).First(&submission).Error; err != nil {
			return ErrManualNotFound
		}
		if submission.Status != paymentModels.ManualStatusSubmitted {
			return ErrManualInvalid
		}
		now := time.Now().UTC()
		submission.Status = paymentModels.ManualStatusRejected
		submission.ReviewedAt = &now
		submission.ReviewedByNumber = actor
		submission.ReviewNote = strings.TrimSpace(note)
		return tx.Save(&submission).Error
	}); err != nil {
		return nil, err
	}
	view := toManualView(submission)
	auditServices.RecordFact(auditServices.AuditRecord{
		OrganisationNumber: submission.OrganisationNumber, ActorUserNumber: actor,
		Action: "manual_payment.rejected", EntityType: "manual_payment", EntityID: submission.Number,
		NewValues: map[string]interface{}{"order": submission.OrderNumber, "note": submission.ReviewNote},
	})
	paymentNotifier.NotifyReviewed(&submission, false, "", "", "")
	return &view, nil
}

func markManualOrderPaid(tx *gorm.DB, submission paymentModels.ManualPaymentSubmission) error {
	switch submission.Purpose {
	case paymentModels.PurposeSetup:
		return tx.Model(&commercialModels.SetupOrder{}).
			Where("number = ? AND status = ?", submission.OrderNumber, "AWAITING_PAYMENT").
			Update("status", "PAID").Error
	case paymentModels.PurposeSubscription:
		return tx.Model(&commercialModels.SubscriptionOrder{}).
			Where("number = ? AND status = ?", submission.OrderNumber, "AWAITING_PAYMENT").
			Update("status", "PAID").Error
	default:
		return ErrManualInvalid
	}
}

func ClaimManualPayment(number, organisationNumber, actor string) (string, *ManualPaymentView, error) {
	var submission paymentModels.ManualPaymentSubmission
	if err := config.DB.Where("number = ?", number).First(&submission).Error; err != nil {
		return "", nil, ErrManualNotFound
	}
	if organisationNumber != "" && submission.OrganisationNumber != organisationNumber {
		return "", nil, ErrManualNotFound
	}
	if organisationNumber == "" && submission.OrganisationNumber != "" {
		return "", nil, ErrManualNotFound
	}
	if submission.Status == paymentModels.ManualStatusRejected {
		return "", nil, ErrManualRejected
	}
	if submission.Status != paymentModels.ManualStatusApproved {
		return "", nil, ErrManualNotApproved
	}
	if submission.TokenClaimedAt != nil {
		view := toManualView(submission)
		return "", &view, ErrManualAlready
	}
	result, err := paymentFulfilment.FulfilOrder(
		submission.Purpose, submission.OrderNumber, actor, submission.Number,
		submission.Method, submission.OrganisationNumber,
	)
	if err != nil {
		if errors.Is(err, paymentFulfilment.ErrAlreadyFulfilled) {
			return "", nil, ErrManualAlready
		}
		return "", nil, err
	}
	now := time.Now().UTC()
	submission.TokenClaimedAt = &now
	submission.ActivationTokenNumber = result.TokenNumber
	if err := config.DB.Save(&submission).Error; err != nil {
		return "", nil, err
	}
	view := toManualView(submission)
	return result.RawToken, &view, nil
}

func ManualOrderReference(number string) (string, error) {
	var submission paymentModels.ManualPaymentSubmission
	if err := config.DB.Select("order_number").Where("number = ?", number).First(&submission).Error; err != nil {
		return "", ErrManualNotFound
	}
	return submission.OrderNumber, nil
}
