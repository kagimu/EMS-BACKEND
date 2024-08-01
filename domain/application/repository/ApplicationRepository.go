package repository

import (
	applicationModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/pagination"
	"github.com/gin-gonic/gin"
)

type ApplicationRepository interface {
	FindAllApplications() ([]applicationModels.Application, error)
	FindAllApplicationsPaginated(ctx *gin.Context) (pagination.PaginatedResponse, error)
}
