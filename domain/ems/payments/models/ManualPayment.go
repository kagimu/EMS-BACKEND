package models

import "time"

const (
	ManualMethodMTNMoMo      = "MTN_MOMO"
	ManualMethodBankTransfer = "BANK_TRANSFER"

	ManualStatusSubmitted = "SUBMITTED"
	ManualStatusApproved  = "APPROVED"
	ManualStatusRejected  = "REJECTED"
)

// ManualPaymentSubmission records an off-line payment claim and the uploaded
// proof. The proof file is private and is only served through the owner
// console. Approval does not itself issue an activation token; the customer
// claims the token after approval, keeping raw-token handling unchanged.
type ManualPaymentSubmission struct {
	ID                    uint       `json:"id" gorm:"primaryKey"`
	Number                string     `json:"number" gorm:"uniqueIndex;size:36;not null"`
	OrderNumber           string     `json:"order_number" gorm:"index;size:36;not null"`
	Purpose               string     `json:"purpose" gorm:"index;size:24;not null"`
	OrganisationNumber    string     `json:"organisation_number" gorm:"index;size:36"`
	CustomerName          string     `json:"customer_name" gorm:"size:160"`
	CustomerEmail         string     `json:"customer_email" gorm:"size:160"`
	CustomerPhone         string     `json:"customer_phone" gorm:"size:32"`
	RequesterName         string     `json:"requester_name" gorm:"size:160;not null"`
	Method                string     `json:"method" gorm:"index;size:24;not null"`
	AmountUGX             int64      `json:"amount_ugx"`
	Currency              string     `json:"currency" gorm:"size:8;default:UGX"`
	PayerName             string     `json:"payer_name" gorm:"size:160"`
	PayerPhone            string     `json:"payer_phone" gorm:"size:32"`
	PayerReference        string     `json:"payer_reference" gorm:"size:160"`
	CustomerNote          string     `json:"customer_note" gorm:"type:text"`
	Status                string     `json:"status" gorm:"index;size:24;not null"`
	ProofPath             string     `json:"-" gorm:"type:text"`
	ProofName             string     `json:"proof_name" gorm:"size:255"`
	ProofMime             string     `json:"proof_mime" gorm:"size:100"`
	ProofSize             int64      `json:"proof_size"`
	SubmittedAt           time.Time  `json:"submitted_at"`
	ReviewedAt            *time.Time `json:"reviewed_at"`
	ReviewedByNumber      string     `json:"reviewed_by_number" gorm:"size:36"`
	ReviewNote            string     `json:"review_note" gorm:"type:text"`
	TokenClaimedAt        *time.Time `json:"token_claimed_at"`
	ActivationTokenNumber string     `json:"activation_token_number" gorm:"size:36"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

func (ManualPaymentSubmission) TableName() string { return "ems_manual_payment_submissions" }
