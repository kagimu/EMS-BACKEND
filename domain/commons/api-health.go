package commons

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func APIHealth(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"message": "Systems are up and running"})
}
