package repositories

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/authority/models"
	permissionsModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/permissions/models"
)

// FindAllByOrganisation returns authority rows of one organisation, optionally
// filtered by decision type and/or status.
func FindAllByOrganisation(organisationNumber, decisionType, status string) ([]models.DecisionAuthority, error) {
	var authorities []models.DecisionAuthority

	query := config.DB.Where("organisation_number = ?", organisationNumber)
	if decisionType != "" {
		query = query.Where("decision_type = ?", decisionType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Order("decision_type asc, approval_level asc").Find(&authorities).Error; err != nil {
		return nil, err
	}

	return authorities, nil
}

// FindByNumberAndOrganisation returns an authority row only when it belongs to
// the given organisation, which prevents cross-organisation (IDOR) disclosure.
func FindByNumberAndOrganisation(number, organisationNumber string) (*models.DecisionAuthority, error) {
	var authority models.DecisionAuthority
	if err := config.DB.
		Where("number = ? AND organisation_number = ?", number, organisationNumber).
		First(&authority).Error; err != nil {
		return nil, err
	}
	return &authority, nil
}

// FindActiveForResolution loads every ACTIVE, currently-valid authority row for
// one organisation and decision type in a SINGLE query. Unit-scope filtering is
// applied in memory by the service so the resolver never issues a query per row.
func FindActiveForResolution(organisationNumber, decisionType string, now interface{}) ([]models.DecisionAuthority, error) {
	var authorities []models.DecisionAuthority

	if err := config.DB.
		Where("organisation_number = ? AND decision_type = ? AND status = ?",
			organisationNumber, decisionType, constants.StatusActive).
		Where("valid_from IS NULL OR valid_from <= ?", now).
		Where("valid_to IS NULL OR valid_to >= ?", now).
		Order("approval_level asc").
		Find(&authorities).Error; err != nil {
		return nil, err
	}

	return authorities, nil
}

// RoleExistsInOrganisation reports whether a system role belongs to the organisation.
func RoleExistsInOrganisation(roleNumber, organisationNumber string) (bool, error) {
	var count int64
	if err := config.DB.Model(&permissionsModel.Role{}).
		Where("number = ? AND entity_number = ?", roleNumber, organisationNumber).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func Create(authority *models.DecisionAuthority) error {
	return config.DB.Create(authority).Error
}

func Save(authority *models.DecisionAuthority) error {
	return config.DB.Save(authority).Error
}
