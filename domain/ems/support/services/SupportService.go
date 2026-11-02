package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	licensingResponses "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/api/responses"
	licensingModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/models"
	licensingServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/services"
	supportRequests "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/support/api/requests"
	supportResponses "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/support/api/responses"
	supportModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/support/models"
	eosModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

var (
	ErrTicketNotFound      = errors.New("support ticket not found")
	ErrTicketValidation    = errors.New("invalid support ticket")
	ErrTicketTransition    = errors.New("invalid ticket status transition")
	ErrTicketResolution    = errors.New("a resolution summary is required when resolving a ticket")
	ErrTicketMessage       = errors.New("invalid ticket message")
	ErrTicketAuthorisation = errors.New("you cannot access this support ticket")
)

func normalisePage(page, perPage int) (int, int) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}
	return page, perPage
}

func validateOrganisation(number string) error {
	var count int64
	if err := config.DB.Model(&eosModel.Entity{}).Where("number = ?", number).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return errors.New("organisation not found")
	}
	return nil
}

func validateUser(number string) error {
	var count int64
	if err := config.DB.Model(&authModel.User{}).Where("number = ?", number).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return errors.New("user not found")
	}
	return nil
}

func validateTicketRelation(organisation, licenceNumber, installationNumber string) error {
	if licenceNumber != "" {
		var licence licensingModels.Licence
		if err := config.DB.Where("number = ? AND organisation_number = ?", licenceNumber, organisation).First(&licence).Error; err != nil {
			return errors.New("licence does not belong to this organisation")
		}
	}
	if installationNumber != "" {
		var installation licensingModels.Installation
		if err := config.DB.Where("number = ? AND organisation_number = ?", installationNumber, organisation).First(&installation).Error; err != nil {
			return errors.New("installation does not belong to this organisation")
		}
	}
	return nil
}

func toSummary(ticket supportModels.Ticket) supportResponses.TicketSummary {
	summary := supportResponses.TicketSummary{
		Number:              ticket.Number,
		OrganisationNumber:  ticket.OrganisationNumber,
		RequestedByNumber:   ticket.RequestedByNumber,
		Subject:             ticket.Subject,
		Category:            ticket.Category,
		Priority:            ticket.Priority,
		Status:              ticket.Status,
		ResolvedAt:          ticket.ResolvedAt,
		CreatedAt:           ticket.CreatedAt,
		UpdatedAt:           ticket.UpdatedAt,
		FirstResponseAt:     ticket.FirstResponseAt,
		SLAFirstResponseDue: ticket.SLAFirstResponseDue,
		SLAResolutionDue:    ticket.SLAResolutionDue,
		Escalated:           ticket.Escalated,
		EscalatedAt:         ticket.EscalatedAt,
		EscalationReason:    ticket.EscalationReason,
	}
	now := time.Now()
	summary.FirstResponseOverdue = firstResponseOverdue(ticket, now)
	summary.ResolutionOverdue = resolutionOverdue(ticket, now)
	if ticket.AssignedToNumber != nil {
		summary.AssignedToNumber = *ticket.AssignedToNumber
	}
	if ticket.LicenceNumber != nil {
		summary.LicenceNumber = *ticket.LicenceNumber
	}
	if ticket.InstallationNumber != nil {
		summary.InstallationNumber = *ticket.InstallationNumber
	}
	return summary
}

func toMessage(message supportModels.Message) supportResponses.MessageResponse {
	return supportResponses.MessageResponse{
		ID:             message.ID,
		TicketNumber:   message.TicketNumber,
		AuthorNumber:   message.AuthorNumber,
		AuthorRole:     message.AuthorRole,
		Visibility:     message.Visibility,
		Body:           message.Body,
		DeliveryStatus: message.DeliveryStatus,
		CreatedAt:      message.CreatedAt,
	}
}

func toAudit(audit supportModels.Audit) supportResponses.AuditResponse {
	return supportResponses.AuditResponse{
		ID: audit.ID, TicketNumber: audit.TicketNumber, ActorNumber: audit.ActorNumber,
		Action: audit.Action, BeforeState: audit.BeforeState, AfterState: audit.AfterState,
		Details: audit.Details, CreatedAt: audit.CreatedAt,
	}
}

func encodeState(value any) string {
	encoded, _ := json.Marshal(value)
	if string(encoded) == "null" {
		return ""
	}
	return string(encoded)
}

func recordAudit(tx *gorm.DB, ticket, actor, action string, before, after any, details string) error {
	return tx.Create(&supportModels.Audit{
		TicketNumber: ticket,
		ActorNumber:  actor,
		Action:       action,
		BeforeState:  encodeState(before),
		AfterState:   encodeState(after),
		Details:      details,
	}).Error
}

