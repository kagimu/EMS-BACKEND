package services

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/repositories"
)

func SaveAdmissionLetterTemplateService(temp models.AdmissionLetterTemplate) error {
	repo := repositories.NewAdmissionTemplateRepository(config.DB)
	if err := repo.SaveAdmissionTemplate(temp); err != nil {
		return err
	}

	return nil
}

func GetAdmissionLetterTemplateService(schoolNumber string) (models.AdmissionLetterTemplate, error) {
	repo := repositories.NewAdmissionTemplateRepository(config.DB)
	admissionTemplate, err := repo.FindAdmissionTemplateBySchoolNumber(schoolNumber)
	if err != nil {
		return models.AdmissionLetterTemplate{}, err
	}
	return admissionTemplate, nil
}

func DeleteAdmissionLetterTemplateService(templateNumber string) error {
	repo := repositories.NewAdmissionTemplateRepository(config.DB)
	if err := repo.DeleteAdmissionTemplateByNumber(templateNumber); err != nil {
		return err
	}
	return nil
}
