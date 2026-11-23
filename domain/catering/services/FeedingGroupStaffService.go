package services

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/repositories"
	staffMembersModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staffMembers/models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func CreateFeedingGroupStaff(staffNumber, campusNumber, feedingGroupNumber, status string, txn *gorm.DB) error {
	checkExists, err := repositories.GetFeedingGroupStaffByStaffNumber(staffNumber)
	if err != nil {
		return errors.New("Checking if staff already added to a group error: " + err.Error())
	}

	if checkExists.Number != "" {
		if checkExists.FeedingGroupNumber != feedingGroupNumber || checkExists.CampusNumber != campusNumber {
			checkExists.FeedingGroupNumber = feedingGroupNumber
			if campusNumber != "" {
				checkExists.CampusNumber = campusNumber
			}
			checkExists.Status = status
			if err := repositories.SaveFeedingGroupStaff(checkExists, txn); err != nil {
				return errors.New("failed to update feeding group staff")
			}
		}
		return nil
	}

	saveFeedingGroupStaff := models.FeedingGroupStaff{
		Number:             uuid.NewV4().String(),
		FeedingGroupNumber: feedingGroupNumber,
		Status:             status,
		StaffNumber:        staffNumber,
		CampusNumber:       campusNumber,
	}

	if err := repositories.SaveFeedingGroupStaff(saveFeedingGroupStaff, txn); err != nil {
		return errors.New("failed to create feeding group staff")
	}

	return nil
}

func UpdateFeedingGroupStaffByNumber(number, campusNumber, feedingGroupNumber, status string, txn *gorm.DB) error {
	fgs, err := repositories.GetFeedingGroupStaffByNumber(number)
	if err != nil {
		return errors.New("finding staff assignment error: " + err.Error())
	}
	if fgs.Number == "" {
		return errors.New("staff feeding group assignment not found")
	}

	if feedingGroupNumber != "" {
		fgs.FeedingGroupNumber = feedingGroupNumber
	}
	if campusNumber != "" {
		fgs.CampusNumber = campusNumber
	}
	if status != "" {
		fgs.Status = status
	}

	if err := repositories.SaveFeedingGroupStaff(fgs, txn); err != nil {
		return errors.New("failed to update staff feeding group assignment")
	}
	return nil
}

func UpdateFeedingGroupStaffByStaffNumber(staffNumber, campusNumber, feedingGroupNumber, status string, txn *gorm.DB) error {
	fgs, err := repositories.GetFeedingGroupStaffByStaffNumber(staffNumber)
	if err != nil {
		return errors.New("finding staff assignment error: " + err.Error())
	}
	if fgs.Number == "" {
		if feedingGroupNumber == "" {
			return errors.New("feeding group number is required to assign staff")
		}
		if status == "" {
			status = "ACTIVE"
		}
		return CreateFeedingGroupStaff(staffNumber, campusNumber, feedingGroupNumber, status, txn)
	}

	if feedingGroupNumber != "" {
		fgs.FeedingGroupNumber = feedingGroupNumber
	}
	if campusNumber != "" {
		fgs.CampusNumber = campusNumber
	}
	if status != "" {
		fgs.Status = status
	}

	if err := repositories.SaveFeedingGroupStaff(fgs, txn); err != nil {
		return errors.New("failed to update staff feeding group assignment")
	}
	return nil
}

func DeleteFeedingGroupStaff(number, staffNumber string) error {
	return repositories.DeleteFeedingGroupStaff(number, staffNumber)
}

func MealStaffs(schoolNumber string, campusNumber string, groupNumber string) ([]responses.FeedingStaffGroupResponse, error) {
	var response = make([]responses.FeedingStaffGroupResponse, 0)

	staffs, err := repositories.GetFeedingGroupStaffByGroupNumber(schoolNumber, campusNumber, groupNumber)
	if err != nil {
		return nil, errors.New("Getting group staffs error: " + err.Error())
	}

	for _, s := range staffs {
		res := createFeedingStaffResponse(s)
		response = append(response, res)
	}

	return response, nil
}

func createFeedingStaffResponse(s models.FeedingGroupStaff) responses.FeedingStaffGroupResponse {
	dep := s.StaffDeployment
	app := dep.JobOffer.JobApplication

	name := app.FirstName
	if app.LastName != "" {
		name += " " + app.LastName
	}
	if name == "" {
		name = "Staff Member"
	}

	gender := ""
	var personal staffMembersModels.StaffBioDataPersonal
	if err := config.DB.Where("job_offer_number = ?", dep.JobOfferNumber).First(&personal).Error; err == nil {
		gender = personal.Sex
	}
	if gender == "" && dep.UserNumber != "" {
		var user authModels.User
		if err := config.DB.Where("number = ?", dep.UserNumber).First(&user).Error; err == nil {
			gender = user.Gender
		}
	}
	if gender == "" && app.UserNumber != "" {
		var user authModels.User
		if err := config.DB.Where("number = ?", app.UserNumber).First(&user).Error; err == nil {
			gender = user.Gender
		}
	}
	if gender == "" {
		gender = "N/A"
	}

	position := ""
	if dep.JobOffer.JobInterview.JobVacancy.JobPosition.PostName != "" {
		position = dep.JobOffer.JobInterview.JobVacancy.JobPosition.PostName
	}

	campusName := s.Campus.Name
	if campusName == "" && dep.JobOffer.JobInterview.JobVacancy.Campus.Name != "" {
		campusName = dep.JobOffer.JobInterview.JobVacancy.Campus.Name
	}

	return responses.FeedingStaffGroupResponse{
		Number:             s.Number,
		StaffNumber:        s.StaffNumber,
		UserNumber:         s.StaffNumber,
		Name:               name,
		Gender:             gender,
		STN:                dep.StaffIdNumber,
		ADN:                "",
		Email:              app.Email,
		PhoneNumber:        app.PhoneNumber,
		JobPosition:        position,
		FeedingGroupNumber: s.FeedingGroupNumber,
		FeedingGroupName:   s.FeedingGroup.Name,
		CampusNumber:       s.CampusNumber,
		CampusName:         campusName,
		SchoolNumber:       s.FeedingGroup.SchoolNumber,
		Status:             s.Status,
	}
}
