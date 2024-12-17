package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Unit struct {
	ID         uuid.UUID `db:"id" json:"id"`
	Label      string    `db:"label" json:"label"`
	StorageID  uuid.UUID `db:"storage_id" json:"storageID"`
	IsOccupied bool      `db:"is_occupied" json:"isOccupied"`
	Capacity   int       `db:"capacity" json:"capacity"`
	CreatedAt  time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt  time.Time `db:"updated_at" json:"updatedAt"`
}

type UnitInput struct {
	Label      string    `json:"label" validate:"required,min=1,max=50"`
	StorageID  uuid.UUID `json:"storageID" validate:"required"`
	IsOccupied bool      `json:"isOccupied"`
	Capacity   int       `json:"capacity" validate:"min=0"`
}

func (input *UnitInput) Sanitize() {
	input.Label = strings.TrimSpace(input.Label)
}

func (input *UnitInput) ToUnit() Unit {
	return Unit{
		Label:      input.Label,
		StorageID:  input.StorageID,
		IsOccupied: input.IsOccupied,
		Capacity:   input.Capacity,
	}
}
