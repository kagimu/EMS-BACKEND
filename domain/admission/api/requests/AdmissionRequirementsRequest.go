package requests

import (
	"errors"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/models"
)

type AdmissionRequirementsRequest struct {
	SchoolNumber       string   `json:"school_number" binding:"required"`
	Items              []string `json:"item_numbers" binding:"required"`
	ClassNumber        string   `json:"class_number" binding:"required"`
	StudyPeriodNumber  string   `json:"study_period_number" binding:"required"`
	AcademicYearNumber string   `json:"academic_year_number" binding:"required"`
}

func (a *AdmissionRequirementsRequest) Validate() error {
	if a.SchoolNumber == "" {
		return errors.New("study period number is required")
	}
	if a.StudyPeriodNumber == "" {
		return errors.New("study period number is required")
	}
	if a.ClassNumber == "" {
		return errors.New("class number is required")
	}
	if a.AcademicYearNumber == "" {
		return errors.New("academic year number is required")
	}
	if len(a.Items) == 0 {
		return errors.New("item numbers is required")
	}
	return nil
}

func (a *AdmissionRequirementsRequest) ToModel() models.AdmissionRequirement {
	return models.AdmissionRequirement{
		SchoolNumber:       a.SchoolNumber,
		CampusClassNumber:  a.ClassNumber,
		StudyPeriodNumber:  a.StudyPeriodNumber,
		AcademicYearNumber: a.AcademicYearNumber,
	}
}
