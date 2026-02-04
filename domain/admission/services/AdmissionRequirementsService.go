package services

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/repositories"
)

func SaveAdmissionRequirementsService(temp models.AdmissionRequirement) error {
	repo := repositories.NewAdmissionRequirementsRepository(config.DB)
	if err := repo.SaveAdmissionRequirements(temp); err != nil {
		return err
	}

	return nil
}
