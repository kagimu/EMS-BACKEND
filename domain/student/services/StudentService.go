package services

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/api/request"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/repositories"
	uuid "github.com/satori/go.uuid"
)

func CreateMedicalRecord(req request.MedicalRecordRequest) error {
	medicalRecordRepository := repositories.NewMedicalRecordRepository(config.DB)
	medicalRecord := req.ToModel()
	if err := medicalRecordRepository.CreateMedicalRecord(medicalRecord); err != nil {
		return err
	}
	return nil
}

func CreateFamilyStructure(payload request.FamilyStructureRequest) error {

	tx := config.DB.Begin()
	familyStructureRepository := repositories.NewFamilyMemberRepository(tx)
	languageRepository := repositories.NewLanguageRepository(tx)

	for _, familyMember := range payload.FamilyMembers {
		var familyMemberModel models.FamilyMember
		familyMemberModel.Number = uuid.NewV4().String()
		familyMemberModel.StudentNumber = payload.StudentNumber
		familyMemberModel.Surname = familyMember.Surname
		familyMemberModel.OtherName = familyMember.OtherName
		familyMemberModel.PlaceOfEmployment = familyMember.PlaceOfEmployment
		familyMemberModel.Relationship = familyMember.Relationship

		if err := familyStructureRepository.CreateFamilyMember(familyMemberModel); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, language := range payload.Languages {
		var languageModel models.StudentLanguage
		languageModel.StudentNumber = payload.StudentNumber
		languageModel.Name = language.Name
		languageModel.Code = language.Code

		if err := languageRepository.CreateLanguage(languageModel); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func CreateEmergencyInformation(payload request.EmergencyInformationRequest) error {
	emergencyInformationRepository := repositories.NewEmergencyInformationRepository(config.DB)
	emergencyInformation := payload.ToModel()
	if err := emergencyInformationRepository.CreateEmergencyInformation(emergencyInformation); err != nil {
		return err
	}
	return nil
}
