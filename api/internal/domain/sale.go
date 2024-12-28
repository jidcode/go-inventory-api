package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Sale struct {
	Id              uuid.UUID  `db:"id" json:"id"`
	InventoryId     uuid.UUID  `db:"inventory_id" json:"inventoryId"`
	CustomerName    string     `db:"customer_name" json:"customerName"`
	CustomerContact string     `db:"customer_contact" json:"customerContact"`
	TotalAmount     float64    `db:"total_amount" json:"totalAmount"`
	CreatedAt       time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updatedAt"`
	Items           []SaleItem `json:"items,omitempty"`
}

type SaleItem struct {
	Id        uuid.UUID `db:"id" json:"id"`
	SaleId    uuid.UUID `db:"sale_id" json:"saleId"`
	ProductId uuid.UUID `db:"product_id" json:"productId"`
	Quantity  int       `db:"quantity" json:"quantity"`
	UnitPrice float64   `db:"unit_price" json:"unitPrice"`
	Subtotal  float64   `db:"subtotal" json:"subtotal"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

// DTOs
type SaleRequest struct {
	InventoryId     uuid.UUID         `json:"inventoryId" validate:"required"`
	CustomerName    string            `json:"customerName" validate:"required"`
	CustomerContact string            `json:"customerContact"`
	TotalAmount     float64           `json:"totalAmount"`
	Items           []SaleItemRequest `json:"items" validate:"required,min=1,dive"`
}

type SaleItemRequest struct {
	SaleId    uuid.UUID `json:"saleId" validate:"required"`
	ProductId uuid.UUID `json:"productId" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,min=1"`
	UnitPrice float64   `json:"unitPrice" validate:"required"`
}

func (req *SaleRequest) ToSale() *Sale {
	sale := &Sale{
		Id:              uuid.New(),
		InventoryId:     req.InventoryId,
		CustomerName:    strings.TrimSpace(req.CustomerName),
		CustomerContact: strings.TrimSpace(req.CustomerContact),
		TotalAmount:     req.TotalAmount,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Items:           make([]SaleItem, len(req.Items)),
	}

	// Convert item requests to items
	for i, item := range req.Items {
		sale.Items[i] = SaleItem{
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		}
	}

	return sale
}

func (req *SaleRequest) ToCreateSaleRequest() *Sale {
	return &Sale{
		Id:              uuid.New(),
		InventoryId:     req.InventoryId,
		CustomerName:    req.CustomerName,
		CustomerContact: req.CustomerContact,
		TotalAmount:     req.TotalAmount,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func (req *SaleRequest) ToUpdateSaleRequest(existing *Sale) *Sale {
	existing.InventoryId = req.InventoryId
	existing.CustomerName = req.CustomerName
	existing.CustomerContact = req.CustomerContact
	existing.TotalAmount = req.TotalAmount
	existing.UpdatedAt = time.Now()
	return existing
}

func (req *SaleRequest) Sanitize() {
	req.CustomerName = strings.TrimSpace(req.CustomerName)
	req.CustomerContact = strings.TrimSpace(req.CustomerContact)
}

// ////////
type SaleItemCreateRequest struct {
	ProductId uuid.UUID `json:"productId" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,min=1"`
	UnitPrice float64   `json:"unitPrice" validate:"required"`
}

type SaleCreateRequest struct {
	InventoryId     uuid.UUID               `json:"inventoryId" validate:"required"`
	CustomerName    string                  `json:"customerName" validate:"required"`
	CustomerContact string                  `json:"customerContact"`
	TotalAmount     float64                 `json:"totalAmount"`
	Items           []SaleItemCreateRequest `json:"items" validate:"required,min=1,dive"`
}

func (req *SaleCreateRequest) ToSale() *Sale {
	now := time.Now()
	sale := &Sale{
		Id:              uuid.New(),
		InventoryId:     req.InventoryId,
		CustomerName:    strings.TrimSpace(req.CustomerName),
		CustomerContact: strings.TrimSpace(req.CustomerContact),
		TotalAmount:     req.TotalAmount,
		CreatedAt:       now,
		UpdatedAt:       now,
		Items:           make([]SaleItem, len(req.Items)),
	}

	// Convert item requests to items
	for i, item := range req.Items {
		sale.Items[i] = SaleItem{
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		}
	}

	return sale
}
