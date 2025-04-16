package request

import "fmt"

type EmergencyInformationRequest struct {
	StudentNumber string `json:"student_number" gorm:"size:36;not null"`
}

func (req *EmergencyInformationRequest) Validate() error {
	if req.StudentNumber == "" {
		return fmt.Errorf("student_number is required")
	}
	return nil
}