func orgNames(numbers []string) map[string]string {
	result := map[string]string{}
	if len(numbers) == 0 {
		return result
	}
	var entities []eosModel.Entity
	config.DB.Where("number IN ?", numbers).Find(&entities)
	for _, entity := range entities {
		result[entity.Number] = entity.Name
	}
	return result
}

// CreateTicket creates an owner or client ticket. Client callers pass the
// server-derived active organisation and actor; owner callers may target another
// organisation and requester.
func CreateTicket(actorNumber, organisationNumber string, request supportRequests.CreateTicketRequest, owner bool) (*supportResponses.TicketDetail, error) {
	organisationNumber = strings.TrimSpace(organisationNumber)
	if organisationNumber == "" {
		return nil, ErrTicketValidation
	}
	if err := validateOrganisation(organisationNumber); err != nil {
		return nil, ErrTicketValidation
	}

	requestedBy := strings.TrimSpace(request.RequestedByNumber)
	explicitRequester := owner && requestedBy != ""
	if !owner || requestedBy == "" {
		requestedBy = actorNumber
	}
	if err := validateUser(requestedBy); err != nil {
		return nil, ErrTicketValidation
	}
	if explicitRequester {
		var membershipCount int64
		if err := config.DB.Model(&eosModel.EntityUser{}).
			Where("user_number = ? AND entity_number = ?", requestedBy, organisationNumber).
			Count(&membershipCount).Error; err != nil {
			return nil, err
		}
		if membershipCount == 0 {
			return nil, ErrTicketValidation
		}
	}

	subject := strings.TrimSpace(request.Subject)
	description := strings.TrimSpace(request.Description)
	if subject == "" || description == "" {
		return nil, ErrTicketValidation
	}
	category, err := supportModels.NormaliseCategory(request.Category)
	if err != nil {
		return nil, ErrTicketValidation
	}
	priority := supportModels.PriorityMedium
	if strings.TrimSpace(request.Priority) != "" {
		priority, err = supportModels.NormalisePriority(request.Priority)
		if err != nil {
			return nil, ErrTicketValidation
		}
	}
	if err := validateTicketRelation(organisationNumber, strings.TrimSpace(request.LicenceNumber), strings.TrimSpace(request.InstallationNumber)); err != nil {
		return nil, ErrTicketValidation
	}

	ticket := supportModels.Ticket{
		Number: uuid.NewV4().String(), OrganisationNumber: organisationNumber,
		RequestedByNumber: requestedBy, Subject: subject, Description: description,
		Category: category, Priority: priority, Status: supportModels.StatusOpen,
	}
	if value := strings.TrimSpace(request.LicenceNumber); value != "" {
		ticket.LicenceNumber = &value
	}
	if value := strings.TrimSpace(request.InstallationNumber); value != "" {
		ticket.InstallationNumber = &value
	}
	applySLA(&ticket)

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&ticket).Error; err != nil {
			return err
		}
		return recordAudit(tx, ticket.Number, actorNumber, "TICKET_CREATED", nil, ticket, "Ticket created")
	})
	if err != nil {
		return nil, err
	}
	_ = notifySupport(ticket, ticket.RequestedByNumber, "support_ticket_received", "ticket:created:"+ticket.Number+":"+ticket.RequestedByNumber, "Support request received", ticket.Subject)
	if ticket.Priority == supportModels.PriorityHigh || ticket.Priority == supportModels.PriorityCritical {
		notifyManagers(ticket, "support_high_priority", "ticket:high:"+ticket.Number, "High-priority support ticket", ticket.Subject)
	}
	return GetTicket(ticket.Number, organisationNumber, true)
}

