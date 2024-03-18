package commons

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RouteNotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"message": "Route not found"})
}
