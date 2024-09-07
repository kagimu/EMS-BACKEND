package services

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/repositories"
)

type SubjectService interface {
	AllSubjects() ([]responses.SubjectResponse, error)
}

type subjectService struct {
}

func (subjectService *subjectService) AllSubjects() ([]responses.SubjectResponse, error) {
	var subjectsResponse = make([]responses.SubjectResponse, 0)

	subjects, err := repositories.AllSubjects()
	if err != nil {
		return subjectsResponse, errors.New("Error getting all subjects " + err.Error())
	}

	for _, subject := range subjects {
		subjectsResponse = append(subjectsResponse, responses.SubjectResponse{
			Name:          subject.Name,
			SubjectNumber: subject.Number,
		})
	}

	return subjectsResponse, nil
}

func NewSubjectService() SubjectService {
	return &subjectService{}
}
