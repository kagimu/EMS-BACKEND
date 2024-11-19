package services

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	userModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	repository2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/repository"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/repositories"
	repositories3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/pagination"
	repositories4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/repositories"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/importCurriculum/models"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/repositories"
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

	curriculum, _ := models2.GetCurriculumByNumber(streamTeacher.CurriculumNumber)
	campus, _ := repositories.GetCampusByNumber(streamTeacher.CampusNumber)
	school, _ := repositories2.GetSchoolByNumber(streamTeacher.SchoolNumber)
	stream, _ := repositories3.GetStreamByNumber(streamTeacher.StreamNumber)
	class, _ := repositories3.GetClassByCampusClassNumber(streamTeacher.CampusClassNumber)
	assessmentArea, _ := repositories4.GetAssessmentAreaByNumber(streamTeacher.AssessmentAreaNumber)

	return responses.TeacherResponse{
		Number:               user.Number,
		FirstName:            user.FirstName,
		LastName:             user.LastName,
		SecondName:           user.SecondName,
		Gender:               user.Gender,
		Type:                 user.Type,
		Status:               user.Status,
		CurriculumNumber:     streamTeacher.CurriculumNumber,
		CurriculumName:       curriculum.Name,
		CampusNumber:         streamTeacher.CampusNumber,
		CampusName:           campus.Name,
		SchoolNumber:         streamTeacher.SchoolNumber,
		SchoolName:           school.Name,
		StreamNumber:         streamTeacher.StreamNumber,
		StreamName:           stream.Name,
		CampusClassNumber:    streamTeacher.CampusClassNumber,
		CampusClassName:      class.Name,
		AssessmentAreaNumber: streamTeacher.AssessmentAreaNumber,
		AssessmentAreaName:   assessmentArea.Name,
	}
}
