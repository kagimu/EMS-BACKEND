package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	conceptModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/concept/models"
	"gorm.io/gorm"
)

type ConceptRepository interface {
	FindConceptsBySubTopicNumber(subTopicNumber string) ([]conceptModels.Concept, error)
	UseTx(tx *gorm.DB) *gorm.DB
}

type conceptRepositoryImp struct {
	Db *gorm.DB
}

func (conceptRepository *conceptRepositoryImp) UseTx(tx *gorm.DB) *gorm.DB {
	conceptRepository.Db = tx
	return conceptRepository.Db
}

func (conceptRepository *conceptRepositoryImp) FindConceptsBySubTopicNumber(subTopicNumber string) ([]conceptModels.Concept, error) {
	db := conceptRepository.Db

	var concepts = make([]conceptModels.Concept, 0)

	if err := db.Where("sub_topic_number = ?", subTopicNumber).Find(&concepts).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return concepts, nil
		}

		return concepts, err
	}

	return concepts, nil
}

func NewConceptRepository() ConceptRepository {
	return &conceptRepositoryImp{
		Db: config.DB,
	}
}
