package requests

import (
	"errors"
	"fmt"
	"strings"
)

type BulkLearnerItem struct {
	Surname          string `json:"surname"`
	SecondName       string `json:"second_name"`
	OtherName        string `json:"other_name"`
	Gender           string `json:"gender"`
	DateOfBirth      string `json:"date_of_birth"`
	Adn              string `json:"adn"`
	Lin              string `json:"lin"`
	ParentSurname    string `json:"parent_surname"`
	ParentSecondName string `json:"parent_second_name"`
	ParentPhone      string `json:"parent_phone"`
	ParentEmail      string `json:"parent_email"`
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
		r.Learners[i].Surname = strings.TrimSpace(r.Learners[i].Surname)
		r.Learners[i].SecondName = strings.TrimSpace(r.Learners[i].SecondName)
		r.Learners[i].OtherName = strings.TrimSpace(r.Learners[i].OtherName)
		r.Learners[i].Gender = strings.TrimSpace(r.Learners[i].Gender)
		if strings.EqualFold(r.Learners[i].Gender, "m") {
			r.Learners[i].Gender = "Male"
		} else if strings.EqualFold(r.Learners[i].Gender, "f") {
			r.Learners[i].Gender = "Female"
		} else if strings.EqualFold(r.Learners[i].Gender, "male") {
			r.Learners[i].Gender = "Male"
		} else if strings.EqualFold(r.Learners[i].Gender, "female") {
			r.Learners[i].Gender = "Female"
		}

		r.Learners[i].DateOfBirth = strings.TrimSpace(r.Learners[i].DateOfBirth)
		r.Learners[i].Adn = strings.TrimSpace(r.Learners[i].Adn)
		r.Learners[i].Lin = strings.TrimSpace(r.Learners[i].Lin)
		r.Learners[i].ParentSurname = strings.TrimSpace(r.Learners[i].ParentSurname)
		r.Learners[i].ParentSecondName = strings.TrimSpace(r.Learners[i].ParentSecondName)
		r.Learners[i].ParentPhone = strings.TrimSpace(r.Learners[i].ParentPhone)
		r.Learners[i].ParentEmail = strings.TrimSpace(r.Learners[i].ParentEmail)
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
