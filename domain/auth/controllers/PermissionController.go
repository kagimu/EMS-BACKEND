package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
)

func FetchPermissionGroups(c *gin.Context) {
	permissions := commons.PermissionGroups()
	c.JSON(http.StatusOK, gin.H{"data": permissions})

}
