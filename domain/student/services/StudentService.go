package services

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/api/request"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/repositories"
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

	//familyStructureRepository := repositories.NewFamilyStructureRepository(config.DB)
	//familyStructure := payload.ToModel()
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
