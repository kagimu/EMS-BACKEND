package responses

import (
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/api/responses"
	authResponses "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
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

type SingleStudent struct {
	responses.EnrolledStudent
	Assessments []responses2.AssessmentResponse
}
type AssessmentAndStudentResponse struct {
	Headers  []SingleAssessment `json:"headers"`
	Students []SingleStudent    `json:"students"`
}

type SingleAssessment struct {
	AssessmentScheduleNumber string `json:"assessment_schedule_number"`
	AssessmentScheduleName   string `json:"assessment_schedule_name"`
}

type EmergencyInformationResponse struct {
	StudentNumber                 string `json:"student_number" `
	ContactPerson                 string `json:"contact_person" `
	ContactPersonPhone            string `json:"contact_person_phone" `
	ContactPersonAlternative      string `json:"contact_person_alternative" `
	ContactPersonAlternativePhone string `json:"contact_person_alternative_phone" `
	ContactDoctorName             string `json:"contact_doctor_name" `
	ContactDoctorPhone            string `json:"contact_doctor_phone" `
}

type MedicalRecordResponse struct {
	StudentNumber     string `json:"student_number"`
	BloodTypeGroup    string `json:"blood_type_group"`
	Allergies         string `json:"allergies" `
	ExistingCondition string `json:"existing_condition"`
	RegularMedication string `json:"regular_medication"`
	DoctorName        string `json:"doctor_name"`
	DoctorPhone       string `json:"doctor_phone"`
}

type FamilyMemberResponse struct {
	Number                 string                       `json:"number"`
	Surname                string                       `json:"surname"`
	IsLiving               string                       `json:"is_living"`
	OtherName              string                       `json:"other_name"`
	PlaceOfEmployment      string                       `json:"place_of_employment"`
	PhysicalContactAddress string                       `json:"physical_contact_address"`
	Contacts               []models.FamilyMemberContact `json:"contacts"`
	Relationship           string                       `json:"relationship"`
}

type Language struct {
	Name string
	Code int64
}

type FamilyStructureResponse struct {
	StudentNumber string                 `json:"student_number" gorm:"size:36;not null"`
	FamilyMembers []FamilyMemberResponse `json:"family_members" `
	Languages     []Language             `json:"languages"`
}
