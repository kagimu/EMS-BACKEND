package services

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportTemplate/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportTemplate/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportTemplate/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportTemplate/repositories"
	uuid "github.com/satori/go.uuid"
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

func PostNewReportTemplate(newTemplateRequest requests.NewTemplateRequest, currentUserNumber string) (string, error) {

	reportTemplate, err := repositories.QueryByNameAndSchool(newTemplateRequest.Name, newTemplateRequest.SchoolNumber)

	if err != nil {
		return "", err
	}

	if reportTemplate.Number != "" {
		return "Report with name already exists", nil
	}

	reportTemplate = models.ReportTemplate{}
	reportTemplate.Name = newTemplateRequest.Name
	reportTemplate.SchoolNumber = newTemplateRequest.SchoolNumber
	reportTemplate.UserNumber = currentUserNumber
	reportTemplate.Number = uuid.NewV4().String()

	err = repositories.Save(reportTemplate)

	if err != nil {
		return "", err
	}

	return "", nil
}

func FetchAllClassReportTemplate(templateNumber string) ([]responses.ClassReportTemplateResponse, error) {

	var response []responses.ClassReportTemplateResponse

	rawResponses, errorResponse := repositories.QueryClassReportTemplates(templateNumber)

	if errorResponse != nil {
		return nil, errorResponse
	}

	for _, rawResponse := range rawResponses {

		reportTemplate := rawResponse.ReportTemplate

		campusClass := rawResponse.CampusClass

		user := rawResponse.User
		userNames := user.FirstName + " " + user.LastName

		resp := responses.ClassReportTemplateResponse{
			TemplateName: reportTemplate.Name,
			ClassName:    campusClass.ClassName,
			Number:       rawResponse.Number,
			AddedBy:      userNames,
		}

		response = append(response, resp)

	}

	return response, nil

}
