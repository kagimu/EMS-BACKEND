package services

import (
	"errors"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/repositories"
	userModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	models3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	repositories3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/services"
	repositories5 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	repositories8 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/repository"
	models4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/models"
	repositories4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/repositories"
	services2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/services"
	repositories6 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/repositories"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/repositories"
	repositories7 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/repositories"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type EnrolmentService interface {
	NewEnrolment(enrolmentRequest requests.EnrollmentRequest, currentUserNumber string) error
	Enrolments(search, schoolNumber, streamNumber, campusNumber, campusClassNumber, studyPeriodNumber, academicYear, status string, user *userModel.User, page func(db *gorm.DB) *gorm.DB) ([]responses.Enrollment, error)
	GetEnrollmentDetailsByNumber(number string) (responses.EnrollmentDetails, error)
	UpdateEnrollmentItem(id uint, quantity float64, brought bool) error
}

type enrolmentService struct {
}

func (enrolment *enrolmentService) Enrolments(search, schoolNumber, streamNumber, campusNumber, campusClassNumber, studyPeriodNumber, academicYear, status string, user *userModel.User, page func(db *gorm.DB) *gorm.DB) ([]responses.Enrollment, error) {

	paginatedResponse, err := enrolmentRepository.PaginatedEnrolments(search, schoolNumber, streamNumber, campusNumber, campusClassNumber, studyPeriodNumber, academicYear, status, user, page)
	students := make([]responses.Enrollment, 0)

	if err != nil {
		return students, err
	}

	for _, enrollment := range paginatedResponse {

		campusLevel, _ := repositories6.GetCampusLevelByNumber(enrollment.CampusClass.CampusLevelNumber)

		curriculumLevel, _ := repositories6.GetCurriculumLevelByNumber(campusLevel.CurriculumLevelNumber)

		studyPeriodName, _ := repositories7.GetStudyPeriodNameByNumber(enrollment.StudyPeriod.StudyPeriodNameNumber)

		students = append(students, responses.Enrollment{
			LevelNumber:            curriculumLevel.LevelNumber,
			CampusCurriculumNumber: campusLevel.CampusCurriculumNumber,
			StudentSurname:         enrollment.Student.Surname,
			StudentSecondName:      enrollment.Student.SecondName,
			StudentOtherName:       enrollment.Student.OtherName,
			StudentNumber:          enrollment.Student.Number,
			Photo:                  commons.GenerateDownloadPath(enrollment.Student.Photo, true),
			Number:                 enrollment.Number,
			School:                 enrollment.School.Name,
			Campus:                 enrollment.Campus.Name,
			StudyPeriod:            studyPeriodName.Name,
			Status:                 enrollment.Status,
			EnrollmentDate:         enrollment.EnrollmentDate,
			ContinuingStudent:      enrollment.ContinuingStudent,
			Class:                  enrollment.CampusClass.ClassName + "(" + enrollment.CampusClass.ClassAbbreviation + ")",
			Adn:                    enrollment.Student.ADN,
			UserNumber:             enrollment.Student.UserNumber,
			CampusClassNumber:      enrollment.CampusClassNumber,
			CampusNumber:           enrollment.CampusNumber,
		})
	}

	return students, nil
}

