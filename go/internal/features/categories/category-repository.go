package categories

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/ventry/internal/domain/models"
)

type CategoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(data *sqlx.DB) *CategoryRepository {
	return &CategoryRepository{db: data}
}

// LIST
func (repo *CategoryRepository) ListCategories(inventoryID uuid.UUID) ([]models.Category, error) {
	categories := []models.Category{}
	query := `SELECT * FROM categories WHERE inventory_id = $1`
	if err := repo.db.Select(&categories, query, inventoryID); err != nil {
		return nil, err
	}
	return categories, nil
}

// GET
func (repo *CategoryRepository) GetCategory(categoryID uuid.UUID) (models.Category, error) {
	var category models.Category
	query := `SELECT * FROM categories WHERE id = $1`
	err := repo.db.Get(&category, query, categoryID)
	if err != nil {
		return models.Category{}, err
	}

	productsQuery := `
        SELECT p.* 
        FROM products p
        JOIN product_categories pc ON p.id = pc.product_id
        WHERE pc.category_id = $1
    `
	products := []models.Product{}
	err = repo.db.Select(&products, productsQuery, categoryID)
	if err != nil {
		return models.Category{}, err
	}

	// category.Products = products

	return category, nil
}

// ADD
func (repo *CategoryRepository) CreateCategory(category *models.Category) error {
	category.ID = uuid.New()
	category.CreatedAt = time.Now()
	category.UpdatedAt = category.CreatedAt

	query := `INSERT INTO categories (id, name, description, inventory_id, created_at, updated_at)
			  VALUES (:id, :name, :description, :inventory_id, :created_at, :updated_at)`
	_, err := repo.db.NamedExec(query, category)
	return err
}

// EDIT
func (repo *CategoryRepository) UpdateCategory(category *models.Category) error {
	category.UpdatedAt = time.Now() // Update the updatedAt field

	query := `UPDATE categories
              SET name = :name, description = :description, inventory_id = :inventory_id,
                  updated_at = :updated_at
              WHERE id = :id`
	_, err := repo.db.NamedExec(query, category)
	return err
}

// DELETE
func (repo *CategoryRepository) DeleteCategory(categoryID uuid.UUID) error {
	_, err := repo.db.Exec(`DELETE FROM categories WHERE id = $1`, categoryID)
	return err
}
