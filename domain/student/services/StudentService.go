package services

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	applicationModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/repository"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/api/request"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/repositories"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
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

		if len(familyMember.OtherName) > 50 {
			tx.Rollback()
			return errors.New(familyMember.OtherName + " other name is  too long maximum 50 characters")
		}

		if len(familyMember.Surname) > 50 {
			tx.Rollback()
			return errors.New(familyMember.Surname + " surname is  too long maximum 50 characters")
		}
		if len(familyMember.PhysicalContactAddress) > 100 {
			tx.Rollback()
			return errors.New(familyMember.PhysicalContactAddress + " Physical Contact Address is  too long maximum 50 characters")
		}

		if len(familyMember.PlaceOfEmployment) > 100 {
			tx.Rollback()
			return errors.New(familyMember.PlaceOfEmployment + "place of employment is  too long maximum 50 characters")
		}

		if err := familyStructureRepository.CreateFamilyMember(familyMemberModel); err != nil {
			tx.Rollback()
			return err
		}

		for _, contact := range familyMember.Contacts {
			var familyMemberContactModel models.FamilyMemberContact
			familyMemberContactModel.Number = uuid.NewV4().String()
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
		languageModel.Number = uuid.NewV4().String()
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
		familyMemberResponse.Contacts = make([]responses.FamilyMemberContactResponse, 0)
		for _, contact := range contacts {
			familyMemberResponse.Contacts = append(familyMemberResponse.Contacts, responses.FamilyMemberContactResponse{
				Number:             contact.Number,
				PhoneNumber:        contact.PhoneNumber,
				Email:              contact.Email,
				FamilyMemberNumber: contact.FamilyMemberNumber,
				CreatedAt:          contact.CreatedAt.Format(time.DateTime),
				UpdatedAt:          contact.UpdatedAt.Format(time.DateTime),
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
			Name:   lang.Name,
			Code:   lang.Code,
			Number: lang.Number,
		})
	}

	response := responses.FamilyStructureResponse{
		StudentNumber: studentNumber,
		FamilyMembers: familyMemberResponses,
		Languages:     languageResponses,
	}

	return response, nil
}

