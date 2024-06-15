package services

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/grades/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/grades/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/grades/repository"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateGrade(ctx *gin.Context, request requests.GradeDataRequest) error {
	// create grade
	var grade = models.Grade{
		Number:     uuid.NewV4().String(),
		Name:       request.Name,
		Code:       request.Code,
		UserNumber: ctx.GetHeader("User-Number"),
	}

	if err := config.DB.Create(&grade).Error; err != nil {
		return err
	}
	return nil
}

func FetchGrades(ctx *gin.Context) ([]models.Grade, error) {
	var grades []models.Grade
	grades, err := repository.FindAllGrades()
	if err != nil {
		return nil, err
	}
	return grades, nil
}

func UpdateGrade(ctx *gin.Context, request requests.EditGradeDataRequest) error {
	grade, err := repository.FindGradeByCode(request.Code)
	if err != nil {
		return err
	}

	grade.Name = request.Name
	grade.Code = request.Code
	grade.UserNumber = ctx.GetHeader("User-Number")

	if err := repository.UpdateGrade(grade); err != nil {
		return err
	}
	return nil
}
