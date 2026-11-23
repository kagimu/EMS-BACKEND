package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	commercialModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commercial/models"
	auditServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/audit/services"
	paymentResponses "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/payments/api/responses"
	paymentEvents "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/payments/events"
	paymentFulfilment "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/payments/fulfilment"
	paymentModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/payments/models"
	paymentNotifier "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/payments/notifier"
	paymentProviders "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/payments/providers"
	eosModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrPaymentNotFound     = errors.New("payment not found")
	ErrOrderNotPayable     = errors.New("this order is not awaiting payment")
	ErrPaymentInProgress   = errors.New("another payment is already in progress for this order")
	ErrInvalidPhone        = errors.New("enter a valid Ugandan mobile-money number")
	ErrInvalidAmount       = errors.New("payment amount is not valid")
	ErrProviderNotEnabled  = errors.New("the selected mobile-money provider is not enabled")
	ErrProviderUnavailable = errors.New("the mobile-money provider could not be reached")
	ErrPaymentNotPaid      = errors.New("the payment has not been confirmed as successful")
	ErrTokenAlreadyClaimed = errors.New("the activation token for this payment has already been claimed")
	ErrPaymentNotAllowed   = errors.New("this payment is not available to the current caller")
)

type InitiateInput struct {
	Purpose              string
	Reference            string
	Provider             string
	Phone                string
	ActorUserNumber      string
	OrganisationNumber   string
	AllowOnboardingOrder bool
	IdempotencyKey       string
}

type ListFilter struct {
	Search             string
	Provider           string
	Status             string
	OrganisationNumber string
	From               *time.Time
	To                 *time.Time
	Page               int
	PerPage            int
}

type orderFacts struct {
	Reference          string
	Purpose            string
	AmountUGX          int64
	Currency           string
	OrganisationNumber string
	ContactPhone       string
	Status             string
}

func newUUID() string { return uuid.NewV4().String() }

