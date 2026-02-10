package requests

import (
	"errors"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/models"
	uuid "github.com/satori/go.uuid"
)

type AdmissionRequirementsRequest struct {
	ItemNumbers        []string `json:"item_numbers" binding:"required"`
	ClassNumber        string   `json:"class_number" binding:"required"`
	CampusNumber       string   `json:"campus_number" binding:"required"`
	StudyPeriodNumber  string   `json:"study_period_number" binding:"required"`
	AcademicYearNumber string   `json:"academic_year_number" binding:"required"`
	SchoolNumber       string   `json:"school_number" binding:"required"`
}

func (a *AdmissionRequirementsRequest) Validate() error {
	if a.SchoolNumber == "" {
		return errors.New("school number is required")
	}
	if a.StudyPeriodNumber == "" {
		return errors.New("study period number is required")
	}
	if a.ClassNumber == "" {
		return errors.New("class number is required")
	}
	if a.CampusNumber == "" {
		return errors.New("campus number is required")
	}
	if a.AcademicYearNumber == "" {
		return errors.New("academic year number is required")
	}
	if len(a.ItemNumbers) == 0 {
		return errors.New("at least one item number is required")
	}
	return nil
}

// ToModels builds one AdmissionRequirement per item_number, with shared class/period/year/school and given user.
func (a *AdmissionRequirementsRequest) ToModels(userNumber string) []models.AdmissionRequirement {
	out := make([]models.AdmissionRequirement, 0, len(a.ItemNumbers))
	for _, itemNumber := range a.ItemNumbers {
		out = append(out, models.AdmissionRequirement{
			Number:             uuid.NewV4().String(),
			ItemNumber:         itemNumber,
			CampusClassNumber:  a.ClassNumber,
			CampusNumber:       a.CampusNumber,
			StudyPeriodNumber:  a.StudyPeriodNumber,
			AcademicYearNumber: a.AcademicYearNumber,
			SchoolNumber:       a.SchoolNumber,
			UserNumber:         userNumber,
		})
	}
	return out
}
