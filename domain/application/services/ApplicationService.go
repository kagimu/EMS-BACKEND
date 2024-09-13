package services

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/repository"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	parentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/parent/models"
	repositories3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/parent/repositories"
	previousSchoolModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/previousSchool/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/previousSchoolResults/services"
	repositories4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/repositories"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"os"
)

var previousSchoolResultsService = services.NewPreviousSchoolResultsService()

func GetApplicationById(id string) (responses.ApplicationResponse, error) {

	applicationRepo := repository.NewApplicationRepositoryImpl(config.DB)
	application, err := applicationRepo.FindApplicationByNumber(id)
	if err != nil {
		return responses.ApplicationResponse{}, err

	}

	var student studentModel.Student
	var previousSchools []previousSchoolModel.PreviousSchool

	student, err = applicationRepo.GetStudentByApplicationNumber(id)
	if err != nil {
		return BuildSingleApplicationResponse(application, &student, previousSchools), nil
	}

	previousSchools, err = applicationRepo.FindPreviousSchoolsByApplicationNumber(id)
	if err != nil {
		return BuildSingleApplicationResponse(application, &student, previousSchools), nil

	}

	return BuildSingleApplicationResponse(application, &student, previousSchools), nil
}

func BuildSingleApplicationResponse(application *models.Application, student *studentModel.Student, previousSchools []previousSchoolModel.PreviousSchool) responses.ApplicationResponse {
	var studentResponse responses.StudentResponse
	var previousSchoolsResponse []responses.PreviousSchoolResponse

	applicationResponse := responses.ApplicationResponse{
		ApplicationFeePayment:     application.ApplicationFeePayment,
		ApplicationFeeInvalidDate: application.ApplicationFeeInvalidDate,
		InterviewMode:             application.InterviewMode,
		Number:                    application.Number,
		Status:                    application.Status,
		UserNumber:                application.UserNumber,
	}

	studentResponse = CreateStudentResponse(student)

	for _, previousSchool := range previousSchools {
		ps := responses.PreviousSchoolResponse{
			Name: previousSchool.Name,
		}

		for _, class := range previousSchool.Classes {
			ps.Classes = append(ps.Classes, responses.PreviousSchoolClassResponse{
				PreviousSchoolNumber: class.PreviousSchoolNumber,
				Year:                 class.Year,
				Class:                class.Class,
			})
		}
		previousSchoolsResponse = append(previousSchoolsResponse, ps)
	}

	parent, err := repositories3.GetParentByApplicationNumber(application.Number)
	if err != nil {
		return applicationResponse
	}

	applicationResponse.Student = studentResponse
	applicationResponse.Parent = createParentResponse(parent)
	applicationResponse.PreviousSchools = previousSchoolsResponse

	prevSchoolResults, _ := previousSchoolResultsService.PreviousSchoolResults(application.Number)

	applicationResponse.PreviousSchoolResult = prevSchoolResults

	return applicationResponse
}

func CreateStudentResponse(student *studentModel.Student) responses.StudentResponse {

	preferredCampus, err := repositories2.GetCampusByNumber(student.PreferredCampusNumber)
	if err != nil {
		return responses.StudentResponse{}
	}

	campusClass, err := repositories.GetCampusClassByNumber(student.CampusClassNumber)
	if err != nil {
		return responses.StudentResponse{}
	}

	class, err := repositories.GetClassByCampusClassNumber(student.CampusClassNumber)
	if err != nil {
		return responses.StudentResponse{}
	}

	campusLevel, err := repositories4.GetCampusLevelByNumber(campusClass.CampusLevelNumber)
	if err != nil {
		return responses.StudentResponse{}
	}

	curriculumLevel, err := repositories4.GetCurriculumLevelByNumber(campusLevel.CurriculumLevelNumber)
	if err != nil {
		return responses.StudentResponse{}
	}

	if student.Photo != "" {
		student.Photo = os.Getenv("BASE_URL") + "attachments/download/" + student.Photo
	}
	return responses.StudentResponse{
		Lin:                       student.Lin,
		CurriculumLevelNumber:     curriculumLevel.Number,
		LevelNumber:               curriculumLevel.LevelNumber,
		CampusCurriculumNumber:    campusLevel.CampusCurriculumNumber,
		CampusLevelNumber:         campusLevel.Number,
		Number:                    student.Number,
		Surname:                   student.Surname,
		SecondName:                student.SecondName,
		OtherName:                 student.OtherName,
		ADN:                       student.ADN,
		APN:                       student.APN,
		APNStatus:                 student.APNStatus,
		ADNStatus:                 student.ADNStatus,
		Photo:                     student.Photo,
		AcademicReport:            student.AcademicReport,
		Status:                    student.Status,
		SourceOfSchoolInfo:        student.SourceOfSchoolInfo,
		ReasonForChoosingSchool:   student.ReasonForChoosingSchool,
		ReasonForAdmissionDenial:  student.ReasonForAdmissionDenial,
		AttentionNeeded:           student.AttentionNeeded,
		DateOfBirth:               student.DateOfBirth,
		PreferredTypeOfEnrollment: student.PreferredTypeOfEnrollment,
		CampusClass:               class.Name,
		PreferredCampus:           preferredCampus.Name,
		CampusClassNumber:         student.CampusClassNumber,
		Gender:                    student.Gender,
		SchoolNumber:              student.SchoolNumber,
		PreferredCampusNumber:     student.PreferredCampusNumber,
		ClassNumber:               class.Number,
		ModeOfStudyNumber:         student.ModeOfStudyNumber,
	}
}

func createParentResponse(parent parentModel.Parent) responses.ParentResponse {
	return responses.ParentResponse{
		Surname:           parent.Surname,
		SecondName:        parent.SecondName,
		AlternativeNumber: parent.AlternativeTelephoneNumber,
		TelephoneNumber:   parent.MainTelephoneNumber,
		Email:             parent.Email,
		AlternativeEmail:  parent.AlternativeEmail,
		OtherNames:        parent.OtherName,
	}
}
