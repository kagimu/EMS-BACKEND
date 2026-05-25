package services

import (
	"errors"
	"time"

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
