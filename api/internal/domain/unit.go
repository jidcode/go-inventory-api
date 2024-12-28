package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type StorageUnit struct {
	Id         uuid.UUID  `db:"id" json:"id"`
	Name       string     `db:"name" json:"name"`
	Label      string     `db:"label" json:"label"`
	StorageId  uuid.UUID  `db:"storage_id" json:"storageId"`
	IsOccupied bool       `db:"is_occupied" json:"isOccupied"`
	CreatedAt  time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updatedAt"`
	Items      []UnitItem `json:"items,omitempty"`
}

type UnitItem struct {
	Id            uuid.UUID `db:"id" json:"id"`
	StorageUnitId uuid.UUID `db:"storage_unit_id" json:"storageUnitId"`
	ProductId     uuid.UUID `db:"product_id" json:"productId"`
	Quantity      int       `db:"quantity" json:"quantity"`
	CreatedAt     time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt     time.Time `db:"updated_at" json:"updatedAt"`
	Product       *Product  `json:"product,omitempty"`
}

// DTO
type StorageUnitRequest struct {
	Label      string    `json:"label" validate:"required,min=2,max=50"`
	StorageId  uuid.UUID `json:"storageId" validate:"required"`
	IsOccupied bool      `json:"isOccupied"`
}

type UnitItemRequest struct {
	StorageUnitId uuid.UUID `json:"storageUnitId" validate:"required"`
	ProductId     uuid.UUID `json:"productId" validate:"required"`
	Quantity      int       `json:"quantity" validate:"required,min=1"`
}

func (req *StorageUnitRequest) ToCreateStorageUnitRequest() *StorageUnit {
	return &StorageUnit{
		Id:         uuid.New(),
		Label:      req.Label,
		StorageId:  req.StorageId,
		IsOccupied: false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

func (req *StorageUnitRequest) ToUpdateStorageUnitRequest(existing *StorageUnit) *StorageUnit {
	existing.Label = req.Label
	existing.IsOccupied = req.IsOccupied
	existing.UpdatedAt = time.Now()
	return existing
}

func (req *UnitItemRequest) ToCreateUnitItemRequest() *UnitItem {
	return &UnitItem{
		Id:            uuid.New(),
		StorageUnitId: req.StorageUnitId,
		ProductId:     req.ProductId,
		Quantity:      req.Quantity,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

func (req *UnitItemRequest) ToUpdateUnitItemRequest(existing *UnitItem) *UnitItem {
	existing.Quantity = req.Quantity
	existing.UpdatedAt = time.Now()
	return existing
}

func (req *StorageUnitRequest) Sanitize() {
	req.Label = strings.TrimSpace(req.Label)
}
