package responses

// ModuleResponse is the API representation of an addon/module.
//
// Type is resolved from the in-code catalogue because it is not persisted.
type ModuleResponse struct {
	Code      string   `json:"code"`
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	Version   string   `json:"version"`
	Status    string   `json:"status"`
	Enabled   bool     `json:"enabled"`
	DependsOn []string `json:"depends_on"`
}
