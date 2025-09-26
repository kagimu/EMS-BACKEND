package services

import (
	"strings"
	"unicode"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/repository"
	repositories4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/repositories"
	repositories3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/repositories"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/repositories"
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

func DecodeAdmissionLetterTemplateService(studentNumber string) (responses.DecodedAdmissionLetter, error) {
	repo := repositories.NewAdmissionTemplateRepository(config.DB)
	student, err := repositories2.FindStudentByNumber(studentNumber)
	if err != nil {
		return responses.DecodedAdmissionLetter{}, err
	}

	user, err := repository.GetUserByNumber(student.UserNumber)
	if err != nil {
		return responses.DecodedAdmissionLetter{}, err
	}

	var admissionRepo = repositories.NewAdmissionRepository()

	admission, err := admissionRepo.FindAdmissionByStudentNumber(student.Number)
	if err != nil {
		return responses.DecodedAdmissionLetter{}, err
	}

	admissionTemplate, err := repo.FindAdmissionTemplateBySchoolNumber(admission.SchoolNumber)
	if err != nil {
		return responses.DecodedAdmissionLetter{}, err
	}

	school, err := repositories3.GetSchoolByNumber(admission.SchoolNumber)
	if err != nil {
		return responses.DecodedAdmissionLetter{}, err
	}

	campus, err := repositories4.GetCampusByNumber(admission.CampusNumber)
	if err != nil {
		return responses.DecodedAdmissionLetter{}, err
	}

	var resp responses.DecodedAdmissionLetter
	resp.AdmissionLetterNumber = admissionTemplate.Number

	// find and replace {{STUDENT_NAME}} with student's full name in the phrase template
	var phrase = admissionTemplate.Phrase
	capitalizedFirstName := capitalizeFirstLetter(user.FirstName)
	capitalizedLastName := capitalizeFirstLetter(user.LastName)
	capitalizedSecondName := capitalizeFirstLetter(user.SecondName)
	phrase = replace(phrase, "{{STUDENT_NAME}}", capitalizedFirstName+" "+capitalizedLastName+" "+capitalizedSecondName)
	phrase = replace(phrase, "{{SCHOOL_NAME}}", school.Name)
	phrase = replace(phrase, "{{ADMISSION_NUMBER}}", student.ADN)
	phrase = replace(phrase, "{{ADMISSION_DATE}}", admission.CreatedAt.Format("02-01-2006"))
	//phrase = replace(phrase, "{{ADMISSION_CLASS}}", admission.CampusClassNumber)
	phrase = replace(phrase, "{{CAMPUS_NAME}}", campus.Name)
	phrase = replace(phrase, "{{PRINCIPAL_NAME}}", school.PrincipalName)
	phrase = replace(phrase, "{{PRINCIPAL_TITLE}}", school.PrincipalTitle)
	//phrase = replace(phrase, "{{REPORTING_DUE_DATE}}", admission.CampusClassNumber)

	resp.Phrase = phrase

	return resp, nil
}

func replace(phrase, placeholder, value string) string {
	if phrase == "" {
		return ""
	}
	// find the placeholder in the phrase and replace it with the value
	if placeholder == "" {
		return phrase
	}
	if value == "" {
		return phrase
	}
	phrase = strings.ReplaceAll(phrase, placeholder, value)
	return phrase
}

func capitalizeFirstLetter(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func DeleteAdmissionLetterTemplateService(templateNumber string) error {
	repo := repositories.NewAdmissionTemplateRepository(config.DB)
	if err := repo.DeleteAdmissionTemplateByNumber(templateNumber); err != nil {
		return err
	}
	return nil
}
