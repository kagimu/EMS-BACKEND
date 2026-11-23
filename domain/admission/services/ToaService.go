package services

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/repositories"
)

func SaveToaService(toa models.TermsOfAdmission) error {
	repo := repositories.NewToaRepository(config.DB)
	return repo.UpsertToa(toa)
}

func GetToaService(schoolNumber string) (models.TermsOfAdmission, error) {
	repo := repositories.NewToaRepository(config.DB)
	return repo.FindToaBySchoolNumber(schoolNumber)
}