func maxPaymentAmount() int64 {
	value := strings.TrimSpace(os.Getenv("PAYMENTS_MAX_AMOUNT_UGX"))
	if value == "" {
		return 100000000
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsed <= 0 {
		return 100000000
	}
	return parsed
}

func pendingPaymentAge() time.Duration {
	value := strings.TrimSpace(os.Getenv("PAYMENTS_PENDING_STALE_MINUTES"))
	if value == "" {
		return 30 * time.Minute
	}
	minutes, err := strconv.Atoi(value)
	if err != nil || minutes <= 0 {
		return 30 * time.Minute
	}
	return time.Duration(minutes) * time.Minute
}

func providerTimeout() time.Duration {
	value := strings.TrimSpace(os.Getenv("PAYMENTS_PROVIDER_TIMEOUT_SECONDS"))
	if value == "" {
		return 35 * time.Second
	}
	seconds, err := strconv.Atoi(value)
	if err != nil || seconds <= 0 {
		return 35 * time.Second
	}
	return time.Duration(seconds) * time.Second
}

func simulatedPaymentsEnabled() bool {
	if strings.EqualFold(strings.TrimSpace(os.Getenv("ENVIRONMENT")), "production") {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(os.Getenv("DEV_SIMULATED_PAYMENTS_ENABLED")), "true")
}

// Options returns the providers a client may choose and whether the dev-only
// simulator is available. It never exposes credentials.
func Options() paymentResponses.PaymentOptions {
	infos := paymentProviders.EnabledProviders()
	providers := make([]paymentResponses.ProviderInfo, 0, len(infos))
	for _, info := range infos {
		providers = append(providers, paymentResponses.ProviderInfo{
			Name:          info.Name,
			DisplayName:   info.DisplayName,
			Environment:   info.Environment,
			TestMode:      info.TestMode,
			ContractReady: info.ContractReady,
			Currency:      info.Currency,
		})
	}
	return paymentResponses.PaymentOptions{
		Providers:         providers,
		SimulatedEnabled:  simulatedPaymentsEnabled(),
		Currency:          "UGX",
		MaxAmountUGX:      maxPaymentAmount(),
		HasActiveProvider: len(providers) > 0,
	}
}

func normalisePhone(raw string) (string, error) {
	value := strings.TrimSpace(raw)
	replacer := strings.NewReplacer(" ", "", "-", "", "(", "", ")", "", ".", "")
	value = replacer.Replace(value)
	value = strings.TrimPrefix(value, "+")
	value = strings.TrimPrefix(value, "00")
	if strings.HasPrefix(value, "0") && len(value) == 10 {
		value = "256" + value[1:]
	}
	if (strings.HasPrefix(value, "7") || strings.HasPrefix(value, "3")) && len(value) == 9 {
		value = "256" + value
	}
	if len(value) != 12 || !strings.HasPrefix(value, "256") {
		return "", ErrInvalidPhone
	}
	for _, character := range value {
		if character < '0' || character > '9' {
			return "", ErrInvalidPhone
		}
	}
	return value, nil
}

func maskPhone(phone string) string {
	value := strings.TrimSpace(phone)
	if len(value) <= 6 {
		return "***"
	}
	return value[:4] + strings.Repeat("*", len(value)-7) + value[len(value)-3:]
}

func factsForOrder(tx *gorm.DB, purpose, reference, organisationNumber string, allowOnboarding bool) (orderFacts, error) {
	switch purpose {
	case paymentModels.PurposeSetup:
		var order commercialModels.SetupOrder
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("number = ?", reference).First(&order).Error; err != nil {
			return orderFacts{}, ErrPaymentNotFound
		}
		if order.Status != "AWAITING_PAYMENT" {
			return orderFacts{}, ErrOrderNotPayable
		}
		return orderFacts{
			Reference: order.Number, Purpose: purpose, AmountUGX: order.SetupTotalUGX,
			Currency: firstNonEmpty(order.Currency, "UGX"), ContactPhone: order.ContactPhone,
			Status: order.Status,
		}, nil
	case paymentModels.PurposeSubscription:
		var order commercialModels.SubscriptionOrder
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("number = ?", reference).First(&order).Error; err != nil {
			return orderFacts{}, ErrPaymentNotFound
		}
		if order.Status != "AWAITING_PAYMENT" {
			return orderFacts{}, ErrOrderNotPayable
		}
		if order.OrganisationNumber != "" && order.OrganisationNumber != organisationNumber {
			return orderFacts{}, ErrPaymentNotFound
		}
		if order.OrganisationNumber == "" && !allowOnboarding {
			return orderFacts{}, ErrPaymentNotAllowed
		}
		return orderFacts{
			Reference: order.Number, Purpose: purpose, AmountUGX: order.TotalUGX,
			Currency: firstNonEmpty(order.Currency, "UGX"), OrganisationNumber: order.OrganisationNumber,
			Status: order.Status,
		}, nil
	default:
		return orderFacts{}, ErrPaymentNotFound
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func customerProfileNumber(tx *gorm.DB, organisationNumber string) string {
	if strings.TrimSpace(organisationNumber) == "" {
		return ""
	}
	var profile struct {
		Number string
	}
	if err := tx.Table("ems_customer_profiles").Select("number").
		Where("organisation_number = ?", organisationNumber).Take(&profile).Error; err != nil {
		return ""
	}
	return profile.Number
}

func stringPointer(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}

func markOrderPaidTx(tx *gorm.DB, payment *paymentModels.MobileMoneyPayment) {
	switch payment.Purpose {
	case paymentModels.PurposeSetup:
		tx.Model(&commercialModels.SetupOrder{}).
			Where("number = ? AND status = ?", payment.Reference, "AWAITING_PAYMENT").
			Update("status", "PAID")
	case paymentModels.PurposeSubscription:
		tx.Model(&commercialModels.SubscriptionOrder{}).
			Where("number = ? AND status = ?", payment.Reference, "AWAITING_PAYMENT").
			Update("status", "PAID")
	}
}

func eventTypeForStatus(status string) string {
	switch status {
	case paymentModels.StatusSuccessful:
		return paymentModels.EventPaymentSuccessful
	case paymentModels.StatusFailed:
		return paymentModels.EventPaymentFailed
	case paymentModels.StatusExpired:
		return paymentModels.EventPaymentExpired
	case paymentModels.StatusCancelled:
		return paymentModels.EventPaymentCancelled
	case paymentModels.StatusUnknown:
		return paymentModels.EventPaymentUnknown
	case paymentModels.StatusPending:
		return paymentModels.EventPaymentPending
	default:
		return paymentModels.EventPaymentInitiated
	}
}

func normaliseProviderStatus(status string) string {
	switch strings.ToUpper(strings.TrimSpace(status)) {
	case paymentModels.StatusSuccessful:
		return paymentModels.StatusSuccessful
	case paymentModels.StatusFailed, "REJECTED":
		return paymentModels.StatusFailed
	case paymentModels.StatusUnknown, "AMBIGUOUS":
		return paymentModels.StatusUnknown
	case paymentModels.StatusExpired, paymentModels.StatusCancelled:
		return strings.ToUpper(strings.TrimSpace(status))
	default:
		return paymentModels.StatusPending
	}
}

func truncateReason(value string) string {
	value = strings.TrimSpace(value)
	if len(value) > 500 {
		return value[:500]
	}
	return value
}

func applyProviderResult(number string, result paymentProviders.Result, actor string) (*paymentModels.MobileMoneyPayment, error) {
	var payment paymentModels.MobileMoneyPayment
	becameSuccessful := false
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("number = ?", number).First(&payment).Error; err != nil {
			return ErrPaymentNotFound
		}
		if paymentModels.IsFinal(payment.Status) {
			return nil
		}
		next := normaliseProviderStatus(result.Status)
		failure := truncateReason(result.FailureReason)
		if next == paymentModels.StatusPending && payment.ExpiresAt != nil && time.Now().UTC().After(*payment.ExpiresAt) {
			next = paymentModels.StatusExpired
			failure = "payment request expired before provider confirmation"
		}
		if next == paymentModels.StatusSuccessful && result.ProviderAmount > 0 &&
			(result.ProviderAmount != payment.AmountUGX || (result.ProviderCurrency != "" && !strings.EqualFold(result.ProviderCurrency, payment.Currency))) {
			next = paymentModels.StatusUnknown
			failure = "provider amount or currency mismatch; verification required"
		}
		if !paymentModels.CanTransition(payment.Status, next) {
			next = paymentModels.StatusUnknown
			failure = "provider returned a status that cannot be applied to the current payment state"
		}
		now := time.Now().UTC()
		payment.ProviderStatus = result.ProviderStatus
		payment.ProviderTransactionID = firstNonEmpty(result.ProviderTransactionID, payment.ProviderTransactionID)
		payment.ProviderAmount = result.ProviderAmount
		payment.ProviderCurrency = firstNonEmpty(result.ProviderCurrency, payment.ProviderCurrency)
		payment.ResponseJSON = result.RawResponse
		payment.LastCheckedAt = &now
		payment.CheckCount++
		payment.Status = next
		payment.StatusVersion++
		if next == paymentModels.StatusSuccessful {
			becameSuccessful = true
			if payment.PaidAt == nil {
				payment.PaidAt = &now
			}
			payment.UnknownAt = nil
			payment.FailureReason = ""
			payment.FailureCode = ""
			if payment.FulfilmentStatus == "" || payment.FulfilmentStatus == paymentModels.FulfilmentNotRequired {
				payment.FulfilmentStatus = paymentModels.FulfilmentPending
			}
			markOrderPaidTx(tx, &payment)
		} else if next == paymentModels.StatusUnknown {
			payment.UnknownAt = &now
			payment.FailureReason = failure
			payment.FailureCode = "PROVIDER_OUTCOME_UNKNOWN"
		} else if next == paymentModels.StatusFailed {
			payment.FailureReason = failure
			payment.FailureCode = "PROVIDER_DECLINED"
		} else if next == paymentModels.StatusExpired {
			payment.FailureReason = failure
			payment.FailureCode = "EXPIRED"
		}
		if err := tx.Save(&payment).Error; err != nil {
			return err
		}
		if err := paymentEvents.Record(tx, payment.Number, eventTypeForStatus(next),
			payment.Provider, payment.ProviderReference, actor, map[string]interface{}{
				"provider_status": payment.ProviderStatus, "failure": payment.FailureReason,
			}); err != nil {
			return err
		}
		if next == paymentModels.StatusSuccessful && payment.FulfilmentStatus == paymentModels.FulfilmentPending {
			_ = paymentEvents.Record(tx, payment.Number, paymentModels.EventFulfilmentRequested,
				payment.Provider, payment.ProviderReference, actor, map[string]interface{}{
					"purpose": payment.Purpose, "order": payment.Reference,
				})
		}
		return nil
	})
	if err == nil && becameSuccessful {
		paymentNotifier.NotifyMobileMoneySuccessful(&payment)
	}
	return &payment, err
}

// Initiate starts a collection request. The amount always comes from the
// server-side order. The payment row is committed before the provider call so
// an ambiguous timeout can be reconciled without creating a second charge.
func Initiate(input InitiateInput) (*paymentModels.MobileMoneyPayment, error) {
	providerName := strings.ToUpper(strings.TrimSpace(input.Provider))
	switch providerName {
	case paymentProviders.MTNMoMo, paymentProviders.AirtelMoney, paymentProviders.Simulated:
	default:
		return nil, ErrProviderNotEnabled
	}
	provider, err := paymentProviders.Get(providerName)
	if err != nil || !provider.Enabled() {
		return nil, ErrProviderNotEnabled
	}
	phone, err := normalisePhone(input.Phone)
	if err != nil {
		return nil, err
	}
	idempotencyKey := strings.TrimSpace(input.IdempotencyKey)

	var payment paymentModels.MobileMoneyPayment
	var facts orderFacts
	var reused bool
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		var existing paymentModels.MobileMoneyPayment
		if idempotencyKey != "" {
			if err := tx.Where("idempotency_key = ?", idempotencyKey).First(&existing).Error; err == nil {
				if existing.Reference != input.Reference || existing.Purpose != input.Purpose {
					return ErrPaymentInProgress
				}
				payment = existing
				reused = true
				return nil
			} else if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}
		var err error
		facts, err = factsForOrder(tx, input.Purpose, input.Reference, input.OrganisationNumber, input.AllowOnboardingOrder)
		if err != nil {
			return err
		}
		if facts.AmountUGX < 100 || facts.AmountUGX > maxPaymentAmount() {
			return ErrInvalidAmount
		}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("reference = ? AND purpose = ? AND status IN ?", facts.Reference, facts.Purpose,
				[]string{paymentModels.StatusInitiated, paymentModels.StatusPending, paymentModels.StatusUnknown}).
			Order("created_at desc").First(&existing).Error; err == nil {
			if existing.Provider == providerName && existing.PayerPhone == phone {
				payment = existing
				reused = true
				return nil
			}
			return ErrPaymentInProgress
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		now := time.Now().UTC()
		expiry := now.Add(pendingPaymentAge())
		number := newUUID()
		payment = paymentModels.MobileMoneyPayment{
			Number:                number,
			Provider:              providerName,
			Purpose:               facts.Purpose,
			Reference:             facts.Reference,
			CommercialOrderNumber: facts.Reference,
			OrganisationNumber:    facts.OrganisationNumber,
			CustomerProfileNumber: customerProfileNumber(tx, facts.OrganisationNumber),
			PayerPhone:            phone,
			AmountUGX:             facts.AmountUGX,
			Currency:              firstNonEmpty(facts.Currency, "UGX"),
			Status:                paymentModels.StatusInitiated,
			StatusVersion:         1,
			ProviderReference:     newUUID(),
			MerchantReference:     stringPointer(number),
			IdempotencyKey:        stringPointer(idempotencyKey),
			InitiatedByNumber:     input.ActorUserNumber,
			CreatedByNumber:       input.ActorUserNumber,
			ExpiresAt:             &expiry,
			FulfilmentStatus:      paymentModels.FulfilmentNotRequired,
		}
		requestPayload, _ := json.Marshal(map[string]interface{}{
			"provider": providerName, "purpose": facts.Purpose, "reference": facts.Reference,
			"amount_ugx": facts.AmountUGX, "currency": payment.Currency, "payer_phone": phone,
		})
		payment.RequestJSON = string(requestPayload)
		if err := tx.Create(&payment).Error; err != nil {
			if idempotencyKey != "" {
				if findErr := tx.Where("idempotency_key = ?", idempotencyKey).First(&existing).Error; findErr == nil {
					payment = existing
					reused = true
					return nil
				}
			}
			return err
		}
		if err := paymentEvents.Record(tx, payment.Number, paymentModels.EventPaymentInitiated,
			payment.Provider, payment.ProviderReference, input.ActorUserNumber, map[string]interface{}{
				"purpose": payment.Purpose, "order": payment.Reference, "amount": payment.AmountUGX,
				"currency": payment.Currency,
			}); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if reused {
		return &payment, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), providerTimeout())
	defer cancel()
	result, providerErr := provider.Initiate(ctx, paymentProviders.Request{
		PaymentNumber: payment.ProviderReference,
		Reference:     payment.Reference,
		Amount:        payment.AmountUGX,
		Currency:      payment.Currency,
		Phone:         phone,
		Description:   paymentDescription(facts),
	})
	if providerErr != nil {
		target := paymentModels.StatusFailed
		if paymentProviders.IsAmbiguous(providerErr) {
			target = paymentModels.StatusUnknown
		}
		updated, updateErr := applyProviderResult(payment.Number, paymentProviders.Result{
			Status:            target,
			ProviderStatus:    "PROVIDER_ERROR",
			ProviderReference: payment.ProviderReference,
			FailureReason:     truncateReason(providerErr.Error()),
			RawResponse:       result.RawResponse,
		}, input.ActorUserNumber)
		if updateErr != nil {
			return &payment, updateErr
		}
		return updated, nil
	}
	updated, err := applyProviderResult(payment.Number, paymentProviders.Result{
		Status:                result.Status,
		ProviderStatus:        result.ProviderStatus,
		ProviderReference:     firstNonEmpty(result.ProviderReference, payment.ProviderReference),
		ProviderTransactionID: result.ProviderTransactionID,
		ProviderAmount:        result.ProviderAmount,
		ProviderCurrency:      result.ProviderCurrency,
		FailureReason:         result.FailureReason,
		RawResponse:           result.RawResponse,
	}, input.ActorUserNumber)
	if err != nil {
		return &payment, err
	}
	return updated, nil
}

func paymentDescription(facts orderFacts) string {
	if facts.Purpose == paymentModels.PurposeSetup {
		return "EMS installation setup " + facts.Reference
	}
	return "EMS software subscription " + facts.Reference
}

func markOrderPaid(payment *paymentModels.MobileMoneyPayment) {
	switch payment.Purpose {
	case paymentModels.PurposeSetup:
		config.DB.Model(&commercialModels.SetupOrder{}).
			Where("number = ? AND status = ?", payment.Reference, "AWAITING_PAYMENT").
			Update("status", "PAID")
	case paymentModels.PurposeSubscription:
		config.DB.Model(&commercialModels.SubscriptionOrder{}).
			Where("number = ? AND status = ?", payment.Reference, "AWAITING_PAYMENT").
			Update("status", "PAID")
	}
}

// Refresh re-queries the telecom provider. Callbacks are never trusted as
// payment proof; this is the only function that advances a payment from a
// provider response.
func Refresh(number string) (*paymentModels.MobileMoneyPayment, error) {
	var payment paymentModels.MobileMoneyPayment
	if err := config.DB.Where("number = ?", number).First(&payment).Error; err != nil {
		return nil, ErrPaymentNotFound
	}
	if paymentModels.IsFinal(payment.Status) {
		return &payment, nil
	}
	if payment.ProviderReference == "" {
		return &payment, nil
	}
	provider, err := paymentProviders.Get(payment.Provider)
	if err != nil || !provider.Enabled() {
		return &payment, fmt.Errorf("%w: provider is not configured", ErrProviderUnavailable)
	}
	ctx, cancel := context.WithTimeout(context.Background(), providerTimeout())
	defer cancel()
	result, err := provider.VerifyCallback(ctx, payment.ProviderReference)
	if err != nil && payment.Provider == paymentProviders.AirtelMoney && payment.ProviderTransactionID != "" {
		result, err = provider.VerifyCallback(ctx, payment.ProviderTransactionID)
	}
	if err != nil {
		return &payment, fmt.Errorf("%w: %v", ErrProviderUnavailable, err)
	}
	return applyProviderResult(payment.Number, result, "")
}

// Claim refreshes the payment and asks the separate fulfilment service to
// issue the one-time activation token. Payment success is never reversed by a
// fulfilment failure. The raw token is returned only from this call.
func Claim(number, organisationNumber string) (string, *paymentModels.MobileMoneyPayment, error) {
	payment, err := Refresh(number)
	if err != nil && !errors.Is(err, ErrPaymentNotFound) {
		return "", payment, err
	}
	if err != nil {
		return "", nil, err
	}
	if payment.Status != paymentModels.StatusSuccessful {
		return "", payment, ErrPaymentNotPaid
	}
	if organisationNumber != "" && payment.OrganisationNumber != organisationNumber {
		return "", nil, ErrPaymentNotFound
	}
	if organisationNumber == "" && payment.OrganisationNumber != "" {
		return "", nil, ErrPaymentNotFound
	}

	fulfilment, err := paymentFulfilment.Fulfil(number, "")
	if err != nil {
		switch {
		case errors.Is(err, paymentFulfilment.ErrAlreadyFulfilled):
			return "", payment, ErrTokenAlreadyClaimed
		case errors.Is(err, paymentFulfilment.ErrPaymentNotSuccessful):
			return "", payment, ErrPaymentNotPaid
		default:
			return "", payment, err
		}
	}
	if fulfilment.RawToken == "" {
		return "", payment, ErrTokenAlreadyClaimed
	}
	now := time.Now().UTC()
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		var locked paymentModels.MobileMoneyPayment
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("number = ?", number).First(&locked).Error; err != nil {
			return ErrPaymentNotFound
		}
		if locked.TokenClaimedAt != nil {
			return ErrTokenAlreadyClaimed
		}
		locked.TokenClaimedAt = &now
		locked.ActivationTokenNumber = fulfilment.TokenNumber
		if err := tx.Save(&locked).Error; err != nil {
			return err
		}
		payment = &locked
		return nil
	})
	if err != nil {
		return "", payment, err
	}
	auditPayment("payment.token.claimed", payment, "", map[string]interface{}{"token": payment.ActivationTokenNumber})
	return fulfilment.RawToken, payment, nil
}

