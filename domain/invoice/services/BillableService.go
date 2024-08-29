package services

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	repositories4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/repositories"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/invoice/repositories"
	repositories5 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/repositories"
	repositories3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/repositories"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type BillableService interface {
	NewBillable(billableRequest requests.BillableRequest, currentUserNumber string) error
	FindAllBillable(schoolNumber, campusNumber, billType, classNumber, studyPeriodNumber, academicYearNumber, status string) ([]responses.BillableResponse, error)
	BillableDetails(billableNumber string) (responses.DetailedBillableResponse, error)
}

type billableService struct {
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
	billResponse.BillableResponse = createBillResponse(billable)

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
		billsResponse = append(billsResponse, createBillResponse(bill))
	}

	return billsResponse, nil
}

func createBillResponse(billable models.Billable) responses.BillableResponse {

	studyPeriod, _ := repositories3.GetStudyPeriodByNumber(billable.StudyPeriodNumber)
	studyPeriodName, _ := repositories3.GetStudyPeriodNameByNumber(studyPeriod.StudyPeriodNameNumber)
	class, _ := repositories2.GetClassByNumber(billable.ClassNumber)
	campus, _ := repositories4.GetCampusByNumber(billable.CampusNumber)

	return responses.BillableResponse{
		BillableNumber:    billable.Number,
		Status:            billable.Status,
		ClassName:         class.Name,
		ClassNumber:       class.Number,
		StudyPeriodNumber: billable.StudyPeriodNumber,
		StudyPeriodName:   studyPeriodName.Name,
		CampusNumber:      campus.Number,
		CampusName:        campus.Name,
		Type:              billable.BillType,
	}
}

func (billableService *billableService) NewBillable(billableRequest requests.BillableRequest, currentUserNumber string) error {
	billable, err := repositories.GetDuplicatedBill(billableRequest.StudyPeriodNumber, billableRequest.ClassNumber, billableRequest.CampusNumber, billableRequest.BillType, billableRequest.ModeOfStudyNumber, constants.StatusActive)
	if err != nil {
		return errors.New("Checking duplicates error " + err.Error())
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
		checkDuplicate, err := repositories.GetDuplicatedBillableItem(billable.Number, itemRequest.ItemNumber)
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
