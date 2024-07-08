package services

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/reporsitory"
)

func CreateHeadTeacherComment(payload requests.HeadTeacherCommentRequest, reportSettingNumber, userNumber string) error {
	var reportSummary models.ReportSetting
	if err := config.DB.Where("number = ?", reportSettingNumber).First(&reportSummary).Error; err != nil {

		return errors.New("failed to fetch given  report setting")
	}

	var headTeacherComment models.HeadTeacherComment
	if err := config.DB.Where("report_setting_number = ?", reportSettingNumber).First(&headTeacherComment).Error; err == nil {
		headTeacherComment.Comment = payload.Comment
		headTeacherComment.UserNumber = userNumber

		err := reporsitory.SaveHeadTeacherComment(headTeacherComment)
		if err != nil {
			return errors.New("failed to update head teacher comment")
		}

	} else {
		headTeacherComment = models.HeadTeacherComment{
			MinAggregate:        payload.MinAggregate,
			MaxAggregate:        payload.MaxAggregate,
			ReportSettingNumber: reportSettingNumber,
			Comment:             payload.Comment,
			UserNumber:          userNumber,
		}

		err := reporsitory.CreateHeadTeacherComment(headTeacherComment)
		if err != nil {
			return errors.New("failed to create head teacher comment")
		}
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
	return responses.HeadTeacherCommentResponse{
		MinAggregate: comment.MinAggregate,
		MaxAggregate: comment.MaxAggregate,
		Comment:      comment.Comment,
		AcademicYear: comment.ReportSetting.AcademicYear.AcademicYearName,
		StudyPeriod:  comment.ReportSetting.StudyPeriod.StudyPeriodName,
		Class:        comment.ReportSetting.CampusClass.ClassName,
	}
}
