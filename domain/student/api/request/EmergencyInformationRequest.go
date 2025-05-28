package request

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
)

type EmergencyInformationRequest struct {
	StudentNumber                 string `json:"student_number" gorm:"size:36;not null"`
	ContactPerson                 string `json:"contact_person" gorm:"size:50;not null"`
	ContactPersonPhone            string `json:"contact_person_phone" gorm:"size:50;not null"`
	ContactPersonAlternative      string `json:"contact_person_alternative" gorm:"size:50;not null"`
	ContactPersonAlternativePhone string `json:"contact_person_alternative_phone" gorm:"size:50;not null"`
	ContactDoctorName             string `json:"contact_doctor_name" gorm:"size:50;not null"`
	ContactDoctorPhone            string `json:"contact_doctor_phone" gorm:"size:50;not null"`
}

func (req *EmergencyInformationRequest) Validate() error {
	if req.StudentNumber == "" {
		return fmt.Errorf("student_number is required")
	}
	return nil
}

func (req *EmergencyInformationRequest) ToModel() models.EmergencyInformation {
	return models.EmergencyInformation{
		StudentNumber:                 req.StudentNumber,
		ContactPerson:                 req.ContactPerson,
		ContactPersonPhone:            req.ContactPersonPhone,
		ContactPersonAlternative:      req.ContactPersonAlternative,
		ContactPersonAlternativePhone: req.ContactPersonAlternativePhone,
		ContactDoctorName:             req.ContactDoctorName,
		ContactDoctorPhone:            req.ContactDoctorPhone,
	}
}
