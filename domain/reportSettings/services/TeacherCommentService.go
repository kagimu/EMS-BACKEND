package services

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/repositories"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/repository"
	repositories4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	repositories3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/reporsitory"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/repositories"
	uuid "github.com/satori/go.uuid"
	"time"
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

func SaveClassTeacherComment(payload requests.ClassTeacherCommentRequest, user *authModel.User) error {
	date, err := commons.ConvertStringToDate(payload.Date)
	if err != nil {
		return errors.New("invalid date format: " + err.Error())
	}
	classTeacherComment := models.ClassTeacherComment{
		Comment:           payload.Comment,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		UserNumber:        user.Number,
		StudyPeriodNumber: payload.StudyPeriodNumber,
		Date:              date,
		ClassNumber:       payload.CampusClassNumber,
		Number:            uuid.NewV4().String(),
	}

	repo := reporsitory.NewClassTeacherCommentsRepository(config.DB)

	err = repo.CreatClasTeacherComment(classTeacherComment)
	if err != nil {
		return errors.New("failed to save class teacher comment " + err.Error())
	}

	return nil
}

func FetchClassTeacherComments() ([]responses.ClassTeacherCommentResponse, error) {
	repo := reporsitory.NewClassTeacherCommentsRepository(config.DB)

	classTeacherComments, err := repo.FindClassTeacherComments()
	if err != nil {
		return nil, err
	}

	var comments = make([]responses.ClassTeacherCommentResponse, 0)

	for _, comment := range classTeacherComments {

		studyPeriod, err := repositories.GetStudyPeriodByNumber(comment.StudyPeriodNumber)
		if err != nil {
			return nil, errors.New("failed to get study period: " + err.Error())
		}

		studyPeriodName, err := repositories.GetStudyPeriodNameByNumber(studyPeriod.StudyPeriodNameNumber)
		if err != nil {
			return nil, errors.New("failed to get study period name: " + err.Error())
		}

		academicYear, err := repositories2.GetAcademicYearByNumber(studyPeriod.AcademicYearNumber)
		if err != nil {
			return nil, errors.New("failed to get academic year: " + err.Error())
		}

		user, err := repository.GetUserByNumber(comment.UserNumber)

		if err != nil {
			return nil, errors.New("failed to get user by number: " + err.Error())
		}

		class, err := repositories4.GetClassByNumber(comment.ClassNumber)

		if err != nil {
			return nil, errors.New("failed to get class by number: " + err.Error())
		}

		comments = append(comments, responses.ClassTeacherCommentResponse{
			Number:             comment.Number,
			Comment:            comment.Comment,
			UserNumber:         comment.UserNumber,
			StudyPeriod:        studyPeriodName.Name,
			StudyPeriodNumber:  comment.StudyPeriodNumber,
			CampusClassNumber:  comment.ClassNumber,
			ClassName:          class.Name,
			AcademicYear:       academicYear.AcademicYearName,
			AcademicYearNumber: academicYear.Number,
			TeacherName:        user.FirstName + " " + user.LastName + " " + user.SecondName,
		})
	}

	return comments, nil
}

func DeleteClassTeacherComment(number string) error {
	repo := reporsitory.NewClassTeacherCommentsRepository(config.DB)

	err := repo.DeleteClassTeacherComment(number)
	if err != nil {
		return errors.New("failed to delete class teacher comment")
	}

	return nil
}

func SaveSubjectTeacherComment(payload requests.SubjectTeacherCommentRequest, user *authModel.User) error {
	repo := reporsitory.NewClassTeacherCommentsRepository(config.DB)
	subjectTeacherComment := models.SubjectTeacherComment{
		Comment:              payload.Comment,
		UserNumber:           user.Number,
		StudyPeriodNumber:    payload.StudyPeriodNumber,
		AssessmentAreaNumber: payload.AssessmentAreaNumber,
	}

	err := repo.CreateSubjectTeacherComment(subjectTeacherComment)
	if err != nil {
		return errors.New("failed to save subject teacher comment")
	}
	return nil
}

func FetchSubjectTeacherComments() ([]responses.SubjectTeacherCommentResponse, error) {
	repo := reporsitory.NewClassTeacherCommentsRepository(config.DB)
	var err error
	subjectTeacherComments, err := repo.FindSubjectTeacherComments()
	if err != nil {
		return nil, err
	}

	var comments = make([]responses.SubjectTeacherCommentResponse, 0)

	for _, comment := range subjectTeacherComments {

		studyPeriod, err := repositories.GetStudyPeriodByNumber(comment.StudyPeriodNumber)
		if err != nil {
			return nil, errors.New("failed to get study period: " + err.Error())
		}

		studyPeriodName, err := repositories.GetStudyPeriodNameByNumber(studyPeriod.StudyPeriodNameNumber)
		if err != nil {
			return nil, errors.New("failed to get study period name: " + err.Error())
		}

		academicYear, err := repositories2.GetAcademicYearByNumber(studyPeriod.AcademicYearNumber)
		if err != nil {
			return nil, errors.New("failed to get academic year: " + err.Error())
		}

		user, err := repository.GetUserByNumber(comment.UserNumber)

		if err != nil {
			return nil, errors.New("failed to get user by number: " + err.Error())
		}

		assessmentArea, err := repositories3.GetAssessmentAreaByNumber(comment.AssessmentAreaNumber)

		if err != nil {
			return nil, errors.New("failed to get assessment area: " + err.Error())
		}

		comments = append(comments, responses.SubjectTeacherCommentResponse{
			Number:               comment.Number,
			Comment:              comment.Comment,
			UserNumber:           comment.UserNumber,
			StudyPeriod:          studyPeriodName.Name,
			AcademicYearName:     academicYear.AcademicYearName,
			StudyPeriodNumber:    comment.StudyPeriodNumber,
			AcademicYearNumber:   academicYear.Number,
			TeacherName:          user.FirstName + " " + user.LastName + " " + user.SecondName,
			AssessmentAreaNumber: comment.AssessmentAreaNumber,
			AssessmentAreaName:   assessmentArea.Name,
		})
	}

	return comments, nil
}

func DeleteSubjectTeacherComment(number string) error {
	repo := reporsitory.NewClassTeacherCommentsRepository(config.DB)

	err := repo.DeleteSubjectTeacherComment(number)
	if err != nil {
		return errors.New("failed to delete subject teacher comment")
	}

	return nil
}
