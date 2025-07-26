package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	LevelModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/models"
	"gorm.io/gorm"
)

type LevelRepository interface {
	FindLevelByNumber(levelNumber string) (LevelModels.Level, error)
	FindLevelByCode(code int64, entityNumber string) (LevelModels.Level, error)
	FindLevelByName(name, entityNumber string) (LevelModels.Level, error)
	FindLevelWithClassesBeforeThisLevelCode(entityNumber string, currentLevelCode int64) (LevelModels.Level, error)
	FindLevelWithClassesAfterThisLevelCode(entityNumber string, currentLevelCode int64) (LevelModels.Level, error)
	SaveLevel(level LevelModels.Level) (LevelModels.Level, error)
}

type levelRepositoryImpl struct {
}

func (l levelRepositoryImpl) FindLevelWithClassesAfterThisLevelCode(entityNumber string, currentLevelCode int64) (LevelModels.Level, error) {
	db := config.DB

	var class LevelModels.Level

	err := db.
		Joins("JOIN classes ON classes.level_number = classes.number").
		Where("classes.entity_number = ? AND levels.code > ?", entityNumber, currentLevelCode).
		Order("levels.code ASC").
		First(&class).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return class, nil
	}
	return class, err
}

func (l levelRepositoryImpl) FindLevelWithClassesBeforeThisLevelCode(entityNumber string, currentLevelCode int64) (LevelModels.Level, error) {
	db := config.DB

	var class LevelModels.Level

	err := db.
		Joins("JOIN classes ON classes.level_number = classes.number").
		Where("classes.entity_number = ? AND levels.code < ?", entityNumber, currentLevelCode).
		Order("levels.code DESC").
		First(&class).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return class, nil
	}
	return class, err
}

func (l levelRepositoryImpl) SaveLevel(level LevelModels.Level) (LevelModels.Level, error) {

	if err := config.DB.Save(&level).Error; err != nil {
		return level, err
	}

	return level, nil
}

func (l levelRepositoryImpl) FindLevelByName(name, entityNumber string) (LevelModels.Level, error) {
	db := config.DB
	var level LevelModels.Level

	if err := db.Where("name = ? AND entity_number = ?", name, entityNumber).First(&level).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return level, nil
		}
		return level, err
	}
	return level, nil
}

func (l levelRepositoryImpl) FindLevelByCode(code int64, entityNumber string) (LevelModels.Level, error) {
	db := config.DB
	var level LevelModels.Level

	if err := db.Where("code = ? AND entity_number = ?", code, entityNumber).First(&level).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return level, nil
		}
		return level, err
	}
	return level, nil
}

func (l levelRepositoryImpl) FindLevelByNumber(levelNumber string) (LevelModels.Level, error) {
	db := config.DB
	var level LevelModels.Level

	if err := db.Where("number = ?", levelNumber).First(&level).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return level, nil
		}
		return level, err
	}
	return level, nil
}

func NewLevelRepository() LevelRepository {
	return &levelRepositoryImpl{}
}
