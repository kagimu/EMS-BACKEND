package services

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/api/requests"
	LevelModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/repositories"
	uuid "github.com/satori/go.uuid"
	"time"
)

type LevelService struct {
}

var levelRepository = repositories.NewLevelRepository()

func (ls LevelService) EditLevel(userNumber string, levelInfo requests.LevelInfo) (LevelModels.Level, error) {

	level, err := levelRepository.FindLevelByNumber(levelInfo.LevelNumber)
	if err != nil {
		return level, err
	}

	if level.Number == "" {
		return level, errors.New("Invalid level number ")
	}

	checkDup, err := levelRepository.FindLevelByCode(levelInfo.Code, levelInfo.EntityNumber)
	if err != nil {
		return checkDup, err
	}

	if checkDup.Number != "" && checkDup.Number != level.Number {
		return checkDup, errors.New("Provided code already used ")
	}

	checkDup, err = levelRepository.FindLevelByName(levelInfo.Name, levelInfo.EntityNumber)
	if err != nil {
		return checkDup, err
	}

	if checkDup.Number != "" && checkDup.Number != level.Number {
		return checkDup, errors.New("Provided name already used ")
	}

	level.UserNumber = userNumber
	level.Name = levelInfo.Name
	level.Code = levelInfo.Code
	level.EntityNumber = levelInfo.EntityNumber
	level.UpdatedAt = time.Now()

	saved, err := levelRepository.SaveLevel(level)
	if err != nil {
		return saved, err
	}

	return saved, nil
}

func (ls LevelService) SaveLevel(userNumber string, levelInfo requests.LevelInfo) (LevelModels.Level, error) {

	checkDup, err := levelRepository.FindLevelByCode(levelInfo.Code, levelInfo.EntityNumber)
	if err != nil {
		return checkDup, err
	}

	if checkDup.Number != "" {
		return checkDup, errors.New("Provided code already used ")
	}

	checkDup, err = levelRepository.FindLevelByName(levelInfo.Name, levelInfo.EntityNumber)
	if err != nil {
		return checkDup, err
	}

	if checkDup.Number != "" {
		return checkDup, errors.New("Provided name already used ")
	}

	level := LevelModels.Level{
		Number:       uuid.NewV4().String(),
		UserNumber:   userNumber,
		Name:         levelInfo.Name,
		Code:         levelInfo.Code,
		EntityNumber: levelInfo.EntityNumber,
	}

	saved, err := levelRepository.SaveLevel(level)
	if err != nil {
		return saved, err
	}

	return saved, nil
}

func NewLevelService() LevelService {
	return LevelService{}
}
