package controllers

import (
	"errors"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	licenseRequests "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/api/requests"
	licenseServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/services"
	"github.com/gin-gonic/gin"
)

func flexibleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, licenseServices.ErrLicenceNotFound):
		c.JSON(http.StatusNotFound, commons.Feedback("Licence not found"))
	case errors.Is(err, licenseServices.ErrFlexibleValidation),
		errors.Is(err, licenseServices.ErrScopeQuantityBelowUsage),
		errors.Is(err, licenseServices.ErrUnitCapacityReached),
		errors.Is(err, licenseServices.ErrUnitDuplicate),
		errors.Is(err, licenseServices.ErrUnitCrossOrganisation),
		errors.Is(err, licenseServices.ErrScopeNotFound):
		c.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
	default:
		c.JSON(http.StatusInternalServerError, commons.Feedback("The flexible licensing request could not be completed"))
	}
}

func FetchFlexibleLicence(c *gin.Context) {
	view, err := licenseServices.FlexibleLicenceView(c.Param("number"))
	if err != nil {
		flexibleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Flexible licence fetched", "data": view})
}

func ConfigureFlexibleLicence(c *gin.Context) {
	actor, ok := currentUserNumber(c)
	if !ok {
		return
	}
	var request licenseRequests.ConfigureFlexibleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}
	view, err := licenseServices.ConfigureFlexible(actor, c.Param("number"), request)
	if err != nil {
		flexibleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Licence scopes and entitlements updated", "data": view})
}

func RegisterLicensedUnit(c *gin.Context) {
	actor, ok := currentUserNumber(c)
	if !ok {
		return
	}
	var request licenseRequests.RegisterUnitRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}
	unit, err := licenseServices.RegisterLicensedUnit(actor, c.Param("number"), request)
	if err != nil {
		flexibleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Licensed unit registered", "data": unit})
}

func ReleaseLicensedUnit(c *gin.Context) {
	actor, ok := currentUserNumber(c)
	if !ok {
		return
	}
	var request licenseRequests.ReleaseUnitRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}
	if err := licenseServices.ReleaseLicensedUnit(actor, c.Param("unit"), request.Reason); err != nil {
		flexibleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Licensed unit released"})
}
