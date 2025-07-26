package services

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/sourceOfSchoolInfo/controllers/api"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/sourceOfSchoolInfo/controllers/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/sourceOfSchoolInfo/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/sourceOfSchoolInfo/repositories"
	uuid "github.com/satori/go.uuid"
	"time"
)

type SourceOfSchoolService struct {
}

func (s SourceOfSchoolService) NewSourceOfSchool(sourceOfSchoolInfoRequest api.SourceOfSchoolRequest, userNumber, entityNumber string) (responses.SourceOfSchoolResponse, error) {

	checkDuplicate, err := repositories.FindSourceOfSchoolInfoByName(entityNumber, sourceOfSchoolInfoRequest.Name)
	if err != nil {
		return responses.SourceOfSchoolResponse{}, err
	}

	if checkDuplicate.Number != "" {
		return responses.SourceOfSchoolResponse{}, errors.New("Source name already used ")
	}

	data := models.SourceOfSchoolInfo{
		Name:         sourceOfSchoolInfoRequest.Name,
		Number:       uuid.NewV4().String(),
		EntityNumber: entityNumber,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		UserNumber:   userNumber,
	}

	err = repositories.SaveSourceOfSchoolInfo(data)

	if err != nil {
		return responses.SourceOfSchoolResponse{}, err
	}

	return responses.SourceOfSchoolResponse{
		Name:         sourceOfSchoolInfoRequest.Name,
		Number:       data.Number,
		EntityNumber: data.EntityNumber,
	}, nil
}

func (s SourceOfSchoolService) FindAllSources(entityNumber string) ([]responses.SourceOfSchoolResponse, error) {

	var result = make([]responses.SourceOfSchoolResponse, 0)

	data, err := repositories.FindAllSourceOfSchoolInfo(entityNumber)
	if err != nil {
		return []responses.SourceOfSchoolResponse{}, err
	}

	for _, v := range data {
		result = append(result, responses.SourceOfSchoolResponse{
			Name:         v.Name,
			Number:       v.Number,
			UserNumber:   v.UserNumber,
			EntityNumber: v.EntityNumber,
		})
	}

	return result, nil
}

func (s SourceOfSchoolService) Update(sourceOfSchoolInfo api.SourceOfSchoolRequest, sourceOfSchoolInfoNumber, userNumber string) (responses.SourceOfSchoolResponse, error) {
	sourceOfSchool, err := repositories.FindSourceOfSchoolInfoByNumber(sourceOfSchoolInfoNumber)
	if err != nil {
		return responses.SourceOfSchoolResponse{}, err
	}

	if sourceOfSchool.Number == "" {
		return responses.SourceOfSchoolResponse{}, errors.New("Invalid number to edit ")
	}

	checkByName, err := repositories.FindSourceOfSchoolInfoByName(sourceOfSchool.EntityNumber, sourceOfSchoolInfo.Name)
	if err != nil {
		return responses.SourceOfSchoolResponse{}, err
	}

	if checkByName.Number != "" && checkByName.Number != sourceOfSchool.Number {
		return responses.SourceOfSchoolResponse{}, errors.New("Source name already used ")
	}

	sourceOfSchool.Name = sourceOfSchoolInfo.Name
	sourceOfSchool.UserNumber = userNumber
	sourceOfSchool.UpdatedAt = time.Now()

	err = repositories.SaveSourceOfSchoolInfo(sourceOfSchool)

	if err != nil {
		return responses.SourceOfSchoolResponse{}, err
	}

	return responses.SourceOfSchoolResponse{
		Name:         sourceOfSchool.Name,
		Number:       sourceOfSchool.Number,
		EntityNumber: sourceOfSchool.EntityNumber,
	}, nil
}

func (s SourceOfSchoolService) DeleteSource(entityNumber, sourceOfSchoolInfoNumber string) error {

	data, err := repositories.FindSourceOfSchoolInfoByNumber(sourceOfSchoolInfoNumber)
	if err != nil {
		return err
	}

	if data.Number == "" {
		return errors.New("Source name does not exist ")
	}

	if data.EntityNumber != entityNumber {
		return errors.New("You cannot delete this source ")
	}

	err = repositories.RemoveSourceOfSchoolInfoByNumber(data)
	if err != nil {
		return err
	}

	return nil
}

func NewSourceOfSchoolService() *SourceOfSchoolService {
	return &SourceOfSchoolService{}
}
