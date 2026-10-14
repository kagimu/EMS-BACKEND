package services

import (
	"errors"
	"strings"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/addons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/addons/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/addons/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/addons/repositories"
	"gorm.io/gorm"
)

// SyncCatalog ensures every catalogue entry has a persisted ems_modules row.
//
// It is idempotent and deterministic:
//   - a missing row is created from the catalogue defaults;
//   - an existing row keeps its Enabled and Status values, so synchronization
//     never unexpectedly disables an addon;
//   - descriptive fields (name, version, dependencies) are refreshed from the
//     catalogue, which is the source of truth for descriptions;
//   - repeated calls never create duplicates because Code is unique.
func SyncCatalog() error {
	tx := config.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	for _, def := range addons.Catalog {
		var existing models.Module
		err := tx.Where("code = ?", def.Code).First(&existing).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			status := constants.StatusInactive
			if def.DefaultEnabled {
				status = constants.StatusActive
			}

			module := models.Module{
				Code:      def.Code,
				Name:      def.Name,
				Version:   def.Version,
				DependsOn: strings.Join(def.DependsOn, ","),
				Enabled:   def.DefaultEnabled,
				Status:    status,
			}

			if err := tx.Create(&module).Error; err != nil {
				tx.Rollback()
				return err
			}

			continue
		}

		if err != nil {
			tx.Rollback()
			return err
		}

		// Existing row: refresh descriptive fields only.
		// Enabled and Status are intentionally left untouched.
		existing.Name = def.Name
		existing.Version = def.Version
		existing.DependsOn = strings.Join(def.DependsOn, ",")

		if err := tx.Save(&existing).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// ListModules synchronises the catalogue, then returns every module with its
// descriptive fields resolved from the catalogue (including Type).
func ListModules() ([]responses.ModuleResponse, error) {
	if err := SyncCatalog(); err != nil {
		return nil, err
	}

	modules, err := repositories.FindAllModules()
	if err != nil {
		return nil, err
	}

	result := make([]responses.ModuleResponse, 0, len(modules))
	for _, module := range modules {
		result = append(result, toResponse(module))
	}

	return result, nil
}

// FindModule returns a single module by code, synchronising the catalogue first.
// It returns gorm.ErrRecordNotFound for a code that is not in the catalogue.
func FindModule(code string) (*responses.ModuleResponse, error) {
	if _, ok := addons.FindModuleDef(code); !ok {
		return nil, gorm.ErrRecordNotFound
	}

	if err := SyncCatalog(); err != nil {
		return nil, err
	}

	module, err := repositories.FindModuleByCode(code)
	if err != nil {
		return nil, err
	}

	response := toResponse(*module)
	return &response, nil
}

// SetModuleEnabled records enable/disable state for an addon.
//
// Foundation A: this only persists registry state. It does not gate, disable,
// or otherwise change Education behaviour, and no licensing enforcement occurs.
func SetModuleEnabled(code string, enabled bool) (*responses.ModuleResponse, error) {
	if _, ok := addons.FindModuleDef(code); !ok {
		return nil, gorm.ErrRecordNotFound
	}

	if err := SyncCatalog(); err != nil {
		return nil, err
	}

	module, err := repositories.FindModuleByCode(code)
	if err != nil {
		return nil, err
	}

	module.Enabled = enabled
	if enabled {
		module.Status = constants.StatusActive
	} else {
		module.Status = constants.StatusInactive
	}

	if err := repositories.SaveModule(module); err != nil {
		return nil, err
	}

	response := toResponse(*module)
	return &response, nil
}

func toResponse(module models.Module) responses.ModuleResponse {
	var dependsOn []string
	if module.DependsOn != "" {
		dependsOn = strings.Split(module.DependsOn, ",")
	}

	moduleType := string(addons.ModuleTypeAddon)
	if def, ok := addons.FindModuleDef(module.Code); ok {
		moduleType = string(def.Type)
	}

	return responses.ModuleResponse{
		Code:      module.Code,
		Name:      module.Name,
		Type:      moduleType,
		Version:   module.Version,
		Status:    module.Status,
		Enabled:   module.Enabled,
		DependsOn: dependsOn,
	}
}
