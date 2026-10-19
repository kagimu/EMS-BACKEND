package controllers

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/permissions/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/permissions/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/permissions/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"net/http"
)

func CreateRole(c *gin.Context) {
	organisation, actor, ok := roleOrganisationContext(c)
	if !ok {
		return
	}

	var payload requests.NewRoleRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to validate input"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return

	}

	if payload.EntityNumber != organisation {
		c.JSON(http.StatusForbidden, gin.H{"message": "You can only manage roles for your current organisation"})
		return
	}

	if err := validatePermissionCatalogue(payload.Permissions); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	allowed, err := actorMayDelegate(actor, organisation, payload.Permissions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to verify role authority"})
		return
	}
	if !allowed {
		c.JSON(http.StatusForbidden, gin.H{"message": "You cannot grant permissions you do not hold"})
		return
	}

	roleNumber := uuid.NewV4().String()
	saveRole := models.Role{
		Number:       roleNumber,
		Name:         payload.Name,
		Status:       constants.StatusActive,
		UserNumber:   actor.Number,
		Description:  payload.Description,
		EntityNumber: payload.EntityNumber,
	}

	tx := config.DB.Begin()

	if err := tx.Create(&saveRole).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save role"})
		tx.Rollback()
		return
	}

	// save the permissions
	for _, permission := range payload.Permissions {
		savePermission := models.RolePermission{
			RoleNumber:    roleNumber,
			AddedByNumber: actor.Number,
			Permission:    permission,
		}

		if err := tx.Create(&savePermission).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save role permissions"})
			tx.Rollback()
			return
		}

	}

	tx.Commit()

	c.JSON(http.StatusOK, commons.Feedback("Role has been created"))
}

func FetchRoles(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	entity_number, _ := commons.GetUserEntityNumber(c)
	if entity_number == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Entity number is required"})
		return
	}

	var rolesResponse []responses.RolesList

	var roles []models.Role
	if err := config.DB.Scopes(filters.ByEntityNumber(entity_number)).Find(&roles).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch roles"})
		return
	}

	// fetch the permissions
	for _, role := range roles {
		var permissions []models.RolePermission
		// select only the permission column in the role_permissions table
		if err := config.DB.Where("role_number = ?", role.Number).Find(&permissions).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"message": "Role not found"})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch role details"})
			return
		}

		//prepare array of permissions
		var pms []string
		for _, permission := range permissions {
			pms = append(pms, permission.Permission)
		}

		rolesResponse = append(rolesResponse, responses.RolesList{
			Number:      role.Number,
			Name:        role.Name,
			Description: role.Description,
			Permissions: pms,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": rolesResponse})
}

func AssignRolesToAUser(c *gin.Context) {
	organisation, actor, ok := roleOrganisationContext(c)
	if !ok {
		return
	}

	type RoleAssignment struct {
		RoleNumber string `json:"role_number"`
		ExpiryDate string `json:"expiry_date"`
	}
	type Input struct {
		UserNumber string           `json:"user_number"`
		Roles      []RoleAssignment `json:"roles"`
	}

	var payload Input
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input- " + err.Error()})
		return
	}

	if len(payload.Roles) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "At least one role is required"})
		return
	}

	if payload.UserNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User number is required"})
		return
	}

	var targetUser authModel.User
	if err := config.DB.Where("number = ?", payload.UserNumber).First(&targetUser).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	member, err := userIsMemberOfOrganisation(payload.UserNumber, organisation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to verify user membership"})
		return
	}
	if !member {
		c.JSON(http.StatusBadRequest, gin.H{"message": "The target user is not a member of this organisation"})
		return
	}

	for _, roleAssignment := range payload.Roles {
		role, ok := loadRoleInOrganisation(c, roleAssignment.RoleNumber, organisation)
		if !ok {
			return
		}

		rolePermissions, err := rolePermissionStrings(role.Number)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to verify role permissions"})
			return
		}

		allowed, err := actorMayDelegate(actor, organisation, rolePermissions)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to verify role authority"})
			return
		}
		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"message": "You cannot grant permissions you do not hold"})
			return
		}
	}

	tx := config.DB.Begin()

	for _, role := range payload.Roles {

		// convert the expiry date to a time object
		date, err := commons.ConvertStringToDate(role.ExpiryDate)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid date format"})
			return
		}
		saveRole := models.UserRole{
			RoleNumber: role.RoleNumber,
			UserNumber: payload.UserNumber,
			AddedBy:    actor.Number,
			ExpiryDate: date,
		}

		// validate if the role and user exists

		var userRole models.UserRole
		err = config.DB.Preload("Role").Where("user_number = ?", payload.UserNumber).Where("role_number = ?", role.RoleNumber).First(&userRole).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to assign a role to the specified user" + err.Error()})
			return
		}

		if err == nil {
			// if the role exists update the expiry date
			userRole.ExpiryDate = date
			if err := tx.Save(&userRole).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update role assignment"})
				return
			}

		} else {

			if err := tx.Create(&saveRole).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to assign a role to the specified user" + err.Error()})
				return
			}
		}

	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to assign roles to the user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User has been updated"})
}

