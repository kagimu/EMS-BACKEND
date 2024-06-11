package services

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/repositories"
)

func SaveCurrentEntity(entityNumber string, currentUserNumber string) (responses.EntityResponse, error) {

	var entity responses.EntityResponse

	entityUser, err := repositories.FindEntityUser(currentUserNumber, entityNumber)
	if err != nil {
		return entity, errors.New("checking entity with user error: " + err.Error())
	}

	if entityUser.UserNumber == "" {
		return entity, errors.New("you can not switch to this entity")
	}

	var activeEntity models.ActiveUserEntity

	activeEntity, err = repositories.FindActiveUserEntityByUserNumber(currentUserNumber)
	if err != nil {
		return entity, errors.New("checking for any active entity error: " + err.Error())
	}

	if activeEntity.UserNumber != "" {
		err = repositories.DeleteActiveUserEntity(activeEntity)
		if err != nil {
			return entity, errors.New("removing previous active entity error: " + err.Error())
		}
	}

	err = repositories.SaveActiveUserEntity(models.ActiveUserEntity{
		EntityNumber: entityNumber,
		UserNumber:   currentUserNumber,
	})

	if err != nil {
		return entity, errors.New("setting current entity error: " + err.Error())
	}

	entityDetails, err := repositories.FindEntityByNumber(entityNumber)
	if err != nil {
		return entity, errors.New("getting entity details error: " + err.Error())
	}

	if entityDetails.Number == "" {
		return entity, errors.New("entity details not found")
	}

	entity = createEntityResponse(entityDetails)

	return entity, nil
}

func createEntityResponse(details models.Entity) responses.EntityResponse {
	return responses.EntityResponse{
		Name:    details.Name,
		Number:  details.Number,
		Email:   details.Email,
		Phone:   details.Phone,
		Address: details.Address,
	}
}

func GetCurrentEntity(currentUserNumber string) (responses.EntityResponse, error) {

	activeUserEntity, err := repositories.FindActiveUserEntityByUserNumber(currentUserNumber)
	if err != nil {
		return responses.EntityResponse{}, errors.New("getting current active entity error: " + err.Error())
	}

	if activeUserEntity.UserNumber == "" {
		return responses.EntityResponse{}, nil
	}

	return createEntityResponse(activeUserEntity.Entity), nil
}
