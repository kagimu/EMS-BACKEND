package requests

import "errors"

type BillableRequest struct {
	SchoolNumber      string               `json:"school_number"`
	ModeOfStudyNumber string               `json:"mode_of_study_number"`
	CampusNumber      string               `json:"campus_number"`
	StudyPeriodNumber string               `json:"study_period_number"`
	ClassNumber       string               `json:"class_number"`
	BillableItems     []InvoiceItemRequest `json:"billable_items"`
}

type UpdateBillableRequest struct {
	BillableRequest
	BillableNumber string `json:"billable_number"`
}

func (billableRequest *BillableRequest) Validate() error {

	if billableRequest.SchoolNumber == "" {
		return errors.New("Choose a school ")
	}

	if billableRequest.ModeOfStudyNumber == "" {
		return errors.New("Choose mode of study ")
	}

	if billableRequest.CampusNumber == "" {
		return errors.New("Choose campus ")
	}

	if billableRequest.StudyPeriodNumber == "" {
		return errors.New("Choose study period ")
	}

	if billableRequest.ClassNumber == "" {
		return errors.New("Choose class ")
	}

	if len(billableRequest.BillableItems) == 0 {
		return errors.New("Add items on this bill ")
	}

	return ValidateItems(billableRequest.BillableItems)
}
func (updateBillableRequest *UpdateBillableRequest) Validate() error {

	if updateBillableRequest.BillableNumber == "" {
		return errors.New("Choose bill being edited ")
	}

	return updateBillableRequest.BillableRequest.Validate()
}
