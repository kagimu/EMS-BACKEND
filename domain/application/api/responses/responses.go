package responses

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/api/responses"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
)

type ApplicationResponse struct {
	Number                    string                   `json:"number"`
	Status                    string                   `json:"status"`
	ApplicationFeePayment     bool                     `json:"application_fee_payment"`
	ApplicationFeeInvalidDate string                   `json:"application_fee_invalid_date"`
	UserNumber                string                   `json:"user_number"`
	InterviewMode             string                   `json:"interview_mode"`
	Student                   StudentApplicant         `json:"student"`
	Parent                    ParentResponse           `json:"parent"`
	PreviousSchools           []PreviousSchoolResponse `json:"previous_schools"`
}

type StudentApplicant struct {
	Number                    string                   `json:"number"`
	APN                       string                   `json:"apn"`
	ADN                       string                   `json:"adn"`
	APNStatus                 string                   `json:"apn_status"`
	ADNStatus                 string                   `json:"adn_status"`
	Surname                   string                   `json:"surname" `
	SecondName                string                   `json:"second_name"`
	OtherName                 string                   `json:"other_name"`
	Photo                     string                   `json:"photo"`
	AcademicReport            string                   `json:"academic_report"`
	Gender                    string                   `json:"gender"`
	ReasonForChoosingSchool   string                   `json:"reason_for_choosing_school"`
	ReasonForAdmissionDenial  string                   `json:"reason_for_admission_denial"`
	AttentionNeeded           string                   `json:"attention_needed"`
	DateOfBirth               string                   `json:"date_of_birth"`
	PreferredTypeOfEnrollment string                   `json:"preferred_type_of_enrollment"`
	PreferredCampusNumber     string                   `json:"preferred_campus_number"`
	PreferredCampus           responses.CampusResponse `json:"preferred_campus"`
	CampusClassNumber         string                   `json:"campus_class_number"`
	CampusClass               classModels.CampusClass  `json:"campus_class"`
	SourceOfSchoolInfo        string                   `json:"source_of_school_info"`
	UserNumber                string                   `json:"user_number"`
	User                      models.User              `json:"user"`
	Status                    string                   `json:"status"`
}
type PreviousSchoolResponse struct {
}

type ParentResponse struct {
}
