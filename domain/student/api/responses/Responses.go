package responses

import (
	authResponses "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/api/responses"
)

type DetailedStudentProfile struct {
	APN                      string                 `gorm:"unique;size:100;" json:"apn"`
	ADN                      string                 `gorm:"unique;size:100;default:null" json:"adn"`
	APNStatus                string                 `gorm:"size:20;default:Inactive" json:"apn_status"`
	ADNStatus                string                 `gorm:"size:20;" json:"adn_status"`
	Surname                  string                 `json:"surname" `
	SecondName               string                 `json:"second_name"`
	OtherName                string                 `json:"other_name"`
	Photo                    string                 `json:"photo"`
	AcademicReport           string                 `json:"academic_report"`
	Gender                   string                 `json:"gender" `
	ReasonForChoosingSchool  string                 `json:"reason_for_choosing_school"`
	ApplicationNumber        string                 `json:"application_number"`
	ReasonForAdmissionDenial string                 `json:"reason_for_admission_denial"`
	AttentionNeeded          string                 `json:"attention_needed"`
	DateOfBirth              string                 `json:"date_of_birth"`
	CurrentEnrollment        responses.Enrollment   `json:"current_enrollment"`
	UserProfile              authResponses.UserView `json:"user_profile"`
	UserNumber               string                 `json:"user_number"`
}
