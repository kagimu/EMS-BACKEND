package services

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	userModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staff/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staff/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staff/repository"
)

var repo = repository.NewStaffRepository(config.DB)

func AttachStaffToStream(streamTeacher *models.StreamTeacher) error {
	return repo.AttachStaffToStream(streamTeacher)
}

func DetachStaffFromStream(staffNumber string, streamNumber string) error {
	return repo.DetachStaffFromStream(staffNumber, streamNumber)
}

func GetStaffOnStream(streamNumber string) ([]responses.StaffResponse, error) {
	staff, err := repo.GetStaffOnStream(streamNumber)
	if err != nil {
		return nil, err
	}
	streams := make([]string, 0)
	classes := make([]string, 0)
	var staffResponses = make([]responses.StaffResponse, 0)
	for _, v := range staff {
		//streams, classes := GetStaffStreamsAndClasses(v)
		staffResponses = append(staffResponses, BuildStaffResponse(v, streams, classes))
	}

	return staffResponses, nil
}

func BuildStaffResponse(user userModel.User, streams, classes []string) responses.StaffResponse {
	return responses.StaffResponse{
		Number:  user.Number,
		Streams: streams,
		Classes: classes,
		UserResponse: responses2.UserResponse{
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			SecondName: user.SecondName,
			Gender:     user.Gender,
			Type:       user.Type,
			Status:     user.Status,
		},
	}
}
