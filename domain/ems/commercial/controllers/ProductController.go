package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	commercialRequests "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/commercial/api/requests"
	commercialServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/commercial/services"
	"github.com/gin-gonic/gin"
)

func productPage(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	return page, perPage
}

func productActor(c *gin.Context) (string, bool) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil || user == nil || user.Number == "" {
		c.JSON(http.StatusUnauthorized, commons.Feedback("Unauthorized"))
		return "", false
	}
	return user.Number, true
}

func writeProductError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, commercialServices.ErrProductNotFound), errors.Is(err, commercialServices.ErrComponentNotFound):
		c.JSON(http.StatusNotFound, commons.Feedback(err.Error()))
	case errors.Is(err, commercialServices.ErrProductExists), errors.Is(err, commercialServices.ErrComponentExists):
		c.JSON(http.StatusConflict, commons.Feedback(err.Error()))
	case errors.Is(err, commercialServices.ErrProductValidation),
		errors.Is(err, commercialServices.ErrProductTransition),
		errors.Is(err, commercialServices.ErrComponentValidation):
		c.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
	default:
		if containsText(err.Error(), "invalid") {
			c.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, commons.Feedback("The commercial product request could not be completed"))
	}
}

func FetchCommercialProducts(c *gin.Context) {
	page, perPage := productPage(c)
	result, err := commercialServices.ListProducts(c.Query("search"), c.Query("status"), page, perPage)
	if err != nil {
		writeProductError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Commercial products fetched", "data": result})
}

func FetchCommercialProduct(c *gin.Context) {
	result, err := commercialServices.GetProduct(c.Param("number"))
	if err != nil {
		writeProductError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Commercial product fetched", "data": result})
}

func CreateCommercialProduct(c *gin.Context) {
	actor, ok := productActor(c)
	if !ok {
		return
	}
	var request commercialRequests.CommercialProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}
	result, err := commercialServices.CreateProduct(actor, request)
	if err != nil {
		writeProductError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Commercial product created", "data": result})
}

func UpdateCommercialProduct(c *gin.Context) {
	actor, ok := productActor(c)
	if !ok {
		return
	}
	var request commercialRequests.CommercialProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}
	result, err := commercialServices.UpdateProduct(actor, c.Param("number"), request)
	if err != nil {
		writeProductError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Commercial product updated", "data": result})
}

func SetCommercialProductStatus(c *gin.Context) {
	actor, ok := productActor(c)
	if !ok {
		return
	}
	var request commercialRequests.CommercialProductStatusRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}
	result, err := commercialServices.SetProductStatus(actor, c.Param("number"), request)
	if err != nil {
		writeProductError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Commercial product status updated", "data": result})
}

func FetchCommercialProductComponents(c *gin.Context) {
	result, err := commercialServices.ListComponents(c.Param("number"))
	if err != nil {
		writeProductError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Commercial product components fetched", "data": result})
}

func CreateCommercialProductComponent(c *gin.Context) {
	actor, ok := productActor(c)
	if !ok {
		return
	}
	var request commercialRequests.CommercialComponentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}
	result, err := commercialServices.CreateComponent(actor, c.Param("number"), request)
	if err != nil {
		writeProductError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Commercial product component added", "data": result})
}

func UpdateCommercialProductComponent(c *gin.Context) {
	actor, ok := productActor(c)
	if !ok {
		return
	}
	var request commercialRequests.CommercialComponentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}
	result, err := commercialServices.UpdateComponent(actor, c.Param("number"), c.Param("componentNumber"), request)
	if err != nil {
		writeProductError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Commercial product component updated", "data": result})
}

func SetCommercialProductComponentStatus(c *gin.Context) {
	actor, ok := productActor(c)
	if !ok {
		return
	}
	var request commercialRequests.CommercialComponentStatusRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}
	result, err := commercialServices.SetComponentStatus(actor, c.Param("number"), c.Param("componentNumber"), request)
	if err != nil {
		writeProductError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Commercial product component status updated", "data": result})
}
