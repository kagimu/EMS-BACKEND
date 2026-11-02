package services

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	commercialModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commercial/models"
	addonResponses "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/addons/api/responses"
	addonServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/addons/services"
	licensingModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/models"
	licensingServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/services"
	ownerResponses "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/owner/api/responses"
	eosModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	"gorm.io/gorm"
)

var ErrOrganisationNotFound = errors.New("organisation not found")

const nearingExpiryDays = 30

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

func effectiveStatus(licence licensingModels.Licence, now time.Time) string {
	switch licence.Status {
	case licensingModels.StatusSuspended:
		return licensingModels.EffectiveStatusSuspended
	case licensingModels.StatusRevoked:
		return licensingModels.EffectiveStatusRevoked
	}
	if now.Before(licence.ValidFrom) {
		return "SCHEDULED"
	}
	if now.After(licence.ValidTo.Add(time.Duration(licence.GraceDays) * 24 * time.Hour)) {
		return licensingModels.EffectiveStatusExpired
	}
	if now.After(licence.ValidTo) {
		return licensingModels.EffectiveStatusInGrace
	}
	return licensingModels.EffectiveStatusActive
}

func customerStatus(licences []licensingModels.Licence, hasOrders bool, now time.Time) string {
	if len(licences) == 0 {
		if hasOrders {
			return "ORDER_ONLY"
		}
		return "UNLICENSED"
	}
	latest := licences[0]
	for _, licence := range licences[1:] {
		if licence.CreatedAt.After(latest.CreatedAt) {
			latest = licence
		}
	}
	return effectiveStatus(latest, now)
}

func decodeAddons(raw string) []string {
	var codes []string
	if err := json.Unmarshal([]byte(raw), &codes); err != nil {
		return []string{}
	}
	return codes
}

func organisationNames(numbers []string) map[string]string {
	result := make(map[string]string, len(numbers))
	if len(numbers) == 0 {
		return result
	}
	var entities []eosModel.Entity
	if err := config.DB.Where("number IN ?", numbers).Find(&entities).Error; err != nil {
		return result
	}
	for _, entity := range entities {
		result[entity.Number] = entity.Name
	}
	return result
}

func ownerNames(organisationNumber string) []string {
	var rows []struct {
		FirstName  string
		SecondName string
		LastName   string
	}
	config.DB.Table("ems_stakeholders AS s").
		Select("u.first_name, u.second_name, u.last_name").
		Joins("JOIN users AS u ON u.number = s.user_number").
		Where("s.organisation_number = ? AND s.stakeholder_type = ? AND s.status = ? AND s.deleted_at IS NULL", organisationNumber, "owner", "ACTIVE").
		Scan(&rows)
	result := make([]string, 0, len(rows))
	for _, row := range rows {
		name := strings.TrimSpace(strings.Join([]string{row.FirstName, row.SecondName, row.LastName}, " "))
		if name != "" {
			result = append(result, name)
		}
	}
	return result
}

func licenceSummaries(licences []licensingModels.Licence, now time.Time) []ownerResponses.LicenceSummary {
	numbers := make([]string, 0, len(licences))
	for _, licence := range licences {
		numbers = append(numbers, licence.Number)
	}
	capacity := make(map[string]int64, len(numbers))
	if len(numbers) > 0 {
		var rows []struct {
			LicenceNumber string
			Count         int64
		}
		config.DB.Model(&licensingModels.Installation{}).
			Select("licence_number, COUNT(*) AS count").
			Where("licence_number IN ? AND status = ?", numbers, licensingModels.InstallationStatusActive).
			Group("licence_number").
			Scan(&rows)
		for _, row := range rows {
			capacity[row.LicenceNumber] = row.Count
		}
	}
	result := make([]ownerResponses.LicenceSummary, 0, len(licences))
	for _, licence := range licences {
		result = append(result, ownerResponses.LicenceSummary{
			Number:             licence.Number,
			OrganisationNumber: licence.OrganisationNumber,
			ValidFrom:          licence.ValidFrom,
			ValidTo:            licence.ValidTo,
			GraceDays:          licence.GraceDays,
			MaxInstallations:   licence.MaxInstallations,
			CapacityUsed:       capacity[licence.Number],
			LicensedAddons:     decodeAddons(licence.LicensedAddons),
			Status:             licence.Status,
			EffectiveStatus:    effectiveStatus(licence, now),
			RenewalRequestedAt: licence.RenewalRequestedAt,
			CreatedAt:          licence.CreatedAt,
		})
	}
	return result
}

