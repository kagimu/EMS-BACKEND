package services

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	repository2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/repository"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/pagination"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculumdelivery/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculumdelivery/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculumdelivery/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculumdelivery/repository"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/repositories"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func SaveCurriculumDeliverySession(ctx *gin.Context, req requests.CurriculumDeliveryRequest) error {

	db := config.DB.Begin()
	repo := repository.NewCurriculumDeliveryRepositoryImpl(db)

	for _, val := range req.StepNumbers {

		date, err := commons.ConvertStringToDate(req.Date)
		if err != nil {
			return err
		}
		// Save session
		curriculumDelivery := models.CurriculumDelivery{
			Number:               uuid.NewV4().String(),
			LessonPlanStepNumber: val,
			StreamNumber:         req.StreamNumber,
			Date:                 date,
			Time:                 req.Time,
			TeacherNumber:        req.TeacherNumber,
			Comment:              req.Comment,
		}

		sessions, err := repo.FindCurriculumDeliveryByStreamNumberAndStepNumber(req.StreamNumber, val)
		if err != nil {
			return err
		}

		if len(sessions) > 0 {
			db.Rollback()
			return errors.New("session has already been recorded")
		}

		err = repo.SaveSession(&curriculumDelivery)
		if err != nil {
			db.Rollback()
			return err
		}

		// Save student session
		for _, student := range req.Students {

			sessions, err := repo.FindCurriculumDeliveryByStudentNumberAndStepNumber(student.StudentNumber, val)

			if len(sessions) > 0 {
				db.Rollback()
				return errors.New("session for this student has already been recorded")
			}

			if err != nil {
				db.Rollback()
				return err
			}

			studentCurriculumDelivery := models.StudentCurriculumDelivery{
				StudentNumber:            student.StudentNumber,
				CurriculumDeliveryNumber: curriculumDelivery.Number,
				Comment:                  student.Comment,
			}

			err = repo.SaveStudentSession(&studentCurriculumDelivery)
			if err != nil {
				return err
			}
		}

	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return err
	}

	return nil
}

func GetCurriculumDeliveryByStreamNumberAndStepNumber(ctx *gin.Context, pr *pagination.PaginatedResponse) (*[]responses.CurriculumDeliveryResponse, error) {
	repo := repository.NewCurriculumDeliveryRepositoryImpl(config.DB)
	sessions, err := repo.FindCurriculumDeliveryByFilters(ctx, pr)
	if err != nil {
		return nil, err
	}

	var cds []responses.CurriculumDeliveryResponse
	for _, session := range sessions {
		cd, err := BuildCurriculumDeliveryResponse(session)
		if err != nil {
			return nil, err
		}
		cds = append(cds, *cd)
	}

	return &cds, nil

}

func BuildCurriculumDeliveryResponse(cds *models.CurriculumDelivery) (*responses.CurriculumDeliveryResponse, error) {
	var cd responses.CurriculumDeliveryResponse
	cd.Number = cds.Number

	stream, err := repositories2.GetStreamByNumber(cds.StreamNumber)
	if err != nil {
		return &cd, err
	}

	cd.Stream = stream.Name

	class, err := repositories2.GetClassByCampusClassNumber(stream.CampusClassNumber)
	if err != nil {
		return &cd, err
	}
	cd.Class = class.Name
	cd.Date = cds.Date.String()
	cd.Time = cds.Time
	cd.Comment = cds.Comment

	teacher, err := repository2.GetUserByNumber(cds.TeacherNumber)
	if err != nil {
		return &cd, err
	}
	cd.Teacher = teacher.FirstName + "" + teacher.LastName

	return &cd, nil

}

func GetStudentsBySession(context *gin.Context, p *pagination.PaginatedResponse) ([]responses2.LeanStudent, error) {
	stdts := make([]responses2.LeanStudent, 0)
	repo := repository.NewCurriculumDeliveryRepositoryImpl(config.DB)
	students, err := repo.FindStudentsBySession(context, p)
	if err != nil {
		return stdts, err
	}

	for _, student := range students {
		st, err := repositories.FindStudentByNumber(student.StudentNumber)
		if err != nil {
			return stdts, err
		}
		stdts = append(stdts, responses2.LeanStudent{
			Number:     student.StudentNumber,
			Surname:    st.Surname,
			Gender:     st.Gender,
			SecondName: st.SecondName,
			OtherName:  st.OtherName,
		})
	}

	return stdts, nil

}
