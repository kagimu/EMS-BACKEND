package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	auditServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/audit/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/workflow/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/workflow/services"
	"github.com/gin-gonic/gin"
)

// FetchDeadlineConfig returns the authenticated organisation's workflow deadline
// configuration. It is read-only and never writes.
func FetchDeadlineConfig(c *gin.Context) {
	_, organisationNumber, ok := approvalOrganisationContext(c)
	if !ok {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Workflow deadline configuration fetched",
		"data":    services.GetDeadlineConfig(organisationNumber),
	})
}

// UpdateDeadlineConfig sets the authenticated organisation's default approval
// deadline. It is an administrative workflow action under manage_ems_workflow
// and records an authoritative audit fact (E4-D9) with old and new values.
func UpdateDeadlineConfig(c *gin.Context) {
	actorNumber, organisationNumber, ok := approvalOrganisationContext(c)
	if !ok {
		return
	}

	var request requests.UpdateDeadlineConfigRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commons.Feedback("Invalid request body"))
		return
	}

	updated, previous, err := services.UpdateDeadlineConfig(organisationNumber, request.DefaultDeadlineHours)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	// Foundation G (G1/G2): best-effort authoritative audit fact. It never
	// changes the result, records no secrets, and ordinary reads record nothing.
	auditServices.Record(auditServices.AuditRecord{
		OrganisationNumber: organisationNumber,
		ActorUserNumber:    actorNumber,
		Action:             "workflow.deadline_config.update",
		EntityType:         "ems_workflow_deadline_config",
		EntityID:           organisationNumber,
		OldValues:          map[string]interface{}{"default_deadline_hours": previous},
		NewValues:          map[string]interface{}{"default_deadline_hours": updated.DefaultDeadlineHours},
		IP:                 c.ClientIP(),
		UserAgent:          c.Request.UserAgent(),
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "Workflow deadline configuration updated",
		"data":    updated,
	})
}
