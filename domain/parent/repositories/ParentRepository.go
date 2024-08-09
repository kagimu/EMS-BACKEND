package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/parent/models"
	"gorm.io/gorm"
)

type ParentRepository interface {
	Create(parent models.Parent) error
	GetParentByApplicationNumber(applicationNumber string) (models.Parent, error)
	FindAll() ([]models.Parent, error)
	Delete(ParentId string) error
}

type ParentRepositoryImpl struct {
	DB *gorm.DB
}

func (p ParentRepositoryImpl) Create(parent models.Parent) error {
	//TODO implement me
	panic("implement me")
}

func (p *ParentRepositoryImpl) GetParentByApplicationNumber(applicationNumber string) (models.Parent, error) {
	var parent models.Parent
	if err := p.DB.Where("application_number = ?", applicationNumber).First(&parent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Parent{}, nil
		}
		return models.Parent{}, err
	}
	return parent, nil
}

func (p ParentRepositoryImpl) FindAll() ([]models.Parent, error) {
	//TODO implement me
	panic("implement me")
}

func (p ParentRepositoryImpl) Delete(ParentId string) error {
	return p.DB.Delete(ParentId).Error
}

func NewParentRepositoryImpl(db *gorm.DB) ParentRepository {
	return &ParentRepositoryImpl{DB: db}
}

func GetParentByApplicationNumber(applicationNumber string, tx ...*gorm.DB) (models.Parent, error) {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var parent models.Parent

	if err := db.Where("application_number = ?", applicationNumber).First(&parent).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Parent{}, nil
		}

		return models.Parent{}, err
	}

	return parent, nil
}
