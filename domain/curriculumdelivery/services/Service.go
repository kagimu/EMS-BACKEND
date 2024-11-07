package services

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculumdelivery/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculumdelivery/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculumdelivery/repository"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func SaveCurriculumDeliverySession(ctx *gin.Context, req requests.CurriculumDeliveryRequest) error {

	db := config.DB.Begin()
	repo := repository.NewCurriculumDeliveryRepositoryImpl(db)

	for _, val := range req.LessonPlanStepNumbers {

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

		err = repo.SaveSession(&curriculumDelivery)
		if err != nil {
			return err
		}

		// Save student session
		for _, student := range req.Students {
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
