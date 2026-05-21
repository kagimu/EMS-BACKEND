package services

import (
	"errors"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/models"
	invoiceRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/repositories"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// DepositPrepaid records a deposit (money brought in) and adds it to the user's prepaid balance.
func DepositPrepaid(payload requests.AccountTopUpRequest, addedByNumber string) (responses.PrepaidBalanceResponse, error) {
	var out responses.PrepaidBalanceResponse
	datePaid, err := time.Parse(time.DateOnly, payload.DatePaid)
	if err != nil {
		// try common alternative
		if parsed, e := commons.ConvertStringToDate(payload.DatePaid); e == nil {
			datePaid = parsed
		} else {
			return out, errors.New("invalid date_paid: use YYYY-MM-DD")
		}
	}
	desc := ""
	if payload.Description != nil {
		desc = *payload.Description
	}
	clientRef := ""
	if payload.ClientReference != nil {
		clientRef = *payload.ClientReference
	}
	deposit := models.PrepaidDeposit{
		Number:          uuid.NewV4().String(),
		UserNumber:      payload.UserNumber,
		Amount:          payload.Amount,
		DatePaid:        datePaid,
		Description:     desc,
		ClientReference: clientRef,
		AddedByNumber:   addedByNumber,
		CreatedAt:       time.Now(),
	}
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		if e := invoiceRepositories.SavePrepaidDeposit(deposit, tx); e != nil {
			return e
		}
		return invoiceRepositories.EnsurePrepaidBalanceRow(payload.UserNumber, payload.Amount, tx)
	})
	if err != nil {
		return out, err
	}
	bal, _ := invoiceRepositories.FindPrepaidBalanceByUserNumber(payload.UserNumber)
	out.UserNumber = bal.UserNumber
	out.Balance = bal.Balance
	return out, nil
}

// GetPrepaidBalance returns the current prepaid balance for a user, optionally with deposit history.
func GetPrepaidBalance(userNumber string, includeDeposits bool) (responses.PrepaidBalanceWithDepositsResponse, error) {
	var out responses.PrepaidBalanceWithDepositsResponse
	if userNumber == "" {
		return out, errors.New("user_number is required")
	}
	bal, err := invoiceRepositories.FindPrepaidBalanceByUserNumber(userNumber)
	if err != nil {
		return out, err
	}
	out.UserNumber = bal.UserNumber
	out.Balance = bal.Balance
	if includeDeposits {
		deposits, err := invoiceRepositories.FindDepositsByUserNumber(userNumber)
		if err != nil {
			return out, err
		}
		for _, d := range deposits {
			out.Deposits = append(out.Deposits, responses.PrepaidDepositResponse{
				Number:          d.Number,
				UserNumber:      d.UserNumber,
				Amount:          d.Amount,
				DatePaid:        d.DatePaid.Format(time.DateOnly),
				Description:     d.Description,
				ClientReference: d.ClientReference,
				CreatedAt:       d.CreatedAt.Format(time.RFC3339),
			})
		}
	}
	return out, nil
}
