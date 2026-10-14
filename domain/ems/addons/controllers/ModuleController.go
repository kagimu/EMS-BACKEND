package controllers

import (
	"errors"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/addons/services"
	auditHooks "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/audit/hooks"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FetchModules returns the EMS addon/module catalogue together with persisted
// state. It lazily and idempotently synchronises the catalogue first, so no
// change to main.go is required.
func FetchModules(c *gin.Context) {
	modules, err := services.ListModules()
	if err != nil {
		c.JSON(http.StatusInternalServerError, commons.Feedback("Failed to fetch modules"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Modules fetched",
		"data":    modules,
	})
}

func FetchModule(c *gin.Context) {
	module, err := services.FindModule(c.Param("code"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, commons.Feedback("Module not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, commons.Feedback("Failed to fetch module"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Module fetched",
		"data":    module,
	})
}

func EnableModule(c *gin.Context) {
	setModuleEnabled(c, true)
}

func DisableModule(c *gin.Context) {
	setModuleEnabled(c, false)
}

func setModuleEnabled(c *gin.Context, enabled bool) {
	module, err := services.SetModuleEnabled(c.Param("code"), enabled)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, commons.Feedback("Module not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, commons.Feedback("Failed to update module"))
		return
	}

	// Foundation G (G2): best-effort audit fact. Never changes the result.
	action := "module.disable"
	if enabled {
		action = "module.enable"
	}
	auditHooks.Record(c, "", action, "ems_module", c.Param("code"), module)

	c.JSON(http.StatusOK, gin.H{
		"message": "Module updated",
		"data":    module,
	})
}
