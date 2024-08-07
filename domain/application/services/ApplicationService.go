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
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
)

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
	var studentResponse responses.StudentApplicant
	var previousSchoolsResponse []responses.PreviousSchoolResponse

	Classes, err := repositories.GetClassByCampusClassNumber(student.CampusClassNumber)
	if err != nil {
		return responses.ApplicationResponse{}
	}

	preferredCampus, err := repositories2.GetCampusByNumber(student.PreferedCampusNumber)
	if err != nil {
		return responses.ApplicationResponse{
			ApplicationFeePayment:     application.ApplicationFeePayment,
			ApplicationFeeInvalidDate: application.ApplicationFeeInvalidDate,
			InterviewMode:             application.InterviewMode,
			Number:                    application.Number,
			Status:                    application.Status,
			UserNumber:                application.UserNumber,
		}
	}

	parent, err := repositories3.GetParentByApplicationNumber(application.Number)
	if err != nil {
		return responses.ApplicationResponse{}
	}

	studentResponse = responses.StudentApplicant{
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
		CampusClass:               Classes.Name,
		PreferredCampus:           preferredCampus.Name,
		CampusClassNumber:         student.CampusClassNumber,
		Gender:                    student.Gender,
		SchoolNumber:              student.SchoolNumber,
	}

	for _, previousSchool := range previousSchools {
		ps := responses.PreviousSchoolResponse{
			Name: previousSchool.Name,
		}

		for _, class := range previousSchool.Classes {
			ps.Classes = append(ps.Classes, responses.ClassResponse{
				PreviousSchoolNumber: class.PreviousSchoolNumber,
				Year:                 class.Year,
				Class:                class.Class,
			})
		}
		previousSchoolsResponse = append(previousSchoolsResponse, ps)
	}

	return responses.ApplicationResponse{
		Student:                   studentResponse,
		Parent:                    createParentResponse(parent),
		PreviousSchools:           previousSchoolsResponse,
		ApplicationFeePayment:     application.ApplicationFeePayment,
		ApplicationFeeInvalidDate: application.ApplicationFeeInvalidDate,
		InterviewMode:             application.InterviewMode,
		Number:                    application.Number,
		Status:                    application.Status,
		UserNumber:                application.UserNumber,
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