func auditPayment(action string, payment *paymentModels.MobileMoneyPayment, actor string, values map[string]interface{}) {
	auditServices.RecordFact(auditServices.AuditRecord{
		OrganisationNumber: payment.OrganisationNumber,
		ActorUserNumber:    actor,
		Action:             action,
		EntityType:         "mobile_money_payment",
		EntityID:           payment.Number,
		NewValues:          values,
	})
}

func organisationNames(numbers []string) map[string]string {
	result := map[string]string{}
	if len(numbers) == 0 {
		return result
	}
	var entities []eosModel.Entity
	if err := config.DB.Where("number IN ?", numbers).Find(&entities).Error; err != nil {
		return result
	}
	for _, entity := range entities {
		result[entity.Number] = entity.Name
	}
	return result
}

func customerForPayment(payment paymentModels.MobileMoneyPayment, names map[string]string) (string, string) {
	if payment.OrganisationNumber != "" {
		return names[payment.OrganisationNumber], ""
	}
	switch payment.Purpose {
	case paymentModels.PurposeSetup:
		var order commercialModels.SetupOrder
		if err := config.DB.Where("number = ?", payment.Reference).First(&order).Error; err == nil {
			return order.ContactName, order.ContactEmail
		}
	case paymentModels.PurposeSubscription:
		var subscription commercialModels.SubscriptionOrder
		if err := config.DB.Where("number = ?", payment.Reference).First(&subscription).Error; err == nil && subscription.SetupOrderNumber != "" {
			var setup commercialModels.SetupOrder
			if err := config.DB.Where("number = ?", subscription.SetupOrderNumber).First(&setup).Error; err == nil {
				return setup.ContactName, setup.ContactEmail
			}
		}
	}
	return "", ""
}

