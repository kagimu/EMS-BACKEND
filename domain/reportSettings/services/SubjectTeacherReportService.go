package services

import (
	"errors"

	userModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	repositories3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/reporsitory"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/repositories"
	uuid "github.com/satori/go.uuid"
)

func CreateSubjectTeacherReport(req requests.SubjectTeacherReportRequest, user *userModel.User) error {

	report := models.SubjectTeacherReport{
		Number:               uuid.NewV4().String(),
		SubjectTeacherName:   user.FirstName + " " + user.LastName,
		Status:               constants.StatusDraft,
		ClassTeacherComment:  "",
		ClassNumber:          req.ClassNumber,
		AssessmentAreaNumber: req.AssessmentAreaNumber,
		StudyPeriodNumber:    req.StudyPeriodNumber,
		Stage:                constants.StageSubjectTeacher,
	}

	return reporsitory.SaveSubjectTeacherReport(report)
}

func GetSubjectTeacherReports(studyPeriodNumber, stage, status, classNumber, assessmentAreaNumber string) ([]responses.SubjectTeacherReportResponse, error) {
	all, err := reporsitory.FetchAllSubjectTeacherReports(studyPeriodNumber, stage, status, classNumber, assessmentAreaNumber)
	resp := make([]responses.SubjectTeacherReportResponse, 0)
	for _, report := range all {

		studyPeriod, err := repositories.GetStudyPeriodNameByStudyPeriodNumber(report.StudyPeriodNumber)
		if err != nil {
			return nil, err
		}

		assessmentArea, err := repositories2.GetAssessmentAreaByNumber(report.AssessmentAreaNumber)

		if err != nil {
			return nil, err
		}

		class, err := repositories3.GetClassByNumber(report.ClassNumber)
		if err != nil {
			return nil, err
		}

		resp = append(resp, responses.SubjectTeacherReportResponse{
			Number:         report.Number,
			SubjectTeacher: report.SubjectTeacherName,
			Status:         report.Status,
			Class:          class.Name,
			AssessmentArea: assessmentArea.Name,
			Stage:          report.Stage,
			StudyPeriod:    studyPeriod.Name,
			SubmittedOn:    report.CreatedAt.String(),
			UpdatedOn:      report.UpdatedAt.String(),
		})
	}

	return resp, nil
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
