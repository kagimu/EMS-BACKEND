package services

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/previousSchoolResults/repository"
)

type PreviousSchoolResultsService interface {
	PreviousSchoolResults(applicationNumber string) (responses.PreviousSchoolResult, error)
}

type previousSchoolResultsRepositoryService struct {
}

var previousSchoolResultsRepository = repository.NewPreviousSchoolResultsRepository()

func (previousSchoolResultsService *previousSchoolResultsRepositoryService) PreviousSchoolResults(applicationNumber string) (responses.PreviousSchoolResult, error) {

	var previousSchoolResultsResponse responses.PreviousSchoolResult

	previousSchoolResults, err := previousSchoolResultsRepository.FetchSchoolResultsByApplicationNumber(applicationNumber)
	if err != nil {
		return responses.PreviousSchoolResult{}, errors.New("Getting school results error: " + err.Error())
	}

	var previousSchoolSubjects = make([]responses.PreviousSchoolSubject, 0)

	if len(previousSchoolResults) > 0 {
		previousSchoolResultsResponse.ReportFile = previousSchoolResults[0].Attachment
	}

	for _, previousSchoolResult := range previousSchoolResults {
		previousSchoolSubjects = append(previousSchoolSubjects, responses.PreviousSchoolSubject{
			Name:  previousSchoolResult.Subject,
			Score: previousSchoolResult.Score,
		})
	}

	previousSchoolResultsResponse.Subjects = previousSchoolSubjects

	return previousSchoolResultsResponse, nil
}

func NewPreviousSchoolResultsService() PreviousSchoolResultsService {
	return &previousSchoolResultsRepositoryService{}
}
