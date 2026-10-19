package services

import (
	"errors"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/repositories"
	"time"
)

type EntityService struct{}

var (
	// ErrNoActiveOrganisation means the caller has no valid active organisation.
	ErrNoActiveOrganisation = errors.New("no active organisation is set")
	// ErrEntityNotInActiveOrganisation means the client-supplied entity is not
	// the caller's active organisation. Permissions are always checked against
	// the active organisation, so a different target must be refused.
	ErrEntityNotInActiveOrganisation = errors.New("you can only modify your active organisation")
	// ErrEntityNotActive means the target organisation does not exist or has
	// already been deactivated.
	ErrEntityNotActive = errors.New("organisation not found or already deactivated")
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

	// A soft-deleted (deactivated) organisation must never be selectable. The
	// preload returns an empty Entity for a soft-deleted row.
	if entityUser.Entity.Number == "" {
		return entity, ErrEntityNotActive
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
		Name:      details.Name,
		Number:    details.Number,
		Email:     details.Email,
		Phone:     details.Phone,
		Address:   details.Address,
		Facebook:  details.Facebook,
		Instagram: details.Instagram,
		Twitter:   details.Twitter,
		Website:   details.Website,
		Whatsapp:  details.Whatsapp,
		Linkedin:  details.Linkedin,
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

	// Stage 0C-2: an active_user_entities row alone is not sufficient. The user
	// must still be a member of a live (non-deactivated) organisation. This makes
	// stale active rows fail closed for every caller, including Stage 0B's
	// organisation-scoped permission middleware.
	membership, err := repositories.FindEntityUser(currentUserNumber, activeUserEntity.EntityNumber)
	if err != nil {
		return responses.EntityResponse{}, errors.New("checking active organisation membership error: " + err.Error())
	}
	if membership.UserNumber == "" || membership.Entity.Number == "" {
		return responses.EntityResponse{}, nil
	}

	return createEntityResponse(membership.Entity), nil
}

func (entity *EntityService) UpdateEntity(updateEntityRequest requests.UpdateEntityRequest, currentUserNumber string) error {
	active, err := entity.GetCurrentEntity(currentUserNumber)
	if err != nil {
		return err
	}
	if active.Number == "" {
		return ErrNoActiveOrganisation
	}
	if updateEntityRequest.EntityNumber != active.Number {
		return ErrEntityNotInActiveOrganisation
	}

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
	update.Email = updateEntityRequest.Email
	update.Facebook = updateEntityRequest.Facebook
	update.Instagram = updateEntityRequest.Instagram
	update.Twitter = updateEntityRequest.Twitter
	update.Website = updateEntityRequest.Website
	update.Whatsapp = updateEntityRequest.Whatsapp
	update.Linkedin = updateEntityRequest.Linkedin

	err = repositories.SaveEntity(update)
	if err != nil {
		return errors.New("updating entity error: " + err.Error())
	}

	return nil
}

// DeactivateEntity removes an organisation from active use without destroying
// membership history or business records.
//
// It soft-deletes the entity (GORM DeletedAt) and removes the active-organisation
// selections for that entity. entity_users rows are deliberately retained: they
// are the membership-history record and are required for a future authorised
// re-activation. EMS, Education, licensing, audit and notification rows reference
// the organisation by value and are retained.
//
// Only the caller's active organisation may be deactivated; the active
// organisation is resolved server-side (membership + live entity), never from
// the client-supplied identifier alone.
func DeactivateEntity(entityNumber string, currentUserNumber string) error {
	active, err := (&EntityService{}).GetCurrentEntity(currentUserNumber)
	if err != nil {
		return err
	}
	if active.Number == "" {
		return ErrNoActiveOrganisation
	}
	if entityNumber != active.Number {
		return ErrEntityNotInActiveOrganisation
	}

	entity, err := repositories.FindEntityByNumber(entityNumber)
	if err != nil {
		return errors.New("checking organisation error: " + err.Error())
	}
	if entity.Number == "" {
		return ErrEntityNotActive
	}

	activeEntries, err := repositories.FindAllActiveEntries(entityNumber)
	if err != nil {
		return errors.New("getting active organisation selections: " + err.Error())
	}

	tx := config.DB.Begin()
	for _, activeEntry := range activeEntries {
		if err := repositories.DeleteActiveUserEntity(activeEntry, tx); err != nil {
			tx.Rollback()
			return errors.New("clearing active organisation selection failed")
		}
	}

	// The entity soft-delete must be inside the same transaction as the active
	// selection cleanup. entity_users are intentionally not deleted.
	if err := repositories.DeleteEntity(entity, tx); err != nil {
		tx.Rollback()
		return errors.New("deactivating organisation failed")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errors.New("deactivating organisation failed, try again")
	}

	return nil
}

func FindCalendarEvents(entityNumber string) ([]responses.Calendar, error) {

	var calendars = make([]responses.Calendar, 0)

	endDateAfter := time.Now()

	calendarEvents, err := repositories.FindCalendarEvent(entityNumber, endDateAfter)
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
