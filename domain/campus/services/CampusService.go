package services

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/api/responses"
	campusModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
)

func CreateCampusResponse(campus campusModel.Campus) responses.CampusResponse {
	return responses.CampusResponse{
		Name:             campus.Name,
		Phone:            campus.Phone,
		Email:            campus.Email,
		Location:         campus.Location,
		SchoolNumber:     campus.SchoolNumber,
		CommencementDate: campus.CommencementDate,
	}
}
