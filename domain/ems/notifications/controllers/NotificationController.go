package controllers

import (
	"errors"
	"net/http"
	"strconv"

	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	notificationServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/notifications/services"
	"github.com/gin-gonic/gin"
)

// FetchNotifications lists the authenticated recipient's in-app notifications
// (recipient-scoped, organisation-scoped).
func FetchNotifications(c *gin.Context) {
	userNumber, ok := currentUserNumber(c)
	if !ok {
		return
	}
	limit, offset := paging(c)
	notifications, total, err := notificationServices.ListRecipientNotifications(
		userNumber, c.Query("type"), c.Query("unread") == "true", limit, offset)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notifications fetched", "data": notifications, "total": total})
}

// FetchNotification returns one of the authenticated recipient's notifications.
func FetchNotification(c *gin.Context) {
	userNumber, ok := currentUserNumber(c)
	if !ok {
		return
	}
	id, valid := parseID(c)
	if !valid {
		return
	}
	notification, err := notificationServices.GetRecipientNotification(userNumber, id)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification fetched", "data": notification})
}

// MarkNotificationRead marks one of the authenticated recipient's notifications
// read. It updates only read_at and never edits the notification content.
func MarkNotificationRead(c *gin.Context) {
	userNumber, ok := currentUserNumber(c)
	if !ok {
		return
	}
	id, valid := parseID(c)
	if !valid {
		return
	}
	notification, err := notificationServices.MarkRecipientNotificationRead(userNumber, id)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked read", "data": notification})
}

// currentUserNumber resolves the JWT-derived actor from the request context.
func currentUserNumber(c *gin.Context) (string, bool) {
	value, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, commons.Feedback("Unauthorized"))
		return "", false
	}
	user, ok := value.(authModels.User)
	if !ok || user.Number == "" {
		c.JSON(http.StatusUnauthorized, commons.Feedback("Unauthorized"))
		return "", false
	}
	return user.Number, true
}

// paging reads bounded limit/offset query parameters.
func paging(c *gin.Context) (int, int) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if err != nil || limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}
	return limit, offset
}

// parseID reads and validates the :number path parameter.
func parseID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("number"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid notification identifier"))
		return 0, false
	}
	return uint(id), true
}

// writeServiceError maps notification service errors to HTTP status codes.
func writeServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, notificationServices.ErrNoActiveOrganisation):
		c.JSON(http.StatusConflict, commons.Feedback(err.Error()))
	case errors.Is(err, notificationServices.ErrNotificationNotFound):
		c.JSON(http.StatusNotFound, commons.Feedback(err.Error()))
	case errors.Is(err, notificationServices.ErrInvalidNotification):
		c.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
	default:
		c.JSON(http.StatusInternalServerError, commons.Feedback("Something went wrong"))
	}
}
