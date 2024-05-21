package services

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/statistics/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/statistics/repository"
	"github.com/gin-gonic/gin"
)

func FetchBasicStatistics(ctx *gin.Context) (responses.BasicStatistics, error) {
	return repository.FetchBasicStatistics(ctx)

}
