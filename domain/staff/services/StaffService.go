package services

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	userModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	repository2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/repository"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/pagination"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staff/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staff/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staff/repository"
	"github.com/gin-gonic/gin"
)

func AttachStaffToStream(streamTeacher *models.StreamTeacher) error {
	var repo = repository.NewStaffRepository(config.DB)
	return repo.AttachStaffToStream(streamTeacher)
}

func DetachStaffFromStream(staffNumber string, streamNumber string) error {
	var repo = repository.NewStaffRepository(config.DB)
	streamTeacher, err := repo.GetTeacherByNumberAndStreamNumber(staffNumber, streamNumber)
	if err != nil {
		return err
	}
	return repo.DetachStaffFromStream(streamTeacher)
}

func GetStaffStreamsAndClasses(ctx *gin.Context) (*pagination.PaginatedResponse, error) {
	var repo = repository.NewStaffRepository(config.DB)
	var pr pagination.PaginatedResponse
	staff, err := repo.GetTeacherStreamDeployments(ctx, &pr)
	if err != nil {
		return nil, err
	}
	cleanResponse := make([]responses.TeacherResponse, 0)
	for _, val := range staff {
		user, err := repository2.GetUserByNumber(val.UserNumber)
		if err != nil {
			return nil, err
		}
		cleanResponse = append(cleanResponse, BuildStaffResponse(user, val))
	}
	pr.Content = cleanResponse

	return &pr, nil
}

func BuildStaffResponse(user userModel.User, streamTeacher models.StreamTeacher) responses.TeacherResponse {

	return responses.TeacherResponse{
		Number:           user.Number,
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		SecondName:       user.SecondName,
		Gender:           user.Gender,
		Type:             user.Type,
		Status:           user.Status,
		CurriculumNumber: streamTeacher.CurriculumNumber,
		//CurriculumName:       streamTeacher.CurriculumName,
		CampusNumber: streamTeacher.CampusNumber,
		//CampusName:           streamTeacher.CampusName,
		SchoolNumber: streamTeacher.SchoolNumber,
		//SchoolName:           streamTeacher.SchoolName,
		StreamNumber: streamTeacher.StreamNumber,
		//StreamName:           streamTeacher.StreamName,
		CampusClassNumber: streamTeacher.CampusClassNumber,
		//CampusClassName:      streamTeacher.CampusClassName,
		AssessmentAreaNumber: streamTeacher.AssessmentAreaNumber,
		//AssessmentAreaName:   streamTeacher.AssessmentAreaName,
	}
}
