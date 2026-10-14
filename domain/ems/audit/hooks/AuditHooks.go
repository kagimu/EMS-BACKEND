package hooks

import (
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	auditServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/audit/services"
	eosRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/repositories"
	"github.com/gin-gonic/gin"
)

// Record appends a best-effort audit fact for a business action (Foundation G
// §18/§6). It is additive and non-semantic: the caller has already executed the
// authoritative business operation, and a recording failure never changes the
// business result (Foundation G §3.5/§7).
//
// The actor is taken from the authenticated request context set by
// AuthorizeJWT — never from a request parameter or a client-supplied header.
// The organisation is the caller-supplied authoritative organisation, or,
// where the caller does not have one (for example a global registry action),
// it is resolved from the authenticated actor through the existing eos
// mechanism. IP and user-agent come from the request context.
func Record(c *gin.Context, organisationNumber, action, entityType, entityID string, newValues interface{}) {
	actorUserNumber := trustedActor(c)
	if organisationNumber == "" && actorUserNumber != "" {
		if active, err := eosRepositories.FindActiveUserEntityByUserNumber(actorUserNumber); err == nil {
			organisationNumber = active.EntityNumber
		}
	}

	auditServices.RecordFact(auditServices.AuditRecord{
		OrganisationNumber: organisationNumber,
		ActorUserNumber:    actorUserNumber,
		Action:             action,
		EntityType:         entityType,
		EntityID:           entityID,
		NewValues:          newValues,
		IP:                 c.ClientIP(),
		UserAgent:          c.Request.UserAgent(),
	})
}

// trustedActor returns the JWT-derived actor set by AuthorizeJWT. It never
// consults request parameters or client-controlled headers.
func trustedActor(c *gin.Context) string {
	value, exists := c.Get("user")
	if !exists {
		return ""
	}
	user, ok := value.(authModels.User)
	if !ok {
		return ""
	}
	return user.Number
}
