package requests

import "errors"

type EnrollmentRequest struct {
	StudentNumber      string           `json:"student_number"`
	StudyPeriodNumber  string           `json:"study_period_number"`
	CampusNumber       string           `json:"campus_number"`
	SchoolNumber       string           `json:"school_number"`
	EnrollmentDate     string           `json:"enrollment_date"`
	Status             string           `json:"status"`
	CampusClassNumber  string           `json:"campus_class_number"`
	ContinuingStudent  bool             `json:"continuing_student"`
	StreamNumber       string           `json:"stream_number"`
	Items              []EnrollmentItem `json:"items"`
	FeedingGroupNumber string           `json:"feeding_group_number"`
	ModeOfStudyNumber  string           `json:"mode_of_study_number"`
	BillableNumber     string           `json:"billable_number"`
	ClientReference    string           `json:"client_reference"`
	Description        string           `json:"description"`
	PaidAmount         float64          `json:"paid_amount"`
}

type EnrollmentItem struct {
	ItemNumber string  `json:"item_number"`
	Quantity   float64 `json:"quantity"`
	Brought    bool    `json:"brought"`
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
	if e.Status == "" {
		return errors.New("status is required")
	}

	if e.ModeOfStudyNumber == "" {
		return errors.New("Choose mode of study ")
	}
	if e.Status != "active" && e.Status != "inactive" {
		return errors.New("status must be active or inactive")
	}
	if e.StreamNumber == "" {
		return errors.New("stream number is required")
	}

	return nil
}
