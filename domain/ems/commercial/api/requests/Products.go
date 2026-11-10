package requests

type CommercialProductRequest struct {
	ProductCode  string `json:"product_code"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	ProductType  string `json:"product_type"`
	DisplayOrder int    `json:"display_order"`
	Reason       string `json:"reason"`
}

type CommercialProductStatusRequest struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}

type CommercialComponentRequest struct {
	ComponentCode string `json:"component_code"`
	IsIncluded    bool   `json:"is_included"`
	IsOptional    bool   `json:"is_optional"`
	DisplayOrder  int    `json:"display_order"`
	Reason        string `json:"reason"`
}

type CommercialComponentStatusRequest struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}
