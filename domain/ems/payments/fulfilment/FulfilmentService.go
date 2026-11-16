package fulfilment

import (
	"errors"
	"strings"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	commercialModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commercial/models"
	commercialServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commercial/services"
	auditServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/audit/services"
	paymentEvents "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/payments/events"
	paymentModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/payments/models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrPaymentNotSuccessful = errors.New("payment is not successful")
	ErrAlreadyFulfilled     = errors.New("payment has already been fulfilled")
	ErrUnsupportedOrder     = errors.New("this payment type has no supported fulfilment operation")
)

type Result struct {
	RawToken    string
	TokenNumber string
	Payment     *paymentModels.MobileMoneyPayment
}

// Fulfil performs the idempotent order-fulfilment step after a verified
// successful payment. It never changes payment success and never calls a
// provider. A retry after a failure is safe.
func Fulfil(paymentNumber, actor string) (Result, error) {
	var result Result
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		var payment paymentModels.MobileMoneyPayment
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("number = ?", paymentNumber).First(&payment).Error; err != nil {
			return err
		}
		result.Payment = &payment
		if payment.Status != paymentModels.StatusSuccessful {
			return ErrPaymentNotSuccessful
		}
		if payment.FulfilmentStatus == paymentModels.FulfilmentCompleted {
			return ErrAlreadyFulfilled
		}
		if payment.Purpose != paymentModels.PurposeSetup && payment.Purpose != paymentModels.PurposeSubscription {
			return ErrUnsupportedOrder
		}

		purpose := commercialServices.TokenPurposeSetup
		orderColumn := "setup_order_number"
		if payment.Purpose == paymentModels.PurposeSubscription {
			purpose = commercialServices.TokenPurposeSubscription
			orderColumn = "subscription_order_number"
		}
		var existing commercialModels.ActivationToken
		if err := tx.Where(orderColumn+" = ? AND purpose = ?", payment.Reference, purpose).First(&existing).Error; err == nil {
			now := time.Now().UTC()
			payment.FulfilmentStatus = paymentModels.FulfilmentCompleted
			payment.FulfilmentCompletedAt = &now
			payment.FulfilmentError = ""
			if err := tx.Save(&payment).Error; err != nil {
				return err
			}
			return ErrAlreadyFulfilled
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		raw, err := commercialServices.GenerateRawToken()
		if err != nil {
			return err
		}
		now := time.Now().UTC()
		token := commercialModels.ActivationToken{
			Number:    uuid.NewV4().String(),
			Purpose:   purpose,
			TokenHash: commercialServices.HashToken(raw),
			Status:    commercialServices.TokenStatusIssued,
			IssuedAt:  &now,
		}
		if purpose == commercialServices.TokenPurposeSetup {
			token.SetupOrderNumber = payment.Reference
			if err := tx.Model(&commercialModels.SetupOrder{}).
				Where("number = ? AND status IN ?", payment.Reference, []string{"AWAITING_PAYMENT", "PAID"}).
				Update("status", "TOKEN_ISSUED").Error; err != nil {
				return err
			}
		} else {
			token.SubscriptionOrderNumber = payment.Reference
			if err := tx.Model(&commercialModels.SubscriptionOrder{}).
				Where("number = ? AND status IN ?", payment.Reference, []string{"AWAITING_PAYMENT", "PAID"}).
				Update("status", "TOKEN_ISSUED").Error; err != nil {
				return err
			}
		}
		if err := tx.Create(&token).Error; err != nil {
			return err
		}
		payment.FulfilmentStatus = paymentModels.FulfilmentCompleted
		payment.FulfilmentCompletedAt = &now
		payment.FulfilmentError = ""
		payment.ActivationTokenNumber = token.Number
		if err := tx.Save(&payment).Error; err != nil {
			return err
		}
		if err := paymentEvents.Record(tx, payment.Number, paymentModels.EventFulfilmentCompleted,
			payment.Provider, payment.ProviderReference, actor, map[string]interface{}{
				"purpose": payment.Purpose, "order": payment.Reference, "token": token.Number,
			}); err != nil {
			return err
		}
		result.RawToken = raw
		result.TokenNumber = token.Number
		result.Payment = &payment
		return nil
	})
	if err != nil {
		if !errors.Is(err, ErrAlreadyFulfilled) && !errors.Is(err, ErrPaymentNotSuccessful) && !errors.Is(err, ErrUnsupportedOrder) {
			_ = markFailed(paymentNumber, err)
		}
		return result, err
	}
	auditFulfilment(result.Payment, actor, "payment.fulfilment.completed")
	return result, nil
}

