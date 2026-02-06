package services

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/repositories"
)

// SaveAdmissionRequirementsService creates one admission requirement per item_number and persists them.
func SaveAdmissionRequirementsService(req *requests.AdmissionRequirementsRequest, userNumber string) error {
	repo := repositories.NewAdmissionRequirementsRepository(config.DB)
	requirements := req.ToModels(userNumber)
	return repo.SaveAdmissionRequirementsBatch(requirements)
}

// FetchAdmissionRequirementsService returns admission requirements filtered by non-empty params.
func FetchAdmissionRequirementsService(schoolNumber, campusClassNumber, studyPeriodNumber, academicYearNumber, itemNumber string) ([]models.AdmissionRequirement, error) {
	repo := repositories.NewAdmissionRequirementsRepository(config.DB)
	return repo.FindAdmissionRequirements(schoolNumber, campusClassNumber, studyPeriodNumber, academicYearNumber, itemNumber)
}
