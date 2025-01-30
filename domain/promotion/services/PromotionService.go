package services

import (
	"errors"
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	assessmentGradeResponses "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/assessmentGrade/api/responses"
	repositories6 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/repositories"
	services5 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	repositories5 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/repositories"
	services3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/services"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	services4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/services"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/importCurriculum/api/responses"
	models3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/importCurriculum/models"
	repositories4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/promotion/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/promotion/api/responses"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/promotion/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/promotion/repositories"
	repositories7 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/repositories"
	repositories3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/repositories"
	services2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/services"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type PromotionService interface {
	AllPromotionSettings(studyPeriodNumber, nextCampusClassNumber, settingsType string) ([]responses.PromotionResponse, error)
	NewPromotionSetting(promotionSettingsRequest requests.PromotionSettingRequest, user *models.User) error
	SavePromotionSettingValue(payload requests.PromotionSettingValueRequest, user *models.User) error
	PromotionSettingDetails(promotionSettingNumber string) (responses.PromotionDetailedResponse, error)
	PromotionSettingGrades(promotionSettingNumber string, c *gin.Context) (responses.PromotionSettingGradeResponse, error)
}

type promotionService struct{}

func (p *promotionService) PromotionSettingGrades(promotionSettingNumber string, c *gin.Context) (responses.PromotionSettingGradeResponse, error) {

	response := responses.PromotionSettingGradeResponse{}

	promotionSetting, err := repositories.FindPromotionSettingByNumber(promotionSettingNumber)
	if err != nil {
		return response, errors.New("Error checking promotion setting " + err.Error())
	}

	if promotionSetting.Number == "" {
		return response, errors.New("Invalid promotion setting number ")
	}

	class, err := repositories2.GetClassByCampusClassNumber(promotionSetting.PreviousCampusClassNumber)
	if err != nil {
		return response, err
	}

	classAssessmentGrade, err := repositories6.GetAssessmentGradeByClassNumberAndStatus(class.Number, constants.StatusActive)
	if err != nil {
		return response, err
	}

	if classAssessmentGrade.Number == "" {
		return response, errors.New(class.Name + " does not have assessment grade set")
	}

	assessmentGrades, err := repositories6.FindAssessmentGradeScaleByGradeNumber(classAssessmentGrade.Number)
	if err != nil {
		return response, errors.New("getting assessments grade scale error: " + err.Error())
	}

	curriculum, _ := repositories4.FindCurriculumByCampusClassNumber(promotionSetting.PreviousCampusClassNumber)

	assessmentAreas, err := services4.FindAllAssessmentAreas("", class.Number, curriculum.Number, "", "", c)
	if err != nil {
		return response, errors.New("getting assessments area for class error: " + err.Error())
	}

	var assessmentGradeItemResponses = make([]assessmentGradeResponses.AssessmentGradeItemResponse, 0)
	for _, assessmentGradeItem := range assessmentGrades {
		assessmentGradeItemResponses = append(assessmentGradeItemResponses, services5.CreateAssessmentItemResponse(assessmentGradeItem))
	}

	response.GradeScales = assessmentGradeItemResponses
	response.AssessmentAreas = assessmentAreas

	return response, nil
}

func (p *promotionService) PromotionSettingDetails(promotionSettingNumber string) (responses.PromotionDetailedResponse, error) {
	promotionSetting, err := repositories.FindPromotionSettingByNumber(promotionSettingNumber)
	if err != nil {
		return responses.PromotionDetailedResponse{}, errors.New(fmt.Sprintf("Error '%s' finding promotion setting by number", err.Error()))
	}

	if promotionSetting.Number == "" {
		return responses.PromotionDetailedResponse{}, errors.New("Invalid promotion setting number ")
	}

	detailedResponse := responses.PromotionDetailedResponse{
		PromotionResponse: createPromotionResponse(promotionSetting),
	}

	detailedResponse.NextCampusClassNumber = promotionSetting.NextCampusClassNumber
	detailedResponse.PreviousCampusClassNumber = promotionSetting.PreviousCampusClassNumber

	return detailedResponse, nil
}

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
		promotionResponses = append(promotionResponses, createPromotionResponse(promotionSetting))
	}

	return promotionResponses, nil
}

func createPromotionResponse(promotionSetting models2.PromotionSetting) responses.PromotionResponse {

	nextClass, _ := repositories2.GetClassByCampusClassNumber(promotionSetting.NextCampusClassNumber)
	prevClass, _ := repositories2.GetClassByCampusClassNumber(promotionSetting.PreviousCampusClassNumber)
	studyPeriod, _ := repositories3.GetStudyPeriodByNumber(promotionSetting.StudyPeriodNumber)

	promo := responses.PromotionResponse{
		Name:                    promotionSetting.Name,
		PromotionSettingsNumber: promotionSetting.Number,
		NextClass:               services.CreateClassResponse(nextClass),
		PreviousClass:           services.CreateClassResponse(prevClass),
		Type:                    promotionSetting.Type,
		Status:                  promotionSetting.Status,
		StudyPeriod:             services2.CreateStudyPeriodResponse(studyPeriod),
	}

	campusClass, err := repositories2.GetCampusClassByNumber(promotionSetting.PreviousCampusClassNumber)
	if err != nil {
		fmt.Println(fmt.Sprintf("Campus class error: %s", err.Error()))
		return promo
	}

	campus, err := repositories5.GetCampusByNumber(campusClass.CampusNumber)
	if err == nil {
		promo.Campus = services3.CreateCampusResponse(campus)
	} else {
		fmt.Println(fmt.Sprintf("Campus error: %s", err.Error()))
		return promo
	}

	curriculumClass, err := repositories2.GetCurriculumClass(campusClass.CurriculumClassNumber)
	if err != nil {
		fmt.Println(fmt.Sprintf("Curriculum classs error: %s", err.Error()))
		return promo
	}

	curriculumLevel, err := repositories7.GetCurriculumLevelByNumber(curriculumClass.CurriculumLevelNumber)
	if err != nil {
		fmt.Println(fmt.Sprintf("Curriculum level error: %s", err.Error()))
		return promo
	}

	curriculum, err := models3.GetCurriculumByNumber(curriculumLevel.CurriculumNumber)
	if err == nil {
		promo.Curriculum = responses2.CurriculumResponse{
			CurriculumNumber: curriculum.Number,
			Name:             curriculum.Name,
		}
	} else {
		fmt.Println(fmt.Sprintf("Curriculum error: %s", err.Error()))
	}

	return promo
}

func NewPromotionService() PromotionService {
	return &promotionService{}
}
