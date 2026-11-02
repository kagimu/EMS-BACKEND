package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	ownerResponses "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/owner/api/responses"
	ownerServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/owner/services"
	"github.com/gin-gonic/gin"
)

func pagination(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	return page, perPage
}

func writeOwnerError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ownerServices.ErrOrganisationNotFound):
		c.JSON(http.StatusNotFound, commons.Feedback("Organisation not found"))
	case strings.Contains(strings.ToLower(err.Error()), "invalid"):
		c.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
	default:
		c.JSON(http.StatusInternalServerError, commons.Feedback("The owner request could not be completed"))
	}
}

// Access confirms the authenticated actor is the configured software owner.
func Access(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, commons.Feedback("Unauthorized"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Owner access confirmed", "data": ownerResponses.OwnerAccessResponse{
		Owner:      true,
		UserNumber: user.Number,
	}})
}

// Overview returns owner dashboard aggregates.
func Overview(c *gin.Context) {
	overview, err := ownerServices.Overview()
	if err != nil {
		writeOwnerError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Owner overview fetched", "data": overview})
}

// FetchOrganisations returns the paginated customer/organisation directory.
func FetchOrganisations(c *gin.Context) {
	page, perPage := pagination(c)
	result, err := ownerServices.ListOrganisations(c.Query("search"), c.Query("status"), page, perPage)
	if err != nil {
		writeOwnerError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Organisations fetched", "data": result})
}

// FetchOrganisation returns an owner customer profile.
func FetchOrganisation(c *gin.Context) {
	detail, err := ownerServices.OrganisationDetail(c.Param("number"))
	if err != nil {
		writeOwnerError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Organisation fetched", "data": detail})
}

// FetchLicences returns the paginated owner licence list (never keys).
func FetchLicences(c *gin.Context) {
	page, perPage := pagination(c)
	result, err := ownerServices.ListLicences(c.Query("search"), c.Query("status"), page, perPage)
	if err != nil {
		writeOwnerError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Licences fetched", "data": result})
}

// FetchLicence returns one licence and its event timeline (never the key).
func FetchLicence(c *gin.Context) {
	detail, err := ownerServices.LicenceDetail(c.Param("number"))
	if err != nil {
		writeOwnerError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Licence fetched", "data": detail})
}

// FetchInstallations returns the paginated owner installation list.
func FetchInstallations(c *gin.Context) {
	page, perPage := pagination(c)
	result, err := ownerServices.ListInstallations(c.Query("search"), c.Query("status"), page, perPage)
	if err != nil {
		writeOwnerError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Installations fetched", "data": result})
}

// FetchInstallation returns one installation with its licence and events.
func FetchInstallation(c *gin.Context) {
	detail, err := ownerServices.InstallationDetail(c.Param("number"))
	if err != nil {
		writeOwnerError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Installation fetched", "data": detail})
}

// FetchAddons returns the addon catalogue for licence configuration.
func FetchAddons(c *gin.Context) {
	modules, err := ownerServices.Modules()
	if err != nil {
		writeOwnerError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Addons fetched", "data": modules})
}
