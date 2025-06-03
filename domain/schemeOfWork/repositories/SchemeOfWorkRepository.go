package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/schemeOfWork/models"
	"gorm.io/gorm"
)

func GetSchemeOfWorkSubTopicByNumber(schemeOfWorkSubTopicNumber string, tx ...*gorm.DB) (models.SchemeOfWorkSubTopic, error) {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var schemeOfWorkSubTopic models.SchemeOfWorkSubTopic

	if err := db.Where("number = ?", schemeOfWorkSubTopicNumber).First(&schemeOfWorkSubTopic).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.SchemeOfWorkSubTopic{}, nil
		}

		return models.SchemeOfWorkSubTopic{}, err
	}

	return schemeOfWorkSubTopic, nil
}

func SaveSchemeOfWork(schemeOfWork models.SchemeOfWork, tx ...*gorm.DB) error {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	if err := db.Save(&schemeOfWork).Error; err != nil {
		return err
	}

	return nil
}

func SaveSchemeOfWorkSubTopic(schemeOfWorkSubTopic models.SchemeOfWorkSubTopic, tx ...*gorm.DB) error {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	if err := db.Save(&schemeOfWorkSubTopic).Error; err != nil {
		return err
	}

	return nil
}

func FindSchemeOfWorkSubTopicByNumber(schemeOfWorkSubTopicNumber string, tx ...*gorm.DB) (models.SchemeOfWorkSubTopic, error) {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var schemeOfWork models.SchemeOfWorkSubTopic

	if err := db.Where("number = ?", schemeOfWorkSubTopicNumber).First(&schemeOfWork).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.SchemeOfWorkSubTopic{}, nil
		}

		return models.SchemeOfWorkSubTopic{}, err
	}

	return schemeOfWork, nil
}

func GetSchemeOfWorkByNumber(schemeOfWorkNumber string, tx ...*gorm.DB) (models.SchemeOfWork, error) {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var schemeOfWork models.SchemeOfWork

	if err := db.Where("number = ?", schemeOfWorkNumber).First(&schemeOfWork).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.SchemeOfWork{}, nil
		}

		return models.SchemeOfWork{}, err
	}

	return schemeOfWork, nil
}

func GetSchemeOfWorkBySchemeOfWorkSubTopicNumber(schemeOfWorkSubTopicNumber string, tx ...*gorm.DB) (models.SchemeOfWork, error) {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var schemeOfWork models.SchemeOfWork

	if err := db.
		Joins("JOIN scheme_of_work_sub_topics ON scheme_of_work_sub_topics.scheme_of_work_number = scheme_of_works.number").
		Where("scheme_of_work_sub_topics.number = ?", schemeOfWorkSubTopicNumber).First(&schemeOfWork).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.SchemeOfWork{}, nil
		}

		return models.SchemeOfWork{}, err
	}

	return schemeOfWork, nil
}

func FindSchemeOfWorkByNameAndClassAndSubjectAndStatus(name, classNumber, subjectNumber, status string, tx ...*gorm.DB) (models.SchemeOfWork, error) {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var schemeOfWork models.SchemeOfWork

	if err := db.Where("name = ? AND subject_number = ? AND campus_class_number = ? AND status = ?", name, subjectNumber, classNumber, status).First(&schemeOfWork).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.SchemeOfWork{}, nil
		}

		return models.SchemeOfWork{}, err
	}

	return schemeOfWork, nil
}
