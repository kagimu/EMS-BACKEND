package services

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	supportResponses "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/support/api/responses"
	supportModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/support/models"
)

func ListStaffTickets(role, userNumber, search, status, category, priority, organisation string, page, perPage int) (*supportResponses.Page[supportResponses.TicketSummary], error) {
	assignedTo := ""
	if role == supportModels.SupportRoleAgent {
		assignedTo = userNumber
	}
	return ListTickets("", organisation, search, status, category, priority, assignedTo, page, perPage)
}

func GetStaffTicket(role, userNumber, number string) (*supportResponses.TicketDetail, error) {
	detail, err := GetTicket(number, "", true)
	if err != nil {
		return nil, err
	}
	if role == supportModels.SupportRoleAgent {
		if detail.Ticket.AssignedToNumber != userNumber {
			return nil, ErrTicketNotFound
		}
	}
	return detail, nil
}

func StaffOverview(role, userNumber string) (*supportResponses.Overview, error) {
	response := &supportResponses.Overview{}
	query := config.DB.Model(&supportModels.Ticket{})
	if role == supportModels.SupportRoleAgent {
		query = query.Where("assigned_to_number = ?", userNumber)
	}
	query.Count(&response.Total)
	countStatus := func(status string, target *int64) {
		q := config.DB.Model(&supportModels.Ticket{}).Where("status = ?", status)
		if role == supportModels.SupportRoleAgent {
			q = q.Where("assigned_to_number = ?", userNumber)
		}
		q.Count(target)
	}
	countStatus(supportModels.StatusOpen, &response.Open)
	countStatus(supportModels.StatusInProgress, &response.InProgress)
	countStatus(supportModels.StatusWaitingForCustomer, &response.WaitingForCustomer)
	countStatus(supportModels.StatusWaitingForSupport, &response.WaitingForSupport)
	countStatus(supportModels.StatusResolved, &response.Resolved)
	countStatus(supportModels.StatusClosed, &response.Closed)
	countPriority := func(priority string, target *int64) {
		q := config.DB.Model(&supportModels.Ticket{}).Where("priority = ?", priority)
		if role == supportModels.SupportRoleAgent {
			q = q.Where("assigned_to_number = ?", userNumber)
		}
		q.Count(target)
	}
	countPriority(supportModels.PriorityHigh, &response.HighPriority)
	countPriority(supportModels.PriorityCritical, &response.Critical)
	return response, nil
}