func installationSummaries(installations []licensingModels.Installation, names map[string]string) []ownerResponses.InstallationSummary {
	result := make([]ownerResponses.InstallationSummary, 0, len(installations))
	for _, installation := range installations {
		summary := ownerResponses.InstallationSummary{
			Number:               installation.Number,
			InstallationIdentity: installation.InstallationIdentity,
			OrganisationNumber:   installation.OrganisationNumber,
			OrganisationName:     names[installation.OrganisationNumber],
			LicenceNumber:        installation.LicenceNumber,
			Status:               installation.Status,
			ActivatedAt:          installation.ActivatedAt,
			LastSeenAt:           installation.LastSeenAt,
			ReleasedAt:           installation.ReleasedAt,
		}
		if installation.ReleaseReason != nil {
			summary.ReleaseReason = *installation.ReleaseReason
		}
		result = append(result, summary)
	}
	return result
}

// Overview returns owner-dashboard aggregates using targeted database queries.
func Overview() (*ownerResponses.OverviewResponse, error) {
	now := time.Now()
	windowEnd := now.Add(nearingExpiryDays * 24 * time.Hour)
	response := &ownerResponses.OverviewResponse{}

	if err := config.DB.Model(&eosModel.Entity{}).Count(&response.Organisations).Error; err != nil {
		return nil, err
	}
	if err := config.DB.Model(&licensingModels.Licence{}).Distinct("organisation_number").Count(&response.LicensedOrganisations).Error; err != nil {
		return nil, err
	}
	if err := config.DB.Raw(`
		SELECT COUNT(DISTINCT org) FROM (
			SELECT organisation_number AS org FROM ems_licences WHERE deleted_at IS NULL AND organisation_number <> ''
			UNION
			SELECT organisation_number FROM ems_setup_orders WHERE deleted_at IS NULL AND organisation_number <> ''
			UNION
			SELECT organisation_number FROM ems_subscription_orders WHERE deleted_at IS NULL AND organisation_number <> ''
		) AS commercial`).Scan(&response.CommercialOrganisations).Error; err != nil {
		return nil, err
	}

	active := config.DB.Model(&licensingModels.Licence{}).Where("status = ? AND valid_from <= ? AND valid_to >= ?", licensingModels.StatusActive, now, now)
	if err := active.Count(&response.ActiveLicences).Error; err != nil {
		return nil, err
	}
	config.DB.Model(&licensingModels.Licence{}).Where("status = ? AND valid_from <= ? AND valid_to >= ? AND valid_to <= ?", licensingModels.StatusActive, now, now, windowEnd).Count(&response.NearingExpiry)
	config.DB.Model(&licensingModels.Licence{}).Where("status = ? AND valid_to < ?", licensingModels.StatusActive, now).Count(&response.ExpiredLicences)
	config.DB.Model(&licensingModels.Licence{}).Where("status = ?", licensingModels.StatusSuspended).Count(&response.SuspendedLicences)
	config.DB.Model(&licensingModels.Licence{}).Where("status = ?", licensingModels.StatusRevoked).Count(&response.RevokedLicences)
	config.DB.Model(&licensingModels.Licence{}).Where("renewal_requested_at IS NOT NULL").Count(&response.PendingRenewals)
	config.DB.Model(&licensingModels.Installation{}).Count(&response.TotalInstallations)
	config.DB.Model(&licensingModels.Installation{}).Where("status = ?", licensingModels.InstallationStatusActive).Count(&response.ActiveInstallations)

	events, _, err := licensingServices.ListLicenceEvents("", "", "", 8, 0)
	if err != nil {
		return nil, err
	}
	response.RecentEvents = events

	var attention []licensingModels.Licence
	config.DB.Where("status = ? AND (valid_to <= ? OR valid_to < ?)", licensingModels.StatusActive, windowEnd, now).
		Or("status = ?", licensingModels.StatusSuspended).
		Order("valid_to ASC").Limit(8).Find(&attention)
	attentionSummaries := licenceSummaries(attention, now)
	attentionNames := make([]string, 0, len(attention))
	for _, licence := range attention {
		attentionNames = append(attentionNames, licence.OrganisationNumber)
	}
	names := organisationNames(attentionNames)
	for index := range attentionSummaries {
		attentionSummaries[index].OrganisationName = names[attentionSummaries[index].OrganisationNumber]
	}
	response.AttentionLicences = attentionSummaries

	var recentInstallations []licensingModels.Installation
	config.DB.Order("activated_at DESC").Limit(5).Find(&recentInstallations)
	recentNames := make([]string, 0, len(recentInstallations))
	for _, installation := range recentInstallations {
		recentNames = append(recentNames, installation.OrganisationNumber)
	}
	response.RecentInstallations = installationSummaries(recentInstallations, organisationNames(recentNames))

	return response, nil
}

