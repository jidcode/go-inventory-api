package categories

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/ventry/internal/domain"
)

type CategoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(data *sqlx.DB) *CategoryRepository {
	return &CategoryRepository{db: data}
}

func (repo *CategoryRepository) ListCategories(inventoryId uuid.UUID) (*[]domain.Category, error) {
	categories := []domain.Category{}
	query := `SELECT * FROM categories WHERE inventory_id = $1`

	err := repo.db.Select(&categories, query, inventoryId)
	return &categories, err
}

func (repo *CategoryRepository) GetCategory(categoryId uuid.UUID) (*domain.Category, error) {
	var category domain.Category
	query := `SELECT * FROM categories WHERE id = $1`

	err := repo.db.Get(&category, query, categoryId)
	return &category, err
}

func (repo *CategoryRepository) CreateCategory(category *domain.Category) error {
	query := `INSERT 
				INTO categories (id, name, inventory_id, created_at, updated_at)
				VALUES (:id, :name, :inventory_id, :created_at, :updated_at)`

	_, err := repo.db.NamedExec(query, category)
	return err
}

func (repo *CategoryRepository) EditCategory(category *domain.Category) error {
	query := `UPDATE categories 
			SET name = :name, inventory_id = :inventory_id, updated_at = :updated_at
			WHERE id = :id`

	_, err := repo.db.NamedExec(query, category)
	return err
}

func (repo *CategoryRepository) DeleteCategory(categoryId uuid.UUID) error {
	query := `DELETE FROM categories WHERE id = $1`

	_, err := repo.db.Exec(query, categoryId)
	return err
}
