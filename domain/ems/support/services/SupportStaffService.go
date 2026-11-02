package services

import (
	"errors"
	"strings"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	supportResponses "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/support/api/responses"
	supportModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/support/models"
)

var ErrSupportStaffInvalid = errors.New("invalid support staff request")

func StaffRoleFor(userNumber string) string {
	if commons.IsSoftwareOwner(userNumber) {
		return "OWNER"
	}
	var staff supportModels.Staff
	if err := config.DB.Where("user_number = ? AND status = ?", userNumber, supportModels.SupportStatusActive).First(&staff).Error; err != nil {
		return ""
	}
	return staff.Role
}

func ListSupportStaff() ([]supportResponses.StaffSummary, error) {
	var staff []supportModels.Staff
	if err := config.DB.Order("role ASC, created_at ASC").Find(&staff).Error; err != nil {
		return nil, err
	}
	result := make([]supportResponses.StaffSummary, 0, len(staff))
	for _, member := range staff {
		summary := supportResponses.StaffSummary{UserNumber: member.UserNumber, Role: member.Role, Status: member.Status}
		config.DB.Model(&supportModels.Ticket{}).
			Where("assigned_to_number = ? AND status NOT IN ?", member.UserNumber, []string{supportModels.StatusResolved, supportModels.StatusClosed}).
			Count(&summary.OpenTickets)
		config.DB.Model(&supportModels.Ticket{}).Where("assigned_to_number = ?", member.UserNumber).Count(&summary.AssignedTickets)
		config.DB.Model(&supportModels.Ticket{}).Where("assigned_to_number = ? AND priority IN ? AND status NOT IN ?", member.UserNumber, []string{supportModels.PriorityHigh, supportModels.PriorityCritical}, []string{supportModels.StatusResolved, supportModels.StatusClosed}).Count(&summary.HighPriority)
		config.DB.Model(&supportModels.Ticket{}).
			Where("assigned_to_number = ? AND sla_resolution_due IS NOT NULL AND sla_resolution_due < NOW() AND status NOT IN ?", member.UserNumber, []string{supportModels.StatusResolved, supportModels.StatusClosed}).
			Count(&summary.Overdue)
		result = append(result, summary)
	}
	return result, nil
}

func UpsertSupportStaff(actorNumber, userNumber, role string) (*supportModels.Staff, error) {
	userNumber = strings.TrimSpace(userNumber)
	role = strings.ToUpper(strings.TrimSpace(role))
	if userNumber == "" || !supportModels.IsSupportRole(role) {
		return nil, ErrSupportStaffInvalid
	}
	var userCount int64
	config.DB.Model(&authModel.User{}).Where("number = ?", userNumber).Count(&userCount)
	if userCount == 0 {
		return nil, ErrSupportStaffInvalid
	}
	var staff supportModels.Staff
	err := config.DB.Where("user_number = ?", userNumber).First(&staff).Error
	if err != nil {
		staff = supportModels.Staff{UserNumber: userNumber, Role: role, Status: supportModels.SupportStatusActive, AddedByNumber: actorNumber}
		if err := config.DB.Create(&staff).Error; err != nil {
			return nil, err
		}
		return &staff, nil
	}
	staff.Role = role
	staff.Status = supportModels.SupportStatusActive
	if err := config.DB.Save(&staff).Error; err != nil {
		return nil, err
	}
	return &staff, nil
}

func SetSupportStaffStatus(actorNumber, userNumber, status string) (*supportModels.Staff, error) {
	status = strings.ToUpper(strings.TrimSpace(status))
	if status != supportModels.SupportStatusActive && status != supportModels.SupportStatusInactive {
		return nil, ErrSupportStaffInvalid
	}
	var staff supportModels.Staff
	if err := config.DB.Where("user_number = ?", userNumber).First(&staff).Error; err != nil {
		return nil, ErrSupportStaffInvalid
	}
	staff.Status = status
	if err := config.DB.Save(&staff).Error; err != nil {
		return nil, err
	}
	return &staff, nil
}

func RemoveSupportStaff(userNumber string) error {
	var staff supportModels.Staff
	if err := config.DB.Where("user_number = ?", userNumber).First(&staff).Error; err != nil {
		return ErrSupportStaffInvalid
	}
	return config.DB.Delete(&staff).Error
}

func EligibleSupportStaff() ([]supportModels.Staff, error) {
	var staff []supportModels.Staff
	err := config.DB.Where("status = ?", supportModels.SupportStatusActive).Order("role ASC, created_at ASC").Find(&staff).Error
	return staff, err
}

func IsEligibleSupportStaff(userNumber string) bool {
	if commons.IsSoftwareOwner(userNumber) {
		return true
	}
	var count int64
	config.DB.Model(&supportModels.Staff{}).Where("user_number = ? AND status = ?", userNumber, supportModels.SupportStatusActive).Count(&count)
	return count > 0
}

func SupportWorkload() (*supportResponses.WorkloadSummary, error) {
	staff, err := ListSupportStaff()
	if err != nil {
		return nil, err
	}
	var unassigned int64
	config.DB.Model(&supportModels.Ticket{}).
		Where("assigned_to_number IS NULL AND status NOT IN ?", []string{supportModels.StatusResolved, supportModels.StatusClosed}).
		Count(&unassigned)
	return &supportResponses.WorkloadSummary{Staff: staff, Unassigned: unassigned, OverloadedThreshold: 10}, nil
}
