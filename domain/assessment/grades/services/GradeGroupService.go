package services

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/grades/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/grades/repository"
)

func CreateGradeGroup(gradeGroup *models.GradeGroup) error {
	// Implement the logic to create a new grade group

	repo := repository.NewGradeGroupRepositoryImpl(config.DB)

	// find grade group by name

	gradeGroupByName, err := repo.FindGradeGroupByName(gradeGroup.Name, gradeGroup.SchoolNumber)
	if err != nil {
		return err

	}

	if gradeGroupByName.Name != "" {
		return errors.New("grade group with name already exists")

	}

	err = repo.Create(gradeGroup)
	if err != nil {
		return err
	}

	return nil

}

func FetchGradeGroups(schoolNumber string) ([]models.GradeGroup, error) {
	// Implement the logic to fetch all grade groups

	repo := repository.NewGradeGroupRepositoryImpl(config.DB)

	gradeGroups, err := repo.FindAllBySchoolNumber(schoolNumber)
	if err != nil {
		return nil, err
	}

	return gradeGroups, nil
}

func FetchGradeGroupByNumber(gradeGroupNumber string) (models.GradeGroup, error) {
	// Implement the logic to fetch a grade group by number

	repo := repository.NewGradeGroupRepositoryImpl(config.DB)

	gradeGroup, err := repo.FindGradeGroupByNumber(gradeGroupNumber)
	if err != nil {
		return models.GradeGroup{}, err
	}

	return gradeGroup, nil

}

func UpdateGradeGroup(gradeGroup *models.GradeGroup) error {
	// Implement the logic to update a grade group

	repo := repository.NewGradeGroupRepositoryImpl(config.DB)

	err := repo.Update(gradeGroup)
	if err != nil {
		return err
	}

	return nil
}

func DeleteGradeGroup(gradeGroupID string) error {
	// Implement the logic to delete a grade group

	repo := repository.NewGradeGroupRepositoryImpl(config.DB)

	gradeGroup, err := repo.FindGradeGroupByNumber(gradeGroupID)
	if err != nil {
		return err
	}

	err = repo.Delete(&gradeGroup)
	if err != nil {
		return err
	}

	return nil
}
