package services

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/api/request"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/repositories"
	uuid "github.com/satori/go.uuid"
)

func CreateMedicalRecord(req request.MedicalRecordRequest) error {
	db := config.DB.Begin()
	medicalRecordRepository := repositories.NewMedicalRecordRepository(db)
	err := medicalRecordRepository.DeleteMedicalRecordByStudentNumber(req.StudentNumber)
	if err != nil {
		db.Rollback()
		return err
	}
	medicalRecord := req.ToModel()
	if err = medicalRecordRepository.CreateMedicalRecord(medicalRecord); err != nil {
		db.Rollback()
		return err
	}
	if err := db.Commit().Error; err != nil {
		db.Rollback()
	}
	return nil
}

func CreateFamilyStructure(payload request.FamilyStructureRequest) error {

	tx := config.DB.Begin()
	familyStructureRepository := repositories.NewFamilyMemberRepository(tx)
	languageRepository := repositories.NewLanguageRepository(tx)

	err := familyStructureRepository.DeleteFamilyMemberByStudentNumber(payload.StudentNumber)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = languageRepository.DeleteLanguageByStudentNumber(payload.StudentNumber)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, familyMember := range payload.FamilyMembers {
		var familyMemberModel models.FamilyMember
		familyMemberModel.Number = uuid.NewV4().String()
		familyMemberModel.StudentNumber = payload.StudentNumber
		familyMemberModel.Surname = familyMember.Surname
		familyMemberModel.OtherName = familyMember.OtherName
		familyMemberModel.PlaceOfEmployment = familyMember.PlaceOfEmployment
		familyMemberModel.Relationship = familyMember.Relationship
		familyMemberModel.PhysicalContactAddress = familyMember.PhysicalContactAddress
		familyMemberModel.IsLiving = familyMember.IsLiving

		if err := familyStructureRepository.CreateFamilyMember(familyMemberModel); err != nil {
			tx.Rollback()
			return err
		}

		for _, contact := range familyMember.Contacts {
			var familyMemberContactModel models.FamilyMemberContact
			familyMemberContactModel.PhoneNumber = contact.PhoneNumber
			familyMemberContactModel.Email = contact.Email
			familyMemberContactModel.FamilyMemberNumber = familyMemberModel.Number

			if err := familyStructureRepository.CreateFamilyMemberContact(familyMemberContactModel); err != nil {
				tx.Rollback()
				return err
			}
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
	tx := config.DB.Begin()
	emergencyInformationRepository := repositories.NewEmergencyInformationRepository(tx)
	emergencyInformation := models.EmergencyInformation{
		StudentNumber:                 payload.StudentNumber,
		ContactPerson:                 payload.ContactPerson,
		ContactPersonPhone:            payload.ContactPersonPhone,
		ContactPersonAlternative:      payload.ContactPersonAlternative,
		ContactPersonAlternativePhone: payload.ContactPersonAlternativePhone,
		ContactDoctorName:             payload.ContactDoctorName,
		ContactDoctorPhone:            payload.ContactDoctorPhone,
	}
	if err := emergencyInformationRepository.DeleteEmergencyInformationByStudentNumber(payload.StudentNumber); err != nil {
		tx.Rollback()
		return err
	}
	if err := emergencyInformationRepository.CreateEmergencyInformation(emergencyInformation); err != nil {
		tx.Rollback()
		return err
	}

	err := tx.Commit().Error

	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func FetchEmergencyInformation(studentNumber string) (responses.EmergencyInformationResponse, error) {

	emergencyInformationRepository := repositories.NewEmergencyInformationRepository(config.DB)
	emergencyInformation, err := emergencyInformationRepository.GetEmergencyInformationByStudentNumber(studentNumber)
	if err != nil {
		return responses.EmergencyInformationResponse{}, err
	}
	resp := responses.EmergencyInformationResponse{
		StudentNumber:                 emergencyInformation.StudentNumber,
		ContactPerson:                 emergencyInformation.ContactPerson,
		ContactPersonPhone:            emergencyInformation.ContactPersonPhone,
		ContactPersonAlternative:      emergencyInformation.ContactPersonAlternative,
		ContactPersonAlternativePhone: emergencyInformation.ContactPersonAlternativePhone,
		ContactDoctorName:             emergencyInformation.ContactDoctorName,
		ContactDoctorPhone:            emergencyInformation.ContactDoctorPhone,
	}
	return resp, nil
}

func FetchMedicalRecord(studentNumber string) (responses.MedicalRecordResponse, error) {
	medicalRecordRepository := repositories.NewMedicalRecordRepository(config.DB)
	medicalRecord, err := medicalRecordRepository.GetMedicalRecordByStudentNumber(studentNumber)
	var medicalRecordResponse responses.MedicalRecordResponse
	if err != nil {
		return medicalRecordResponse, err
	}
	medicalRecordResponse = responses.MedicalRecordResponse{
		StudentNumber:     medicalRecord.StudentNumber,
		Allergies:         medicalRecord.Allergies,
		DoctorName:        medicalRecord.DoctorName,
		DoctorPhone:       medicalRecord.DoctorPhone,
		RegularMedication: medicalRecord.RegularMedication,
		ExistingCondition: medicalRecord.ExistingCondition,
		BloodTypeGroup:    medicalRecord.BloodTypeGroup,
	}
	return medicalRecordResponse, nil
}

func FetchFamilyStructure(studentNumber string) (responses.FamilyStructureResponse, error) {

	familyMemberRepository := repositories.NewFamilyMemberRepository(config.DB)
	familyMembers, err := familyMemberRepository.FindFamilyMembersByStudentNumber(studentNumber)
	if err != nil {
		return responses.FamilyStructureResponse{}, err
	}
	familyMemberResponses := make([]responses.FamilyMemberResponse, 0)
	for _, member := range familyMembers {
		var familyMemberResponse responses.FamilyMemberResponse
		familyMemberResponse.Number = member.Number
		familyMemberResponse.Surname = member.Surname
		familyMemberResponse.OtherName = member.OtherName
		familyMemberResponse.PlaceOfEmployment = member.PlaceOfEmployment
		familyMemberResponse.PhysicalContactAddress = member.PhysicalContactAddress
		familyMemberResponse.IsLiving = member.IsLiving
		//familyMemberResponse.Contacts = member.Contacts
		familyMemberResponse.Relationship = member.Relationship
		contacts, err := familyMemberRepository.FindFamilyMemberContactsByFamilyMemberNumber(member.Number)
		if err != nil {
			return responses.FamilyStructureResponse{}, err
		}
		familyMemberResponse.Contacts = make([]models.FamilyMemberContact, 0)
		for _, contact := range contacts {
			familyMemberResponse.Contacts = append(familyMemberResponse.Contacts, models.FamilyMemberContact{
				PhoneNumber: contact.PhoneNumber,
				Email:       contact.Email,
			})
		}

		familyMemberResponses = append(familyMemberResponses, familyMemberResponse)
	}

	languageRepository := repositories.NewLanguageRepository(config.DB)
	languages, err := languageRepository.FindLanguagesByStudentNumber(studentNumber)

	if err != nil {
		return responses.FamilyStructureResponse{}, err
	}

	languageResponses := make([]responses.Language, 0)

	for _, lang := range languages {
		languageResponses = append(languageResponses, responses.Language{
			Name: lang.Name,
			Code: lang.Code,
		})
	}

	response := responses.FamilyStructureResponse{
		StudentNumber: studentNumber,
		FamilyMembers: familyMemberResponses,
		Languages:     languageResponses,
	}

	return response, nil
}