// ListTickets returns a paginated list. organisationFilter is empty for the
// owner and forced to the active organisation for clients.
func ListTickets(organisationFilter, organisation, search, status, category, priority, assignedTo string, page, perPage int) (*supportResponses.Page[supportResponses.TicketSummary], error) {
	page, perPage = normalisePage(page, perPage)
	query := config.DB.Model(&supportModels.Ticket{})
	if organisationFilter != "" {
		query = query.Where("organisation_number = ?", organisationFilter)
	}
	if organisation != "" {
		query = query.Where("organisation_number = ?", organisation)
	}
	if search = strings.TrimSpace(search); search != "" {
		like := "%" + search + "%"
		query = query.Where("number LIKE ? OR subject LIKE ? OR description LIKE ?", like, like, like)
	}
	if status = strings.ToUpper(strings.TrimSpace(status)); status != "" && status != "ALL" {
		query = query.Where("status = ?", status)
	}
	if category = strings.ToUpper(strings.TrimSpace(category)); category != "" {
		query = query.Where("category = ?", category)
	}
	if priority = strings.ToUpper(strings.TrimSpace(priority)); priority != "" {
		query = query.Where("priority = ?", priority)
	}
	if assignedTo = strings.TrimSpace(assignedTo); assignedTo != "" {
		query = query.Where("assigned_to_number = ?", assignedTo)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	var tickets []supportModels.Ticket
	if err := query.Order("updated_at DESC").Offset((page - 1) * perPage).Limit(perPage).Find(&tickets).Error; err != nil {
		return nil, err
	}
	numbers := make([]string, 0, len(tickets))
	for _, ticket := range tickets {
		numbers = append(numbers, ticket.OrganisationNumber)
	}
	names := orgNames(numbers)
	items := make([]supportResponses.TicketSummary, 0, len(tickets))
	for _, ticket := range tickets {
		summary := toSummary(ticket)
		summary.OrganisationName = names[ticket.OrganisationNumber]
		items = append(items, summary)
	}
	return &supportResponses.Page[supportResponses.TicketSummary]{Items: items, Total: total, Page: page, PerPage: perPage}, nil
}

// GetTicket returns one ticket. Client callers pass their active organisation
// and cannot see internal messages or audit records.
func GetTicket(number, organisationFilter string, owner bool) (*supportResponses.TicketDetail, error) {
	var ticket supportModels.Ticket
	if err := config.DB.Where("number = ?", number).First(&ticket).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTicketNotFound
		}
		return nil, err
	}
	if organisationFilter != "" && ticket.OrganisationNumber != organisationFilter {
		return nil, ErrTicketNotFound
	}

	messageQuery := config.DB.Where("ticket_number = ?", ticket.Number)
	if !owner {
		messageQuery = messageQuery.Where("visibility = ?", supportModels.VisibilityCustomer)
	}
	var messages []supportModels.Message
	if err := messageQuery.Order("created_at ASC").Find(&messages).Error; err != nil {
		return nil, err
	}
	messageResponses := make([]supportResponses.MessageResponse, 0, len(messages))
	for _, message := range messages {
		messageResponses = append(messageResponses, toMessage(message))
	}

	detail := &supportResponses.TicketDetail{
		Ticket:      toSummary(ticket),
		Description: ticket.Description,
		Resolution:  ticket.ResolutionSummary,
		Messages:    messageResponses,
	}
	detail.Attachments, _ = ListSupportAttachments(ticket.Number, owner)
	detail.Ticket.OrganisationName = orgNames([]string{ticket.OrganisationNumber})[ticket.OrganisationNumber]

	if licenceNumber := ticket.LicenceNumber; licenceNumber != nil {
		if licence, err := licensingServices.GetLicence(*licenceNumber); err == nil {
			detail.Technical.LicenceStatus = licence.EffectiveStatus
			detail.Technical.LicenceValidFrom = &licence.ValidFrom
			detail.Technical.LicenceValidTo = &licence.ValidTo
			detail.Technical.LicenceAddons = licence.LicensedAddons
			if owner {
				events, _, _ := licensingServices.ListLicenceEvents(*licenceNumber, "", "", 50, 0)
				detail.Events = events
			}
		}
	}
	if installationNumber := ticket.InstallationNumber; installationNumber != nil {
		if installation, err := licensingServices.GetInstallation(*installationNumber); err == nil && installation.OrganisationNumber == ticket.OrganisationNumber {
			detail.Technical.InstallationStatus = installation.Status
			detail.Technical.InstallationActivatedAt = &installation.ActivatedAt
			detail.Technical.InstallationLastSeenAt = installation.LastSeenAt
		}
	}
	if owner {
		var audits []supportModels.Audit
		config.DB.Where("ticket_number = ?", ticket.Number).Order("created_at ASC").Find(&audits)
		for _, audit := range audits {
			detail.Audit = append(detail.Audit, toAudit(audit))
		}
	}
	return detail, nil
}