// ListOrganisations returns the paginated owner directory.
func ListOrganisations(search, status string, page, perPage int) (*ownerResponses.Page[ownerResponses.OrganisationSummary], error) {
	page, perPage = normalisePage(page, perPage)
	now := time.Now()
	windowEnd := now.Add(nearingExpiryDays * 24 * time.Hour)

	query := config.DB.Model(&eosModel.Entity{})
	search = strings.TrimSpace(search)
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("name LIKE ? OR email LIKE ? OR phone LIKE ?", like, like, like)
	}

	switch strings.ToLower(strings.TrimSpace(status)) {
	case "", "all":
	case "licensed":
		query = query.Where("number IN (?)", config.DB.Model(&licensingModels.Licence{}).Select("organisation_number").Where("status = ? AND valid_to >= ?", licensingModels.StatusActive, now))
	case "expiring":
		query = query.Where("number IN (?)", config.DB.Model(&licensingModels.Licence{}).Select("organisation_number").Where("status = ? AND valid_from <= ? AND valid_to BETWEEN ? AND ?", licensingModels.StatusActive, now, now, windowEnd))
	case "expired":
		query = query.Where("number IN (?)", config.DB.Model(&licensingModels.Licence{}).Select("organisation_number").Where("status = ? AND valid_to < ?", licensingModels.StatusActive, now))
	case "suspended":
		query = query.Where("number IN (?)", config.DB.Model(&licensingModels.Licence{}).Select("organisation_number").Where("status = ?", licensingModels.StatusSuspended))
	case "unlicensed":
		query = query.Where("number NOT IN (?)", config.DB.Model(&licensingModels.Licence{}).Select("organisation_number"))
	case "order_only":
		query = query.Where("number IN (?)", config.DB.Model(&commercialModel.SetupOrder{}).Select("organisation_number")).
			Where("number NOT IN (?)", config.DB.Model(&licensingModels.Licence{}).Select("organisation_number"))
	default:
		return nil, errors.New("invalid customer status")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	var entities []eosModel.Entity
	if err := query.Order("name ASC").Offset((page - 1) * perPage).Limit(perPage).Find(&entities).Error; err != nil {
		return nil, err
	}

	numbers := make([]string, 0, len(entities))
	for _, entity := range entities {
		numbers = append(numbers, entity.Number)
	}
	licencesByOrg := make(map[string][]licensingModels.Licence)
	installationCounts := make(map[string]int64)
	orderOrgs := make(map[string]bool)
	if len(numbers) > 0 {
		var licences []licensingModels.Licence
		config.DB.Where("organisation_number IN ?", numbers).Find(&licences)
		for _, licence := range licences {
			licencesByOrg[licence.OrganisationNumber] = append(licencesByOrg[licence.OrganisationNumber], licence)
		}
		var rows []struct {
			OrganisationNumber string
			Count              int64
		}
		config.DB.Model(&licensingModels.Installation{}).
			Select("organisation_number, COUNT(*) AS count").
			Where("organisation_number IN ?", numbers).
			Group("organisation_number").Scan(&rows)
		for _, row := range rows {
			installationCounts[row.OrganisationNumber] = row.Count
		}
		var orderRows []struct{ OrganisationNumber string }
		config.DB.Model(&commercialModel.SetupOrder{}).Select("DISTINCT organisation_number").Where("organisation_number IN ?", numbers).Scan(&orderRows)
		for _, row := range orderRows {
			orderOrgs[row.OrganisationNumber] = true
		}
	}

	items := make([]ownerResponses.OrganisationSummary, 0, len(entities))
	for _, entity := range entities {
		licences := licencesByOrg[entity.Number]
		var activeCount int64
		var latestExpiry *time.Time
		for _, licence := range licences {
			if licence.Status == licensingModels.StatusActive {
				activeCount++
			}
			if latestExpiry == nil || licence.ValidTo.After(*latestExpiry) {
				expiry := licence.ValidTo
				latestExpiry = &expiry
			}
		}
		items = append(items, ownerResponses.OrganisationSummary{
			Number:              entity.Number,
			Name:                entity.Name,
			Email:               entity.Email,
			Phone:               entity.Phone,
			CreatedAt:           entity.CreatedAt,
			CustomerStatus:      customerStatus(licences, orderOrgs[entity.Number], now),
			LicenceCount:        int64(len(licences)),
			ActiveLicenceCount:  activeCount,
			InstallationCount:   installationCounts[entity.Number],
			OwnerNames:          ownerNames(entity.Number),
			LatestLicenceExpiry: latestExpiry,
		})
	}

	return &ownerResponses.Page[ownerResponses.OrganisationSummary]{Items: items, Total: total, Page: page, PerPage: perPage}, nil
}

