package requests

import (
	"errors"
	"fmt"
	"strings"
)

// BulkFamilyContact holds a single phone/email pair for a family member.
type BulkFamilyContact struct {
	Phone string `json:"phone"`
	Email string `json:"email"`
}

type BulkLearnerItem struct {
	Surname     string `json:"surname"`
	SecondName  string `json:"second_name"`
	OtherName   string `json:"other_name"`
	Gender      string `json:"gender"`
	DateOfBirth string `json:"date_of_birth"`
	Adn         string `json:"adn"`
	Lin         string `json:"lin"`

	// Legacy generic parent fields — kept for backward compatibility with old templates
	ParentSurname    string `json:"parent_surname"`
	ParentSecondName string `json:"parent_second_name"`
	ParentPhone      string `json:"parent_phone"`
	ParentEmail      string `json:"parent_email"`

	// Medical
	BloodTypeGroup    string `json:"blood_type_group"`
	Allergies         string `json:"allergies"`
	ExistingCondition string `json:"existing_condition"`
	RegularMedication string `json:"regular_medication"`
	DoctorName        string `json:"doctor_name"`
	DoctorPhone       string `json:"doctor_phone"`

	// Emergency
	EmergencyContactPerson                 string `json:"emergency_contact_person"`
	EmergencyContactPersonPhone            string `json:"emergency_contact_person_phone"`
	EmergencyContactPersonAlternative      string `json:"emergency_contact_person_alternative"`
	EmergencyContactPersonAlternativePhone string `json:"emergency_contact_person_alternative_phone"`
	EmergencyContactDoctorName             string `json:"emergency_contact_doctor_name"`
	EmergencyContactDoctorPhone            string `json:"emergency_contact_doctor_phone"`

	// Language
	PrimaryLanguage   string `json:"primary_language"`
	SecondaryLanguage string `json:"secondary_language"`
	TertiaryLanguage  string `json:"tertiary_language"`

	// Father — surname, second name, other name + up to 3 phone/email contacts
	FatherSurname           string              `json:"father_surname"`
	FatherSecondName        string              `json:"father_second_name"`
	FatherOtherName         string              `json:"father_other_name"`
	FatherIsLiving          string              `json:"father_is_living"`
	FatherPlaceOfEmployment string              `json:"father_place_of_employment"`
	FatherContacts          []BulkFamilyContact `json:"father_contacts"`

	// Mother — surname, second name, other name + up to 3 phone/email contacts
	MotherSurname           string              `json:"mother_surname"`
	MotherSecondName        string              `json:"mother_second_name"`
	MotherOtherName         string              `json:"mother_other_name"`
	MotherIsLiving          string              `json:"mother_is_living"`
	MotherPlaceOfEmployment string              `json:"mother_place_of_employment"`
	MotherContacts          []BulkFamilyContact `json:"mother_contacts"`

	// Guardian (optional) — surname, second name, other name + up to 3 phone/email contacts
	GuardianSurname           string              `json:"guardian_surname"`
	GuardianSecondName        string              `json:"guardian_second_name"`
	GuardianOtherName         string              `json:"guardian_other_name"`
	GuardianIsLiving          string              `json:"guardian_is_living"`
	GuardianPlaceOfEmployment string              `json:"guardian_place_of_employment"`
	GuardianContacts          []BulkFamilyContact `json:"guardian_contacts"`
}

type BulkDirectAdmitRequest struct {
	SchoolNumber       string            `json:"school_number"`
	CampusNumber       string            `json:"campus_number"`
	CampusClassNumber  string            `json:"campus_class_number"`
	StreamNumber       string            `json:"stream_number"`
	AcademicYearNumber string            `json:"academic_year_number"`
	StudyPeriodNumber  string            `json:"study_period_number"`
	Learners           []BulkLearnerItem `json:"learners"`
}

