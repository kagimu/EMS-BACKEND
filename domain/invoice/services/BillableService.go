package services

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	repositories6 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/repositories"
	repositories4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/repositories"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	repositories8 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/repositories"
	repositories5 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/repositories"
	repositories7 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/repositories"
	repositories3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/repositories"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type BillableService interface {
	NewBillable(billableRequest requests.BillableRequest, currentUserNumber string) error
	FindAllBillable(schoolNumber, campusNumber, billType, classNumber, studyPeriodNumber, academicYearNumber, status string) ([]responses.BillableResponse, error)
	BillableDetails(billableNumber string) (responses.DetailedBillableResponse, error)
	GetBill(studentNumber, campusNumber, classNumber, modeOfStudyNumber, studyPeriodNumber, campusCurriculumNumber string) (responses.DetailedBillableResponse, error)
	DeleteBillable(billableNumber string) error
	UpdateBillable(updateBillableRequest requests.UpdateBillableRequest, currentUserNumber string) error
}

type billableService struct {
}

func (billableService *billableService) UpdateBillable(updateBillableRequest requests.UpdateBillableRequest, currentUserNumber string) error {
	billable, err := repositories.GetBillableByNumber(updateBillableRequest.BillableNumber)
	if err != nil {
		return errors.New("Checking billable number error " + err.Error())
	}

	if billable.Number == "" {
		return errors.New("Invalid billable number ")
	}

	invoices, err := repositories.GetBillableInvoiceByBillableNumber(billable.Number, []string{constants.StatusPaid, constants.StatusPending, constants.StatusPartiallyPaid})
	if err != nil {
		return errors.New("Checking invoices from this bill error " + err.Error())
	}

	if len(invoices) != 0 {
		return errors.New("This bill created invoices, it can not be edited ")
	}

	billable.ClassNumber = updateBillableRequest.ClassNumber
	billable.CampusNumber = updateBillableRequest.CampusNumber
	billable.Description = updateBillableRequest.Description
	billable.UserNumber = currentUserNumber
	billable.SchoolNumber = updateBillableRequest.SchoolNumber
	billable.StudyPeriodNumber = updateBillableRequest.StudyPeriodNumber
	billable.BillType = updateBillableRequest.BillType
	billable.ModeOfStudyNumber = updateBillableRequest.ModeOfStudyNumber

	if updateBillableRequest.Status != "" {
		billable.Status = updateBillableRequest.Status
	}

	billItems, err := repositories.FindBillableItems(billable.Number)
	if err != nil {
		return errors.New("Checking billable items error " + err.Error())
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {

		for _, billableItem := range billItems {
			err := repositories.RemoveBillableItem(billableItem, tx)
			if err != nil {
				return errors.New("Removing billable item error " + err.Error())
			}
		}

		err := repositories.SaveBillable(billable, tx)
		if err != nil {
			return errors.New("Error updating bill " + err.Error())
		}

		err = processBillItems(billable, updateBillableRequest.BillableItems, tx)

		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (billableService *billableService) DeleteBillable(billableNumber string) error {
	bill, err := repositories.GetBillableByNumber(billableNumber)
	if err != nil {
		return errors.New("Checking billable by number error " + err.Error())
	}

	if bill.Number == "" {
		return errors.New("Invalid billable number ")
	}

	billableInvoice, err := repositories.GetBillableInvoiceByBillableNumber(billableNumber, []string{constants.StatusPaid, constants.StatusPending, constants.StatusPartiallyPaid})
	if err != nil {
		return errors.New("Checking if bill has ever created an invoice error " + err.Error())
	}

	if len(billableInvoice) != 0 {
		return errors.New("You can not remove this bill, but you can make it Inactive ")
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		err := repositories.RemoveBillable(bill, tx)

		if err != nil {
			return errors.New("Removing bill error " + err.Error())
		}

		return nil
	})

	return err
}

func (billableService *billableService) GetBill(studentNumber, campusNumber, classNumber, modeOfStudyNumber, studyPeriodNumber, campusCurriculumNumber string) (responses.DetailedBillableResponse, error) {

	billableResponse := responses.DetailedBillableResponse{}

	admission, err := admissionRepository.FindAdmissionByStudentNumber(studentNumber)
	if err != nil {
		return billableResponse, errors.New("Error checking for student admission " + err.Error())
	}

	if admission.Number == "" {
		return billableResponse, errors.New("This student has no any admission! ")
	}

	billType := constants.BillTypeContinuing
	if admission.Status == constants.StatusPending {
		billType = constants.BillTypeNewApplication
	}

	curriculum, err := repositories8.GetCurriculumByCampusCurriculumNumber(campusCurriculumNumber, campusNumber)
	if err != nil {
		return billableResponse, errors.New("Checking curriculum and campus error " + err.Error())
	}

	if curriculum.Number == "" {
		return billableResponse, errors.New("Invalid Curriculum on campus ")
	}

	billable, err := repositories.GetEnrollmentBill(studyPeriodNumber, classNumber, campusNumber, billType, modeOfStudyNumber, curriculum.Number, constants.StatusActive)
	if err != nil {
		return billableResponse, errors.New("Checking billable error " + err.Error())
	}

	if billable.Number == "" {
		return billableResponse, errors.New("This enrollment can not be invoiced, check billable configurations ")
	}

	return createDetailedBillableResponse(billable)
}

func (billableService *billableService) BillableDetails(billableNumber string) (responses.DetailedBillableResponse, error) {

	var detailedResponse responses.DetailedBillableResponse

	bill, err := repositories.GetBillableByNumber(billableNumber)
	if err != nil {
		return detailedResponse, errors.New("Checking billable number error " + err.Error())
	}

	if bill.Number == "" {
		return detailedResponse, errors.New("Invalid billable number ")
	}

	return createDetailedBillableResponse(bill)
}

func createDetailedBillableResponse(billable models.Billable) (responses.DetailedBillableResponse, error) {
	billResponse := responses.DetailedBillableResponse{}
	bResponse, err := createBillResponse(billable)
	if err != nil {
		return billResponse, err
	}
	billResponse.BillableResponse = bResponse

	var itemsResponse = make([]responses.InvoiceItemResponse, 0)

	billItems, err := repositories.FindBillableItems(billable.Number)
	if err != nil {
		return billResponse, errors.New("Getting billable items error " + err.Error())
	}

	for _, item := range billItems {
		itemsResponse = append(itemsResponse, createItemResponse(item))
	}

	billResponse.BillableItems = itemsResponse

	return billResponse, nil
}

func createItemResponse(billableItem models.BillableItem) responses.InvoiceItemResponse {

	item, _ := repositories5.GetItemByNumber(billableItem.ItemNumber)

	return responses.InvoiceItemResponse{
		ItemNumber: billableItem.ItemNumber,
		ItemName:   item.Name,
		UnitCost:   billableItem.UnitPrice,
		Quantity:   billableItem.Quantity,
		TotalCost:  billableItem.Quantity * billableItem.UnitPrice,
	}
}

func (billableService *billableService) FindAllBillable(schoolNumber, campusNumber, billType, classNumber, studyPeriodNumber, academicYearNumber, status string) ([]responses.BillableResponse, error) {

	var billsResponse = make([]responses.BillableResponse, 0)

	bills, err := repositories.FindBillable(schoolNumber, studyPeriodNumber, classNumber, campusNumber, academicYearNumber, status, billType)
	if err != nil {
		return billsResponse, errors.New("Error getting bills " + err.Error())
	}

	for _, bill := range bills {
		bResponse, err := createBillResponse(bill)
		if err != nil {
			return billsResponse, err
		}
		billsResponse = append(billsResponse, bResponse)
	}

	return billsResponse, nil
}

func createBillResponse(billable models.Billable) (responses.BillableResponse, error) {

	studyPeriod, _ := repositories3.GetStudyPeriodByNumber(billable.StudyPeriodNumber)
	studyPeriodName, _ := repositories3.GetStudyPeriodNameByNumber(studyPeriod.StudyPeriodNameNumber)
	class, _ := repositories2.GetClassByNumber(billable.ClassNumber)
	campus, _ := repositories4.GetCampusByNumber(billable.CampusNumber)

	amount, err := repositories.GetBillableTotal(billable.Number)

	if err != nil {
		return responses.BillableResponse{}, errors.New("Getting billable total amount error: " + err.Error())
	}

	academicYear, _ := repositories6.GetAcademicYearByNumber(studyPeriod.AcademicYearNumber)

	campusLevel, _ := repositories7.GetCampusLevelByNumber(academicYear.CampusLevelNumber)

	return responses.BillableResponse{
		BillableNumber:         billable.Number,
		Status:                 billable.Status,
		ClassName:              class.Name,
		ClassNumber:            class.Number,
		StudyPeriodNumber:      billable.StudyPeriodNumber,
		StudyPeriodName:        studyPeriodName.Name,
		CampusNumber:           campus.Number,
		CampusName:             campus.Name,
		BillType:               billable.BillType,
		TotalAmount:            amount,
		AcademicYearNumber:     academicYear.Number,
		CampusCurriculumNumber: campusLevel.CampusCurriculumNumber,
		AcademicYearName:       academicYear.AcademicYearName,
		CampusLevelNumber:      academicYear.CampusLevelNumber,
		Description:            billable.Description,
		ModeOfStudyNumber:      billable.ModeOfStudyNumber,
	}, nil
}

func (billableService *billableService) NewBillable(billableRequest requests.BillableRequest, currentUserNumber string) error {

	curriculum, err := repositories8.GetCurriculumByCampusCurriculumNumber(billableRequest.CampusCurriculumNumber, billableRequest.CampusNumber)
	if err != nil {
		return errors.New("Checking curriculum and campus error " + err.Error())
	}

	if curriculum.Number == "" {
		return errors.New("Invalid Curriculum on campus ")
	}

	billable, err := repositories.GetEnrollmentBill(billableRequest.StudyPeriodNumber, billableRequest.ClassNumber, billableRequest.CampusNumber, billableRequest.BillType, billableRequest.ModeOfStudyNumber, constants.StatusActive, curriculum.Number)
	if err != nil {
		return errors.New("Checking bill duplicates error at creation " + err.Error())
	}

	if billable.Number != "" {
		return errors.New("You already have an active similar bill ")
	}

	newBillable := models.Billable{
		ModeOfStudyNumber: billableRequest.ModeOfStudyNumber,
		ClassNumber:       billableRequest.ClassNumber,
		StudyPeriodNumber: billableRequest.StudyPeriodNumber,
		UserNumber:        currentUserNumber,
		Number:            uuid.NewV4().String(),
		CampusNumber:      billableRequest.CampusNumber,
		SchoolNumber:      billableRequest.SchoolNumber,
		Description:       billableRequest.Description,
		BillType:          billableRequest.BillType,
		CurriculumNumber:  curriculum.Number,
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {

		err = repositories.SaveBillable(newBillable, tx)
		if err != nil {
			return errors.New("Error saving bill " + err.Error())
		}

		err = processBillItems(newBillable, billableRequest.BillableItems, tx)

		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func processBillItems(billable models.Billable, billItems []requests.InvoiceItemRequest, tx *gorm.DB) error {

	for _, itemRequest := range billItems {
		checkDuplicate, err := repositories.GetDuplicatedBillableItem(billable.Number, itemRequest.ItemNumber, tx)
		if err != nil {
			return errors.New("Checking duplicate item error " + err.Error())
		}

		if checkDuplicate.Number != "" {
			return errors.New("You have added an item multiple times ")
		}

		err = repositories.SaveBillableItem(models.BillableItem{
			ItemNumber:     itemRequest.ItemNumber,
			Quantity:       itemRequest.Quantity,
			UnitPrice:      itemRequest.UnitCost,
			BillableNumber: billable.Number,
			Number:         uuid.NewV4().String(),
		}, tx)

		if err != nil {
			return errors.New("Saving billable item error: " + err.Error())
		}
	}

	return nil
}

func NewBillableService() BillableService {
	return &billableService{}
}
