package services

import (
	"errors"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	notificationServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/notifications/services"
	supportResponses "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/support/api/responses"
	supportModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/support/models"
)

var defaultSLATargets = map[string]supportModels.SLATarget{
	supportModels.PriorityCritical: {Priority: supportModels.PriorityCritical, FirstResponseMinutes: 60, ResolutionMinutes: 480},
	supportModels.PriorityHigh:     {Priority: supportModels.PriorityHigh, FirstResponseMinutes: 240, ResolutionMinutes: 1440},
	supportModels.PriorityMedium:   {Priority: supportModels.PriorityMedium, FirstResponseMinutes: 480, ResolutionMinutes: 4320},
	supportModels.PriorityLow:      {Priority: supportModels.PriorityLow, FirstResponseMinutes: 1440, ResolutionMinutes: 7200},
}

func SLATargetFor(priority string) supportModels.SLATarget {
	var target supportModels.SLATarget
	if err := config.DB.Where("priority = ?", priority).First(&target).Error; err == nil {
		return target
	}
	target = defaultSLATargets[priority]
	if target.Priority == "" {
		target = defaultSLATargets[supportModels.PriorityMedium]
	}
	config.DB.Create(&target)
	return target
}

func ListSLATargets() ([]supportResponses.SLATargetResponse, error) {
	result := make([]supportResponses.SLATargetResponse, 0, len(supportModels.Priorities))
	for _, priority := range supportModels.Priorities {
		target := SLATargetFor(priority)
		result = append(result, supportResponses.SLATargetResponse{
			Priority: target.Priority, FirstResponseMinutes: target.FirstResponseMinutes, ResolutionMinutes: target.ResolutionMinutes,
		})
	}
	return result, nil
}

func SaveSLATargets(actorNumber string, targets []supportResponses.SLATargetResponse) error {
	for _, target := range targets {
		if !supportModels.IsPriority(target.Priority) || target.FirstResponseMinutes < 1 || target.ResolutionMinutes < 1 {
			return errors.New("invalid SLA target")
		}
		var existing supportModels.SLATarget
		if err := config.DB.Where("priority = ?", target.Priority).First(&existing).Error; err != nil {
			config.DB.Create(&supportModels.SLATarget{Priority: target.Priority, FirstResponseMinutes: target.FirstResponseMinutes, ResolutionMinutes: target.ResolutionMinutes, UpdatedByNumber: actorNumber})
		} else {
			existing.FirstResponseMinutes = target.FirstResponseMinutes
			existing.ResolutionMinutes = target.ResolutionMinutes
			existing.UpdatedByNumber = actorNumber
			config.DB.Save(&existing)
		}
	}
	return nil
}

func applySLA(ticket *supportModels.Ticket) {
	target := SLATargetFor(ticket.Priority)
	first := ticket.CreatedAt.Add(time.Duration(target.FirstResponseMinutes) * time.Minute)
	resolution := ticket.CreatedAt.Add(time.Duration(target.ResolutionMinutes) * time.Minute)
	ticket.SLAFirstResponseDue = &first
	ticket.SLAResolutionDue = &resolution
}

func firstResponseOverdue(ticket supportModels.Ticket, now time.Time) bool {
	return ticket.FirstResponseAt == nil && ticket.SLAFirstResponseDue != nil && now.After(*ticket.SLAFirstResponseDue)
}

func resolutionOverdue(ticket supportModels.Ticket, now time.Time) bool {
	return ticket.ResolvedAt == nil && ticket.SLAResolutionDue != nil && now.After(*ticket.SLAResolutionDue)
}