func CreateNewStudentAndApplication(payload request.StudentRequest, userNumber string) (string, error) {
	campusClass, err := services.ResolveCampusClass(payload.ClassNumber, payload.CampusNumber, payload.CurriculumNumber, payload.CampusCurriculumNumber)
	if err != nil {
		return "", errors.New("Error checking class and campus: " + err.Error())
	}

	if campusClass.Number == "" {
		return "", errors.New("campus class is not active on the selected campus")
	}

	/*var campusClass classModels.CampusClass
	if err := config.DB.Where("number = ?", campusClass.Number).First(&campusClass).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Campus class not found"})
		return
	}*/

	tx := config.DB.Begin()

	applicationRepo := repository.NewApplicationRepositoryImpl(tx)

	saveApplication := applicationModel.Application{
		Number:     uuid.NewV4().String(),
		UserNumber: userNumber,
		Type:       payload.Type,
	}

	if payload.ApplicationNumber == "" {
		if err := applicationRepo.SaveApplication(saveApplication); err != nil {
			tx.Rollback()
			return "", errors.New("failed to create application " + err.Error())
		}

	}

	studentAPN, err := CreateAPN(tx)
	if err != nil {
		tx.Rollback()
		return "", errors.New("failed to create APN " + err.Error())
	}

	studentNumber := uuid.NewV4().String()

	// find Student by application number

	student, err := applicationRepo.GetStudentByApplicationNumber(saveApplication.Number)
	if err != nil {
		tx.Rollback()
		return "", errors.New("failed to get student " + err.Error())
	}

	// NOTE: To avoid duplicate records, check if the student number is empty else, update the existing student
	if student.Number != "" {

		// update student details
		updateStudent := models.Student{
			Surname:                   payload.Surname,
			OtherName:                 payload.OtherName,
			Gender:                    payload.Gender,
			DateOfBirth:               payload.DateOfBirth,
			ADNStatus:                 payload.ADNStatus,
			ReasonForChoosingSchool:   payload.ReasonForChoosingSchool,
			PreferredTypeOfEnrollment: payload.PreferredTypeOfEnrollment,
			SecondName:                payload.SecondName,
			PreferredCampusNumber:     payload.PreferredCampusNumber,
			SourceOfSchoolInfo:        payload.SourceOfSchoolInfo,
			ApplicationNumber:         saveApplication.Number,
			CampusClassNumber:         campusClass.Number,
			SchoolNumber:              payload.SchoolNumber,
			Lin:                       payload.Lin,
			ModeOfStudyNumber:         payload.ModeOfStudyNumber,
			AddedBy:                   userNumber,
		}

		if payload.ADN != "" {
			updateStudent.ADN = payload.ADN
		}

		if err := applicationRepo.UpdateStudentByApplicationNumber(payload.ApplicationNumber, updateStudent); err != nil {
			return "", errors.New("failed to update student " + err.Error())
		}

		updateUser := authModel.User{
			FirstName:  payload.SecondName,
			LastName:   payload.Surname,
			SecondName: payload.OtherName,
		}

		if err := tx.Where("number = ?", student.UserNumber).Updates(updateUser).Error; err != nil {
			return "", errors.New("failed to update user")
		}

		if err := tx.Commit().Error; err != nil {
			return "", errors.New("failed to commit transaction")
		}

	}

	var createdUser = authModel.User{
		Number:     uuid.NewV4().String(),
		FirstName:  payload.SecondName,
		LastName:   payload.Surname,
		SecondName: payload.OtherName,
		Gender:     payload.Gender,
		Type:       constants.TypeStudent,
		Status:     constants.StatusInactive,
	}

	if err := applicationRepo.CreateUser(createdUser); err != nil {
		tx.Rollback()
		return "", errors.New("failed to create user " + err.Error())
	}

	savStudent := models.Student{
		Number:                    studentNumber,
		APN:                       studentAPN,
		Surname:                   payload.Surname,
		OtherName:                 payload.OtherName,
		Gender:                    payload.Gender,
		DateOfBirth:               payload.DateOfBirth,
		ADNStatus:                 payload.ADNStatus,
		ReasonForChoosingSchool:   payload.ReasonForChoosingSchool,
		PreferredTypeOfEnrollment: payload.PreferredTypeOfEnrollment,
		SecondName:                payload.SecondName,
		PreferredCampusNumber:     payload.PreferredCampusNumber,
		UserNumber:                createdUser.Number,
		SourceOfSchoolInfo:        payload.SourceOfSchoolInfo,
		ApplicationNumber:         saveApplication.Number,
		CampusClassNumber:         campusClass.Number,
		SchoolNumber:              payload.SchoolNumber,
		Lin:                       payload.Lin,
		ModeOfStudyNumber:         payload.ModeOfStudyNumber,
	}

	if payload.ADN != "" {
		savStudent.ADN = payload.ADN
	}

	if err := applicationRepo.CreateStudent(savStudent); err != nil {
		tx.Rollback()
		return "", errors.New("failed to save student " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		return "", errors.New("failed to commit transaction " + err.Error())
	}

	var apn = payload.ApplicationNumber
	if apn == "" {
		apn = saveApplication.Number
	}

	return apn, nil
}

func CreateAPN(tx *gorm.DB) (string, error) {
	year := strconv.Itoa(time.Now().Year())
	month := time.Now().Format("01")

	saveAPN := models.APN{
		Year:  year,
		Month: month,
	}

	if tx == nil {
		tx = config.DB
	}

	if err := tx.Create(&saveAPN).Error; err != nil {
		return "", err
	}

	savedID := saveAPN.ID + 1
	result := strings.Join([]string{saveAPN.Year, "01", fmt.Sprintf("%04s", strconv.FormatUint(uint64(savedID), 10))}, "")

	return result, nil
}
