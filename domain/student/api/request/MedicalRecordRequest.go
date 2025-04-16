package request

import "fmt"

type MedicalRecordRequest struct {
	StudentNumber     string `json:"student_number" gorm:"size:36;not null"`
	BloodTypeGroup    string `json:"blood_type_group" gorm:"size:50;not null"`
	Allergies         string `json:"allergies" gorm:"size:50;not null"`
	ExistingCondition string `json:"existing_condition" gorm:"size:100;not null"`
	RegularMedication string `json:"regular_medication" gorm:"size:100;not null"`
	DoctorName        string `json:"doctor_name" gorm:"size:100;not null"`
	DoctorPhone       string `json:"doctor_phone" gorm:"size:100;not null"`
}

func (req *MedicalRecordRequest) Validate() error {
	if req.StudentNumber == "" {
		return fmt.Errorf("student_number is required")
	}
	return nil
}
