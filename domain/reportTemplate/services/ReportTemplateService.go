package services

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportTemplate/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportTemplate/repositories"
)

func FetchAllReportTemplates(schoolNumber string, search string) ([]responses.ReportTemplateResponse, error) {

	var response []responses.ReportTemplateResponse

	rawResponses, errorResponse := repositories.QueryReportTemplates(schoolNumber, search)

	if errorResponse != nil {

		return nil, errorResponse
	}

	for _, rawResponse := range rawResponses {

		user := rawResponse.User

		userNames := user.FirstName + " " + user.LastName

		resp := responses.ReportTemplateResponse{
			Name:             rawResponse.Name,
			ReportTemplateId: rawResponse.Number,
			Status:           rawResponse.Status,
			AddedBy:          userNames,
		}

		response = append(response, resp)

	}

	return response, nil
}
