package products

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/ventry/internal/domain/models"
)

type ImagesRepository struct {
	db *sqlx.DB
}

func NewImagesRepository(data *sqlx.DB) *ImagesRepository {
	return &ImagesRepository{db: data}
}

func (repo *ImagesRepository) ListImages(productID uuid.UUID) ([]models.Image, error) {
	images := []models.Image{}
	query := `SELECT * FROM images WHERE product_id = $1`
	err := repo.db.Select(&images, query, productID)
	return images, err
}

func (repo *ImagesRepository) GetImage(id uuid.UUID) (models.Image, error) {
	var image models.Image
	query := `SELECT * FROM images WHERE id = $1`
	err := repo.db.Get(&image, query, id)
	return image, err
}

func (repo *ImagesRepository) CreateImage(image *models.Image) error {
	image.ID = uuid.New()
	image.CreatedAt = time.Now()
	image.UpdatedAt = time.Now()

	query := `INSERT INTO images (id, url, product_id, is_primary, created_at, updated_at)
			  VALUES (:id, :url, :product_id, :is_primary, :created_at, :updated_at)`
	_, err := repo.db.NamedExec(query, image)
	return err
}

func (repo *ImagesRepository) UpdateImage(image *models.Image) error {
	image.UpdatedAt = time.Now()
	query := `UPDATE images
              SET url = :url, product_id = :product_id, is_primary = :is_primary, updated_at = :updated_at
              WHERE id = :id`
	_, err := repo.db.NamedExec(query, image)
	return err
}

func (repo *ImagesRepository) DeleteImage(id uuid.UUID) error {
	query := `DELETE FROM images WHERE id = $1`
	_, err := repo.db.Exec(query, id)
	return err
}