// OrganisationDetail returns a customer profile using existing records.
func OrganisationDetail(number string) (*ownerResponses.OrganisationDetail, error) {
	var entity eosModel.Entity
	if err := config.DB.Where("number = ?", number).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrganisationNotFound
		}
		return nil, err
	}
	licences, err := licensingServices.ListLicences(number, "")
	if err != nil {
		return nil, err
	}
	installations, err := licensingServices.ListInstallations("", number, "")
	if err != nil {
		return nil, err
	}
	events, _, err := licensingServices.ListLicenceEvents("", number, "", 25, 0)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	var licenceModels []licensingModels.Licence
	config.DB.Where("organisation_number = ?", number).Find(&licenceModels)
	var activeCount int64
	for _, licence := range licenceModels {
		if licence.Status == licensingModels.StatusActive {
			activeCount++
		}
	}
	var setupOrders, subscriptionOrders int64
	config.DB.Model(&commercialModel.SetupOrder{}).Where("organisation_number = ?", number).Count(&setupOrders)
	config.DB.Model(&commercialModel.SubscriptionOrder{}).Where("organisation_number = ?", number).Count(&subscriptionOrders)

	return &ownerResponses.OrganisationDetail{
		Organisation: ownerResponses.OrganisationSummary{
			Number:             entity.Number,
			Name:               entity.Name,
			Email:              entity.Email,
			Phone:              entity.Phone,
			CreatedAt:          entity.CreatedAt,
			CustomerStatus:     customerStatus(licenceModels, setupOrders+subscriptionOrders > 0, now),
			LicenceCount:       int64(len(licences)),
			ActiveLicenceCount: activeCount,
			InstallationCount:  int64(len(installations)),
			OwnerNames:         ownerNames(entity.Number),
		},
		Licences:           licences,
		Installations:      installations,
		Events:             events,
		SetupOrders:        setupOrders,
		SubscriptionOrders: subscriptionOrders,
	}, nil
}

