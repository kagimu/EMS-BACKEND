package repositories

import (
	"errors"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"gorm.io/gorm"
)

func GetFeedingGroupStaffByGroupNumber(schoolNumber string, campusNumber string, groupNumber string) ([]models.FeedingGroupStaff, error) {
	db := config.DB
	var feedingGroupStaffs []models.FeedingGroupStaff

	if err := db.Preload("Campus").Preload("FeedingGroup").
		Preload("StaffDeployment").
		Preload("StaffDeployment.JobOffer").
		Preload("StaffDeployment.JobOffer.JobApplication").
		Preload("StaffDeployment.JobOffer.JobInterview").
		Preload("StaffDeployment.JobOffer.JobInterview.JobVacancy").
		Preload("StaffDeployment.JobOffer.JobInterview.JobVacancy.JobPosition").
		Preload("StaffDeployment.JobOffer.JobInterview.JobVacancy.Campus").
		Joins("JOIN feeding_groups ON feeding_group_staffs.feeding_group_number = feeding_groups.number").
		Scopes(filters.ByAnyField("feeding_groups.school_number", schoolNumber), filters.ByCampusNumber(campusNumber), filters.ByAnyField("feeding_group_number", groupNumber)).
		Find(&feedingGroupStaffs).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return feedingGroupStaffs, nil
		}

		return feedingGroupStaffs, err
	}

	return feedingGroupStaffs, nil
}

func GetFeedingGroupStaffByStaffNumber(staffNumber string) (models.FeedingGroupStaff, error) {
	db := config.DB
	var feedingGroupStaff models.FeedingGroupStaff

	if err := db.Where("staff_number = ?", staffNumber).First(&feedingGroupStaff).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return feedingGroupStaff, nil
		}
		return feedingGroupStaff, err
	}

	return feedingGroupStaff, nil
}

func SaveFeedingGroupStaff(fgs models.FeedingGroupStaff, txn ...*gorm.DB) error {
	var db *gorm.DB
	if len(txn) > 0 {
		db = txn[0]
	} else {
		db = config.DB
	}

	return db.Save(&fgs).Error
}

func GetFeedingGroupStaffByNumber(number string) (models.FeedingGroupStaff, error) {
	db := config.DB
	var feedingGroupStaff models.FeedingGroupStaff

	if err := db.Where("number = ?", number).First(&feedingGroupStaff).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return feedingGroupStaff, nil
		}
		return feedingGroupStaff, err
	}

	return feedingGroupStaff, nil
}

func DeleteFeedingGroupStaff(number, staffNumber string) error {
	db := config.DB
	query := db.Model(&models.FeedingGroupStaff{})
	if number != "" {
		query = query.Where("number = ?", number)
	}
	if staffNumber != "" {
		query = query.Where("staff_number = ?", staffNumber)
	}
	return query.Delete(&models.FeedingGroupStaff{}).Error
}
