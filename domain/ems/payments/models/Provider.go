package models

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type PaymentProvider struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	ProviderCode  string    `json:"provider_code" gorm:"uniqueIndex;size:32;not null"`
	ProviderName  string    `json:"provider_name" gorm:"size:80;not null"`
	CountryCode   string    `json:"country_code" gorm:"size:8;not null"`
	CurrencyCode  string    `json:"currency_code" gorm:"size:8;not null"`
	Environment   string    `json:"environment" gorm:"size:24;not null"`
	IsEnabled     bool      `json:"is_enabled" gorm:"not null;default:false"`
	ContractReady bool      `json:"contract_ready" gorm:"not null;default:false"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (PaymentProvider) TableName() string { return "ems_payment_providers" }

// EnsurePaymentProviders seeds non-secret provider metadata. It never stores
// credentials and never enables a provider; runtime enablement remains
// environment-controlled.
func EnsurePaymentProviders(db *gorm.DB) error {
	rows := []PaymentProvider{
		{
			ProviderCode: ProviderMTNMoMo, ProviderName: "MTN MoMo",
			CountryCode: "UG", CurrencyCode: "UGX", Environment: "sandbox",
			IsEnabled: false, ContractReady: true,
		},
		{
			ProviderCode: ProviderAirtelMoney, ProviderName: "Airtel Money",
			CountryCode: "UG", CurrencyCode: "UGX", Environment: "uat",
			IsEnabled: false, ContractReady: false,
		},
		{
			ProviderCode: ProviderSimulated, ProviderName: "Simulated",
			CountryCode: "UG", CurrencyCode: "UGX", Environment: "test",
			IsEnabled: false, ContractReady: true,
		},
	}
	for _, row := range rows {
		var existing PaymentProvider
		err := db.Where("provider_code = ?", row.ProviderCode).First(&existing).Error
		if err == nil {
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if err := db.Create(&row).Error; err != nil && !strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			return err
		}
	}
	return nil
}
