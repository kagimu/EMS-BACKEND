// Package addons defines the EMS addon/module catalogue.
//
// EMS Core must never depend on an Education domain. This package is generic:
// it describes addons (including Education) without importing any of them, so
// new addons can be declared here without changing the Core architecture.
package addons

// ModuleType distinguishes the platform itself from installed addons.
type ModuleType string

const (
	ModuleTypeCore  ModuleType = "core"
	ModuleTypeAddon ModuleType = "addon"
)

// ModuleDef is the in-code description of an addon.
//
// The catalogue is the source of truth for descriptive fields; the ems_modules
// table stores state (enabled/status). Type is descriptive only and is not
// persisted, to keep the table minimal.
type ModuleDef struct {
	Code           string
	Name           string
	Type           ModuleType
	Version        string
	DependsOn      []string
	DefaultEnabled bool
}

// Catalog lists every addon EMS Core knows about.
//
// Keep this list minimal and deterministic. "core" is the platform itself;
// "education" is the first and currently only real addon.
//
// Foundation A: do NOT add Finance, HR, Inventory, Healthcare or Construction.
var Catalog = []ModuleDef{
	{
		Code:           "core",
		Name:           "EMS Core",
		Type:           ModuleTypeCore,
		Version:        "1.0.0",
		DefaultEnabled: true,
	},
	{
		Code:           "education",
		Name:           "Education",
		Type:           ModuleTypeAddon,
		Version:        "1.0.0",
		DependsOn:      []string{"core"},
		DefaultEnabled: true,
	},
}

// FindModuleDef returns the catalogue entry for a code.
func FindModuleDef(code string) (ModuleDef, bool) {
	for _, def := range Catalog {
		if def.Code == code {
			return def, true
		}
	}
	return ModuleDef{}, false
}