func (enrolment *enrolmentService) GetEnrollmentDetailsByNumber(number string) (responses.EnrollmentDetails, error) {
	enrollment, err := enrolmentRepository.GetEnrollmentByNumber(number)
	if err != nil {
		return responses.EnrollmentDetails{}, err
	}
	if enrollment.Number == "" {
		return responses.EnrollmentDetails{}, errors.New("enrollment not found")
	}
	studyPeriodName, _ := repositories7.GetStudyPeriodNameByNumber(enrollment.StudyPeriod.StudyPeriodNameNumber)
	streamName := ""
	if enrollment.Stream.Number != "" {
		streamName = enrollment.Stream.Name
	}
	classDisplay := ""
	if enrollment.CampusClass.Number != "" {
		classDisplay = enrollment.CampusClass.ClassName + "(" + enrollment.CampusClass.ClassAbbreviation + ")"
	}
	items := make([]responses.EnrollmentItemWithName, 0, len(enrollment.EnrollmentItems))
	for _, ei := range enrollment.EnrollmentItems {
		itemName := ""
		if ei.Item.Number != "" {
			itemName = ei.Item.Name
		}
		items = append(items, responses.EnrollmentItemWithName{
			ID:               ei.ID,
			EnrollmentNumber: ei.EnrollmentNumber,
			ItemNumber:       ei.ItemNumber,
			ItemName:         itemName,
			Quantity:         ei.Quantity,
			Brought:          ei.Brought,
		})
	}
	return responses.EnrollmentDetails{
		Number:            enrollment.Number,
		StudentNumber:     enrollment.StudentNumber,
		Adn:               enrollment.Student.ADN,
		StudentSurname:    enrollment.Student.Surname,
		StudentSecondName: enrollment.Student.SecondName,
		StudentOtherName:  enrollment.Student.OtherName,
		Class:             classDisplay,
		Stream:            streamName,
		School:            enrollment.School.Name,
		Campus:            enrollment.Campus.Name,
		StudyPeriod:       studyPeriodName.Name,
		Status:            enrollment.Status,
		EnrollmentDate:    enrollment.EnrollmentDate,
		ContinuingStudent: enrollment.ContinuingStudent,
		CampusClassNumber: enrollment.CampusClassNumber,
		CampusNumber:      enrollment.CampusNumber,
		StudyPeriodNumber: enrollment.StudyPeriodNumber,
		SchoolNumber:      enrollment.SchoolNumber,
		StreamNumber:      enrollment.StreamNumber,
		EnrollmentItems:   items,
	}, nil
}

func (enrolment *enrolmentService) UpdateEnrollmentItem(id uint, quantity float64, brought bool) error {
	item, err := repository.GetEnrollmentItemByID(id)
	if err != nil {
		return err
	}
	if item.ID == 0 {
		return errors.New("enrollment item not found")
	}
	item.Quantity = quantity
	item.Brought = brought
	return repository.SaveEnrollmentItem(item)
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

	emergencyInformationRepository := repositories2.NewEmergencyInformationRepository(config.DB)
	emergencyInformation, err := emergencyInformationRepository.GetEmergencyInformationByStudentNumber(student.Number)
	if err != nil {
		return errors.New("Checking emergency information error: " + err.Error())
	}

	if emergencyInformation.StudentNumber != student.Number {
		return errors.New("The student profile is incomplete ")
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
		return errors.New("Checking student current enrollment error: " + err.Error())
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

	class, _ := repositories5.GetClassByCampusClassNumber(enrolmentRequest.CampusClassNumber)

	billType := constants.BillTypeContinuing
	if admission.Status == constants.StatusPending {
		billType = constants.BillTypeNewApplication
	}

	var billable models4.Billable

	if enrolmentRequest.BillableNumber != "" {

		billable, err = repositories4.GetBillableByNumber(enrolmentRequest.BillableNumber)
		if err != nil {
			return errors.New("Checking billable by number error " + err.Error())
		}

		if billable.Number == "" {
			return errors.New("Invalid bill number ")
		}

		if billable.Status != constants.StatusActive {
			return errors.New("Selected bill is no longer active ")
		}

	} else {

		curriculum, err := repositories8.GetCurriculumByCampusCurriculumNumber(enrolmentRequest.CampusCurriculumNumber, enrolmentRequest.CampusNumber)
		if err != nil {
			return errors.New("Checking curriculum and campus error " + err.Error())
		}

		if curriculum.Number == "" {
			return errors.New("Invalid Curriculum on campus ")
		}

		billable, err = repositories4.GetEnrollmentBill(enrolmentRequest.StudyPeriodNumber, class.Number, enrolmentRequest.CampusNumber, billType, enrolmentRequest.ModeOfStudyNumber, curriculum.Number, constants.StatusActive)
		if err != nil {
			return errors.New("Checking billable error " + err.Error())
		}
	}

	if billable.Number == "" {
		return errors.New("This enrollment can not be invoiced, check billable configurations ")
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
		ModeOfStudyNumber: enrolmentRequest.ModeOfStudyNumber,
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

		invoice, err := services2.CreateInvoiceFromBillable(billable, student, enrolmentRequest.ClientReference, enrolmentRequest.Description, currentUserNumber, enrolmentRequest.PaidAmount, enrolmentRequest.DatePaid, tx)
		if err != nil {
			return err
		}

		err = repositories4.SaveEnrollmentInvoice(models4.EnrollmentInvoice{
			InvoiceNumber:    invoice.Number,
			EnrollmentNumber: saveEnrollment.Number,
		}, tx)

		if err != nil {
			return errors.New("Attaching enrollment to promissory note error: " + err.Error())
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
