package services

import (
	"errors"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/reporsitory"
	uuid "github.com/satori/go.uuid"
)

func CreateSubjectTeacherReport(req requests.SubjectTeacherReportRequest) error {
	report := models.SubjectTeacherReport{
		Number:               uuid.NewV4().String(),
		SubjectTeacherName:   req.SubjectTeacherName,
		Status:               constants.StatusDraft,
		ClassTeacherComment:  req.ClassTeacherComment,
		AssessmentAreaNumber: req.AssessmentAreaNumber,
		StudyPeriodNumber:    req.StudyPeriodNumber,
		Stage:                constants.StageSubjectTeacher,
	}

	return reporsitory.SaveSubjectTeacherReport(report)
}

func GetSubjectTeacherReports(studyPeriodNumber, stage, status string) ([]models.SubjectTeacherReport, error) {
	return reporsitory.FetchAllSubjectTeacherReports(studyPeriodNumber, stage, status)
}

func UpdateSubjectTeacherReport(number string, req requests.UpdateSubjectTeacherReportRequest) error {
	report, err := reporsitory.GetSubjectTeacherReportByNumber(number)
	if err != nil {
		return err
	}
	if report.Number == "" {
		return errors.New("report not found")
	}

	if req.SubjectTeacherName != "" {
		report.SubjectTeacherName = req.SubjectTeacherName
	}
	if req.Status != "" {
		report.Status = req.Status
	}
	if req.ClassTeacherComment != "" {
		report.ClassTeacherComment = req.ClassTeacherComment
	}
	if req.AssessmentAreaNumber != "" {
		report.AssessmentAreaNumber = req.AssessmentAreaNumber
	}
	if req.StudyPeriodNumber != "" {
		report.StudyPeriodNumber = req.StudyPeriodNumber
	}
	if req.Stage != "" {
		report.Stage = req.Stage
	}

	return reporsitory.SaveSubjectTeacherReport(report)
}

func DeleteSubjectTeacherReport(number string) error {
	report, err := reporsitory.GetSubjectTeacherReportByNumber(number)
	if err != nil {
		return err
	}
	if report.Number == "" {
		return errors.New("report not found")
	}

	return reporsitory.DeleteSubjectTeacherReport(report)
}

func ForwardSubjectTeacherReport(req requests.ForwardSubjectTeacherReportRequest, commenterName string) error {
	report, err := reporsitory.GetSubjectTeacherReportByNumber(req.SubjectTeacherReportNumber)
	if err != nil {
		return err
	}
	if report.Number == "" {
		return errors.New("report not found")
	}

	report.Stage = req.NextStage

	report.Status = constants.StatusDraft
	if err := reporsitory.SaveSubjectTeacherReport(report); err != nil {
		return err
	}

	comment := models.SubjectTeacherReportComment{
		Number:                     uuid.NewV4().String(),
		SubjectTeacherReportNumber: req.SubjectTeacherReportNumber,
		CommenterName:              commenterName,
		Comment:                    req.Comment,
	}
	return reporsitory.SaveSubjectTeacherReportComment(comment)
}

func ApproveSubjectTeacherReport(number string) error {
	report, err := reporsitory.GetSubjectTeacherReportByNumber(number)
	if err != nil {
		return err
	}
	if report.Number == "" {
		return errors.New("report not found")
	}

	report.Status = constants.StatusFinal
	return reporsitory.SaveSubjectTeacherReport(report)
}

func GetSubjectTeacherReportComments(reportNumber string) ([]models.SubjectTeacherReportComment, error) {
	return reporsitory.FetchSubjectTeacherReportComments(reportNumber)
}
