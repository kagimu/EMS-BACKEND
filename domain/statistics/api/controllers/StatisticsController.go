package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/statistics/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func FetchStatistics(c *gin.Context) {

	bs, err := services.FetchBasicStatistics(c)
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to fetch statistics"})
		return
	}
	c.JSON(http.StatusOK, bs)
}
