package services

import (
	"errors"
	userModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	uuid "github.com/satori/go.uuid"
)

type UnitOfMeasureService interface {
	NewUnitOfMeasure(unitOfMeasureRequest requests.UnitOfMeasureRequest) error
	UpdateUnitOfMeasure(editUnitOfMeasureRequest requests.EditUnitOfMeasureRequest) error
	FindUnitOfMeasures(schoolNumber string) ([]responses.UnitOfMeasureResponse, error)
	RemoveUnitOfMeasures(unitOfMeasureNumber string) error
}

var unitOfMeasureRepo = repositories.NewUnitOfMeasureRepository()

type unitOfMeasureService struct {
	User *userModel.User
}

func (u *unitOfMeasureService) RemoveUnitOfMeasures(unitOfMeasureNumber string) error {
	unitOfMeasure, err := unitOfMeasureRepo.GetUnitOfMeasureByNumber(unitOfMeasureNumber)
	if err != nil {
		return errors.New("checking unit of measure number error: " + err.Error())
	}

	if unitOfMeasure.Number == "" {
		return errors.New("unit of measure not found")
	}

	err = unitOfMeasureRepo.RemoveUnitOfMeasure(unitOfMeasure)
	if err != nil {
		return errors.New("removing unit of measure error: " + err.Error())
	}

	return nil
}

func (u *unitOfMeasureService) FindUnitOfMeasures(schoolNumber string) ([]responses.UnitOfMeasureResponse, error) {

	var unitOfMeasureResponses = make([]responses.UnitOfMeasureResponse, 0)

	unitOfMeasures, err := unitOfMeasureRepo.FindUnitOfMeasuresBySchoolNumber(schoolNumber)
	if err != nil {
		return unitOfMeasureResponses, errors.New("Unable to find unit of measures: " + err.Error())
	}

	for _, unitOfMeasure := range unitOfMeasures {
		unitOfMeasureResponses = append(unitOfMeasureResponses, createUnitOfMeasureResponse(unitOfMeasure))
	}

	return unitOfMeasureResponses, nil
}

func createUnitOfMeasureResponse(unitOfMeasure models.UnitOfMeasure) responses.UnitOfMeasureResponse {
	return responses.UnitOfMeasureResponse{
		UnitOfMeasureNumber: unitOfMeasure.Number,
		Name:                unitOfMeasure.Name,
		Code:                unitOfMeasure.Code,
		Abbreviation:        unitOfMeasure.Abbreviation,
		Status:              unitOfMeasure.Status,
		SchoolNumber:        unitOfMeasure.SchoolNumber,
	}
}

func (u *unitOfMeasureService) NewUnitOfMeasure(unitOfMeasureRequest requests.UnitOfMeasureRequest) error {

	unitOfMeasure, err := unitOfMeasureRepo.FindUnitOfMeasureByCodeAndSchoolNumber(unitOfMeasureRequest.Code, unitOfMeasureRequest.SchoolNumber)
	if err != nil {
		return err
	}

	if unitOfMeasure.Number != "" {
		return errors.New("unit of measure with same code already exists")
	}

	unitOfMeasure, err = unitOfMeasureRepo.FindUnitOfMeasureByNameAndSchoolNumber(unitOfMeasureRequest.Name, unitOfMeasureRequest.SchoolNumber)
	if err != nil {
		return err
	}

	if unitOfMeasure.Number != "" {
		return errors.New("unit of measure with same name already exists")
	}

	unitOfMeasure = models.UnitOfMeasure{
		Number:       uuid.NewV4().String(),
		Name:         unitOfMeasureRequest.Name,
		SchoolNumber: unitOfMeasureRequest.SchoolNumber,
		Code:         unitOfMeasureRequest.Code,
		Abbreviation: unitOfMeasureRequest.Abbreviation,
		UserNumber:   u.User.Number,
		Status:       constants.StatusActive,
	}

	err = unitOfMeasureRepo.CreateUnitOfMeasure(unitOfMeasure)
	if err != nil {
		return errors.New("save unit of measure fail: " + err.Error())
	}

	return nil
}

func (u *unitOfMeasureService) UpdateUnitOfMeasure(editUnitOfMeasureRequest requests.EditUnitOfMeasureRequest) error {

	baseUnitOfMeasure, err := unitOfMeasureRepo.GetUnitOfMeasureByNumber(editUnitOfMeasureRequest.UnitOfMeasureNumber)
	if err != nil {
		return err
	}

	if baseUnitOfMeasure.Number == "" {
		return errors.New("invalid unit of measure number")
	}

	unitOfMeasure, err := unitOfMeasureRepo.FindUnitOfMeasureByCodeAndSchoolNumber(editUnitOfMeasureRequest.Code, editUnitOfMeasureRequest.SchoolNumber)
	if err != nil {
		return err
	}

	if unitOfMeasure.Number != "" && unitOfMeasure.Number != baseUnitOfMeasure.Number {
		return errors.New("unit of measure with same code already exists")
	}

	unitOfMeasure, err = unitOfMeasureRepo.FindUnitOfMeasureByNameAndSchoolNumber(editUnitOfMeasureRequest.Name, editUnitOfMeasureRequest.SchoolNumber)
	if err != nil {
		return err
	}

	if unitOfMeasure.Number != "" && unitOfMeasure.Number != baseUnitOfMeasure.Number {
		return errors.New("unit of measure with same name already exists")
	}

	unitOfMeasure.Name = editUnitOfMeasureRequest.Name
	unitOfMeasure.SchoolNumber = editUnitOfMeasureRequest.SchoolNumber
	unitOfMeasure.Code = editUnitOfMeasureRequest.Code
	unitOfMeasure.Abbreviation = editUnitOfMeasureRequest.Abbreviation
	unitOfMeasure.UserNumber = u.User.Number
	unitOfMeasure.Status = editUnitOfMeasureRequest.Status

	err = unitOfMeasureRepo.UpdateUnitOfMeasure(unitOfMeasure)
	if err != nil {
		return errors.New("save unit of measure fail: " + err.Error())
	}

	return nil
}

func NewUnitOfMeasureService(user *userModel.User) UnitOfMeasureService {
	return &unitOfMeasureService{
		User: user,
	}
}