func toView(payment paymentModels.MobileMoneyPayment, names map[string]string, includeCustomer bool) paymentResponses.PaymentView {
	providerName := payment.Provider
	if provider, err := paymentProviders.Get(payment.Provider); err == nil {
		providerName = provider.DisplayName()
	}
	customerName, customerEmail := "", ""
	if includeCustomer {
		customerName, customerEmail = customerForPayment(payment, names)
	}
	return paymentResponses.PaymentView{
		Number:                payment.Number,
		Provider:              payment.Provider,
		ProviderName:          providerName,
		Purpose:               payment.Purpose,
		Reference:             payment.Reference,
		OrganisationNumber:    payment.OrganisationNumber,
		OrganisationName:      names[payment.OrganisationNumber],
		CommercialOrderNumber: firstNonEmpty(payment.CommercialOrderNumber, payment.Reference),
		CustomerProfileNumber: payment.CustomerProfileNumber,
		CustomerName:          customerName,
		CustomerEmail:         customerEmail,
		PayerPhoneMasked:      maskPhone(payment.PayerPhone),
		AmountUGX:             payment.AmountUGX,
		Currency:              payment.Currency,
		Status:                payment.Status,
		ProviderStatus:        payment.ProviderStatus,
		ProviderReference:     payment.ProviderReference,
		ProviderTransactionID: payment.ProviderTransactionID,
		FailureCode:           payment.FailureCode,
		FailureReason:         payment.FailureReason,
		FulfilmentStatus:      payment.FulfilmentStatus,
		FulfilmentError:       payment.FulfilmentError,
		ExpiresAt:             payment.ExpiresAt,
		UnknownAt:             payment.UnknownAt,
		InitiatedAt:           payment.CreatedAt,
		LastCheckedAt:         payment.LastCheckedAt,
		PaidAt:                payment.PaidAt,
		VerifiedAt:            payment.VerifiedAt,
		TokenClaimedAt:        payment.TokenClaimedAt,
		TokenClaimed:          payment.TokenClaimedAt != nil,
	}
}

