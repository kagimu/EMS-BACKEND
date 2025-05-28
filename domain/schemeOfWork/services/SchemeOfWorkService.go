package services

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/schemeOfWork/repositories"
)

type SchemeOfWorkService struct {
}

func (s SchemeOfWorkService) InvalidateSchemeOfWork(schemeOfWorkNumber, userNumber string) error {
	schemeOfWork, err := repositories.GetSchemeOfWorkByNumber(schemeOfWorkNumber)
	if err != nil {
		return err
	}

	if schemeOfWork.Number == "" {
		return errors.New("Invalid schemeOfWork number ")
	}

	if schemeOfWork.Status == constants.StatusInactive {
		return errors.New("This scheme is already invalidated ")
	}

	schemeOfWork.Status = constants.StatusInactive
	schemeOfWork.UserNumber = userNumber

	err = repositories.SaveSchemeOfWork(schemeOfWork)
	if err != nil {
		return err
	}

	return nil
}

func (s SchemeOfWorkService) ActivateSchemeOfWork(schemeOfWorkNumber, userNumber string) error {
	schemeOfWork, err := repositories.GetSchemeOfWorkByNumber(schemeOfWorkNumber)
	if err != nil {
		return err
	}

	if schemeOfWork.Number == "" {
		return errors.New("Invalid schemeOfWork number ")
	}

	if schemeOfWork.Status == constants.StatusActive {
		return errors.New("This scheme is already active ")
	}

	schemeOfWork.Status = constants.StatusActive
	schemeOfWork.UserNumber = userNumber

	err = repositories.SaveSchemeOfWork(schemeOfWork)
	if err != nil {
		return err
	}

	return nil
}

func InitSchemeOfWorkService() *SchemeOfWorkService {
	return &SchemeOfWorkService{}
}
