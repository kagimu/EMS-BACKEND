package controllers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	customfieldModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/customfields/models"
	customfieldServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/customfields/services"
	licensingServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/services"
	"github.com/gin-gonic/gin"
)

func currentOrg(c *gin.Context) (string, string, bool) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, commons.Feedback("Unauthorized"))
		return "", "", false
	}
	org, err := licensingServices.OrganisationNumberForUser(user.Number)
	if err != nil {
		c.JSON(http.StatusConflict, commons.Feedback("No active organisation for this account"))
		return "", "", false
	}
	return org, user.Number, true
}

func fieldError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, customfieldServices.ErrDefinitionNotFound):
		c.JSON(http.StatusNotFound, commons.Feedback(err.Error()))
	case errors.Is(err, customfieldServices.ErrRequiredField),
		errors.Is(err, customfieldServices.ErrInvalidValue),
		errors.Is(err, customfieldServices.ErrInvalidDropdown),
		errors.Is(err, customfieldServices.ErrInvalidFieldType),
		errors.Is(err, customfieldServices.ErrSelectNeedsOptions),
		errors.Is(err, customfieldServices.ErrLabelRequired),
		errors.Is(err, customfieldServices.ErrKeyRequired),
		errors.Is(err, customfieldServices.ErrKeyExists),
		errors.Is(err, customfieldServices.ErrTypeChangeWithValue):
		c.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
	default:
		c.JSON(http.StatusInternalServerError, commons.Feedback("The request could not be completed"))
	}
}

type definitionRequest struct {
	Label        string   `json:"label"`
	Key          string   `json:"key"`
	Type         string   `json:"type"`
	Description  string   `json:"description"`
	Required     bool     `json:"required"`
	DefaultValue string   `json:"default_value"`
	Options      []string `json:"options"`
	DisplayOrder int      `json:"display_order"`
	Active       bool     `json:"active"`
}

// List returns the organisation's custom-field definitions and values.
func List(c *gin.Context) {
	org, _, ok := currentOrg(c)
	if !ok {
		return
	}
	definitions, err := customfieldServices.ListDefinitions(org, customfieldModels.ScopeOrganisation)
	if err != nil {
		fieldError(c, err)
		return
	}
	values, err := customfieldServices.ListValues(org, customfieldModels.ScopeOrganisation)
	if err != nil {
		fieldError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"definitions": definitions, "values": values}})
}

// Create adds a custom-field definition (authorised administrators only).
func Create(c *gin.Context) {
	org, actor, ok := currentOrg(c)
	if !ok {
		return
	}
	var req definitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid field definition"))
		return
	}
	def, err := customfieldServices.CreateDefinition(org, customfieldModels.ScopeOrganisation, actor, toInput(req))
	if err != nil {
		fieldError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": def})
}

// Update edits a custom-field definition.
func Update(c *gin.Context) {
	org, _, ok := currentOrg(c)
	if !ok {
		return
	}
	var req definitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid field definition"))
		return
	}
	def, err := customfieldServices.UpdateDefinition(org, c.Param("number"), toInput(req))
	if err != nil {
		fieldError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": def})
}

type valuesRequest struct {
	Values map[string]string `json:"values"`
}

// SaveValues validates and persists custom-field values for the organisation.
func SaveValues(c *gin.Context) {
	org, actor, ok := currentOrg(c)
	if !ok {
		return
	}
	var req valuesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid values payload"))
		return
	}
	if err := customfieldServices.SaveValues(org, customfieldModels.ScopeOrganisation, actor, req.Values); err != nil {
		fieldError(c, err)
		return
	}
	c.JSON(http.StatusOK, commons.Feedback("Custom field values saved"))
}

func toInput(req definitionRequest) customfieldServices.DefinitionInput {
	options := make([]string, 0, len(req.Options))
	for _, option := range req.Options {
		if trimmed := strings.TrimSpace(option); trimmed != "" {
			options = append(options, trimmed)
		}
	}
	return customfieldServices.DefinitionInput{
		Label:        req.Label,
		Key:          req.Key,
		Type:         req.Type,
		Description:  req.Description,
		Required:     req.Required,
		DefaultValue: req.DefaultValue,
		Options:      options,
		DisplayOrder: req.DisplayOrder,
		Active:       req.Active,
	}
}
