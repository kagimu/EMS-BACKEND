package services

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/repositories"
	models3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	repositories3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/repository"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/repositories"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type EnrolmentService interface {
	NewEnrolment(enrolmentRequest requests.EnrollmentRequest, currentUserNumber string) error
}

type enrolmentService struct {
}

var enrolmentRepository = repository.NewEnrolmentRepository()
var admissionRepository = repositories.NewAdmissionRepository()

func (enrolment *enrolmentService) NewEnrolment(enrolmentRequest requests.EnrollmentRequest, currentUserNumber string) error {

	student, err := repositories2.FindStudentByNumber(enrolmentRequest.StudentNumber)
	if err != nil {
		return errors.New("Checking student number error: " + err.Error())
	}

	if student.Number == "" {
		return errors.New("Invalid student number ")
	}

	admission, err := admissionRepository.FindAdmissionByStudentNumber(enrolmentRequest.StudentNumber)
	if err != nil {
		return errors.New("Error checking for student admission " + err.Error())
	}

	if admission.Number == "" {
		return errors.New("This student has no any admission! ")
	}

	currentEnrolment, err := enrolmentRepository.GetStudentEnrolment(enrolmentRequest.StudentNumber, enrolmentRequest.StudyPeriodNumber, enrolmentRequest.CampusClassNumber, "")
	if err != nil {
		return errors.New("Checking duplicate error: " + err.Error())
	}

	if currentEnrolment.Number != "" {
		return errors.New("Student already enrolled for this study period ")
	}

	allOtherEnrolments, err := enrolmentRepository.AllEnrolments(enrolmentRequest.StudentNumber, "", "", "", constants.StatusActive)
	if err != nil {
		return errors.New("Checking previous enrolments error: " + err.Error())
	}

	feedingGroup := models3.FeedingGroup{}
	groupUser := models3.FeedingGroupUser{}

	if enrolmentRequest.FeedingGroupNumber != "" {
		feedingGroup, err = repositories3.GetFeedingGroupByNumber(enrolmentRequest.FeedingGroupNumber)
		if err != nil {
			return errors.New("Checking feeding group error: " + err.Error())
		}

		if feedingGroup.Number == "" {
			return errors.New("Invalid feeding group number ")
		}

		groupUser, err = repositories3.GetFeedingGroupUserByUserNumber(student.UserNumber)
		if err != nil {
			return errors.New("Checking current user feeding group error: " + err.Error())
		}
	}

	saveEnrollment := models.Enrollment{
		Number:            uuid.NewV4().String(),
		StudentNumber:     student.Number,
		StreamNumber:      enrolmentRequest.StreamNumber,
		CampusNumber:      enrolmentRequest.CampusNumber,
		SchoolNumber:      enrolmentRequest.SchoolNumber,
		EnrollmentDate:    enrolmentRequest.EnrollmentDate,
		StudyPeriodNumber: enrolmentRequest.StudyPeriodNumber,
		CampusClassNumber: enrolmentRequest.CampusClassNumber,
		Status:            constants.StatusActive,
		UserNumber:        currentUserNumber,
		ContinuingStudent: enrolmentRequest.ContinuingStudent,
	}

	var requiredItems []models2.AdmissionCondition

	if admission.Status != constants.StatusEnrolled {
		requiredItems, err = repositories.FindAllAdmissionItem(admission.Number)
		if err != nil {
			return errors.New("Checking required items error: " + err.Error())
		}
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		err = enrolmentRepository.SaveEnrolment(saveEnrollment, tx)
		if err != nil {
			return errors.New("Saving enrolment error: " + err.Error())
		}

		if len(requiredItems) > 0 {

			for _, item := range enrolmentRequest.Items {
				var enrolmentItem models.EnrollmentItem
				enrolmentItem.ItemNumber = item.ItemNumber
				enrolmentItem.EnrollmentNumber = saveEnrollment.Number
				enrolmentItem.Quantity = item.Quantity
				enrolmentItem.Brought = item.Brought
				enrolmentItem.UserNumber = currentUserNumber

				err = repository.SaveEnrollmentItem(enrolmentItem, tx)
				if err != nil {
					return errors.New("Saving enrolment items error: " + err.Error())
				}
			}

		}

		if admission.Status != constants.StatusEnrolled {
			admission.Status = constants.StatusEnrolled
			err = admissionRepository.SaveAdmission(admission, tx)
			if err != nil {
				return errors.New("Updating admission to enrolled error " + err.Error())
			}
		}

		for _, otherEnrolment := range allOtherEnrolments {
			otherEnrolment.Status = constants.StatusInactive
			err = enrolmentRepository.SaveEnrolment(otherEnrolment, tx)
			if err != nil {
				return errors.New("Updating previous enrolment error " + err.Error())
			}
		}

		if feedingGroup.Number != "" {

			if groupUser.Number != "" {
				groupUser.FeedingGroupNumber = feedingGroup.Number
				groupUser.Status = feedingGroup.Status
				groupUser.CampusNumber = enrolmentRequest.CampusNumber
				err = repositories3.SaveFeedingGroupUser(groupUser, tx)
				if err != nil {
					return errors.New("Updating feeding group user error: " + err.Error())
				}
			} else {
				err := services.CreateFeedingGroupUser(student.UserNumber, enrolmentRequest.CampusNumber, feedingGroup.Number, constants.StatusActive, tx)
				if err != nil {
					return errors.New("Failed to create feeding group user " + err.Error())
				}
			}

		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func NewEnrolmentService() EnrolmentService {
	return &enrolmentService{}
}
