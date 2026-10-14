package addons

import "testing"

func TestCatalogContainsCoreAndEducation(t *testing.T) {
	if _, ok := FindModuleDef("core"); !ok {
		t.Fatal("catalogue is missing the core entry")
	}
	if _, ok := FindModuleDef("education"); !ok {
		t.Fatal("catalogue is missing the education entry")
	}
}

func TestEducationIsAnAddonDependingOnCore(t *testing.T) {
	def, ok := FindModuleDef("education")
	if !ok {
		t.Fatal("education is not in the catalogue")
	}

	if def.Type != ModuleTypeAddon {
		t.Errorf("education type = %q, want %q", def.Type, ModuleTypeAddon)
	}
	if len(def.DependsOn) != 1 || def.DependsOn[0] != "core" {
		t.Errorf("education depends_on = %v, want [core]", def.DependsOn)
	}
}

func TestCatalogIsMinimalAndDeterministic(t *testing.T) {
	if len(Catalog) != 2 {
		t.Fatalf("catalogue size = %d, want 2 (core, education only)", len(Catalog))
	}
	if Catalog[0].Code != "core" {
		t.Errorf("catalogue[0] = %q, want core", Catalog[0].Code)
	}
	if Catalog[1].Code != "education" {
		t.Errorf("catalogue[1] = %q, want education", Catalog[1].Code)
	}
	if Catalog[0].Type != ModuleTypeCore {
		t.Errorf("core type = %q, want %q", Catalog[0].Type, ModuleTypeCore)
	}
}

func TestFindModuleDefUnknownCode(t *testing.T) {
	if _, ok := FindModuleDef("finance"); ok {
		t.Fatal("finance must not exist in Foundation A")
	}
}
