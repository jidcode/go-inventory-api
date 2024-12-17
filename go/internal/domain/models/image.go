package models

import (
	"time"

	"github.com/google/uuid"
)

type Image struct {
	ID        uuid.UUID `db:"id" json:"id"`
	URL       string    `db:"url" json:"url"`
	ProductID uuid.UUID `db:"product_id" json:"productId"`
	IsPrimary bool      `db:"is_primary" json:"isPrimary"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

type ImageInput struct {
	URL       string    `json:"url" validate:"required,url"`
	ProductID uuid.UUID `json:"productId" validate:"required"`
	IsPrimary bool      `json:"isPrimary,omitempty"`
}
