package services

import (
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/promotion/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/promotion/repositories"
	repositories3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/repositories"
	services2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/services"
)

type PromotionService interface {
	AllPromotionSettings(studyPeriodNumber, nextCampusClassNumber, settingsType string) ([]responses.PromotionResponse, error)
}

type promotionService struct{}

func (p *promotionService) AllPromotionSettings(studyPeriodNumber, nextCampusClassNumber, settingsType string) ([]responses.PromotionResponse, error) {
	var promotionResponses = make([]responses.PromotionResponse, 0)

	promotionSettings, err := repositories.FindAllPromotionSettings(studyPeriodNumber, nextCampusClassNumber, settingsType)
	if err != nil {
		return promotionResponses, err
	}

	for _, promotionSetting := range promotionSettings {
		nextClass, _ := repositories2.GetClassByCampusClassNumber(promotionSetting.NextCampusClassNumber)
		prevClass, _ := repositories2.GetClassByCampusClassNumber(promotionSetting.PreviousCampusClassNumber)
		studyPeriod, _ := repositories3.GetStudyPeriodByNumber(promotionSetting.StudyPeriodNumber)

		promotionResponses = append(promotionResponses, responses.PromotionResponse{
			Name:                    promotionSetting.Name,
			PromotionSettingsNumber: promotionSetting.Number,
			NextClass:               services.CreateClassResponse(nextClass),
			PreviousClass:           services.CreateClassResponse(prevClass),
			Type:                    promotionSetting.Type,
			Status:                  promotionSetting.Status,
			StudyPeriod:             services2.CreateStudyPeriodResponse(studyPeriod),
		})
	}

	return promotionResponses, nil
}

func NewPromotionService() PromotionService {
	return &promotionService{}
}
