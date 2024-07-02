package services

import (
	"errors"
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/repositories"
)

type EntityService struct{}

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

	entity = createEntityResponse(entityUser.Entity)

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

func (entity *EntityService) GetCurrentEntity(currentUserNumber string) (responses.EntityResponse, error) {

	activeUserEntity, err := repositories.FindActiveUserEntityByUserNumber(currentUserNumber)
	if err != nil {
		return responses.EntityResponse{}, errors.New("getting current active entity error: " + err.Error())
	}

	if activeUserEntity.UserNumber == "" {
		return responses.EntityResponse{}, nil
	}

	return createEntityResponse(activeUserEntity.Entity), nil
}

func (entity *EntityService) UpdateEntity(updateEntityRequest requests.UpdateEntityRequest, currentUserNumber string) error {

	entityUser, err := repositories.FindEntityUser(currentUserNumber, updateEntityRequest.EntityNumber)
	if err != nil {
		return errors.New("checking entity with user error: " + err.Error())
	}

	if entityUser.UserNumber == "" {
		return errors.New("you are not attached to this entity")
	}

	checkName, err := repositories.FindEntityByName(updateEntityRequest.Name)
	if err != nil {
		return errors.New("checking name error: " + err.Error())
	}

	if checkName.Number != "" && checkName.Number != updateEntityRequest.EntityNumber {
		return errors.New("name already used")
	}

	update := entityUser.Entity
	update.Name = updateEntityRequest.Name
	update.Address = updateEntityRequest.Address
	update.Phone = updateEntityRequest.Phone
	update.AddedBy = currentUserNumber

	err = repositories.SaveEntity(update)
	if err != nil {
		return errors.New("updating entity error: " + err.Error())
	}

	return nil
}

func DeleteEntity(entityNumber string, currentUserNumber string) error {

	checkUserOnEntity, err := repositories.FindEntityUser(currentUserNumber, entityNumber)
	if err != nil {
		return errors.New("checking entity with user error: " + err.Error())
	}

	if checkUserOnEntity.UserNumber == "" {
		return errors.New("you are not attached to this entity")
	}

	entityUsers, err := repositories.FindAllEntityUsers(entityNumber)
	if err != nil {
		return errors.New("getting all entity users: " + err.Error())
	}

	activeEntries, err := repositories.FindAllActiveEntries(entityNumber)
	if err != nil {
		return errors.New("getting all entity users: " + err.Error())
	}

	tx := config.DB.Begin()

	for _, activeEntry := range activeEntries {

		err = repositories.DeleteActiveUserEntity(activeEntry, tx)
		if err != nil {
			tx.Rollback()
			fmt.Println(err.Error())
			return errors.New("clearing active entity error")
		}

	}

	for _, entityUser := range entityUsers {

		err = repositories.DeleteEntityUser(entityUser, tx)
		if err != nil {
			tx.Rollback()
			fmt.Println(err.Error())
			return errors.New("removing users from entity failed")
		}

	}

	err = repositories.DeleteEntity(checkUserOnEntity.Entity)
	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		return errors.New("removing entity failed, confirm if this entity has never been used")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		return errors.New("confirm remove error, try again")
	}

	return nil
}

func FindCalendarEvents(entityNumber string) ([]responses.Calendar, error) {

	var calendars = make([]responses.Calendar, 0)

	calendarEvents, err := repositories.FindCalendarEvent(entityNumber)
	if err != nil {
		return calendars, errors.New("finding calendar events error: " + err.Error())
	}

	for _, calendarEvent := range calendarEvents {
		calendars = append(calendars, createCalendarResponse(calendarEvent))
	}

	return calendars, nil
}

func createCalendarResponse(calendarEvent models.CalendarEvent) responses.Calendar {
	return responses.Calendar{
		Name:      calendarEvent.Name,
		Number:    calendarEvent.Number,
		StartDate: calendarEvent.StartDate,
		EndDate:   calendarEvent.EndDate,
	}
}
