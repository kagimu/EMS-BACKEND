package repository

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/grades/models"
)

func FindAllGrades() ([]models.Grade, error) {
	var grades []models.Grade
	if err := config.DB.Find(&grades).Error; err != nil {
		return nil, err
	}
	return grades, nil
}

func FindGradeByCode(code int64) (*models.Grade, error) {
	var grade models.Grade
	if err := config.DB.Where("code = ?", code).First(&grade).Error; err != nil {
		return nil, err
	}
	return &grade, nil
}

func UpdateGrade(grade *models.Grade) error {
	if err := config.DB.Save(grade).Error; err != nil {
		return err
	}
	return nil
}