// AddMessage appends a customer-visible or internal conversation item.
func AddMessage(actorNumber, number, organisationFilter, authorRole, visibility, body string) (*supportResponses.MessageResponse, error) {
	var ticket supportModels.Ticket
	if err := config.DB.Where("number = ?", number).First(&ticket).Error; err != nil {
		return nil, ErrTicketNotFound
	}
	if organisationFilter != "" && ticket.OrganisationNumber != organisationFilter {
		return nil, ErrTicketNotFound
	}
	if ticket.Status == supportModels.StatusClosed {
		return nil, ErrTicketTransition
	}

	body = strings.TrimSpace(body)
	if body == "" {
		return nil, ErrTicketMessage
	}
	visibility = strings.ToUpper(strings.TrimSpace(visibility))
	if visibility == "" {
		visibility = supportModels.VisibilityCustomer
	}
	if visibility != supportModels.VisibilityCustomer && visibility != supportModels.VisibilityInternal {
		return nil, ErrTicketMessage
	}
	if authorRole == supportModels.AuthorCustomer && visibility == supportModels.VisibilityInternal {
		return nil, ErrTicketAuthorisation
	}

	message := supportModels.Message{
		TicketNumber: ticket.Number, AuthorNumber: actorNumber, AuthorRole: authorRole,
		Visibility: visibility, Body: body, DeliveryStatus: "IN_APP",
	}
	oldStatus := ticket.Status
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&message).Error; err != nil {
			return err
		}
		action := "MESSAGE_CUSTOMER"
		if visibility == supportModels.VisibilityInternal {
			action = "MESSAGE_INTERNAL"
		}
		if err := recordAudit(tx, ticket.Number, actorNumber, action, nil, message, ""); err != nil {
			return err
		}
		if authorRole != supportModels.AuthorCustomer && visibility == supportModels.VisibilityCustomer && ticket.FirstResponseAt == nil {
			now := time.Now()
			ticket.FirstResponseAt = &now
			if err := tx.Model(&supportModels.Ticket{}).Where("number = ?", ticket.Number).Update("first_response_at", now).Error; err != nil {
				return err
			}
		}
		if authorRole == supportModels.AuthorCustomer && ticket.Status == supportModels.StatusWaitingForCustomer {
			ticket.Status = supportModels.StatusWaitingForSupport
			if err := tx.Model(&supportModels.Ticket{}).Where("number = ?", ticket.Number).Update("status", ticket.Status).Error; err != nil {
				return err
			}
			return recordAudit(tx, ticket.Number, actorNumber, "TICKET_STATUS_CHANGED", oldStatus, ticket.Status, "Customer replied")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	dedupe := "ticket:message:" + fmt.Sprint(message.ID)
	if authorRole == supportModels.AuthorCustomer {
		if ticket.AssignedToNumber != nil {
			_ = notifySupport(ticket, *ticket.AssignedToNumber, "support_customer_reply", dedupe+":assignee", "Customer replied", ticket.Subject)
		}
		notifyManagers(ticket, "support_customer_reply", dedupe+":manager", "Customer replied", ticket.Subject)
	} else if visibility == supportModels.VisibilityCustomer && ticket.RequestedByNumber != "" {
		_ = notifySupport(ticket, ticket.RequestedByNumber, "support_reply", dedupe+":customer", "Support updated your ticket", ticket.Subject)
	}

	response := toMessage(message)
	return &response, nil
}