// notifySupport records deduplicated delivery status and writes the existing
// in-app notification. No email is sent unless a provider is later configured.
func notifySupport(ticket supportModels.Ticket, recipient, notificationType, dedupeKey, title, body string) error {
	if recipient == "" {
		return nil
	}
	var existing supportModels.Notification
	if err := config.DB.Where("dedupe_key = ?", dedupeKey).First(&existing).Error; err == nil {
		return nil
	}
	record := supportModels.Notification{
		TicketNumber: ticket.Number, RecipientUserNumber: recipient, NotificationType: notificationType,
		DedupeKey: dedupeKey, DeliveryStatus: supportModels.DeliveryInApp, Attempts: 1,
	}
	if err := config.DB.Create(&record).Error; err != nil {
		return err
	}
	if err := notificationServices.Notify(notificationServices.NotificationInput{
		OrganisationNumber: ticket.OrganisationNumber, RecipientUserNumber: recipient,
		Type: notificationType, Title: title, Body: body,
		Payload: map[string]any{"ticket_number": ticket.Number},
	}); err != nil {
		config.DB.Model(&record).Updates(map[string]any{"delivery_status": supportModels.DeliveryFailed, "last_error": err.Error()})
		return err
	}
	return nil
}

func notifyManagers(ticket supportModels.Ticket, notificationType, dedupePrefix, title, body string) {
	var staff []supportModels.Staff
	config.DB.Where("role = ? AND status = ?", supportModels.SupportRoleManager, supportModels.SupportStatusActive).Find(&staff)
	for _, manager := range staff {
		_ = notifySupport(ticket, manager.UserNumber, notificationType, dedupePrefix+":"+manager.UserNumber, title, body)
	}
}

func SupportReport() (*supportResponses.SupportReport, error) {
	report := &supportResponses.SupportReport{
		ByPriority: map[string]int64{}, ByCategory: map[string]int64{}, ByAgent: map[string]int64{},
	}
	config.DB.Model(&supportModels.Ticket{}).Count(&report.TicketsReceived)
	config.DB.Model(&supportModels.Ticket{}).Where("status IN ?", []string{supportModels.StatusResolved, supportModels.StatusClosed}).Count(&report.TicketsResolved)
	config.DB.Model(&supportModels.Ticket{}).Where("status NOT IN ?", []string{supportModels.StatusResolved, supportModels.StatusClosed}).Count(&report.OpenTickets)
	config.DB.Model(&supportModels.Ticket{}).Where("sla_resolution_due IS NOT NULL AND sla_resolution_due < NOW() AND status NOT IN ?", []string{supportModels.StatusResolved, supportModels.StatusClosed}).Count(&report.OverdueTickets)

	var avgFirst, avgResolution float64
	config.DB.Model(&supportModels.Ticket{}).Where("first_response_at IS NOT NULL").
		Select("COALESCE(AVG(TIMESTAMPDIFF(MINUTE, created_at, first_response_at)), 0)").Scan(&avgFirst)
	config.DB.Model(&supportModels.Ticket{}).Where("resolved_at IS NOT NULL").
		Select("COALESCE(AVG(TIMESTAMPDIFF(MINUTE, created_at, resolved_at)), 0)").Scan(&avgResolution)
	report.AverageFirstResponseMinutes = avgFirst
	report.AverageResolutionMinutes = avgResolution

	for _, priority := range supportModels.Priorities {
		var count int64
		config.DB.Model(&supportModels.Ticket{}).Where("priority = ?", priority).Count(&count)
		report.ByPriority[priority] = count
	}
	for _, category := range supportModels.Categories {
		var count int64
		config.DB.Model(&supportModels.Ticket{}).Where("category = ?", category).Count(&count)
		report.ByCategory[category] = count
	}
	var agentRows []struct {
		AssignedToNumber string
		Count            int64
	}
	config.DB.Model(&supportModels.Ticket{}).Select("assigned_to_number, COUNT(*) AS count").Where("assigned_to_number IS NOT NULL").Group("assigned_to_number").Scan(&agentRows)
	for _, row := range agentRows {
		report.ByAgent[row.AssignedToNumber] = row.Count
	}
	return report, nil
}
