package services

import (
	"errors"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	repositories3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/models"
	invoiceRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/repositories"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func CreateCredit(payload requests.CreateCreditRequest, addedByNumber string) (responses.CreditResponse, error) {
	var out responses.CreditResponse

	// Verify the student exists
	student, err := repositories.FindStudentByUserNumber(payload.UserNumber)
	if err != nil {
		return out, err
	}
	if student.Number == "" {
		return out, errors.New("selected user is not a learner")
	}

	campus, err := repositories3.GetCampusByCampusClassNumber(student.CampusClassNumber)
	if err != nil {
		return out, errors.New("determine campus error: " + err.Error())
	}
	if campus.Number == "" {
		return out, errors.New("could not determine campus for student: " + student.Surname)
	}

	date, err := time.Parse(time.DateOnly, payload.Date)
	if err != nil {
		if parsed, e := commons.ConvertStringToDate(payload.Date); e == nil {
			date = parsed
		} else {
			return out, errors.New("invalid date: use YYYY-MM-DD")
		}
	}

	desc := ""
	if payload.Description != nil {
		desc = *payload.Description
	}

	clientRef := ""
	if payload.ClientReference != nil {
		clientRef = *payload.ClientReference
		existingCredit, err := invoiceRepositories.GetCreditByClientReference(clientRef, campus.Number)
		if err != nil {
			return out, errors.New("checking client reference error: " + err.Error())
		}
		if existingCredit.Number != "" {
			return out, errors.New("client reference '" + clientRef + "' already used!")
		}
	}

	credit := models.Credit{
		Number:          uuid.NewV4().String(),
		UserNumber:      payload.UserNumber,
		SchoolNumber:    student.SchoolNumber,
		CampusNumber:    campus.Number,
		Amount:          payload.Amount,
		Balance:         payload.Amount, // initially the balance is same as amount
		Date:            date,
		Status:          constants.StatusActive,
		Description:     desc,
		ClientReference: clientRef,
		AddedByNumber:   addedByNumber,
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		_, err := invoiceRepositories.SaveCredit(credit, tx)
		return err
	})

	if err != nil {
		return out, err
	}

	out = responses.CreditResponse{
		Number:          credit.Number,
		Date:            credit.Date.Format(time.DateOnly),
		Amount:          credit.Amount,
		Balance:         credit.Balance,
		Status:          credit.Status,
		ClientReference: credit.ClientReference,
		Description:     credit.Description,
		SchoolNumber:    credit.SchoolNumber,
		CampusNumber:    credit.CampusNumber,
		UserNumber:      credit.UserNumber,
		AddedByNumber:   credit.AddedByNumber,
		CreatedAt:       credit.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       credit.UpdatedAt.Format(time.RFC3339),
	}

	return out, nil
}

func GetCreditsByUser(userNumber string) ([]responses.CreditResponse, error) {
	var out []responses.CreditResponse
	if userNumber == "" {
		return out, errors.New("user_number is required")
	}

	credits, err := invoiceRepositories.FindCreditsByUserNumber(userNumber)
	if err != nil {
		return out, err
	}

	for _, c := range credits {
		out = append(out, responses.CreditResponse{
			Number:          c.Number,
			Date:            c.Date.Format(time.DateOnly),
			Amount:          c.Amount,
			Balance:         c.Balance,
			Status:          c.Status,
			ClientReference: c.ClientReference,
			Description:     c.Description,
			SchoolNumber:    c.SchoolNumber,
			CampusNumber:    c.CampusNumber,
			UserNumber:      c.UserNumber,
			AddedByNumber:   c.AddedByNumber,
			CreatedAt:       c.CreatedAt.Format(time.RFC3339),
			UpdatedAt:       c.UpdatedAt.Format(time.RFC3339),
		})
	}

	return out, nil
}
