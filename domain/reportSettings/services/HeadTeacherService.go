package services

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/reporsitory"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/repositories"
)

func CreateHeadTeacherComment(headTeacherCommentRequest requests.HeadTeacherCommentRequest, reportSettingNumber, userNumber string) error {
	reportSummary, err := reporsitory.GetReportSettingByNumber(reportSettingNumber)
	if err != nil {
		return errors.New("Checking report settings error: " + err.Error())
	}

	if reportSummary.Number == "" {
		return errors.New("Invalid report settings number ")
	}
	var headTeacherComment models.HeadTeacherComment
	if headTeacherCommentRequest.Number != 0 {
		headTeacherComment, err = reporsitory.GetHeadTeacherCommentByNumber(headTeacherCommentRequest.Number)

		if err != nil {
			return errors.New("failed to get head teacher comment")
		}

		if headTeacherComment.ID == 0 {
			return errors.New("head teacher comment not found")
		}

		headTeacherComment.MinAggregate = headTeacherCommentRequest.MinAggregate
		headTeacherComment.MaxAggregate = headTeacherCommentRequest.MaxAggregate
		headTeacherComment.Comment = headTeacherCommentRequest.Comment
		headTeacherComment.ReportSettingNumber = reportSettingNumber
		headTeacherComment.UserNumber = userNumber

	} else {
		headTeacherComment = models.HeadTeacherComment{
			MinAggregate:        headTeacherCommentRequest.MinAggregate,
			MaxAggregate:        headTeacherCommentRequest.MaxAggregate,
			ReportSettingNumber: reportSettingNumber,
			Comment:             headTeacherCommentRequest.Comment,
			UserNumber:          userNumber,
		}

	}

	err = reporsitory.SaveHeadTeacherComment(headTeacherComment)
	if err != nil {
		return errors.New("failed to create head teacher comment")
	}

	return nil
}

func FetchHeadTeacherComments() ([]responses.HeadTeacherCommentResponse, error) {
	headTeacherComments, err := reporsitory.FindHeadTeacherComments()
	if err != nil {
		return nil, err

	}

	var comments = make([]responses.HeadTeacherCommentResponse, 0)

	for _, comment := range headTeacherComments {
		comments = append(comments, BuildHeadTeacherResponse(comment))
	}
	return comments, nil
}

func BuildHeadTeacherResponse(comment models.HeadTeacherComment) responses.HeadTeacherCommentResponse {

	studyPeriod, _ := repositories.GetStudyPeriodByNumber(comment.ReportSetting.StudyPeriodNumber)

	studyPeriodName, _ := repositories.GetStudyPeriodNameByNumber(studyPeriod.StudyPeriodNameNumber)

	return responses.HeadTeacherCommentResponse{
		ID:                  int64(comment.ID),
		MinAggregate:        comment.MinAggregate,
		MaxAggregate:        comment.MaxAggregate,
		Comment:             comment.Comment,
		AcademicYear:        comment.ReportSetting.AcademicYear.AcademicYearName,
		StudyPeriod:         studyPeriodName.Name,
		Class:               comment.ReportSetting.CampusClass.ClassName,
		StudyPeriodNumber:   comment.ReportSetting.StudyPeriodNumber,
		ReportSettingNumber: comment.ReportSetting.Number,
		AcademicYearNumber:  comment.ReportSetting.AcademicYearNumber,
		CampusClassNumber:   comment.ReportSetting.CampusClassNumber,
	}
}

func DeleteHeadTeacherComment(commentNumber uint) error {
	comment, err := reporsitory.GetHeadTeacherCommentByNumber(commentNumber)

	if err != nil {
		return errors.New("failed to get head teacher comment")
	}

	err = reporsitory.DeleteHeadTeacherComment(comment.ID)
	if err != nil {
		return errors.New("failed to delete head teacher comment")
	}

	return nil
}