func Get(number, organisationNumber string) (*paymentResponses.PaymentView, error) {
	var payment paymentModels.MobileMoneyPayment
	if err := config.DB.Where("number = ?", number).First(&payment).Error; err != nil {
		return nil, ErrPaymentNotFound
	}
	if organisationNumber != "" && payment.OrganisationNumber != organisationNumber {
		return nil, ErrPaymentNotFound
	}
	if organisationNumber == "" && payment.OrganisationNumber != "" {
		return nil, ErrPaymentNotFound
	}
	names := organisationNames([]string{payment.OrganisationNumber})
	view := toView(payment, names, false)
	return &view, nil
}

func applyListFilters(query *gorm.DB, filter ListFilter) *gorm.DB {
	if strings.TrimSpace(filter.Provider) != "" {
		query = query.Where("provider = ?", strings.ToUpper(strings.TrimSpace(filter.Provider)))
	}
	if strings.TrimSpace(filter.Status) != "" {
		if strings.EqualFold(strings.TrimSpace(filter.Status), paymentModels.StatusUnknown) {
			query = query.Where("status = ?", paymentModels.StatusUnknown)
		} else {
			query = query.Where("status = ?", strings.ToUpper(strings.TrimSpace(filter.Status)))
		}
	}
	if strings.TrimSpace(filter.OrganisationNumber) != "" {
		query = query.Where("organisation_number = ?", strings.TrimSpace(filter.OrganisationNumber))
	}
	if filter.From != nil {
		query = query.Where("created_at >= ?", *filter.From)
	}
	if filter.To != nil {
		query = query.Where("created_at <= ?", *filter.To)
	}
	if search := strings.TrimSpace(filter.Search); search != "" {
		like := "%" + search + "%"
		query = query.Where(
			`number LIKE ? OR reference LIKE ? OR organisation_number LIKE ? OR payer_phone LIKE ? OR provider_reference LIKE ? OR provider_transaction_id LIKE ?
			OR reference IN (SELECT number FROM ems_setup_orders WHERE contact_name LIKE ? OR contact_email LIKE ?)
			OR reference IN (
				SELECT number FROM ems_subscription_orders
				WHERE setup_order_number IN (SELECT number FROM ems_setup_orders WHERE contact_name LIKE ? OR contact_email LIKE ?)
			)
			OR organisation_number IN (SELECT number FROM entities WHERE name LIKE ? AND deleted_at IS NULL)`,
			like, like, like, like, like, like, like, like, like, like, like,
		)
	}
	return query
}

