package services

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/repositories"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/repositories"
	repositories3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	repositories4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/repositories"
	repositories5 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/repositories"
	repositories6 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/repositories"
)

// SaveAdmissionRequirementsService creates one admission requirement per item_number and persists them.
func SaveAdmissionRequirementsService(req *requests.AdmissionRequirementsRequest, userNumber string) error {
	repo := repositories.NewAdmissionRequirementsRepository(config.DB)
	requirements := req.ToModels(userNumber)
	return repo.SaveAdmissionRequirementsBatch(requirements)
}

// createAdmissionRequirementResponse maps a model to the API response with resolved names.
func createAdmissionRequirementResponse(m models.AdmissionRequirement, schoolName, campusName, curriculumName, className, studyPeriodName, itemName string) responses.AdmissionRequirementResponse {
	return responses.AdmissionRequirementResponse{
		Number:             m.Number,
		ItemNumber:         m.ItemNumber,
		ItemName:           itemName,
		CampusClassNumber:  m.CampusClassNumber,
		CampusNumber:       m.CampusNumber,
		StudyPeriodNumber:  m.StudyPeriodNumber,
		AcademicYearNumber: m.AcademicYearNumber,
		SchoolNumber:       m.SchoolNumber,
		SchoolName:         schoolName,
		CampusName:         campusName,
		CurriculumName:     curriculumName,
		ClassName:          className,
		StudyPeriodName:    studyPeriodName,
		UserNumber:         m.UserNumber,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}
}

// FetchAdmissionRequirementsService returns admission requirements filtered by non-empty params, as response DTOs.
func FetchAdmissionRequirementsService(schoolNumber, campusClassNumber, studyPeriodNumber, academicYearNumber, itemNumber string) ([]responses.AdmissionRequirementResponse, error) {
	repo := repositories.NewAdmissionRequirementsRepository(config.DB)
	list, err := repo.FindAdmissionRequirements(schoolNumber, campusClassNumber, studyPeriodNumber, academicYearNumber, itemNumber)
	if err != nil {
		return nil, err
	}
	out := make([]responses.AdmissionRequirementResponse, 0, len(list))
	for _, m := range list {
		schoolName := ""
		if m.SchoolNumber != "" {
			if school, err := repositories5.GetSchoolByNumber(m.SchoolNumber); err == nil {
				schoolName = school.Name
			}
		}
		itemName := ""
		if m.ItemNumber != "" {
			if item, err := repositories4.GetItemByNumber(m.ItemNumber); err == nil && item.Number != "" {
				itemName = item.Name
			}
		}
		campusName := ""
		if m.CampusNumber != "" {
			if campus, err := repositories2.GetCampusByNumber(m.CampusNumber); err == nil {
				campusName = campus.Name
			}
		} else if m.CampusClassNumber != "" {
			if campusClass, err := repositories3.GetCampusClassByNumber(m.CampusClassNumber); err == nil && campusClass.CampusNumber != "" {
				if campus, err := repositories2.GetCampusByNumber(campusClass.CampusNumber); err == nil {
					campusName = campus.Name
				}
			}
		}
		curriculumName := ""
		className := ""

		if m.CampusClassNumber != "" {
			class, err := repositories3.GetClassByNumber(m.CampusClassNumber)
			if err == nil {
				className = class.Name
			}
		}

		studyPeriodName := ""
		if m.StudyPeriodNumber != "" {
			if name, err := repositories6.GetStudyPeriodNameByStudyPeriodNumber(m.StudyPeriodNumber); err == nil {
				studyPeriodName = name.Name
			}
		}

		out = append(out, createAdmissionRequirementResponse(m, schoolName, campusName, curriculumName, className, studyPeriodName, itemName))
	}
	return out, nil
}
