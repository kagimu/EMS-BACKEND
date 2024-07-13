package services

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/responses"
)

func FindSchoolClasses(schoolNumber, levelNumber string) ([]responses.ClassResponse, error) {

	var classesResponse = make([]responses.ClassResponse, 0)

	classes, err := repositories.FindClassesBySchoolNumber(schoolNumber, levelNumber)
	if err != nil {
		return classesResponse, errors.New("getting classes error: " + err.Error())
	}

	for _, class := range classes {
		classesResponse = append(classesResponse, responses.ClassResponse{
			Number:       class.Number,
			Name:         class.Name,
			Abbreviation: class.Abbreviation,
		})
	}

	return classesResponse, nil
}
