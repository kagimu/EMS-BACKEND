package services

import (
	"errors"
	"fmt"
	activityReportingRequests "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/activityreporting/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/activityreporting/models"
	activityReportingRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/activityreporting/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
)

func SaveIndicator(indicatorRequest activityReportingRequests.IndicatorRequest, currentUserNumber string) error {

	indicator, err := activityReportingRepositories.FindIndicatorByName(indicatorRequest.Name, indicatorRequest.EntityNumber)
	if err != nil {
		return err
	}

	if indicator.Number != "" {
		return errors.New("the indicator name already exists")
	}

	indicator, err = activityReportingRepositories.FindIndicatorByCode(indicatorRequest.Code, indicatorRequest.EntityNumber)
	if err != nil {
		return err
	}

	if indicator.Number != "" {
		return errors.New("the indicator code already exists")
	}

	err = activityReportingRepositories.SaveIndicator(models.Indicator{
		Name:         indicatorRequest.Name,
		Code:         indicatorRequest.Code,
		EntityNumber: indicatorRequest.EntityNumber,
		Description:  indicatorRequest.Description,
		UserNumber:   currentUserNumber,
		Status:       constants.StatusActive,
		InputType:    indicatorRequest.InputType,
	})

	if err != nil {
		fmt.Println(err)
		return errors.New("failed to save the indicator, try again")
	}

	return nil
}
