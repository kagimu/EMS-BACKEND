package repository

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/landingPage/models"
	"gorm.io/gorm"
)

type LandingPageRepo interface {
	GetLandingPage() (models.LandingPage, error)
	CreateLandingPage(page *models.LandingPage) error
}

type LandingPageRepoImpl struct {
	DB *gorm.DB
}

func (l LandingPageRepoImpl) GetLandingPage() (models.LandingPage, error) {
	var page models.LandingPage
	err := l.DB.Order("created_at desc").First(&page).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return page, nil
		}
		return page, err
	}
	return page, nil
}

func (l LandingPageRepoImpl) CreateLandingPage(page *models.LandingPage) error {
	return l.DB.Create(page).Error
}

func NewLandingPageRepo(db *gorm.DB) LandingPageRepo {
	return &LandingPageRepoImpl{DB: db}
}
