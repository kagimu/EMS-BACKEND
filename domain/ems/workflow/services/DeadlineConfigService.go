package services

import (
	"errors"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/workflow/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/workflow/models"
	workflowRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/workflow/repositories"
)

// ErrInvalidDeadlineHours is returned for a deadline configuration outside the
// accepted range. Zero, negative and above-maximum values are rejected.
var ErrInvalidDeadlineHours = errors.New("default_deadline_hours must be between 1 and 8760")

// deadlineDefaultHours resolves the organisation's configured default deadline in
// hours, falling back to the built-in 7-day default when no configuration exists
// (or an invalid one was somehow stored). It reads configuration only; it never
// writes and never triggers background behaviour.
func deadlineDefaultHours(organisationNumber string) int {
	record, err := workflowRepositories.FindDeadlineConfig(organisationNumber)
	if err != nil || record == nil || record.DefaultDeadlineHours <= 0 {
		return builtInDefaultDeadlineHours
	}
	if record.DefaultDeadlineHours > maxDeadlineHours {
		return builtInDefaultDeadlineHours
	}
	return record.DefaultDeadlineHours
}

// GetDeadlineConfig returns the effective workflow deadline configuration for the
// organisation (the configured value, or the built-in default when unset).
func GetDeadlineConfig(organisationNumber string) *responses.DeadlineConfigResponse {
	return &responses.DeadlineConfigResponse{
		OrganisationNumber:   organisationNumber,
		DefaultDeadlineHours: deadlineDefaultHours(organisationNumber),
	}
}

// UpdateDeadlineConfig sets the organisation's default approval deadline. It
// returns the updated configuration and the previous effective value so the
// caller can record an authoritative audit fact (E4-D9). Changing the default
// never rewrites existing requests: their resolved deadline was persisted at
// creation.
func UpdateDeadlineConfig(organisationNumber string, hours int) (*responses.DeadlineConfigResponse, int, error) {
	if hours <= 0 || hours > maxDeadlineHours {
		return nil, 0, ErrInvalidDeadlineHours
	}

	previous := deadlineDefaultHours(organisationNumber)

	record := models.WorkflowDeadlineConfig{
		OrganisationNumber:   organisationNumber,
		DefaultDeadlineHours: hours,
	}
	if err := workflowRepositories.UpsertDeadlineConfig(&record); err != nil {
		return nil, 0, err
	}

	return &responses.DeadlineConfigResponse{
		OrganisationNumber:   organisationNumber,
		DefaultDeadlineHours: hours,
	}, previous, nil
}
