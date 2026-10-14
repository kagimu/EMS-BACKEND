package controllers

import (
	"errors"
	"net/http"
	"strconv"

	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	auditServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/audit/services"
	"github.com/gin-gonic/gin"
)

// FetchOrganisationAuditLogs lists the authenticated organisation's audit facts
// (client plane, organisation-scoped).
func FetchOrganisationAuditLogs(c *gin.Context) {
	userNumber, ok := currentUserNumber(c)
	if !ok {
		return
	}
	limit, offset := paging(c)
	logs, total, err := auditServices.ListOrganisationAuditLogs(
		userNumber, c.Query("actor"), c.Query("action"), c.Query("entity_type"), c.Query("entity_id"), limit, offset)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Audit records fetched", "data": logs, "total": total})
}

// FetchOrganisationAuditLog returns one of the authenticated organisation's
// audit facts (client plane, organisation-scoped).
func FetchOrganisationAuditLog(c *gin.Context) {
	userNumber, ok := currentUserNumber(c)
	if !ok {
		return
	}
	id, valid := parseID(c)
	if !valid {
		return
	}
	log, err := auditServices.GetOrganisationAuditLog(userNumber, id)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Audit record fetched", "data": log})
}

// FetchAuditRegistryLogs lists audit facts across organisations (control plane).
// It is read-only and never exposes secrets.
func FetchAuditRegistryLogs(c *gin.Context) {
	limit, offset := paging(c)
	logs, total, err := auditServices.ListAuditLogs(
		"", c.Query("organisation"), c.Query("actor"), c.Query("action"), c.Query("entity_type"), c.Query("entity_id"), limit, offset)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Audit records fetched", "data": logs, "total": total})
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
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid audit record identifier"))
		return 0, false
	}
	return uint(id), true
}

// writeServiceError maps audit service errors to HTTP status codes.
func writeServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, auditServices.ErrNoActiveOrganisation):
		c.JSON(http.StatusConflict, commons.Feedback(err.Error()))
	case errors.Is(err, auditServices.ErrAuditLogNotFound):
		c.JSON(http.StatusNotFound, commons.Feedback(err.Error()))
	case errors.Is(err, auditServices.ErrInvalidAuditRecord):
		c.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
	default:
		c.JSON(http.StatusInternalServerError, commons.Feedback("Something went wrong"))
	}
}
