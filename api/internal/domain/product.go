package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	Id           uuid.UUID  `db:"id" json:"id"`
	Name         string     `db:"name" json:"name"`
	Description  *string    `db:"description" json:"description"`
	SKU          string     `db:"sku" json:"sku"`
	Code         *string    `db:"code" json:"code"`
	Quantity     int        `db:"quantity" json:"quantity"`
	RestockLevel int        `db:"restock_level" json:"restockLevel"`
	OptimalLevel int        `db:"optimal_level" json:"optimalLevel"`
	Cost         float64    `db:"cost" json:"cost"`
	Price        float64    `db:"price" json:"price"`
	InventoryId  uuid.UUID  `db:"inventory_id" json:"inventoryId"`
	CreatedAt    time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updatedAt"`
	Categories   []Category `db:"categories" json:"categories"`
	Storages     []Storage  `db:"storages" json:"storages"`
	Images       []Image    `db:"images" json:"images"`
}

// DTOs
type ProductRequest struct {
	Name         string    `json:"name" validate:"required"`
	Description  *string   `json:"description"`
	SKU          string    `json:"sku" validate:"required"`
	Code         *string   `json:"code"`
	Quantity     int       `json:"quantity"`
	RestockLevel int       `json:"restockLevel"`
	OptimalLevel int       `json:"optimalLevel"`
	Cost         float64   `json:"cost"`
	Price        float64   `json:"price"`
	InventoryId  uuid.UUID `json:"inventoryId" validate:"required"`
	Categories   []string  `db:"categories" json:"categories"`
	Storages     []Storage `db:"storages" json:"storages"`
	Images       []string  `db:"images" json:"images"`
}

type ProductResponse struct {
	Id           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Description  *string   `json:"description"`
	SKU          string    `json:"sku"`
	Code         *string   `json:"code"`
	Quantity     int       `json:"quantity"`
	RestockLevel int       `json:"restockLevel"`
	OptimalLevel int       `json:"optimalLevel"`
	Cost         float64   `json:"cost"`
	Price        float64   `json:"price"`
	InventoryId  uuid.UUID `json:"inventoryId"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func (req *ProductRequest) ToCreateProductRequest() *Product {
	return &Product{
		Id:           uuid.New(),
		Name:         req.Name,
		Description:  req.Description,
		SKU:          req.SKU,
		Code:         req.Code,
		Quantity:     req.Quantity,
		RestockLevel: req.RestockLevel,
		OptimalLevel: req.OptimalLevel,
		Cost:         req.Cost,
		Price:        req.Price,
		InventoryId:  req.InventoryId,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func (req *ProductRequest) ToEditProductRequest(existingProduct *Product) *Product {
	existingProduct.Name = req.Name
	existingProduct.Description = req.Description
	existingProduct.SKU = req.SKU
	existingProduct.Code = req.Code
	existingProduct.Quantity = req.Quantity
	existingProduct.RestockLevel = req.RestockLevel
	existingProduct.OptimalLevel = req.OptimalLevel
	existingProduct.Cost = req.Cost
	existingProduct.Price = req.Price
	existingProduct.InventoryId = req.InventoryId
	existingProduct.UpdatedAt = time.Now()

	return existingProduct
}

func (req *ProductRequest) Sanitize() {
	req.Name = strings.TrimSpace(req.Name)
	req.SKU = strings.TrimSpace(req.SKU)
	if req.Code != nil {
		trimmedCode := strings.TrimSpace(*req.Code)
		req.Code = &trimmedCode
	}
	if req.Description != nil {
		trimmedDesc := strings.TrimSpace(*req.Description)
		req.Description = &trimmedDesc
	}
}
