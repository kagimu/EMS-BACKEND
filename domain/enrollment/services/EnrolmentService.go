package services

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/repository"
)

type EnrolmentService interface {
	NewEnrolment(enrolmentRequest requests.EnrollmentRequest, currentUserNumber string) error
}

type enrolmentService struct {
}

var enrolmentRepository = repository.NewEnrolmentRepository()

func (enrolment *enrolmentService) NewEnrolment(enrolmentRequest requests.EnrollmentRequest, currentUserNumber string) error {
	currentEnrolment, err := enrolmentRepository.GetStudentEnrolment(enrolmentRequest.StudentNumber, enrolmentRequest.StudyPeriodNumber)
}

func NewEnrolmentService() EnrolmentService {
	return &enrolmentService{}
}
