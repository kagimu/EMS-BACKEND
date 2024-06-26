package services

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/repositories"
	uuid "github.com/satori/go.uuid"
)

func SaveGradingStandard(gradingStandardRequest requests.GradingStandardRequest, currentUser string) error {

	byCode, err := repositories.FindGradingStandardByCodeAndClassNumber(gradingStandardRequest.Code, gradingStandardRequest.CampusClassNumber)
	if err != nil {
		return errors.New("Checking code error: " + err.Error())
	}

	if byCode.Number != "" {
		return errors.New("code already used: " + byCode.Number)
	}

	byName, err := repositories.FindGradingStandardByNameAndClassNumber(gradingStandardRequest.Name, gradingStandardRequest.CampusClassNumber)
	if err != nil {
		return errors.New("Checking by name error: " + err.Error())
	}

	if byName.Number != "" {
		return errors.New("name already used: " + byName.Number)
	}

	gs := models.GradingStandard{
		Code:              gradingStandardRequest.Code,
		Number:            uuid.NewV4().String(),
		Name:              gradingStandardRequest.Name,
		Description:       gradingStandardRequest.Description,
		Type:              gradingStandardRequest.Type,
		CampusClassNumber: gradingStandardRequest.CampusClassNumber,
		UserNumber:        currentUser,
		CampusNumber:      gradingStandardRequest.CampusNumber,
	}

	err = repositories.SaveGradingStandard(gs)
	if err != nil {
		return errors.New("Saving grading standard error: " + err.Error())
	}

	return nil
}

func UpdateGradingStandard(gradingStandardRequest requests.UpdateGradingStandardRequest, currentUser string) error {

	gradingStandard, err := repositories.GetGradingStandardByNumber(gradingStandardRequest.GradingStandardNumber)
	if err != nil {
		return errors.New("Checking grading standard error: " + err.Error())
	}

	if gradingStandard.Number == "" {
		return errors.New("grading standard number is empty: " + gradingStandardRequest.GradingStandardNumber)
	}

	byCode, err := repositories.FindGradingStandardByCodeAndClassNumber(gradingStandardRequest.Code, gradingStandardRequest.CampusClassNumber)
	if err != nil {
		return errors.New("Checking code error: " + err.Error())
	}

	if byCode.Number != "" && byCode.Number != gradingStandard.Number {
		return errors.New("code already used: " + byCode.Number)
	}

	byName, err := repositories.FindGradingStandardByNameAndClassNumber(gradingStandardRequest.Name, gradingStandardRequest.CampusClassNumber)
	if err != nil {
		return errors.New("Checking by name error: " + err.Error())
	}

	if byName.Number != "" && byCode.Number != gradingStandard.Number {
		return errors.New("name already used: " + byName.Number)
	}

	gs := models.GradingStandard{
		Code:              gradingStandardRequest.Code,
		Number:            uuid.NewV4().String(),
		Name:              gradingStandardRequest.Name,
		Description:       gradingStandardRequest.Description,
		Type:              gradingStandardRequest.Type,
		CampusClassNumber: gradingStandardRequest.CampusClassNumber,
		UserNumber:        currentUser,
		CampusNumber:      gradingStandardRequest.CampusNumber,
	}

	err = repositories.SaveGradingStandard(gs)
	if err != nil {
		return errors.New("Saving grading standard error: " + err.Error())
	}

	return nil
}
