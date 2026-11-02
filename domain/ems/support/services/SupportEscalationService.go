package services

import (
	"errors"
	"strings"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	supportResponses "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/support/api/responses"
	supportModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/support/models"
	"gorm.io/gorm"
)

var ErrEscalationReason = errors.New("an escalation reason is required")

func EscalateTicket(actorNumber, number, reason string) (*supportResponses.TicketDetail, error) {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return nil, ErrEscalationReason
	}
	var ticket supportModels.Ticket
	if err := config.DB.Where("number = ?", number).First(&ticket).Error; err != nil {
		return nil, ErrTicketNotFound
	}
	if ticket.Status == supportModels.StatusClosed {
		return nil, ErrTicketTransition
	}
	before := ticket
	now := time.Now()
	ticket.Escalated = true
	ticket.EscalatedAt = &now
	ticket.EscalatedBy = &actorNumber
	ticket.EscalationReason = reason
	if err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&ticket).Error; err != nil {
			return err
		}
		return recordAudit(tx, ticket.Number, actorNumber, "TICKET_ESCALATED", before, ticket, reason)
	}); err != nil {
		return nil, err
	}
	notifyManagers(ticket, "support_escalated", "ticket:escalated:"+ticket.Number, "Support ticket escalated", reason)
	if ticket.AssignedToNumber != nil {
		_ = notifySupport(ticket, *ticket.AssignedToNumber, "support_escalated", "ticket:escalated:"+ticket.Number+":assignee", "Support ticket escalated", reason)
	}
	return GetTicket(ticket.Number, "", true)
}