func normalisePage(page, perPage int) (int, int) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}
	return page, perPage
}

func List(filter ListFilter) (*paymentResponses.PaymentList, error) {
	page, perPage := normalisePage(filter.Page, filter.PerPage)
	var total int64
	if err := applyListFilters(config.DB.Model(&paymentModels.MobileMoneyPayment{}), filter).Count(&total).Error; err != nil {
		return nil, err
	}
	var rows []paymentModels.MobileMoneyPayment
	if err := applyListFilters(config.DB.Model(&paymentModels.MobileMoneyPayment{}), filter).
		Order("created_at desc").Limit(perPage).Offset((page - 1) * perPage).Find(&rows).Error; err != nil {
		return nil, err
	}
	orgNumbers := make([]string, 0, len(rows))
	for _, row := range rows {
		if row.OrganisationNumber != "" {
			orgNumbers = append(orgNumbers, row.OrganisationNumber)
		}
	}
	names := organisationNames(orgNumbers)
	items := make([]paymentResponses.PaymentView, 0, len(rows))
	for _, row := range rows {
		items = append(items, toView(row, names, true))
	}

	var grouped []struct {
		Provider string
		Status   string
		Count    int64
		Amount   int64
	}
	if err := applyListFilters(config.DB.Model(&paymentModels.MobileMoneyPayment{}), filter).
		Select("provider, status, COUNT(*) AS count, COALESCE(SUM(amount_ugx),0) AS amount").
		Group("provider, status").Scan(&grouped).Error; err != nil {
		return nil, err
	}
	summary := paymentResponses.PaymentSummary{ByProvider: []paymentResponses.ProviderTotal{}}
	byProvider := map[string]*paymentResponses.ProviderTotal{}
	for _, group := range grouped {
		summary.TotalCount += group.Count
		switch group.Status {
		case paymentModels.StatusSuccessful:
			summary.SuccessfulCount += group.Count
			summary.SuccessfulAmountUGX += group.Amount
		case paymentModels.StatusPending, paymentModels.StatusInitiated:
			summary.PendingCount += group.Count
		case paymentModels.StatusUnknown:
			summary.UnknownCount += group.Count
		case paymentModels.StatusFailed, paymentModels.StatusCancelled, paymentModels.StatusExpired:
			summary.FailedCount += group.Count
		}
		entry, ok := byProvider[group.Provider]
		if !ok {
			entry = &paymentResponses.ProviderTotal{Provider: group.Provider}
			byProvider[group.Provider] = entry
		}
		entry.Count += group.Count
		if group.Status == paymentModels.StatusSuccessful {
			entry.AmountUGX += group.Amount
		}
	}
	for _, entry := range byProvider {
		summary.ByProvider = append(summary.ByProvider, *entry)
	}
	status := WorkerStatus()
	return &paymentResponses.PaymentList{
		Items: items, Total: total, Page: page, PerPage: perPage, Summary: summary,
		Worker: paymentResponses.ReconciliationWorkerView{
			Enabled: status.Enabled, Running: status.Running,
			IntervalSeconds: status.IntervalSeconds, MaxAttempts: status.MaxAttempts,
			LastRunAt: status.LastRunAt, LastSuccessAt: status.LastSuccessAt,
			LastError: status.LastError, LastScannedCount: status.LastScannedCount,
			LastUpdatedCount: status.LastUpdatedCount,
		},
	}, nil
}