func FetchRoleDetails(c *gin.Context) {
	organisation, _, ok := roleOrganisationContext(c)
	if !ok {
		return
	}

	role, ok := loadRoleInOrganisation(c, c.Param("role_number"), organisation)
	if !ok {
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": role})
}

func RemoveRoleAssignment(c *gin.Context) {
	organisation, _, ok := roleOrganisationContext(c)
	if !ok {
		return
	}

	type RemoveRoleRequest struct {
		UserNumber string `json:"user_number"`
		RoleNumber string `json:"role_number"`
	}

	var payload RemoveRoleRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input- " + err.Error()})
		return
	}

	if len(payload.RoleNumber) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "At least one role is required"})
		return
	}

	if _, ok := loadRoleInOrganisation(c, payload.RoleNumber, organisation); !ok {
		return
	}

	// delete the role assignment

	if err := config.DB.Where("user_number = ?", payload.UserNumber).Where("role_number = ?", payload.RoleNumber).Delete(&models.UserRole{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to remove role assignment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role assignment has been updated"})
}

func UpdateRole(ctx *gin.Context) {
	organisation, actor, ok := roleOrganisationContext(ctx)
	if !ok {
		return
	}

	role, ok := loadRoleInOrganisation(ctx, ctx.Param("role_number"), organisation)
	if !ok {
		return
	}
	roleNumber := role.Number

	var payload requests.EditRoleRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to validate payload"})
		return
	}

	if err := payload.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := validatePermissionCatalogue(payload.Permissions); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	allowed, err := actorMayDelegate(actor, organisation, payload.Permissions)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to verify role authority"})
		return
	}
	if !allowed {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "You cannot grant permissions you do not hold"})
		return
	}

	role.Name = payload.Name
	role.Description = payload.Description
	role.Status = constants.StatusActive
	role.UserNumber = actor.Number

	if err := config.DB.Save(role).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update role"})
		return

	}

	tx := config.DB.Begin()

	// delete the permissions
	if err := tx.Where("role_number = ?", roleNumber).Delete(&models.RolePermission{}).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update role permissions"})
		return
	}

	// save the permissions
	for _, permission := range payload.Permissions {
		savePermission := models.RolePermission{
			RoleNumber:    roleNumber,
			AddedByNumber: actor.Number,
			Permission:    permission,
		}

		if err := tx.Create(&savePermission).Error; err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save role permissions"})
			return
		}

	}

	tx.Commit()

	ctx.JSON(http.StatusOK, gin.H{"message": "Role has been updated"})

}

func DeleteRole(context *gin.Context) {
	organisation, _, ok := roleOrganisationContext(context)
	if !ok {
		return
	}

	role, ok := loadRoleInOrganisation(context, context.Param("role_number"), organisation)
	if !ok {
		return
	}
	roleNumber := role.Number

	tx := config.DB.Begin()
	// delete associated permissions
	if err := tx.Where("role_number = ?", roleNumber).Delete(&models.RolePermission{}).Error; err != nil {
		tx.Rollback()
		context.JSON(500, gin.H{"message": "Failed to delete role permissions"})
		return
	}

	// delete the role
	if err := tx.Where("number = ?", roleNumber).Delete(&models.Role{}).Error; err != nil {
		tx.Rollback()
		context.JSON(500, gin.H{"message": "Failed to delete role"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		context.JSON(500, gin.H{"message": "Failed to delete role"})
		return
	}

	context.JSON(200, gin.H{"message": "Role has been deleted"})

	//

}
