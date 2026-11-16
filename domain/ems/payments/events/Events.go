package events

import (
	"encoding/json"
	"strings"

	paymentModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/payments/models"
	"gorm.io/gorm"
)

// Record appends a safe payment event. Metadata must contain no credentials,
// PINs, access tokens or raw personal data.
func Record(tx *gorm.DB, transaction, eventType, provider, providerReference, actor string, metadata map[string]interface{}) error {
	encoded := ""
	if len(metadata) > 0 {
		raw, err := json.Marshal(metadata)
		if err != nil {
			return err
		}
		encoded = string(raw)
	}
	return tx.Create(&paymentModels.PaymentEvent{
		TransactionNumber: transaction,
		EventType:         strings.ToUpper(strings.TrimSpace(eventType)),
		ProviderCode:      strings.ToUpper(strings.TrimSpace(provider)),
		ProviderReference: strings.TrimSpace(providerReference),
		ActorUserNumber:   strings.TrimSpace(actor),
		Metadata:          encoded,
	}).Error
}
