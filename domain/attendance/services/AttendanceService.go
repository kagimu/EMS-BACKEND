package services

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/attendance/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/attendance/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/attendance/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/repository"
	campusModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	repositories3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/repositories"
	classModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	repositories4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/pagination"
	schoolModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/repositories"
	"github.com/gin-gonic/gin"
)

func RecordAttendance(attendance *models.Attendance) error {

	repo := repositories.NewAttendanceRepository(config.DB)
	_, err := repo.Save(attendance)
	if err != nil {
		return err
	}

	return nil
}

func FetchAttendance(ctxt *gin.Context) (pagination.PaginatedResponse, error) {
	repo := repositories.NewAttendanceRepository(config.DB)

	response := make([]responses.AttendanceResponse, 0)

	var paginatedResponse pagination.PaginatedResponse
	paginatedResponse.Content = response

	attendances, err := repo.AllAttendance(ctxt, &paginatedResponse)
	if err != nil {
		return pagination.PaginatedResponse{}, err
	}

	for _, attendance := range attendances {

		user, err := repository.GetUserByNumber(attendance.UserNumber)

		if err != nil {
			return pagination.PaginatedResponse{}, err
		}

		var school schoolModels.School
		if attendance.SchoolNumber != "" {
			school, err = repositories2.GetSchoolByNumber(attendance.SchoolNumber)
			if err != nil {
				return pagination.PaginatedResponse{}, err
			}
		}

		var campus campusModel.Campus

		if attendance.CampusNumber != "" {
			campus, err = repositories3.GetCampusByNumber(attendance.CampusNumber)
		}

		var campusClass classModel.CampusClass
		if attendance.CampusClassNumber != "" {
			campusClass, err = repositories4.GetCampusClassByNumber(attendance.CampusClassNumber)
		}

		response = append(response, responses.AttendanceResponse{
			Number:       attendance.Number,
			Date:         attendance.Date.String(),
			CheckInTime:  attendance.CheckInTime.String(),
			CheckOutTime: attendance.CheckOutTime.String(),
			Present:      attendance.Present,
			Class:        campusClass.ClassName,
			School:       school.Name,
			Campus:       campus.Name,
			Name:         user.FirstName + " " + user.LastName,
		})

		paginatedResponse.Content = response
	}

	return paginatedResponse, nil
}
