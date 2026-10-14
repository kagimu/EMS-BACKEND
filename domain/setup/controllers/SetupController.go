package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	licensingModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/models"
	stakeholderModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/stakeholders/models"
	eosModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	"github.com/gin-gonic/gin"
)

// next_action values returned by GET /setup/progress. Each maps to an existing
// application route; no route is invented.
const (
	NextActionCommercialSetup      = "commercial_setup"      // /setup/commercial
	NextActionResumeActivation     = "resume_activation"     // existing activation / reset flow
	NextActionCompleteSetup        = "complete_setup"        // /setup (registration/org/ownership)
	NextActionActivateSubscription = "activate_subscription" // /ems/subscription
	NextActionReady                = "ready"                 // /login then /ems
)

// educationAddonCode is the addon code used by the EMS catalogue.
const educationAddonCode = "education"

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

// Progress reports, as database facts only, which first-installation steps have
// been completed, plus the single next action the operator should take. It is
// public and read-only and never creates, repairs, activates or deletes
// anything.
//
// It exposes only booleans, one derived state and one next_action — never
// personal information, email addresses, passwords, reset tokens, licence keys
// or other secrets. It does not assume an ACTIVE owner implies every step is
// complete: each step is reported from its own records.
func Progress(c *gin.Context) {
	usersExist, err := countExists(&authModel.User{}, "")
	if err != nil {
		unavailable(c)
		return
	}
	adminActive, err := countExists(&authModel.User{}, "type = ? AND status = ?", constants.TypeAdmin, constants.StatusActive)
	if err != nil {
		unavailable(c)
		return
	}
	ownerExists, err := countExists(&stakeholderModel.Stakeholder{}, "stakeholder_type = ? AND status = ?", stakeholderModel.StakeholderTypeOwner, constants.StatusActive)
	if err != nil {
		unavailable(c)
		return
	}
	organisationExists, err := countExists(&eosModel.Entity{}, "")
	if err != nil {
		unavailable(c)
		return
	}

	var licences []licensingModel.Licence
	if err := config.DB.Where("status = ?", licensingModel.StatusActive).Find(&licences).Error; err != nil {
		unavailable(c)
		return
	}
	subscriptionActive := len(licences) > 0
	educationEntitled := false
	for _, licence := range licences {
		if licenceContainsAddon(licence.LicensedAddons, educationAddonCode) {
			educationEntitled = true
			break
		}
	}

	state := SetupStateUninitialized
	if ownerExists {
		state = SetupStateInitialized
	}

	c.JSON(http.StatusOK, gin.H{
		"state": state,
		"steps": gin.H{
			"users_exist":         usersExist,
			"admin_active":        adminActive,
			"owner_exists":        ownerExists,
			"organisation_exists": organisationExists,
			"subscription_active": subscriptionActive,
			"education_entitled":  educationEntitled,
		},
		"next_action": nextInstallationAction(usersExist, adminActive, ownerExists, organisationExists, subscriptionActive),
	})
}

func unavailable(c *gin.Context) {
	c.JSON(http.StatusServiceUnavailable, gin.H{"message": "Setup progress is temporarily unavailable"})
}

// countExists reports whether at least one row matches the optional condition.
// A database error is returned so the caller never reports a step as done when
// it could not be verified.
func countExists(model interface{}, where string, args ...interface{}) (bool, error) {
	query := config.DB.Model(model)
	if strings.TrimSpace(where) != "" {
		query = query.Where(where, args...)
	}
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// nextInstallationAction is the single next step implied by the state facts.
// Each value maps to an existing application route (see the constants above).
func nextInstallationAction(usersExist, adminActive, ownerExists, organisationExists, subscriptionActive bool) string {
	switch {
	case !usersExist:
		return NextActionCommercialSetup
	case !ownerExists && !adminActive:
		return NextActionResumeActivation
	case !ownerExists || !organisationExists:
		return NextActionCompleteSetup
	case !subscriptionActive:
		return NextActionActivateSubscription
	default:
		return NextActionReady
	}
}

// licenceContainsAddon reports whether a licence's licensed_addons JSON list
// contains the given addon code (case-insensitive). Reuses the existing
// licensing representation — no parallel entitlement system.
func licenceContainsAddon(addonsJSON, code string) bool {
	var addons []string
	if err := json.Unmarshal([]byte(addonsJSON), &addons); err != nil {
		return false
	}
	for _, addon := range addons {
		if strings.EqualFold(addon, code) {
			return true
		}
	}
	return false
}
