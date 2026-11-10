package services

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/addons"
	commercialRequests "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/commercial/api/requests"
	commercialResponses "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/commercial/api/responses"
	commercialModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/commercial/models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

var (
	ErrProductNotFound     = errors.New("commercial product not found")
	ErrProductExists       = errors.New("commercial product code already exists")
	ErrProductValidation   = errors.New("invalid commercial product")
	ErrProductTransition   = errors.New("invalid commercial product status transition")
	ErrComponentNotFound   = errors.New("commercial product component not found")
	ErrComponentExists     = errors.New("component already belongs to this product")
	ErrComponentValidation = errors.New("invalid commercial product component")
)

var productCodePattern = regexp.MustCompile(`^[A-Z0-9][A-Z0-9-]{1,58}[A-Z0-9]$`)

func normalisePage(page, perPage int) (int, int) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 50
	}
	if perPage > 100 {
		perPage = 100
	}
	return page, perPage
}

func toProductSummary(product commercialModels.CommercialProduct, count int64) commercialResponses.CommercialProductSummary {
	return commercialResponses.CommercialProductSummary{
		Number: product.Number, ProductCode: product.ProductCode, Name: product.Name, Description: product.Description,
		ProductType: product.ProductType, Status: product.Status, DisplayOrder: product.DisplayOrder,
		ComponentCount: count, CreatedAt: product.CreatedAt, UpdatedAt: product.UpdatedAt,
	}
}

func toComponentResponse(component commercialModels.CommercialProductComponent) commercialResponses.CommercialComponentResponse {
	return commercialResponses.CommercialComponentResponse{
		Number: component.Number, ProductNumber: component.ProductNumber, ComponentCode: component.ComponentCode,
		ComponentName: component.ComponentName, IsIncluded: component.IsIncluded, IsOptional: component.IsOptional,
		Status: component.Status, DisplayOrder: component.DisplayOrder, CreatedAt: component.CreatedAt, UpdatedAt: component.UpdatedAt,
	}
}

func recordProductEvent(tx *gorm.DB, entityType, entityNumber, action, before, after, reason, actor string) error {
	return tx.Create(&commercialModels.CommercialEvent{
		Number: uuid.NewV4().String(), EntityType: entityType, EntityNumber: entityNumber,
		OrganisationNumber: "", Action: action, BeforeState: before, AfterState: after,
		Reason: strings.TrimSpace(reason), ActorUserNumber: actor,
	}).Error
}

func productState(product commercialModels.CommercialProduct) string {
	encoded, _ := json.Marshal(toProductSummary(product, 0))
	return string(encoded)
}

func componentState(component commercialModels.CommercialProductComponent) string {
	encoded, _ := json.Marshal(toComponentResponse(component))
	return string(encoded)
}

func ListProducts(search, status string, page, perPage int) (*commercialResponses.CommercialProductPage, error) {
	page, perPage = normalisePage(page, perPage)
	query := config.DB.Model(&commercialModels.CommercialProduct{})
	search = strings.TrimSpace(search)
	if search != "" {
		like := "%" + strings.ToUpper(search) + "%"
		query = query.Where("product_code LIKE ? OR name LIKE ?", like, "%"+search+"%")
	}
	if status = strings.ToUpper(strings.TrimSpace(status)); status != "" && status != "ALL" {
		if _, err := commercialModels.NormaliseCommercialProductStatus(status); err != nil {
			return nil, err
		}
		query = query.Where("status = ?", status)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	var products []commercialModels.CommercialProduct
	if err := query.Order("display_order ASC, name ASC").Offset((page - 1) * perPage).Limit(perPage).Find(&products).Error; err != nil {
		return nil, err
	}
	numbers := make([]string, 0, len(products))
	for _, product := range products {
		numbers = append(numbers, product.Number)
	}
	counts := map[string]int64{}
	if len(numbers) > 0 {
		var rows []struct {
			ProductNumber string
			Count         int64
		}
		config.DB.Model(&commercialModels.CommercialProductComponent{}).Select("product_number, COUNT(*) AS count").Where("product_number IN ?", numbers).Group("product_number").Scan(&rows)
		for _, row := range rows {
			counts[row.ProductNumber] = row.Count
		}
	}
	items := make([]commercialResponses.CommercialProductSummary, 0, len(products))
	for _, product := range products {
		items = append(items, toProductSummary(product, counts[product.Number]))
	}
	return &commercialResponses.CommercialProductPage{Items: items, Total: total, Page: page, PerPage: perPage}, nil
}

func GetProduct(number string) (*commercialResponses.CommercialProductDetail, error) {
	var product commercialModels.CommercialProduct
	if err := config.DB.Where("number = ?", number).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	var components []commercialModels.CommercialProductComponent
	config.DB.Where("product_number = ?", number).Order("display_order ASC, component_code ASC").Find(&components)
	items := make([]commercialResponses.CommercialComponentResponse, 0, len(components))
	for _, component := range components {
		items = append(items, toComponentResponse(component))
	}
	return &commercialResponses.CommercialProductDetail{Product: toProductSummary(product, int64(len(components))), Components: items}, nil
}

func CreateProduct(actor string, request commercialRequests.CommercialProductRequest) (*commercialResponses.CommercialProductDetail, error) {
	code := strings.ToUpper(strings.TrimSpace(request.ProductCode))
	name := strings.TrimSpace(request.Name)
	productType, err := commercialModels.NormaliseCommercialProductType(request.ProductType)
	if err != nil {
		return nil, err
	}
	if code == "" || !productCodePattern.MatchString(code) || name == "" {
		return nil, ErrProductValidation
	}
	var count int64
	config.DB.Model(&commercialModels.CommercialProduct{}).Where("product_code = ?", code).Count(&count)
	if count > 0 {
		return nil, ErrProductExists
	}
	product := commercialModels.CommercialProduct{
		Number: uuid.NewV4().String(), ProductCode: code, Name: name, Description: strings.TrimSpace(request.Description),
		ProductType: productType, Status: commercialModels.CommercialProductStatusDraft, DisplayOrder: request.DisplayOrder,
		CreatedByNumber: actor, UpdatedByNumber: actor,
	}
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&product).Error; err != nil {
			var exists int64
			tx.Model(&commercialModels.CommercialProduct{}).Where("product_code = ?", code).Count(&exists)
			if exists > 0 {
				return ErrProductExists
			}
			return err
		}
		return recordProductEvent(tx, "COMMERCIAL_PRODUCT", product.Number, "COMMERCIAL_PRODUCT_CREATED", "", productState(product), request.Reason, actor)
	})
	if err != nil {
		return nil, err
	}
	return GetProduct(product.Number)
}

