package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	Name        string     `db:"name" json:"name"`
	Description string     `db:"description" json:"description"`
	Code        string     `db:"code" json:"code"`
	Cost        float64    `db:"cost" json:"cost"`
	Price       float64    `db:"price" json:"price"`
	Quantity    int        `db:"quantity" json:"quantity"`
	InventoryID uuid.UUID  `db:"inventory_id" json:"inventoryID"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updatedAt"`
	Categories  []Category `db:"categories" json:"categories"`
	Storages    []Storage  `db:"storages" json:"storages"`
	Images      []Image    `db:"images" json:"images"`
}

type ProductInput struct {
	Name        string    `json:"name" validate:"required,min=3,max=100"`
	Description string    `json:"description"`
	Code        string    `json:"code"`
	Cost        float64   `json:"cost" validate:"gte=0"`
	Price       float64   `json:"price" validate:"gte=0"`
	Quantity    int       `json:"quantity" validate:"gte=0"`
	InventoryID uuid.UUID `json:"inventoryID" validate:"required"`
	Categories  []string  `json:"categories"`
	Storages    []Storage `json:"storages"`
	Images      []string  `json:"images"`
}

func (p *ProductInput) Sanitize() {
	p.Name = strings.TrimSpace(p.Name)
	p.Description = strings.TrimSpace(p.Description)
	p.Code = strings.TrimSpace(p.Code)
}
