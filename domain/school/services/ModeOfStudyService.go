package services

import (
	"errors"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/repositories"
	uuid "github.com/satori/go.uuid"
)

type ModeOfStudyService interface {
	NewModeOfStudy(modeOfStudyRequest requests.ModeOfStudyRequest, currentUserNumber string) error
	AllModesOfStudy(schoolNumber string) ([]responses.ModeOfStudyResponse, error)
	UpdateModeOfStudy(updateModeOfStudyRequest requests.UpdateModeOfStudyRequest, currentUserNumber string) error
	DeleteModeOfStudy(modeOfStudyNumber string) error
}

type modeOfStudyService struct {
}

func (modeOfStudyService *modeOfStudyService) DeleteModeOfStudy(modeOfStudyNumber string) error {
	modeOfStudy, err := repositories.GetModeOfStudyByNumber(modeOfStudyNumber)

	if err != nil {
		return errors.New("Error checking mode of study " + err.Error())
	}

	if modeOfStudy.Number == "" {
		return errors.New("Invalid mode of study number provided ")
	}

	err = repositories.RemoveModeOfStudy(modeOfStudy)
	if err != nil {
		modeOfStudy.Status = constants.StatusInactive
		err = repositories.SaveModeOfStudy(modeOfStudy)

		if err != nil {
			return errors.New("Deleting mode of study failed, try again later ")
		}

	}

	return nil
}

func (modeOfStudyService *modeOfStudyService) UpdateModeOfStudy(updateModeOfStudyRequest requests.UpdateModeOfStudyRequest, currentUserNumber string) error {
	modeOfStudy, err := repositories.GetModeOfStudyByNumber(updateModeOfStudyRequest.ModeOfStudyNumber)

	if err != nil {
		return errors.New("Error checking mode of study " + err.Error())
	}

	if modeOfStudy.Number == "" {
		return errors.New("Invalid mode of study number provided ")
	}

	checkName, err := repositories.GetModeOfStudyByName(updateModeOfStudyRequest.Name, updateModeOfStudyRequest.SchoolNumber)
	if err != nil {
		return errors.New("Checking name error " + err.Error())
	}

	if checkName.Number != "" && modeOfStudy.Number != checkName.Number {
		return errors.New("This name is already used ")
	}

	modeOfStudy.Name = updateModeOfStudyRequest.Name
	modeOfStudy.SchoolNumber = updateModeOfStudyRequest.SchoolNumber
	modeOfStudy.UserNumber = currentUserNumber

	if updateModeOfStudyRequest.Status != "" {
		modeOfStudy.Status = updateModeOfStudyRequest.Status
	}

	err = repositories.SaveModeOfStudy(modeOfStudy)
	if err != nil {
		return errors.New("Error updating mode of study " + err.Error())
	}

	return nil
}

func (modeOfStudyService *modeOfStudyService) AllModesOfStudy(schoolNumber string) ([]responses.ModeOfStudyResponse, error) {
	modes, err := repositories.FindModeOfStudy(schoolNumber)

	var modeOfStudyResponse = make([]responses.ModeOfStudyResponse, 0)

	if err != nil {
		return nil, errors.New("Fetching modes of study error: " + err.Error())
	}

	for _, modeOfStudy := range modes {
		modeOfStudyResponse = append(modeOfStudyResponse, responses.ModeOfStudyResponse{
			ModeOfStudyNumber: modeOfStudy.Number,
			Name:              modeOfStudy.Name,
			Status:            modeOfStudy.Status,
			CreatedAt:         modeOfStudy.CreatedAt.Format(time.DateTime),
			SchoolNumber:      modeOfStudy.SchoolNumber,
			Code:              modeOfStudy.Code,
		})
	}

	return modeOfStudyResponse, nil
}

func (modeOfStudyService *modeOfStudyService) NewModeOfStudy(modeOfStudyRequest requests.ModeOfStudyRequest, currentUserNumber string) error {
	modeOfStudy, err := repositories.GetModeOfStudyByName(modeOfStudyRequest.Name, modeOfStudyRequest.SchoolNumber)
	if err != nil {
		return errors.New("Error checking duplicate name: " + err.Error())
	}

	if modeOfStudy.Number != "" {
		return errors.New("Name already added ")
	}

	err = repositories.SaveModeOfStudy(models.ModeOfStudy{
		Number:       uuid.NewV4().String(),
		Name:         modeOfStudyRequest.Name,
		SchoolNumber: modeOfStudyRequest.SchoolNumber,
		UserNumber:   currentUserNumber,
		Code:         modeOfStudyRequest.Code,
	})

	if err != nil {
		return errors.New("Saving mode of study error: " + err.Error())
	}

	return nil
}

func NewModeOfStudyService() ModeOfStudyService {
	return &modeOfStudyService{}
}
