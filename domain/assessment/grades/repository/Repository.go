package repository

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/grades/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"gorm.io/gorm"
)

func FindAllGrades(schoolNumber string) ([]models.Grade, error) {
	var grades []models.Grade
	if err := config.DB.Scopes(filters.BySchoolNumber(schoolNumber)).Find(&grades).Error; err != nil {
		return nil, err
	}
	return grades, nil
}

func FindGradeByCode(code int64, schoolNumber string) (models.Grade, error) {
	var grade models.Grade
	if err := config.DB.Where("code = ? AND school_number = ?", code, schoolNumber).First(&grade).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return grade, nil
		}

		return grade, err
	}
	return grade, nil
}

func FindGradeByNameAndShortName(name, shortName, schoolNumber string) (models.Grade, error) {
	var grade models.Grade
	if err := config.DB.Where("name = ? AND short_name = ? AND school_number = ?", name, shortName, schoolNumber).First(&grade).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return grade, nil
		}

		return grade, err
	}
	return grade, nil
}

func UpdateGrade(grade models.Grade) error {
	if err := config.DB.Save(grade).Error; err != nil {
		return err
	}
	return nil
}

func FindGradeByNumber(gradeNumber string) (models.Grade, error) {
	var grade models.Grade
	if err := config.DB.Where("number = ?", gradeNumber).First(&grade).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return grade, nil
		}

		return grade, err
	}
	return grade, nil
}
