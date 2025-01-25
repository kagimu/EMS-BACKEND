package services

import (
	"errors"
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/promotion/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/promotion/api/responses"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/promotion/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/promotion/repositories"
	repositories3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/repositories"
	services2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/services"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type PromotionService interface {
	AllPromotionSettings(studyPeriodNumber, nextCampusClassNumber, settingsType string) ([]responses.PromotionResponse, error)
	NewPromotionSetting(promotionSettingsRequest requests.PromotionSettingRequest, user *models.User) error
	SavePromotionSettingValue(payload requests.PromotionSettingValueRequest, user *models.User) error
}

type promotionService struct{}

func (p *promotionService) SavePromotionSettingValue(payload requests.PromotionSettingValueRequest, user *models.User) error {
	savePromotionSettings := models2.PromotionSettingValue{
		Number:                     uuid.NewV4().String(),
		Score:                      payload.Score,
		PromotionSettingNumber:     payload.PromotionSettingNumber,
		AssessmentGradeScaleNumber: payload.AssessmentGradeScaleNumber,
		MinimumPoints:              payload.MinimumPoints,
		UserNumber:                 user.Number,
	}

	err := config.DB.Transaction(func(tx *gorm.DB) error {

		checkOld, err := repositories.FindPromotionSettingValueBySettingNumber(payload.PromotionSettingNumber, tx)
		if err != nil {
			return errors.New("Checking old value error: " + err.Error())
		}

		if checkOld.Number != "" {
			savePromotionSettings.Number = checkOld.Number
		}

		err = repositories.SavePromotionSettingValue(savePromotionSettings, tx)
		if err != nil {
			return errors.New(fmt.Sprintf("Error saving promotion setting value: %s", err.Error()))
		}

		assessments := payload.AssessmentAreas
		if len(assessments) > 0 {
			err = handleAssessmentAreas(payload.PromotionSettingNumber, assessments, user, tx)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func handleAssessmentAreas(promotionSettingsNumber string, assessments []requests.AssessmentAreaRequestValues, user *models.User, tr *gorm.DB) error {

	for _, assessment := range assessments {
		checkDuplicates, err := repositories.FindPromotionAssessmentAreaValueBySettingsNumberAndAssessmentAreaNumber(promotionSettingsNumber, assessment.AssessmentAreaNumber, tr)
		if err != nil {
			return errors.New("Checking promotion assessment area duplicates error: " + err.Error())
		}

		newValue := models2.PromotionAssessmentAreaValue{
			AssessmentAreaNumber:   assessment.AssessmentAreaNumber,
			PromotionSettingNumber: promotionSettingsNumber,
			GradeNumber:            assessment.GradeNumber,
			Number:                 uuid.NewV4().String(),
			UserNumber:             user.Number,
		}

		if checkDuplicates.Number != "" {
			newValue.Number = checkDuplicates.Number
		}

		err = repositories.SavePromotionAssessmentAreaValue(newValue, tr)
		if err != nil {
			return errors.New(fmt.Sprintf("Error saving promotion assessment area value: %s", err.Error()))
		}

	}

	return nil
}

func (p *promotionService) NewPromotionSetting(promotionSettingsRequest requests.PromotionSettingRequest, user *models.User) error {

	nextCampusClass, err := services.ResolveCampusClass(promotionSettingsRequest.NextClassNumber, promotionSettingsRequest.CampusNumber, promotionSettingsRequest.CurriculumNumber, promotionSettingsRequest.CampusCurriculumNumber)
	if err != nil {
		return err
	}

	prevCampusClass, err := services.ResolveCampusClass(promotionSettingsRequest.PreviousClassNumber, promotionSettingsRequest.CampusNumber, promotionSettingsRequest.CurriculumNumber, promotionSettingsRequest.CampusCurriculumNumber)
	if err != nil {
		return err
	}

	studyPeriod, err := repositories3.GetStudyPeriodByNumber(promotionSettingsRequest.StudyPeriodNumber)
	if err != nil {
		return err
	}

	if studyPeriod.Number == "" {
		return errors.New("Invalid study period number ")
	}

	savePromotionSettings := models2.PromotionSetting{
		Number:                    uuid.NewV4().String(),
		Name:                      promotionSettingsRequest.Name,
		Type:                      promotionSettingsRequest.Type,
		Status:                    promotionSettingsRequest.Status,
		NextCampusClassNumber:     nextCampusClass.Number,
		PreviousCampusClassNumber: prevCampusClass.Number,
		UserNumber:                user.Number,
		StudyPeriodNumber:         studyPeriod.Number,
	}

	err = repositories.SavePromotionSettings(savePromotionSettings)
	if err != nil {
		return err
	}

	return nil
}

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
