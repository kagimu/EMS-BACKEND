package services

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/models"
)

func GetApplicationById(id string) (responses.ApplicationResponse, error) {

	applicationRepo := NewApplicationRepositoryImpl(config.DB)
	bankRepo := repo.NewBankRepositoryImpl(config.DB)
	application, err := applicationRepo.FindById(id)

	if err != nil {
		return responses.ApplicationResponse{}, err
	}

	bank, err := bankRepo.FindById(application.BankId)
	if err != nil {
		return responses.ApplicationResponse{}, err

	}

	return BuildSingleApplicationResponse(application, bank), nil
}

func BuildSingleApplicationResponse(application *models.Application) responses.ApplicationResponse {
	var studentResponse responses.StudentApplicant
	var parentResponse responses.ParentResponse
	var previousSchoolsResponse []responses.PreviousSchoolResponse
	return responses.ApplicationResponse{
		Student:         studentResponse,
		Parent:          parentResponse,
		PreviousSchools: previousSchoolsResponse,
	}
}
