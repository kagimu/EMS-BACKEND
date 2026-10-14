package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	stakeholderModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/stakeholders/models"
	"github.com/gin-gonic/gin"
)

// Setup state values returned by GET /setup/status. They are deliberately the
// only information this endpoint exposes.
const (
	SetupStateInitialized   = "initialized"
	SetupStateUninitialized = "uninitialized"
)

// Status reports whether the installation has completed its initial owner
// setup, so the public frontend root can decide between the first-time setup
// experience and the login experience.
//
// It is public and read-only by design and returns only the minimum routing
// state. It never discloses account counts, identities, contact details,
// tokens, roles, permissions or licence data, and it does not trust any
// client-supplied value.
//
// "Initialized" means at least one ACTIVE owner stakeholder exists
// (ems_stakeholders.stakeholder_type='owner' AND status='ACTIVE'). That is the
// explicit "owner established" marker; it implies an active administrator and a
// current organisation, and the last-owner invariant makes it durable. A fresh
// installation, an inactive initial administrator, or an organisation without an
// owner all report "uninitialized" — they are never treated as complete.
func Status(c *gin.Context) {
	var activeOwners int64
	err := config.DB.
		Model(&stakeholderModel.Stakeholder{}).
		Where("stakeholder_type = ?", stakeholderModel.StakeholderTypeOwner).
		Where("status = ?", constants.StatusActive).
		Count(&activeOwners).Error
	if err != nil {
		// Never guess the state: the caller must show a retry state rather than
		// assume the installation is new or complete.
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": "Setup status is temporarily unavailable"})
		return
	}

	state := SetupStateUninitialized
	if activeOwners > 0 {
		state = SetupStateInitialized
	}

	c.JSON(http.StatusOK, gin.H{"state": state})
}