// FulfilOrder performs the same idempotent token issuance for an approved
// manual payment that has no provider transaction row. Approval has already
// marked the commercial order PAID; this function only creates the one-time
// activation token.
func FulfilOrder(purpose, reference, actor, transactionNumber, provider, organisation string) (Result, error) {
	var result Result
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		tokenPurpose := commercialServices.TokenPurposeSetup
		orderColumn := "setup_order_number"
		if purpose == paymentModels.PurposeSubscription {
			tokenPurpose = commercialServices.TokenPurposeSubscription
			orderColumn = "subscription_order_number"
		} else if purpose != paymentModels.PurposeSetup {
			return ErrUnsupportedOrder
		}

		var existing commercialModels.ActivationToken
		if err := tx.Where(orderColumn+" = ? AND purpose = ?", reference, tokenPurpose).First(&existing).Error; err == nil {
			return ErrAlreadyFulfilled
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		raw, err := commercialServices.GenerateRawToken()
		if err != nil {
			return err
		}
		now := time.Now().UTC()
		token := commercialModels.ActivationToken{
			Number:    uuid.NewV4().String(),
			Purpose:   tokenPurpose,
			TokenHash: commercialServices.HashToken(raw),
			Status:    commercialServices.TokenStatusIssued,
			IssuedAt:  &now,
		}
		if purpose == paymentModels.PurposeSetup {
			var order commercialModels.SetupOrder
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("number = ?", reference).First(&order).Error; err != nil {
				return err
			}
			if order.Status != "PAID" && order.Status != "AWAITING_PAYMENT" && order.Status != "TOKEN_ISSUED" {
				return ErrUnsupportedOrder
			}
			token.SetupOrderNumber = reference
			if err := tx.Model(&commercialModels.SetupOrder{}).Where("number = ?", reference).
				Update("status", "TOKEN_ISSUED").Error; err != nil {
				return err
			}
		} else {
			var order commercialModels.SubscriptionOrder
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("number = ?", reference).First(&order).Error; err != nil {
				return err
			}
			if order.Status != "PAID" && order.Status != "AWAITING_PAYMENT" && order.Status != "TOKEN_ISSUED" {
				return ErrUnsupportedOrder
			}
			token.SubscriptionOrderNumber = reference
			if err := tx.Model(&commercialModels.SubscriptionOrder{}).Where("number = ?", reference).
				Update("status", "TOKEN_ISSUED").Error; err != nil {
				return err
			}
		}
		if err := tx.Create(&token).Error; err != nil {
			return err
		}
		if err := paymentEvents.Record(tx, transactionNumber, paymentModels.EventFulfilmentCompleted,
			provider, transactionNumber, actor, map[string]interface{}{
				"purpose": purpose, "order": reference, "token": token.Number, "manual": true,
			}); err != nil {
			return err
		}
		result.RawToken = raw
		result.TokenNumber = token.Number
		return nil
	})
	if err != nil {
		return result, err
	}
	auditServices.RecordFact(auditServices.AuditRecord{
		OrganisationNumber: organisation, ActorUserNumber: actor,
		Action: "manual_payment.fulfilment.completed", EntityType: "manual_payment",
		EntityID: transactionNumber, NewValues: map[string]interface{}{"order": reference, "token": result.TokenNumber},
	})
	return result, nil
}

func markFailed(paymentNumber string, cause error) error {
	payment := paymentModels.MobileMoneyPayment{}
	if err := config.DB.Where("number = ?", paymentNumber).First(&payment).Error; err != nil {
		return err
	}
	if payment.FulfilmentStatus == paymentModels.FulfilmentCompleted {
		return nil
	}
	payment.FulfilmentStatus = paymentModels.FulfilmentFailed
	payment.FulfilmentError = truncate(cause.Error())
	if err := config.DB.Save(&payment).Error; err != nil {
		return err
	}
	_ = paymentEvents.Record(config.DB, payment.Number, paymentModels.EventFulfilmentFailed,
		payment.Provider, payment.ProviderReference, "", map[string]interface{}{
			"purpose": payment.Purpose, "order": payment.Reference, "error": payment.FulfilmentError,
		})
	auditFulfilment(&payment, "", "payment.fulfilment.failed")
	return nil
}

func auditFulfilment(payment *paymentModels.MobileMoneyPayment, actor, action string) {
	if payment == nil {
		return
	}
	auditServices.RecordFact(auditServices.AuditRecord{
		OrganisationNumber: payment.OrganisationNumber,
		ActorUserNumber:    actor,
		Action:             action,
		EntityType:         "mobile_money_payment",
		EntityID:           payment.Number,
		NewValues:          map[string]interface{}{"fulfilment_status": payment.FulfilmentStatus},
	})
}

func truncate(value string) string {
	value = strings.TrimSpace(value)
	if len(value) > 500 {
		return value[:500]
	}
	return value
}