func (r *BulkDirectAdmitRequest) Prepare() {
	r.SchoolNumber = strings.TrimSpace(r.SchoolNumber)
	r.CampusNumber = strings.TrimSpace(r.CampusNumber)
	r.CampusClassNumber = strings.TrimSpace(r.CampusClassNumber)
	r.StreamNumber = strings.TrimSpace(r.StreamNumber)
	r.AcademicYearNumber = strings.TrimSpace(r.AcademicYearNumber)
	r.StudyPeriodNumber = strings.TrimSpace(r.StudyPeriodNumber)

	for i := range r.Learners {
		l := &r.Learners[i]
		l.Surname = strings.TrimSpace(l.Surname)
		l.SecondName = strings.TrimSpace(l.SecondName)
		l.OtherName = strings.TrimSpace(l.OtherName)
		l.Gender = strings.TrimSpace(l.Gender)
		switch strings.ToLower(l.Gender) {
		case "m", "male":
			l.Gender = "Male"
		case "f", "female":
			l.Gender = "Female"
		}
		l.DateOfBirth = strings.TrimSpace(l.DateOfBirth)
		l.Adn = strings.TrimSpace(l.Adn)
		l.Lin = strings.TrimSpace(l.Lin)

		// Legacy parent fields
		l.ParentSurname = strings.TrimSpace(l.ParentSurname)
		l.ParentSecondName = strings.TrimSpace(l.ParentSecondName)
		l.ParentPhone = strings.TrimSpace(l.ParentPhone)
		l.ParentEmail = strings.TrimSpace(l.ParentEmail)

		// Father
		l.FatherSurname = strings.TrimSpace(l.FatherSurname)
		l.FatherSecondName = strings.TrimSpace(l.FatherSecondName)
		l.FatherOtherName = strings.TrimSpace(l.FatherOtherName)
		l.FatherPlaceOfEmployment = strings.TrimSpace(l.FatherPlaceOfEmployment)
		for j := range l.FatherContacts {
			l.FatherContacts[j].Phone = strings.TrimSpace(l.FatherContacts[j].Phone)
			l.FatherContacts[j].Email = strings.TrimSpace(l.FatherContacts[j].Email)
		}

		// Mother
		l.MotherSurname = strings.TrimSpace(l.MotherSurname)
		l.MotherSecondName = strings.TrimSpace(l.MotherSecondName)
		l.MotherOtherName = strings.TrimSpace(l.MotherOtherName)
		l.MotherPlaceOfEmployment = strings.TrimSpace(l.MotherPlaceOfEmployment)
		for j := range l.MotherContacts {
			l.MotherContacts[j].Phone = strings.TrimSpace(l.MotherContacts[j].Phone)
			l.MotherContacts[j].Email = strings.TrimSpace(l.MotherContacts[j].Email)
		}

		// Guardian
		l.GuardianSurname = strings.TrimSpace(l.GuardianSurname)
		l.GuardianSecondName = strings.TrimSpace(l.GuardianSecondName)
		l.GuardianOtherName = strings.TrimSpace(l.GuardianOtherName)
		l.GuardianPlaceOfEmployment = strings.TrimSpace(l.GuardianPlaceOfEmployment)
		for j := range l.GuardianContacts {
			l.GuardianContacts[j].Phone = strings.TrimSpace(l.GuardianContacts[j].Phone)
			l.GuardianContacts[j].Email = strings.TrimSpace(l.GuardianContacts[j].Email)
		}

		// Backward compat: if new father fields are absent, promote legacy parent fields
		if l.FatherSurname == "" && l.ParentSurname != "" {
			l.FatherSurname = l.ParentSurname
		}
		if l.FatherSecondName == "" && l.ParentSecondName != "" {
			l.FatherSecondName = l.ParentSecondName
		}
		if len(l.FatherContacts) == 0 && (l.ParentPhone != "" || l.ParentEmail != "") {
			l.FatherContacts = []BulkFamilyContact{{Phone: l.ParentPhone, Email: l.ParentEmail}}
		}

		// Default living status
		if l.FatherIsLiving == "" {
			l.FatherIsLiving = "Yes"
		}
		if l.MotherIsLiving == "" {
			l.MotherIsLiving = "Yes"
		}
		if l.GuardianIsLiving == "" {
			l.GuardianIsLiving = "Yes"
		}
	}
}

func (r *BulkDirectAdmitRequest) Validate() error {
	if r.SchoolNumber == "" {
		return errors.New("school number is required")
	}
	if r.CampusNumber == "" {
		return errors.New("campus number is required")
	}
	if r.CampusClassNumber == "" {
		return errors.New("campus class number is required")
	}
	if r.StreamNumber == "" {
		return errors.New("stream number is required")
	}
	if r.AcademicYearNumber == "" {
		return errors.New("academic year number is required")
	}
	if r.StudyPeriodNumber == "" {
		return errors.New("study period number is required")
	}
	if len(r.Learners) == 0 {
		return errors.New("at least one learner is required for bulk admission")
	}

	seenAdns := make(map[string]int)
	for i, l := range r.Learners {
		rowNum := i + 1
		if l.Adn == "" {
			return fmt.Errorf("learner row %d: ADN (Admission Number) is required", rowNum)
		}
		if prevRow, exists := seenAdns[l.Adn]; exists {
			return fmt.Errorf("learner row %d: duplicate ADN '%s' in spreadsheet (already entered in row %d)", rowNum, l.Adn, prevRow)
		}
		seenAdns[l.Adn] = rowNum

		if l.Surname == "" {
			return fmt.Errorf("learner row %d: surname is required", rowNum)
		}
		if l.SecondName == "" {
			return fmt.Errorf("learner row %d: first/second name is required", rowNum)
		}
		if l.Gender == "" {
			return fmt.Errorf("learner row %d: gender is required", rowNum)
		}
		if l.Gender != "Male" && l.Gender != "Female" {
			return fmt.Errorf("learner row %d: gender must be Male or Female", rowNum)
		}
	}

	return nil
}
