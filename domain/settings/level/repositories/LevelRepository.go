package repositories

import LevelModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/models"

type LevelRepository interface {
	FindLevelByNumber(levelNumber string) (LevelModels.Level, error)
	FindLevelByCode(code int64, entityNumber string) (LevelModels.Level, error)
	FindLevelByName(name, entityNumber string) (LevelModels.Level, error)
	SaveLevel(level LevelModels.Level) (LevelModels.Level, error)
}

type levelRepositoryImpl struct {
}

func (l levelRepositoryImpl) SaveLevel(level LevelModels.Level) (LevelModels.Level, error) {
	//TODO implement me
	panic("implement me")
}

func (l levelRepositoryImpl) FindLevelByName(name, entityNumber string) (LevelModels.Level, error) {
	//TODO implement me
	panic("implement me")
}

func (l levelRepositoryImpl) FindLevelByCode(code int64, entityNumber string) (LevelModels.Level, error) {
	//TODO implement me
	panic("implement me")
}

func (l levelRepositoryImpl) FindLevelByNumber(levelNumber string) (LevelModels.Level, error) {
	//TODO implement me
	panic("implement me")
}

func NewLevelRepository() LevelRepository {
	return &levelRepositoryImpl{}
}
