package models

import (
	"time"

	"gorm.io/gorm"
)

type (
	QBSyncQueue struct {
		gorm.Model
		Number            string     `json:"number" gorm:"size:36;unique;not null"`
		TransactionNumber string     `json:"transaction_number" gorm:"size:36;index"`
		ActionType        string     `json:"action_type" gorm:"size:50;not null;default:'SALES_RECEIPT_ADD'"`
		Status            string     `json:"status" gorm:"size:50;not null;default:'QUEUED'"`
		QbXmlPayload      string     `json:"qbxml_payload" gorm:"type:text"`
		ErrorMessage      string     `json:"error_message" gorm:"type:text"`
		SchoolNumber      string     `json:"school_number" gorm:"size:36"`
		UserNumber        string     `json:"user_number" gorm:"size:36"`
		ProcessedAt       *time.Time `json:"processed_at"`
	}
)

func (u *QBSyncQueue) TableName() string {
	return "qb_sync_queue"
}
