package service

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/landingPage/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/landingPage/repository"
)

type LandingPageService struct {
}

func NewLandingPageService() LandingPageService {
	return LandingPageService{}
}

func (l LandingPageService) GetLandingPage() (models.LandingPage, error) {
	var landingPage models.LandingPage
	db := config.DB
	var repo = repository.NewLandingPageRepo(db)
	landingPage, err := repo.GetLandingPage()

	if err != nil {
		return landingPage, err
	}
	return landingPage, nil

}

func (l LandingPageService) CreateLandingPage(page *models.LandingPage) error {
	db := config.DB
	var repo = repository.NewLandingPageRepo(db)
	err := repo.CreateLandingPage(page)
	if err != nil {
		return err
	}
	return nil
}