func orderStatus(payment paymentModels.MobileMoneyPayment) string {
	switch payment.Purpose {
	case paymentModels.PurposeSetup:
		var order commercialModels.SetupOrder
		if err := config.DB.Where("number = ?", payment.Reference).First(&order).Error; err == nil {
			return order.Status
		}
	case paymentModels.PurposeSubscription:
		var order commercialModels.SubscriptionOrder
		if err := config.DB.Where("number = ?", payment.Reference).First(&order).Error; err == nil {
			return order.Status
		}
	}
	return ""
}

func Detail(number string) (*paymentResponses.PaymentDetail, error) {
	var payment paymentModels.MobileMoneyPayment
	if err := config.DB.Where("number = ?", number).First(&payment).Error; err != nil {
		return nil, ErrPaymentNotFound
	}
	names := organisationNames([]string{payment.OrganisationNumber})
	view := toView(payment, names, true)
	events := make([]paymentResponses.PaymentEventView, 0)
	var storedEvents []paymentModels.PaymentEvent
	config.DB.Where("transaction_number = ?", payment.Number).Order("created_at asc").Find(&storedEvents)
	for _, event := range storedEvents {
		events = append(events, paymentResponses.PaymentEventView{
			EventType: event.EventType, ProviderCode: event.ProviderCode,
			ProviderReference: event.ProviderReference, ActorUserNumber: event.ActorUserNumber,
			Metadata: event.Metadata, CreatedAt: event.CreatedAt,
		})
	}
	reconciliations := make([]paymentResponses.ReconciliationView, 0)
	var storedReconciliations []paymentModels.PaymentReconciliation
	config.DB.Where("transaction_number = ?", payment.Number).Order("reconciled_at asc").Find(&storedReconciliations)
	for _, reconciliation := range storedReconciliations {
		reconciliations = append(reconciliations, paymentResponses.ReconciliationView{
			Result: reconciliation.Result, ProviderReference: reconciliation.ProviderReference,
			ProviderAmount: reconciliation.ProviderAmount, EMSAmount: reconciliation.EMSAmount,
			ProviderCurrency: reconciliation.ProviderCurrency, EMSCurrency: reconciliation.EMSCurrency,
			Discrepancy: reconciliation.Discrepancy, AttemptNumber: reconciliation.AttemptNumber,
			ReconciledAt: reconciliation.ReconciledAt,
		})
	}
	return &paymentResponses.PaymentDetail{
		PaymentView:      view,
		PayerPhone:       payment.PayerPhone,
		CheckCount:       payment.CheckCount,
		VerifiedByNumber: payment.VerifiedByNumber,
		OrderStatus:      orderStatus(payment),
		Events:           events,
		Reconciliations:  reconciliations,
	}, nil
}

