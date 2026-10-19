package controllers

import (
	"errors"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	authRepository "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/repository"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	eosRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/repositories"
	permissionModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/permissions/models"
	permissionRepository "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/permissions/repository"
	"github.com/gin-gonic/gin"
)

// roleOrganisationContext resolves the authenticated actor and the actor's
// current organisation. Every role-management handler uses this so organisation
// scoping is always derived server-side from active_user_entities.
func roleOrganisationContext(c *gin.Context) (string, *authModel.User, bool) {
	actor, err := commons.GetUserFromHeader(c)
	if err != nil || actor == nil || actor.Number == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return "", nil, false
	}

	organisation, err := commons.GetUserEntityNumber(c)
	if err != nil || organisation == "" {
		c.JSON(http.StatusConflict, gin.H{"message": "No active organisation is selected"})
		return "", nil, false
	}

	return organisation, actor, true
}

// validatePermissionCatalogue rejects permission strings that are not part of
// the authoritative catalogue (Education catalogue plus EMS Core/EOS strings).
func validatePermissionCatalogue(permissions []string) error {
	known := commons.KnownPermissions()
	for _, permission := range permissions {
		if _, ok := known[permission]; !ok {
			return errors.New("unknown permission: " + permission)
		}
	}
	return nil
}

// actorMayDelegate reports whether the actor may grant the requested
// permissions in the given organisation. An installation-level Admin may grant
// any catalogue permission; every other actor may only grant permissions they
// themselves hold through roles in that same organisation.
func actorMayDelegate(actor *authModel.User, organisation string, permissions []string) (bool, error) {
	if actor.Type == "Admin" {
		return true, nil
	}

	granted, err := authRepository.FetchPermissionsByUserAndEntity(actor.Number, organisation)
	if err != nil {
		return false, err
	}

	allowed := make(map[string]struct{}, len(granted))
	for _, permission := range granted {
		allowed[permission] = struct{}{}
	}
	for _, permission := range permissions {
		if _, ok := allowed[permission]; !ok {
			return false, nil
		}
	}
	return true, nil
}

// loadRoleInOrganisation returns the role only when it belongs to the actor's
// current organisation. A cross-organisation role is reported as not found so
// its existence is not disclosed.
func loadRoleInOrganisation(c *gin.Context, roleNumber string, organisation string) (*permissionModels.Role, bool) {
	if roleNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Role number is required"})
		return nil, false
	}

	role, err := permissionRepository.FindRoleByNumber(roleNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Role not found"})
		return nil, false
	}
	if role.EntityNumber != organisation {
		c.JSON(http.StatusNotFound, gin.H{"message": "Role not found"})
		return nil, false
	}
	return role, true
}

func rolePermissionStrings(roleNumber string) ([]string, error) {
	var permissions []permissionModels.RolePermission
	if err := config.DB.Where("role_number = ?", roleNumber).Find(&permissions).Error; err != nil {
		return nil, err
	}

	result := make([]string, 0, len(permissions))
	for _, permission := range permissions {
		result = append(result, permission.Permission)
	}
	return result, nil
}

func userIsMemberOfOrganisation(userNumber string, organisation string) (bool, error) {
	membership, err := eosRepositories.FindEntityUser(userNumber, organisation)
	if err != nil {
		return false, err
	}
	return membership.UserNumber != "", nil
}