// ListLicences returns a paginated, searchable owner licence list without keys.
func ListLicences(search, status string, page, perPage int) (*ownerResponses.Page[ownerResponses.LicenceSummary], error) {
	page, perPage = normalisePage(page, perPage)
	now := time.Now()
	windowEnd := now.Add(nearingExpiryDays * 24 * time.Hour)
	query := config.DB.Model(&licensingModels.Licence{})
	search = strings.TrimSpace(search)
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("number LIKE ? OR organisation_number LIKE ? OR organisation_number IN (?)",
			like, like, config.DB.Model(&eosModel.Entity{}).Select("number").Where("name LIKE ?", like))
	}
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "", "all":
	case "active":
		query = query.Where("status = ? AND valid_from <= ? AND valid_to >= ?", licensingModels.StatusActive, now, now)
	case "expiring":
		query = query.Where("status = ? AND valid_from <= ? AND valid_to BETWEEN ? AND ?", licensingModels.StatusActive, now, now, windowEnd)
	case "expired":
		query = query.Where("status = ? AND valid_to < ?", licensingModels.StatusActive, now)
	case "in_grace":
		query = query.Where("status = ? AND valid_to < ? AND DATE_ADD(valid_to, INTERVAL grace_days DAY) >= ?", licensingModels.StatusActive, now, now)
	case "suspended":
		query = query.Where("status = ?", licensingModels.StatusSuspended)
	case "revoked":
		query = query.Where("status = ?", licensingModels.StatusRevoked)
	case "scheduled":
		query = query.Where("status = ? AND valid_from > ?", licensingModels.StatusActive, now)
	default:
		return nil, errors.New("invalid licence status")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	var licences []licensingModels.Licence
	if err := query.Order("valid_to ASC").Offset((page - 1) * perPage).Limit(perPage).Find(&licences).Error; err != nil {
		return nil, err
	}
	items := licenceSummaries(licences, now)
	numbers := make([]string, 0, len(licences))
	for _, licence := range licences {
		numbers = append(numbers, licence.OrganisationNumber)
	}
	names := organisationNames(numbers)
	for index := range items {
		items[index].OrganisationName = names[items[index].OrganisationNumber]
	}
	return &ownerResponses.Page[ownerResponses.LicenceSummary]{Items: items, Total: total, Page: page, PerPage: perPage}, nil
}

// ListInstallations returns a paginated, searchable owner installation list.
func ListInstallations(search, status string, page, perPage int) (*ownerResponses.Page[ownerResponses.InstallationSummary], error) {
	page, perPage = normalisePage(page, perPage)
	query := config.DB.Model(&licensingModels.Installation{})
	search = strings.TrimSpace(search)
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("number LIKE ? OR installation_identity LIKE ? OR organisation_number LIKE ? OR licence_number LIKE ? OR organisation_number IN (?)",
			like, like, like, like, config.DB.Model(&eosModel.Entity{}).Select("number").Where("name LIKE ?", like))
	}
	status = strings.ToUpper(strings.TrimSpace(status))
	if status != "" && status != "ALL" {
		query = query.Where("status = ?", status)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	var installations []licensingModels.Installation
	if err := query.Order("activated_at DESC").Offset((page - 1) * perPage).Limit(perPage).Find(&installations).Error; err != nil {
		return nil, err
	}
	numbers := make([]string, 0, len(installations))
	for _, installation := range installations {
		numbers = append(numbers, installation.OrganisationNumber)
	}
	items := installationSummaries(installations, organisationNames(numbers))
	return &ownerResponses.Page[ownerResponses.InstallationSummary]{Items: items, Total: total, Page: page, PerPage: perPage}, nil
}

// LicenceDetail returns a licence and its event timeline without the key.
func LicenceDetail(number string) (*ownerResponses.LicenceDetail, error) {
	licence, err := licensingServices.GetLicence(number)
	if err != nil {
		return nil, err
	}
	events, _, err := licensingServices.ListLicenceEvents(number, "", "", 50, 0)
	if err != nil {
		return nil, err
	}
	return &ownerResponses.LicenceDetail{Licence: *licence, Events: events}, nil
}

// InstallationDetail returns an installation, its licence and event timeline.
func InstallationDetail(number string) (*ownerResponses.InstallationDetail, error) {
	installation, err := licensingServices.GetInstallation(number)
	if err != nil {
		return nil, err
	}
	licence, err := licensingServices.GetLicence(installation.LicenceNumber)
	if err != nil {
		return nil, err
	}
	events, _, err := licensingServices.ListLicenceEvents("", installation.OrganisationNumber, "", 50, 0)
	if err != nil {
		return nil, err
	}
	return &ownerResponses.InstallationDetail{Installation: *installation, Licence: *licence, Events: events}, nil
}

// Modules returns the addon catalogue for licence configuration.
func Modules() ([]addonResponses.ModuleResponse, error) {
	return addonServices.ListModules()
}