// UpdateTicket performs owner-only assignment, priority, status and resolution
// changes with audit records.
func UpdateTicket(actorNumber, number string, request supportRequests.UpdateTicketRequest) (*supportResponses.TicketDetail, error) {
	var ticket supportModels.Ticket
	if err := config.DB.Where("number = ?", number).First(&ticket).Error; err != nil {
		return nil, ErrTicketNotFound
	}
	before := ticket
	now := time.Now()

	if request.Status != nil {
		status, err := supportModels.NormaliseStatus(*request.Status)
		if err != nil {
			return nil, ErrTicketValidation
		}
		if err := supportModels.ValidateTransition(ticket.Status, status); err != nil {
			return nil, ErrTicketTransition
		}
		if status == ticket.Status {
			// no-op
		} else {
			ticket.Status = status
			switch status {
			case supportModels.StatusResolved:
				if request.ResolutionSummary != nil {
					ticket.ResolutionSummary = strings.TrimSpace(*request.ResolutionSummary)
				}
				if strings.TrimSpace(ticket.ResolutionSummary) == "" {
					return nil, ErrTicketResolution
				}
				ticket.ResolvedAt = &now
			case supportModels.StatusOpen:
				ticket.ResolvedAt = nil
			case supportModels.StatusClosed:
				if ticket.ResolvedAt == nil {
					ticket.ResolvedAt = &now
				}
			}
		}
	}
	if request.Priority != nil {
		priority, err := supportModels.NormalisePriority(*request.Priority)
		if err != nil {
			return nil, ErrTicketValidation
		}
		ticket.Priority = priority
	}
	if request.Category != nil {
		category, err := supportModels.NormaliseCategory(*request.Category)
		if err != nil {
			return nil, ErrTicketValidation
		}
		ticket.Category = category
	}
	if request.ResolutionSummary != nil {
		ticket.ResolutionSummary = strings.TrimSpace(*request.ResolutionSummary)
	}
	if request.AssignedToNumber != nil {
		assignee := strings.TrimSpace(*request.AssignedToNumber)
		if assignee != "" {
			if err := validateUser(assignee); err != nil {
				return nil, ErrTicketValidation
			}
			if !IsEligibleSupportStaff(assignee) {
				return nil, ErrTicketValidation
			}
			ticket.AssignedToNumber = &assignee
		} else {
			ticket.AssignedToNumber = nil
		}
	}

	action := "TICKET_UPDATED"
	if before.Status != ticket.Status {
		if ticket.Status == supportModels.StatusResolved {
			action = "TICKET_RESOLVED"
		} else if before.Status == supportModels.StatusResolved && ticket.Status == supportModels.StatusOpen {
			action = "TICKET_REOPENED"
		} else if ticket.Status == supportModels.StatusClosed {
			action = "TICKET_CLOSED"
		} else {
			action = "TICKET_STATUS_CHANGED"
		}
	} else if before.AssignedToNumber == nil && ticket.AssignedToNumber != nil {
		action = "TICKET_ASSIGNED"
	} else if before.AssignedToNumber != nil && ticket.AssignedToNumber != nil && *before.AssignedToNumber != *ticket.AssignedToNumber {
		action = "TICKET_REASSIGNED"
	} else if before.Priority != ticket.Priority {
		action = "TICKET_PRIORITY_CHANGED"
	} else if before.Category != ticket.Category {
		action = "TICKET_CATEGORY_CHANGED"
	}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&ticket).Error; err != nil {
			return err
		}
		if err := recordAudit(tx, ticket.Number, actorNumber, action, before, ticket, strings.TrimSpace(request.AssignmentReason)); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if before.AssignedToNumber == nil && ticket.AssignedToNumber != nil {
		_ = notifySupport(ticket, *ticket.AssignedToNumber, "support_assigned", "ticket:assigned:"+ticket.Number+":"+*ticket.AssignedToNumber, "Support ticket assigned", ticket.Subject)
	} else if before.AssignedToNumber != nil && ticket.AssignedToNumber != nil && *before.AssignedToNumber != *ticket.AssignedToNumber {
		_ = notifySupport(ticket, *ticket.AssignedToNumber, "support_reassigned", "ticket:reassigned:"+ticket.Number+":"+*ticket.AssignedToNumber, "Support ticket reassigned", ticket.Subject)
	}
	if ticket.RequestedByNumber != "" && before.Status != ticket.Status {
		_ = notifySupport(ticket, ticket.RequestedByNumber, "support_status", "ticket:status:"+ticket.Number+":"+ticket.Status, "Support ticket status updated", ticket.Subject)
	}
	return GetTicket(ticket.Number, "", true)
}

// Overview returns real support counts.
func Overview() (*supportResponses.Overview, error) {
	response := &supportResponses.Overview{}
	config.DB.Model(&supportModels.Ticket{}).Count(&response.Total)
	config.DB.Model(&supportModels.Ticket{}).Where("status = ?", supportModels.StatusOpen).Count(&response.Open)
	config.DB.Model(&supportModels.Ticket{}).Where("status = ?", supportModels.StatusInProgress).Count(&response.InProgress)
	config.DB.Model(&supportModels.Ticket{}).Where("status = ?", supportModels.StatusWaitingForCustomer).Count(&response.WaitingForCustomer)
	config.DB.Model(&supportModels.Ticket{}).Where("status = ?", supportModels.StatusWaitingForSupport).Count(&response.WaitingForSupport)
	config.DB.Model(&supportModels.Ticket{}).Where("priority = ?", supportModels.PriorityHigh).Count(&response.HighPriority)
	config.DB.Model(&supportModels.Ticket{}).Where("priority = ?", supportModels.PriorityCritical).Count(&response.Critical)
	config.DB.Model(&supportModels.Ticket{}).Where("status = ?", supportModels.StatusResolved).Count(&response.Resolved)
	config.DB.Model(&supportModels.Ticket{}).Where("status = ?", supportModels.StatusClosed).Count(&response.Closed)
	return response, nil
}

// Events returns licence events for owner technical investigation.
func Events(licenceNumber string) ([]licensingResponses.LicenceEventResponse, error) {
	events, _, err := licensingServices.ListLicenceEvents(licenceNumber, "", "", 100, 0)
	return events, err
}
