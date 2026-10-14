package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/services"
	"github.com/gin-gonic/gin"
)

// ActivateInstallation registers (or reactivates) the caller organisation's
// installation against a licence key. Capacity is enforced transactionally.
func ActivateInstallation(c *gin.Context) {
	actorNumber, ok := currentUserNumber(c)
	if !ok {
		return
	}

	var request requests.ActivateInstallationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}

	response, err := services.ActivateInstallation(actorNumber, request)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Installation activated", "data": response})
}

// FetchClientInstallations lists the authenticated organisation's own
// installations. Never exposes the licence key or another organisation.
func FetchClientInstallations(c *gin.Context) {
	actorNumber, ok := currentUserNumber(c)
	if !ok {
		return
	}

	installations, err := services.ListClientInstallations(actorNumber)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Installations fetched", "data": installations})
}

// DeactivateOwnInstallation performs a permissible client self-deactivation.
// The installation stops being usable but keeps consuming capacity.
func DeactivateOwnInstallation(c *gin.Context) {
	actorNumber, ok := currentUserNumber(c)
	if !ok {
		return
	}

	installation, err := services.DeactivateOwnInstallation(actorNumber, c.Param("number"))
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Installation deactivated", "data": installation})
}

// FetchInstallations lists installations for the control plane (deliberately
// global; guarded by the explicit registry authority).
func FetchInstallations(c *gin.Context) {
	installations, err := services.ListInstallations(
		c.Query("licence"), c.Query("organisation"), c.Query("status"))
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Installations fetched", "data": installations})
}

// FetchInstallation returns one installation for the control plane.
func FetchInstallation(c *gin.Context) {
	installation, err := services.GetInstallation(c.Param("number"))
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Installation fetched", "data": installation})
}

// ReleaseInstallation is the explicit control-plane release — the only action
// that frees capacity. It records actor, timestamp and reason.
func ReleaseInstallation(c *gin.Context) {
	actorNumber, ok := currentUserNumber(c)
	if !ok {
		return
	}

	var request requests.ReleaseInstallationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}

	installation, err := services.ReleaseInstallation(actorNumber, c.Param("number"), request.Reason)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Installation released", "data": installation})
}
