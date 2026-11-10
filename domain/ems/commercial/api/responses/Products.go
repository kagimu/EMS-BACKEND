package responses

import "time"

type CommercialProductSummary struct {
	Number         string    `json:"number"`
	ProductCode    string    `json:"product_code"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	ProductType    string    `json:"product_type"`
	Status         string    `json:"status"`
	DisplayOrder   int       `json:"display_order"`
	ComponentCount int64     `json:"component_count"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CommercialComponentResponse struct {
	Number        string    `json:"number"`
	ProductNumber string    `json:"product_number"`
	ComponentCode string    `json:"component_code"`
	ComponentName string    `json:"component_name"`
	IsIncluded    bool      `json:"is_included"`
	IsOptional    bool      `json:"is_optional"`
	Status        string    `json:"status"`
	DisplayOrder  int       `json:"display_order"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CommercialProductDetail struct {
	Product    CommercialProductSummary      `json:"product"`
	Components []CommercialComponentResponse `json:"components"`
}

type CommercialProductPage struct {
	Items   []CommercialProductSummary `json:"items"`
	Total   int64                      `json:"total"`
	Page    int                        `json:"page"`
	PerPage int                        `json:"per_page"`
}
