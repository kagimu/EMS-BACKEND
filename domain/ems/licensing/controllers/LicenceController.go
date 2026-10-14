package controllers

import (
	"errors"
	"net/http"
	"strconv"

	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/services"
	"github.com/gin-gonic/gin"
)

// IssueLicence creates and assigns a licence (control plane). The opaque
// licence key is returned exactly once, at issuance.
func IssueLicence(c *gin.Context) {
	actorNumber, ok := currentUserNumber(c)
	if !ok {
		return
	}

	var request requests.IssueLicenceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}

	licence, err := services.IssueLicence(actorNumber, request)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Licence issued", "data": licence})
}

// FetchLicences lists licences (control plane; deliberately global).
func FetchLicences(c *gin.Context) {
	licences, err := services.ListLicences(c.Query("organisation"), c.Query("status"))
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Licences fetched", "data": licences})
}

// FetchLicence returns one licence (control plane).
func FetchLicence(c *gin.Context) {
	licence, err := services.GetLicence(c.Param("number"))
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Licence fetched", "data": licence})
}

// ConfigureLicence updates validity, installation allowance and entitlements.
func ConfigureLicence(c *gin.Context) {
	actorNumber, ok := currentUserNumber(c)
	if !ok {
		return
	}

	var request requests.ConfigureLicenceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}

	licence, err := services.ConfigureLicence(actorNumber, c.Param("number"), request)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Licence updated", "data": licence})
}

// RenewLicence extends an existing licence, preserving its identity (F-Q6).
func RenewLicence(c *gin.Context) {
	actorNumber, ok := currentUserNumber(c)
	if !ok {
		return
	}

	var request requests.RenewLicenceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}

	licence, err := services.RenewLicence(actorNumber, c.Param("number"), request)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Licence renewed", "data": licence})
}

// SuspendLicence suspends an ACTIVE licence.
func SuspendLicence(c *gin.Context) {
	actorNumber, ok := currentUserNumber(c)
	if !ok {
		return
	}

	licence, err := services.SuspendLicence(actorNumber, c.Param("number"))
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Licence suspended", "data": licence})
}

// RevokeLicence revokes an ACTIVE or SUSPENDED licence (terminal).
func RevokeLicence(c *gin.Context) {
	actorNumber, ok := currentUserNumber(c)
	if !ok {
		return
	}

	licence, err := services.RevokeLicence(actorNumber, c.Param("number"))
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Licence revoked", "data": licence})
}

// ReplaceLicence issues a new licence that atomically supersedes an ACTIVE
// licence: the old licence is revoked, its capacity-holding installations are
// explicitly released with the mandatory reason, and the new licence becomes
// the organisation's single ACTIVE licence. The new key is returned once, at
// issuance, like ordinary issuance.
func ReplaceLicence(c *gin.Context) {
	actorNumber, ok := currentUserNumber(c)
	if !ok {
		return
	}

	var request requests.ReplaceLicenceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}

	licence, err := services.ReplaceLicence(actorNumber, c.Param("number"), request)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Replacement licence issued", "data": licence})
}

// ClientLicenceStatus returns the authenticated organisation's active licence
// (client plane). The licence key is never exposed.
func ClientLicenceStatus(c *gin.Context) {
	actorNumber, ok := currentUserNumber(c)
	if !ok {
		return
	}

	status, err := services.ClientLicenceStatus(actorNumber)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Licence status fetched", "data": status})
}

// RequestRenewal records the authenticated organisation's renewal request on
// its latest licence. It changes no licence state and never exposes the key.
func RequestRenewal(c *gin.Context) {
	actorNumber, ok := currentUserNumber(c)
	if !ok {
		return
	}

	status, err := services.RequestRenewal(actorNumber)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Renewal requested", "data": status})
}

// FetchLicenceEvents lists append-only licensing facts (control plane). Facts
// are recorded by Foundation F for Foundation G to reference later; this view
// is read-only and never exposes the licence key.
func FetchLicenceEvents(c *gin.Context) {
	limit, offset := paging(c)
	events, total, err := services.ListLicenceEvents(
		c.Query("licence"), c.Query("organisation"), c.Query("type"), limit, offset)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Licence events fetched", "data": events, "total": total})
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

// currentUserNumber resolves the JWT-derived actor from the request context.
// It writes the error response itself and returns ok=false when the request
// must not proceed.
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

// writeServiceError maps service errors to HTTP status codes.
func writeServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrNoActiveOrganisation):
		c.JSON(http.StatusConflict, commons.Feedback(err.Error()))
	case errors.Is(err, services.ErrLicenceNotFound),
		errors.Is(err, services.ErrInstallationNotFound):
		c.JSON(http.StatusNotFound, commons.Feedback(err.Error()))
	case errors.Is(err, services.ErrActiveLicenceExists),
		errors.Is(err, services.ErrInvalidTransition),
		errors.Is(err, services.ErrLicenceNotActive),
		errors.Is(err, services.ErrLicenceNotYetValid),
		errors.Is(err, services.ErrLicenceExpired),
		errors.Is(err, services.ErrCapacityReached),
		errors.Is(err, services.ErrInstallationReleased),
		errors.Is(err, services.ErrInstallationNotActive),
		errors.Is(err, services.ErrInstallationAlreadyDeactivated):
		c.JSON(http.StatusConflict, commons.Feedback(err.Error()))
	case errors.Is(err, services.ErrAddonNotLicensed):
		c.JSON(http.StatusForbidden, commons.Feedback(err.Error()))
	case errors.Is(err, services.ErrOrganisationNotFound),
		errors.Is(err, services.ErrInvalidValidity),
		errors.Is(err, services.ErrInvalidMaxInstallations),
		errors.Is(err, services.ErrUnknownAddon),
		errors.Is(err, services.ErrInvalidStatus),
		errors.Is(err, services.ErrValidToNotExtended),
		errors.Is(err, services.ErrInvalidLicenceKey),
		errors.Is(err, services.ErrInvalidInstallationIdentity),
		errors.Is(err, services.ErrReleaseReasonRequired),
		errors.Is(err, services.ErrInvalidReplacementLink):
		c.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
	default:
		c.JSON(http.StatusInternalServerError, commons.Feedback("Something went wrong"))
	}
}
