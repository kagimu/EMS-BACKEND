package repositories

import (
	"errors"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/sourceOfSchoolInfo/models"
	"gorm.io/gorm"
)

func SaveSourceOfSchoolInfo(sourceOfSchoolInfo models2.SourceOfSchoolInfo, tx ...*gorm.DB) error {

	db := config.DB

	if len(tx) > 0 {
		db = tx[0]
	}

	if err := db.Save(&sourceOfSchoolInfo).Error; err != nil {
		return err
	}

	return nil

}

func FindAllSourceOfSchoolInfo(entityNumber string, tx ...*gorm.DB) ([]models2.SourceOfSchoolInfo, error) {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var sourceOfSchoolInfos []models2.SourceOfSchoolInfo

	if err := db.
		Scopes(filters.ByAnyField("entity_number", entityNumber)).
		Find(&sourceOfSchoolInfos).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sourceOfSchoolInfos, nil
		}

		return sourceOfSchoolInfos, err
	}

	return sourceOfSchoolInfos, nil
}

func FindSourceOfSchoolInfoByName(entityNumber, name string, tx ...*gorm.DB) (models2.SourceOfSchoolInfo, error) {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var sourceOfSchoolInfo models2.SourceOfSchoolInfo

	if err := db.
		Scopes(filters.ByAnyField("entity_number", entityNumber), filters.ByAnyField("name", name)).
		Find(&sourceOfSchoolInfo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sourceOfSchoolInfo, nil
		}

		return sourceOfSchoolInfo, err
	}

	return sourceOfSchoolInfo, nil
}

func FindSourceOfSchoolInfoByCode(entityNumber string, code int64, tx ...*gorm.DB) (models2.SourceOfSchoolInfo, error) {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var sourceOfSchoolInfo models2.SourceOfSchoolInfo

	if err := db.
		Scopes(filters.ByAnyField("entity_number", entityNumber)).
		Where("code = ?", code).
		Find(&sourceOfSchoolInfo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sourceOfSchoolInfo, nil
		}

		return sourceOfSchoolInfo, err
	}

	return sourceOfSchoolInfo, nil
}

func FindSourceOfSchoolInfoByNumber(number string, tx ...*gorm.DB) (models2.SourceOfSchoolInfo, error) {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var sourceOfSchoolInfo models2.SourceOfSchoolInfo

	if err := db.
		Scopes(filters.ByAnyField("number", number)).
		First(&sourceOfSchoolInfo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sourceOfSchoolInfo, nil
		}

		return sourceOfSchoolInfo, err
	}

	return sourceOfSchoolInfo, nil
}

func RemoveSourceOfSchoolInfoByNumber(sourceOfSchoolInfo models2.SourceOfSchoolInfo, tx ...*gorm.DB) error {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	return db.Delete(&sourceOfSchoolInfo).Error
}
