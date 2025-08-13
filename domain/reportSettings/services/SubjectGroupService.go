package services

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/api/responses"
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

func GetSubjectGroupsByReportSettingNumber(reportSettingsNumber string) ([]responses.SubjectGroupResponse, error) {
	sgRepo := reporsitory.NewSubjectGroupRepository(config.DB)
	var resp []responses.SubjectGroupResponse
	var subjectGroups = make([]models.SubjectGroup, 0)
	var err error
	subjectGroups, err = sgRepo.FindSubjectGroupsByReportSettingNumber(reportSettingsNumber)
	if err != nil {
		return resp, err
	}

	for _, subjectGroup := range subjectGroups {
		assessmentAreas, err := sgRepo.FindAssessmentAreasByGroupNumber(subjectGroup.Number)
		if err != nil {
			return resp, err
		}

		var groupAssessmentArea []responses2.AssessmentAreaResponse

		for _, assessmentArea := range assessmentAreas {
			asArea, err := repositories.GetAssessmentAreaByNumber(assessmentArea.AssessmentAreaNumber)
			if err != nil {
				return resp, err
			}

			groupAssessmentArea = append(groupAssessmentArea, responses2.AssessmentAreaResponse{
				Name:   asArea.Name,
				Number: asArea.Number,
			})

		}

		resp = append(resp, responses.SubjectGroupResponse{
			Number:          subjectGroup.Number,
			Name:            subjectGroup.SubjectGroupName,
			AssessmentAreas: groupAssessmentArea,
		})
	}

	return resp, nil

}

func DeleteSubjectGroup(subjectGroupNumber string) error {
	txn := config.DB.Begin()
	sgRepo := reporsitory.NewSubjectGroupRepository(txn)
	subjectGroup, err := sgRepo.FindSubjectGroupByNumber(subjectGroupNumber)
	if err != nil {
		txn.Rollback()
		return err
	}

	err = sgRepo.DeleteGroupAssessmentAreaByGroupNumber(subjectGroup.Number)
	if err != nil {
		txn.Rollback()
		return err
	}

	err = sgRepo.DeleteGroup(subjectGroup.Number)
	if err != nil {
		txn.Rollback()
		return err
	}

	txn.Commit()

	return nil

}

func DeleteSubjectGroupAssessment(assessmentArea string) error {
	txn := config.DB.Begin()
	sgRepo := reporsitory.NewSubjectGroupRepository(txn)

	err := sgRepo.DeleteGroupAssessmentAreaByAssessmentAreaNumber(assessmentArea)
	if err != nil {
		txn.Rollback()
		return err
	}

	txn.Commit()

	return nil

}
