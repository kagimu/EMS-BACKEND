package services

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/repositories"
	services4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/services"
	userModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	repositories3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/repositories"
	services2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/services"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/services"
	repositories5 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/repositories"
	repositories4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/repositories"
	services3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/services"
	"gorm.io/gorm"
)

type AdmissionService interface {
	AllAdmissions(campusNumber, studyPeriodNumber, streamNumber, interviewNumber, schoolNumber, campusClassNumber string, user *userModel.User, paginateFunc func(db *gorm.DB) *gorm.DB) ([]responses.AdmissionResponse, error)
}

type admissionService struct {
}

var admissionRepository = repositories.NewAdmissionRepository()

func (admissionService *admissionService) AllAdmissions(campusNumber, studyPeriodNumber, streamNumber, interviewNumber, schoolNumber, campusClassNumber string, user *userModel.User, paginateFunc func(db *gorm.DB) *gorm.DB) ([]responses.AdmissionResponse, error) {

	var admissionResponse = make([]responses.AdmissionResponse, 0)

	admissions, err := admissionRepository.AllAdmissions(campusNumber, studyPeriodNumber, streamNumber, interviewNumber, schoolNumber, campusClassNumber, user, paginateFunc)
	if err != nil {
		return admissionResponse, errors.New("Error getting admissions " + err.Error())
	}

	for _, admission := range admissions {
		admissionResponse = append(admissionResponse, createAdmissionResponse(admission))
	}

	return admissionResponse, err
}

func createAdmissionResponse(admission models.Admission) responses.AdmissionResponse {
	admissionResponse := responses.AdmissionResponse{
		AdmissionNumber: admission.Number,
		Status:          admission.Status,
	}

	var studentRepository = repositories5.NewStudentRepositoryImpl(config.DB)

	class, _ := repositories2.GetClassByCampusClassNumber(admission.CampusClassNumber)
	campus, _ := repositories3.GetCampusByNumber(admission.CampusNumber)
	studyPeriod, _ := repositories4.GetStudyPeriodByNumber(admission.StudyPeriodNumber)
	student, _ := studentRepository.GetStudentByNumber(admission.StudentNumber)
	stream, _ := repositories2.GetStreamByNumber(admission.StreamNumber)

	admissionResponse.Class = services.CreateClassResponse(class)
	admissionResponse.CampusResponse = services2.CreateCampusResponse(campus)
	admissionResponse.StudyPeriod = services3.CreateStudyPeriod(studyPeriod)
	admissionResponse.Student = services4.CreateStudentResponse(&student)
	admissionResponse.Stream = services.CreateStreamResponse(stream)

	return admissionResponse
}

func NewAdmissionService() AdmissionService {
	return &admissionService{}
}
