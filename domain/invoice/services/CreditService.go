package services

import (
	"errors"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/api/responses"
	invoiceRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/repositories"
)

func GetCreditsByUser(userNumber string) (responses.CreditReportResponse, error) {
	var out responses.CreditReportResponse
	if userNumber == "" {
		return out, errors.New("user_number is required")
	}

	credits, err := invoiceRepositories.FindCreditsByUserNumber(userNumber)
	if err != nil {
		return out, err
	}

	var creditResponses []responses.CreditResponse
	var totalCredits float64
	var netBalance float64

	for _, c := range credits {
		creditResponses = append(creditResponses, responses.CreditResponse{
			Number:                 c.Number,
			Date:                   c.Date.Format(time.DateOnly),
			Amount:                 c.Amount,
			OriginalAmountInCredit: c.OriginalAmountInCredit,
			Status:                 c.Status,
			ClientReference:        c.ClientReference,
			Description:            c.Description,
			SchoolNumber:           c.SchoolNumber,
			CampusNumber:           c.CampusNumber,
			UserNumber:             c.UserNumber,
			AddedByNumber:          c.AddedByNumber,
			CreatedAt:              c.CreatedAt.Format(time.RFC3339),
			UpdatedAt:              c.UpdatedAt.Format(time.RFC3339),
		})
		totalCredits += c.OriginalAmountInCredit
		netBalance += c.Amount
	}

	out = responses.CreditReportResponse{
		Credits:      creditResponses,
		NetBalance:   netBalance,
		TotalCredits: totalCredits,
	}

	return out, nil
}

func EditCredit(payload requests.EditCreditRequest, currentUser string) error {
	credit, err := invoiceRepositories.FindCreditByNumber(payload.CreditNumber)
	if err != nil {
		return errors.New("checking credit error: " + err.Error())
	}
	if credit.Number == "" {
		return errors.New("credit not found")
	}

	// Update fields if provided
	if payload.Description != nil {
		credit.Description = *payload.Description
	}

	if payload.ClientReference != nil {
		credit.ClientReference = *payload.ClientReference
	}

	if payload.Date != nil {
		parsedDate, err := commons.ConvertStringToDate(*payload.Date)
		if err != nil {
			return errors.New("invalid date format: " + err.Error())
		}
		credit.Date = parsedDate
	}

	if payload.Amount != nil {
		newAmount := *payload.Amount
		usedAmount := credit.OriginalAmountInCredit - credit.Amount
		if newAmount < usedAmount {
			return errors.New("cannot edit credit amount to less than what has already been used")
		}
		credit.OriginalAmountInCredit = newAmount
		credit.Amount = newAmount - usedAmount
	}

	_, err = invoiceRepositories.SaveCredit(credit)
	if err != nil {
		return errors.New("failed to save edited credit: " + err.Error())
	}

	return nil
}

func DeleteCredit(creditNumber string, currentUser string) error {
	credit, err := invoiceRepositories.FindCreditByNumber(creditNumber)
	if err != nil {
		return errors.New("checking credit error: " + err.Error())
	}
	if credit.Number == "" {
		return errors.New("credit not found")
	}

	// Check if the credit has been used
	if credit.Amount < credit.OriginalAmountInCredit {
		return errors.New("cannot delete credit that has already been partially or fully used")
	}

	err = invoiceRepositories.DeleteCredit(credit)
	if err != nil {
		return errors.New("failed to delete credit: " + err.Error())
	}

	return nil
}
