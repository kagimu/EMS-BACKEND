package services

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/repository"
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
		return responses.ApplicationResponse{}, err
	}

	previousSchools, err = applicationRepo.FindPreviousSchoolsByApplicationNumber(id)
	if err != nil {
		return responses.ApplicationResponse{}, err

	}

	return BuildSingleApplicationResponse(application, &student, previousSchools), nil
}

func BuildSingleApplicationResponse(application *models.Application, student *studentModel.Student, previousSchools []previousSchoolModel.PreviousSchool) responses.ApplicationResponse {
	var studentResponse responses.StudentApplicant
	var parentResponse responses.ParentResponse
	var previousSchoolsResponse []responses.PreviousSchoolResponse

	studentResponse = responses.StudentApplicant{
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
		Parent:                    parentResponse,
		PreviousSchools:           previousSchoolsResponse,
		ApplicationFeePayment:     application.ApplicationFeePayment,
		ApplicationFeeInvalidDate: application.ApplicationFeeInvalidDate,
		InterviewMode:             application.InterviewMode,
		Number:                    application.Number,
		Status:                    application.Status,
		UserNumber:                application.UserNumber,
	}
}