func UpdateProduct(actor, number string, request commercialRequests.CommercialProductRequest) (*commercialResponses.CommercialProductDetail, error) {
	var product commercialModels.CommercialProduct
	if err := config.DB.Where("number = ?", number).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	if product.Status == commercialModels.CommercialProductStatusInactive {
		return nil, ErrProductTransition
	}
	name := strings.TrimSpace(request.Name)
	productType, err := commercialModels.NormaliseCommercialProductType(request.ProductType)
	if err != nil {
		return nil, err
	}
	if name == "" {
		return nil, ErrProductValidation
	}
	before := productState(product)
	product.Name = name
	product.Description = strings.TrimSpace(request.Description)
	product.ProductType = productType
	product.DisplayOrder = request.DisplayOrder
	product.UpdatedByNumber = actor
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&product).Error; err != nil {
			return err
		}
		return recordProductEvent(tx, "COMMERCIAL_PRODUCT", product.Number, "COMMERCIAL_PRODUCT_UPDATED", before, productState(product), request.Reason, actor)
	})
	if err != nil {
		return nil, err
	}
	return GetProduct(product.Number)
}

func SetProductStatus(actor, number string, request commercialRequests.CommercialProductStatusRequest) (*commercialResponses.CommercialProductDetail, error) {
	var product commercialModels.CommercialProduct
	if err := config.DB.Where("number = ?", number).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	status, err := commercialModels.NormaliseCommercialProductStatus(request.Status)
	if err != nil {
		return nil, err
	}
	if err := commercialModels.ValidateProductStatusTransition(product.Status, status); err != nil {
		return nil, ErrProductTransition
	}
	if product.Status == status {
		return GetProduct(number)
	}
	before := productState(product)
	product.Status = status
	product.UpdatedByNumber = actor
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&product).Error; err != nil {
			return err
		}
		return recordProductEvent(tx, "COMMERCIAL_PRODUCT", product.Number, "COMMERCIAL_PRODUCT_STATUS_CHANGED", before, productState(product), request.Reason, actor)
	})
	if err != nil {
		return nil, err
	}
	return GetProduct(product.Number)
}

func ListComponents(productNumber string) ([]commercialResponses.CommercialComponentResponse, error) {
	var product commercialModels.CommercialProduct
	if err := config.DB.Where("number = ?", productNumber).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	var components []commercialModels.CommercialProductComponent
	config.DB.Where("product_number = ?", productNumber).Order("display_order ASC, component_code ASC").Find(&components)
	result := make([]commercialResponses.CommercialComponentResponse, 0, len(components))
	for _, component := range components {
		result = append(result, toComponentResponse(component))
	}
	return result, nil
}

func validateComponent(product commercialModels.CommercialProduct, componentCode string, included, optional bool) (string, string, error) {
	if product.Status == commercialModels.CommercialProductStatusInactive {
		return "", "", ErrProductTransition
	}
	code := strings.ToLower(strings.TrimSpace(componentCode))
	def, ok := addons.FindModuleDef(code)
	if !ok {
		return "", "", ErrComponentValidation
	}
	if included == optional {
		return "", "", ErrComponentValidation
	}
	return code, def.Name, nil
}

