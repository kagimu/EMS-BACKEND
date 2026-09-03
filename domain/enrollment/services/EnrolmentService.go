package services

import (
	"errors"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/repositories"
	userModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	authRepo "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/repository"
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
	itemRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/repositories"
	repositories6 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/repositories"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/repositories"
	repositories7 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/repositories"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type EnrolmentService interface {
	NewEnrolment(enrolmentRequest requests.EnrollmentRequest, currentUserNumber string) (responses.EnrollmentResponse, error)
	Enrolments(search, schoolNumber, streamNumber, campusNumber, campusClassNumber, studyPeriodNumber, academicYear, status string, user *userModel.User, page func(db *gorm.DB) *gorm.DB) ([]responses.Enrollment, error)
	GetEnrollmentDetailsByNumber(number string) (responses.EnrollmentDetails, error)
	UpdateEnrollmentItem(id uint, quantity float64, brought bool, dateBrought string, remarks string, currentUserNumber string) error
	GetRequirementClearanceLogs(enrollmentNumber, itemNumber, dateFrom, dateTo, staffNumber string) ([]responses.RequirementClearanceLogResponse, error)
	DeleteClearanceLog(logID uint) error
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
			StudyPeriodStartDate:   enrollment.StudyPeriod.StudyPeriodStartDate.Format("2006-01-02"),
			StudyPeriodEndDate:     enrollment.StudyPeriod.StudyPeriodEndDate.Format("2006-01-02"),
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
		reqQty := ei.RequiredQuantity
		if reqQty == 0 {
			reqQty = ei.Quantity
		}
		items = append(items, responses.EnrollmentItemWithName{
			ID:               ei.ID,
			EnrollmentNumber: ei.EnrollmentNumber,
			ItemNumber:       ei.ItemNumber,
			ItemName:         itemName,
			RequiredQuantity: reqQty,
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
		StudyPeriodStartDate: enrollment.StudyPeriod.StudyPeriodStartDate.Format("2006-01-02"),
		StudyPeriodEndDate:   enrollment.StudyPeriod.StudyPeriodEndDate.Format("2006-01-02"),
		ContinuingStudent: enrollment.ContinuingStudent,
		CampusClassNumber: enrollment.CampusClassNumber,
		CampusNumber:      enrollment.CampusNumber,
		StudyPeriodNumber: enrollment.StudyPeriodNumber,
		SchoolNumber:      enrollment.SchoolNumber,
		StreamNumber:      enrollment.StreamNumber,
		EnrollmentItems:   items,
	}, nil
}

func (enrolment *enrolmentService) UpdateEnrollmentItem(id uint, quantity float64, brought bool, dateBrought string, remarks string, currentUserNumber string) error {
	item, err := repository.GetEnrollmentItemByID(id)
	if err != nil {
		return err
	}
	if item.ID == 0 {
		return errors.New("enrollment item not found")
	}

	oldQuantity := item.Quantity
	quantityCleared := quantity - oldQuantity

	item.Quantity = quantity
	item.Brought = brought
	if err := repository.SaveEnrollmentItem(&item); err != nil {
		return err
	}

	var staffName string
	user, userErr := authRepo.GetUserByNumber(currentUserNumber)
	if userErr == nil && user.ID != 0 {
		staffName = user.FirstName + " " + user.SecondName
		if user.LastName != "" {
			staffName += " " + user.LastName
		}
	}

	if dateBrought == "" {
		dateBrought = time.Now().Format("2006-01-02")
	}

	remaining := item.RequiredQuantity - quantity
	if remaining < 0 {
		remaining = 0
	}

	log := models.RequirementClearanceLog{
		EnrollmentNumber:     item.EnrollmentNumber,
		EnrollmentItemID:     item.ID,
		ItemNumber:           item.ItemNumber,
		ItemName:             item.Item.Name,
		DateCleared:          dateBrought,
		QuantityCleared:      quantityCleared,
		TotalQuantityBrought: quantity,
		RequiredQuantity:     item.RequiredQuantity,
		RemainingQuantity:    remaining,
		ClearedByStaffNumber: currentUserNumber,
		ClearedByStaffName:   staffName,
		Remarks:              remarks,
	}
	return repository.SaveRequirementClearanceLog(log)
}

func (enrolment *enrolmentService) GetRequirementClearanceLogs(enrollmentNumber, itemNumber, dateFrom, dateTo, staffNumber string) ([]responses.RequirementClearanceLogResponse, error) {
	logs, err := repository.FindRequirementClearanceLogs(enrollmentNumber, itemNumber, dateFrom, dateTo, staffNumber)
	if err != nil {
		return nil, err
	}
	var res []responses.RequirementClearanceLogResponse
	for _, l := range logs {
		res = append(res, responses.RequirementClearanceLogResponse{
			ID:                   l.ID,
			EnrollmentNumber:     l.EnrollmentNumber,
			EnrollmentItemID:     l.EnrollmentItemID,
			ItemNumber:           l.ItemNumber,
			ItemName:             l.ItemName,
			DateCleared:          l.DateCleared,
			QuantityCleared:      l.QuantityCleared,
			TotalQuantityBrought: l.TotalQuantityBrought,
			RequiredQuantity:     l.RequiredQuantity,
			RemainingQuantity:    l.RemainingQuantity,
			ClearedByStaffNumber: l.ClearedByStaffNumber,
			ClearedByStaffName:   l.ClearedByStaffName,
			Remarks:              l.Remarks,
			CreatedAt:            l.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return res, nil
}

func (enrolment *enrolmentService) DeleteClearanceLog(logID uint) error {
	// 1. Fetch the log to know which enrollment item it belongs to.
	log, err := repository.GetClearanceLogByID(logID)
	if err != nil {
		return errors.New("clearance log not found: " + err.Error())
	}
	if log.ID == 0 {
		return errors.New("clearance log not found")
	}

	// 2. Fetch the parent enrollment item.
	item, err := repository.GetEnrollmentItemByID(log.EnrollmentItemID)
	if err != nil || item.ID == 0 {
		return errors.New("enrollment item not found")
	}

	// 3. Run delete + item recalculation in a single transaction.
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// Delete the log record permanently.
		if err := repository.DeleteRequirementClearanceLogByID(logID, tx); err != nil {
			return errors.New("failed to delete clearance log: " + err.Error())
		}

		// Recalculate the total quantity from the remaining logs for this item.
		newTotal, err := repository.SumClearanceLogQuantityByItemID(log.EnrollmentItemID, tx)
		if err != nil {
			return errors.New("failed to recalculate quantity: " + err.Error())
		}
		if newTotal < 0 {
			newTotal = 0
		}

		// Update the enrollment item with recalculated quantity and brought flag.
		item.Quantity = newTotal
		item.Brought = newTotal >= item.RequiredQuantity && newTotal > 0
		if err := repository.SaveEnrollmentItem(&item, tx); err != nil {
			return errors.New("failed to update enrollment item: " + err.Error())
		}
		return nil
	})
}

var enrolmentRepository = repository.NewEnrolmentRepository()
var admissionRepository = repositories.NewAdmissionRepository()

func (enrolment *enrolmentService) NewEnrolment(enrolmentRequest requests.EnrollmentRequest, currentUserNumber string) (responses.EnrollmentResponse, error) {

	student, err := repositories2.FindStudentByNumber(enrolmentRequest.StudentNumber)
	enrolmentResponse := responses.EnrollmentResponse{}
	if err != nil {
		return enrolmentResponse, errors.New("Checking student number error: " + err.Error())
	}

	if student.Number == "" {
		return enrolmentResponse, errors.New("Invalid student number ")
	}

	emergencyInformationRepository := repositories2.NewEmergencyInformationRepository(config.DB)
	emergencyInformation, err := emergencyInformationRepository.GetEmergencyInformationByStudentNumber(student.Number)
	if err != nil {
		return enrolmentResponse, errors.New("Checking emergency information error: " + err.Error())
	}

	if emergencyInformation.StudentNumber != student.Number {
		return enrolmentResponse, errors.New("The student profile is incomplete ")
	}

	admission, err := admissionRepository.FindAdmissionByStudentNumber(enrolmentRequest.StudentNumber)
	if err != nil {
		return enrolmentResponse, errors.New("Error checking for student admission " + err.Error())
	}

	if admission.Number == "" {
		return enrolmentResponse, errors.New("This student has no any admission! ")
	}

	currentEnrolment, err := enrolmentRepository.GetStudentEnrolment(enrolmentRequest.StudentNumber, enrolmentRequest.StudyPeriodNumber, enrolmentRequest.CampusClassNumber, "")
	if err != nil {
		return enrolmentResponse, errors.New("Checking student current enrollment error: " + err.Error())
	}

	if currentEnrolment.Number != "" {
		return enrolmentResponse, errors.New("Student already enrolled for this study period ")
	}

	allOtherEnrolments, err := enrolmentRepository.AllEnrolments(enrolmentRequest.StudentNumber, "", "", "", constants.StatusActive)
	if err != nil {
		return enrolmentResponse, errors.New("Checking previous enrolments error: " + err.Error())
	}

	feedingGroup := models3.FeedingGroup{}
	groupUser := models3.FeedingGroupUser{}

	if enrolmentRequest.FeedingGroupNumber != "" {
		feedingGroup, err = repositories3.GetFeedingGroupByNumber(enrolmentRequest.FeedingGroupNumber)
		if err != nil {
			return enrolmentResponse, errors.New("Checking feeding group error: " + err.Error())
		}

		if feedingGroup.Number == "" {
			return enrolmentResponse, errors.New("Invalid feeding group number ")
		}

		groupUser, err = repositories3.GetFeedingGroupUserByUserNumber(student.UserNumber)
		if err != nil {
			return enrolmentResponse, errors.New("Checking current user feeding group error: " + err.Error())
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
			return enrolmentResponse, errors.New("Checking billable by number error " + err.Error())
		}

		if billable.Number == "" {
			return enrolmentResponse, errors.New("Invalid bill number ")
		}

		if billable.Status != constants.StatusActive {
			return enrolmentResponse, errors.New("Selected bill is no longer active ")
		}

	} else {

		curriculum, err := repositories8.GetCurriculumByCampusCurriculumNumber(enrolmentRequest.CampusCurriculumNumber, enrolmentRequest.CampusNumber)
		if err != nil {
			return enrolmentResponse, errors.New("Checking curriculum and campus error " + err.Error())
		}

		if curriculum.Number == "" {
			return enrolmentResponse, errors.New("Invalid Curriculum on campus ")
		}

		billable, err = repositories4.GetEnrollmentBill(enrolmentRequest.StudyPeriodNumber, class.Number, enrolmentRequest.CampusNumber, billType, enrolmentRequest.ModeOfStudyNumber, curriculum.Number, constants.StatusActive)
		if err != nil {
			return enrolmentResponse, errors.New("Checking billable error " + err.Error())
		}
	}

	if billable.Number == "" {
		return enrolmentResponse, errors.New("This enrollment can not be invoiced, check billable configurations ")
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

	enrolmentResponse.EnrollmentNumber = saveEnrollment.Number

	if admission.Status != constants.StatusEnrolled {
		_, err = repositories.FindAllAdmissionItem(admission.Number)
		if err != nil {
			return enrolmentResponse, errors.New("Checking required items error: " + err.Error())
		}
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		err = enrolmentRepository.SaveEnrolment(saveEnrollment, tx)
		if err != nil {
			return errors.New("Saving enrolment error: " + err.Error())
		}

		if len(enrolmentRequest.Items) > 0 {
			var staffName string
			user, userErr := authRepo.GetUserByNumber(currentUserNumber)
			if userErr == nil && user.ID != 0 {
				staffName = user.FirstName + " " + user.SecondName
				if user.LastName != "" {
					staffName += " " + user.LastName
				}
			}

			dateBrought := enrolmentRequest.EnrollmentDate
			if dateBrought == "" {
				dateBrought = time.Now().Format("2006-01-02")
			}

			for _, item := range enrolmentRequest.Items {
				var enrolmentItem models.EnrollmentItem
				enrolmentItem.ItemNumber = item.ItemNumber
				enrolmentItem.EnrollmentNumber = saveEnrollment.Number
				reqQty := item.RequiredQuantity
				if reqQty == 0 && item.Quantity > 0 {
					reqQty = item.Quantity
				}
				enrolmentItem.RequiredQuantity = reqQty
				enrolmentItem.Quantity = item.Quantity
				enrolmentItem.Brought = item.Brought
				enrolmentItem.UserNumber = currentUserNumber

				err = repository.SaveEnrollmentItem(&enrolmentItem, tx)
				if err != nil {
					return errors.New("Saving enrolment items error: " + err.Error())
				}

				if item.Quantity > 0 {
					itemObj, _ := itemRepositories.GetItemByNumber(item.ItemNumber, tx)
					remaining := reqQty - item.Quantity
					if remaining < 0 {
						remaining = 0
					}

					log := models.RequirementClearanceLog{
						EnrollmentNumber:     saveEnrollment.Number,
						EnrollmentItemID:     enrolmentItem.ID,
						ItemNumber:           item.ItemNumber,
						ItemName:             itemObj.Name,
						DateCleared:          dateBrought,
						QuantityCleared:      item.Quantity,
						TotalQuantityBrought: item.Quantity,
						RequiredQuantity:     reqQty,
						RemainingQuantity:    remaining,
						ClearedByStaffNumber: currentUserNumber,
						ClearedByStaffName:   staffName,
						Remarks:              "Initial clearance at enrolment",
					}
					if err := repository.SaveRequirementClearanceLog(log, tx); err != nil {
						return errors.New("Saving clearance log error: " + err.Error())
					}
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

		invoice, err := services2.CreateInvoiceFromBillable(billable, student, enrolmentRequest.ClientReference, enrolmentRequest.Description, currentUserNumber, enrolmentRequest.PaidAmount, enrolmentRequest.DatePaid, enrolmentRequest.BalanceDueDate, enrolmentRequest.ItemPayments, tx)
		if err != nil {
			return err
		}

		enrolmentResponse.PromissoryNoteNumber = invoice.Number

		err = repositories4.SaveEnrollmentInvoice(models4.EnrollmentInvoice{
			InvoiceNumber:    invoice.Number,
			EnrollmentNumber: saveEnrollment.Number,
			BalanceDueDate:   enrolmentRequest.BalanceDueDate,
		}, tx)

		if err != nil {
			return errors.New("Attaching enrollment to promissory note error: " + err.Error())
		}

		return nil
	})

	if err != nil {
		return enrolmentResponse, err
	}

	return enrolmentResponse, nil
}

func NewEnrolmentService() EnrolmentService {
	return &enrolmentService{}
}
