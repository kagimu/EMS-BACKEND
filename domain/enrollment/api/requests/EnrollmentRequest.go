package requests

import (
	"errors"
)

type EnrollmentRequest struct {
	StudentNumber     string `json:"student_number"`
	StudyPeriodNumber string `json:"study_period_number"`
	CampusNumber      string `json:"campus_number"`
	SchoolNumber      string `json:"school_number"`
	EnrollmentDate    string `json:"enrollment_date"`
	//Status                 string           `json:"status"`
	CampusClassNumber      string           `json:"campus_class_number"`
	ContinuingStudent      bool             `json:"continuing_student"`
	StreamNumber           string           `json:"stream_number"`
	Items                  []EnrollmentItem `json:"items"`
	FeedingGroupNumber     string           `json:"feeding_group_number"`
	ModeOfStudyNumber      string           `json:"mode_of_study_number"`
	BillableNumber         string           `json:"billable_number"`
	ClientReference        string           `json:"client_reference"`
	Description            string           `json:"description"`
	PaidAmount             float64          `json:"paid_amount"`
	DatePaid               string           `json:"date_paid"`
	CampusCurriculumNumber string           `json:"campus_curriculum_number"`
}

type EnrollmentItem struct {
	ItemNumber string  `json:"item_number"`
	Quantity   float64 `json:"quantity"`
	Brought    bool    `json:"brought"`
}

// UpdateEnrollmentItemRequest is the body for updating an enrollment item (quantity, brought).
type UpdateEnrollmentItemRequest struct {
	Quantity float64 `json:"quantity"`
	Brought  bool    `json:"brought"`
}

func (e *EnrollmentRequest) Validate() error {
	if e.CampusNumber == "" {
		return errors.New("campus number is required")
	}
	if e.SchoolNumber == "" {
		return errors.New("school number is required")
	}
	if e.CampusClassNumber == "" {
		return errors.New("class is required")
	}
	if e.StudentNumber == "" {
		return errors.New("student number is required")
	}
	if e.StudyPeriodNumber == "" {
		return errors.New("study period number is required")
	}
	if e.EnrollmentDate == "" {
		return errors.New("enrollment date is required")
	}
	/*if e.Status == "" {
		return errors.New("status is required")
	}*/

	if e.ModeOfStudyNumber == "" {
		return errors.New("Choose mode of study ")
	}

	if e.BillableNumber == "" {
		return errors.New("No bill selected for this enrollment ")
	}
	if e.PaidAmount < 0 {
		return errors.New("You can not enter amount in negative ")
	}
	if e.ClientReference == "" && e.PaidAmount > 0 {
		return errors.New("Enter payment reference I.E Receipt number or Bank transaction reference ")
	}
	if e.DatePaid == "" && e.PaidAmount > 0 {
		return errors.New("Provide date paid ")
	}
	if e.Description == "" && e.PaidAmount > 0 {
		return errors.New("Describe the payment E.G Bank account, Mobile phone or Cashier ")
	}

	/*if e.Status != constants.StatusActive && e.Status != constants.StatusInactive {
		return errors.New(fmt.Sprintf("status must be %s or %s", constants.StatusActive, constants.StatusInactive))
	}*/
	if e.StreamNumber == "" {
		return errors.New("stream number is required")
	}

	return nil
}
