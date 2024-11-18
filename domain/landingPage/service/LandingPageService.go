package service

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/attachments/services"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/landingPage/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/landingPage/repository"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/landingPage/responses"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

type LandingPageService struct {
}

func NewLandingPageService() LandingPageService {
	return LandingPageService{}
}

func (l LandingPageService) GetLandingPage() (responses.LandingPageResponse, error) {
	var landingPage models.LandingPage
	db := config.DB
	var repo = repository.NewLandingPageRepo(db)
	landingPage, err := repo.GetLandingPage()

	if err != nil {
		return responses.LandingPageResponse{}, err
	}

	return responses.LandingPageResponse{
		Logo:  landingPage.Image,
		Title: landingPage.Title,
	}, nil

}

func (l LandingPageService) LandingPageDetails(c *gin.Context, pageTitle string, user *authModel.User) error {
	logoName, err := services.UploadAttachment(c, "logo")
	if err != nil {
		return errors.New("Failed to upload logo with error: " + err.Error())
	}

	page := models.LandingPage{
		Title:      pageTitle,
		Image:      logoName,
		UserNumber: user.Number,
		Number:     uuid.NewV4().String(),
	}

	db := config.DB
	var repo = repository.NewLandingPageRepo(db)
	err = repo.CreateLandingPage(&page)
	if err != nil {
		return errors.New("Failed to create LandingPage with error: " + err.Error())
	}

	return nil
}
