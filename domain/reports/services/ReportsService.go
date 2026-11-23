package services

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reports/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reports/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reports/repositories"
	"github.com/gin-gonic/gin"
)

func FetchSummary(ctx *gin.Context, filters requests.ReportsFilterRequest) (responses.SummaryResponse, error) {
	return repositories.FetchSummary(ctx, filters)
}

func FetchSchoolStructureReport(ctx *gin.Context, filters requests.ReportsFilterRequest) (responses.SchoolStructureReportResponse, error) {
	return repositories.FetchSchoolStructureReport(ctx, filters)
}
