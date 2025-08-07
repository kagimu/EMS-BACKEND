package services

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/reporsitory"
	uuid "github.com/satori/go.uuid"
)

func CreateSubjectGroup(req requests.CreateSubjectGroupRequest) error {
	txn := config.DB.Begin()

	sgRepo := reporsitory.NewSubjectGroupRepository(txn)

	subjectGroup := models.SubjectGroup{
		Number:                 uuid.NewV4().String(),
		SubjectGroupName:       req.SubjectGroupName,
		ReportSettingsNumber:   req.ReportSettingsNumber,
		ClassNumber:            req.ClassNumber,
		CampusNumber:           req.CampusNumber,
		CurriculumNumber:       req.CurriculumNumber,
		CampusCurriculumNumber: req.CampusCurriculumNumber,
	}

	err := sgRepo.Save(subjectGroup)

	if err != nil {
		txn.Rollback()
		return err
	}

	for _, area := range req.AssessmentAreas {

		var groupAssessmentArea models.GroupAssessmentArea

		groupAssessmentArea.AssessmentAreaNumber = area.AssessmentAreaNumber
		groupAssessmentArea.GroupNumber = subjectGroup.Number

		err := sgRepo.SaveGroupAssessmentArea(groupAssessmentArea)
		if err != nil {
			txn.Rollback()
			return err
		}
	}

	txn.Commit()

	return nil
}
