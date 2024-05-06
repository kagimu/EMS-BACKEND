package services

import (
	"errors"
	classRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportTemplate/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportTemplate/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportTemplate/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportTemplate/repositories"
	uuid "github.com/satori/go.uuid"
)

func FetchAllReportTemplates(schoolNumber string, search string) ([]responses.ReportTemplateResponse, error) {

	var response []responses.ReportTemplateResponse

	rawResponses, errorResponse := repositories.QueryReportTemplates(schoolNumber, search, "")

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

func FetchAllClassReportTemplate(schoolNumber string, reportTemplateNumber string) ([]responses.ClassReportTemplateResponse, error) {

	var response []responses.ClassReportTemplateResponse

	allReportTemplates, err := repositories.QueryReportTemplates(schoolNumber, "", reportTemplateNumber)

	if err != nil {
		return nil, err
	}

	for _, reportTemplate := range allReportTemplates {

		classReportTemplateResponse, err := queryAndSetClassesReportTemplateResponse(reportTemplate)

		if err == nil {
			response = append(response, classReportTemplateResponse)
		}

	}

	return response, nil

}

func queryAndSetClassesReportTemplateResponse(reportTemplate models.ReportTemplate) (responses.ClassReportTemplateResponse, error) {

	var classReportTemplateResponse responses.ClassReportTemplateResponse

	rawResponses, errorResponse := repositories.QueryClassReportTemplates(reportTemplate.Number)

	if errorResponse != nil {
		return classReportTemplateResponse, errorResponse
	}

	classReportTemplateResponse.TemplateName = reportTemplate.Name

	var classesResponse []responses.ClassResponse

	for _, rawResponse := range rawResponses {

		campusClass := rawResponse.CampusClass

		user := rawResponse.User
		userNames := user.FirstName + " " + user.LastName

		resp := responses.ClassResponse{
			ClassName: campusClass.ClassName,
			AddedBy:   userNames,
			Number:    rawResponse.Number,
		}

		classesResponse = append(classesResponse, resp)

	}

	classReportTemplateResponse.Classes = classesResponse

	return classReportTemplateResponse, nil
}

func AddTemplateToClass(classReportTemplateRequest requests.ClassReportTemplateRequest, currentUserNumber string) (string, []string) {

	reportTemplate, err := repositories.QueryByReportTemplateNumber(classReportTemplateRequest.ReportTemplateNumber)

	var errors []string

	if err != nil {
		return "", append(errors, err.Error())
	}

	if reportTemplate.Number == "" {
		return "Invalid report template number", nil
	}

	selectedClasses := classReportTemplateRequest.CampusClasses

	for _, campusClassNumber := range selectedClasses {

		errorMessage := checkAndSaveClass(reportTemplate.Number, campusClassNumber, currentUserNumber)

		errors = append(errors, errorMessage)

	}

	if len(errors) == len(selectedClasses) {
		return "No class added to the template, try again", errors
	}

	return "", errors
}

func checkAndSaveClass(reportTemplateNumber string, campusClassNumber string, currentUserNumber string) string {

	campusClass, err := classRepositories.QueryCampusClassByNumber(campusClassNumber)

	if err != nil {
		return err.Error()
	}

	if campusClass.Number == "" {
		return "Invalid campus class number"
	}

	oldClassReportTemplate, err := repositories.QueryByCampusClassNumber(campusClass.Number)

	if err != nil {
		return err.Error()
	}

	if oldClassReportTemplate.Number == "" {
		return "This class already has a report template"
	}

	classReportTemplate := models.ClassReportTemplate{
		ReportTemplateNumber: reportTemplateNumber,
		UserNumber:           currentUserNumber,
		CampusClassNumber:    campusClass.Number,
		Number:               uuid.NewV4().String(),
	}

	err = repositories.SaveClassReportTemplate(classReportTemplate)

	if err != nil {
		return err.Error()
	}

	return ""
}

func RemoveTemplateFromClass(campusClassNumber string) error {

	classCampusReportTemplate, err := repositories.QueryByCampusClassNumber(campusClassNumber)

	if err != nil {
		return err
	}

	if classCampusReportTemplate.Number == "" {
		return errors.New("this class has no report templated added to it")
	}

	deletionError := repositories.DeleteCampusClassReportTemplate(classCampusReportTemplate)
	if deletionError != nil {
		return deletionError
	}

	return nil
}