func CreateComponent(actor, productNumber string, request commercialRequests.CommercialComponentRequest) (*commercialResponses.CommercialComponentResponse, error) {
	var product commercialModels.CommercialProduct
	if err := config.DB.Where("number = ?", productNumber).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	code, name, err := validateComponent(product, request.ComponentCode, request.IsIncluded, request.IsOptional)
	if err != nil {
		return nil, err
	}
	var count int64
	config.DB.Model(&commercialModels.CommercialProductComponent{}).Where("product_number = ? AND component_code = ?", productNumber, code).Count(&count)
	if count > 0 {
		return nil, ErrComponentExists
	}
	component := commercialModels.CommercialProductComponent{
		Number: uuid.NewV4().String(), ProductNumber: productNumber, ComponentCode: code, ComponentName: name,
		IsIncluded: request.IsIncluded, IsOptional: request.IsOptional, Status: commercialModels.CommercialComponentStatusActive,
		DisplayOrder: request.DisplayOrder, CreatedByNumber: actor, UpdatedByNumber: actor,
	}
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&component).Error; err != nil {
			var exists int64
			tx.Model(&commercialModels.CommercialProductComponent{}).Where("product_number = ? AND component_code = ?", productNumber, code).Count(&exists)
			if exists > 0 {
				return ErrComponentExists
			}
			return err
		}
		return recordProductEvent(tx, "COMMERCIAL_PRODUCT_COMPONENT", component.Number, "COMMERCIAL_PRODUCT_COMPONENT_ADDED", "", componentState(component), request.Reason, actor)
	})
	if err != nil {
		return nil, err
	}
	response := toComponentResponse(component)
	return &response, nil
}

func UpdateComponent(actor, productNumber, componentNumber string, request commercialRequests.CommercialComponentRequest) (*commercialResponses.CommercialComponentResponse, error) {
	var product commercialModels.CommercialProduct
	if err := config.DB.Where("number = ?", productNumber).First(&product).Error; err != nil {
		return nil, ErrProductNotFound
	}
	var component commercialModels.CommercialProductComponent
	if err := config.DB.Where("number = ? AND product_number = ?", componentNumber, productNumber).First(&component).Error; err != nil {
		return nil, ErrComponentNotFound
	}
	if product.Status == commercialModels.CommercialProductStatusInactive {
		return nil, ErrProductTransition
	}
	if request.IsIncluded == request.IsOptional {
		return nil, ErrComponentValidation
	}
	code := strings.ToLower(strings.TrimSpace(request.ComponentCode))
	if code != "" && code != component.ComponentCode {
		return nil, ErrComponentValidation
	}
	before := componentState(component)
	component.IsIncluded = request.IsIncluded
	component.IsOptional = request.IsOptional
	component.DisplayOrder = request.DisplayOrder
	component.UpdatedByNumber = actor
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&component).Error; err != nil {
			return err
		}
		return recordProductEvent(tx, "COMMERCIAL_PRODUCT_COMPONENT", component.Number, "COMMERCIAL_PRODUCT_COMPONENT_UPDATED", before, componentState(component), request.Reason, actor)
	})
	if err != nil {
		return nil, err
	}
	response := toComponentResponse(component)
	return &response, nil
}

func SetComponentStatus(actor, productNumber, componentNumber string, request commercialRequests.CommercialComponentStatusRequest) (*commercialResponses.CommercialComponentResponse, error) {
	var product commercialModels.CommercialProduct
	if err := config.DB.Where("number = ?", productNumber).First(&product).Error; err != nil {
		return nil, ErrProductNotFound
	}
	if product.Status == commercialModels.CommercialProductStatusInactive {
		return nil, ErrProductTransition
	}
	var component commercialModels.CommercialProductComponent
	if err := config.DB.Where("number = ? AND product_number = ?", componentNumber, productNumber).First(&component).Error; err != nil {
		return nil, ErrComponentNotFound
	}
	status, err := commercialModels.NormaliseCommercialComponentStatus(request.Status)
	if err != nil {
		return nil, err
	}
	if component.Status == status {
		response := toComponentResponse(component)
		return &response, nil
	}
	before := componentState(component)
	component.Status = status
	component.UpdatedByNumber = actor
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&component).Error; err != nil {
			return err
		}
		return recordProductEvent(tx, "COMMERCIAL_PRODUCT_COMPONENT", component.Number, "COMMERCIAL_PRODUCT_COMPONENT_STATUS_CHANGED", before, componentState(component), request.Reason, actor)
	})
	if err != nil {
		return nil, err
	}
	response := toComponentResponse(component)
	return &response, nil
}