// Verify is the owner console action: it re-queries the provider, records who
// checked it and returns the refreshed detail. It never marks a payment paid
// without a successful provider response.
func Verify(number, actor string) (*paymentResponses.PaymentDetail, error) {
	payment, err := Refresh(number)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	if payment.Status == paymentModels.StatusSuccessful {
		payment.VerifiedAt = &now
		payment.VerifiedByNumber = actor
		config.DB.Save(payment)
	}
	auditPayment("payment.owner_verified", payment, actor, map[string]interface{}{
		"status": payment.Status, "provider_status": payment.ProviderStatus,
	})
	return Detail(number)
}

// HandleCallback records the provider callback and immediately re-queries the
// provider. A callback body is never sufficient proof of payment.
func HandleCallback(providerName string, body []byte) error {
	provider, err := paymentProviders.Get(providerName)
	if err != nil {
		return nil
	}
	reference, transactionID, err := provider.CallbackReference(body)
	if err != nil || (reference == "" && transactionID == "") {
		return nil
	}
	var payment paymentModels.MobileMoneyPayment
	query := config.DB.Where("provider = ?", strings.ToUpper(strings.TrimSpace(providerName)))
	if reference != "" && transactionID != "" {
		query = query.Where("provider_reference = ? OR provider_transaction_id = ?", reference, transactionID)
	} else if reference != "" {
		query = query.Where("provider_reference = ?", reference)
	} else {
		query = query.Where("provider_transaction_id = ?", transactionID)
	}
	if err := query.First(&payment).Error; err != nil {
		// Airtel callbacks may identify the merchant order rather than our
		// internal provider reference. Fall back to the order reference and
		// choose the latest non-final attempt.
		if reference == "" {
			return nil
		}
		if err := config.DB.Where("provider = ? AND reference = ?", strings.ToUpper(strings.TrimSpace(providerName)), reference).
			Order("created_at desc").First(&payment).Error; err != nil {
			return nil
		}
	}
	safeCallback, _ := json.Marshal(map[string]interface{}{
		"provider": providerName, "provider_reference": reference,
		"provider_transaction_id": transactionID,
	})
	payment.CallbackJSON = string(safeCallback)
	if transactionID != "" {
		payment.ProviderTransactionID = transactionID
	}
	config.DB.Save(&payment)
	_ = paymentEvents.Record(config.DB, payment.Number, paymentModels.EventCallbackReceived,
		payment.Provider, payment.ProviderReference, "", map[string]interface{}{
			"reference": reference, "transaction_id": transactionID,
		})
	if paymentModels.IsFinal(payment.Status) {
		return nil
	}
	_, err = Refresh(payment.Number)
	return err
}
