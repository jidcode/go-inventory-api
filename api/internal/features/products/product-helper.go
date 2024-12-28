package products

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/ventry/internal/domain"
)

func (repo *ProductRepository) UpdateProductQuantity(productId uuid.UUID, quantity int) error {
	query := `UPDATE products SET quantity = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`

	_, err := repo.db.Exec(query, quantity, productId)
	return err
}

func (repo *ProductRepository) GetProductWithRelations(productId uuid.UUID) (domain.Product, error) {
	var product domain.Product

	productQuery := `
        SELECT * FROM products WHERE id = $1
    `
	if err := repo.db.Get(&product, productQuery, productId); err != nil {
		return product, err
	}

	categoriesQuery := `
        SELECT c.* FROM categories c
        INNER JOIN product_categories pc ON c.id = pc.category_id
        WHERE pc.product_id = $1
    `
	if err := repo.db.Select(&product.Categories, categoriesQuery, productId); err != nil {
		return product, err
	}

	storagesQuery := `
        SELECT s.* FROM storages s
        INNER JOIN product_storages ps ON s.id = ps.storage_id
        WHERE ps.product_id = $1
    `
	if err := repo.db.Select(&product.Storages, storagesQuery, productId); err != nil {
		return product, err
	}

	imagesQuery := `
        SELECT * FROM images 
        WHERE product_id = $1 
        ORDER BY is_primary DESC, created_at ASC
    `
	if err := repo.db.Select(&product.Images, imagesQuery, productId); err != nil {
		return product, err
	}

	return product, nil
}

func (repo *ProductRepository) handleProductCategories(tx *sqlx.Tx, productId, inventoryId uuid.UUID, categoryNames []string) error {
	for _, categoryName := range categoryNames {
		var category domain.Category
		query := `SELECT * FROM categories WHERE name = $1 AND inventory_id = $2`

		// Try to get the category
		if err := tx.Get(&category, query, categoryName, inventoryId); err != nil {
			if err == sql.ErrNoRows {
				// Create a new category if not found
				category = domain.Category{
					Id:          uuid.New(),
					Name:        categoryName,
					InventoryId: inventoryId,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}

				insertQuery := `
                    INSERT INTO categories (id, name, inventory_id, created_at, updated_at)
                    VALUES (:id, :name, :inventory_id, :created_at, :updated_at)
                `
				if _, err := tx.NamedExec(insertQuery, &category); err != nil {
					return err
				}
			} else {
				return err
			}
		}

		// Add product-category relationship
		relationQuery := `
            INSERT INTO product_categories (product_id, category_id, created_at)
            VALUES ($1, $2, $3)
            ON CONFLICT DO NOTHING
        `
		if _, err := tx.Exec(relationQuery, productId, category.Id, time.Now()); err != nil {
			return err
		}
	}
	return nil
}

func (repo *ProductRepository) handleProductStorages(tx *sqlx.Tx, productId, inventoryId uuid.UUID, storages []domain.Storage) error {
	for _, storage := range storages {
		// Validate that the storage exists and belongs to the inventory
		var existingStorage domain.Storage
		validationQuery := `
            SELECT * FROM storages WHERE id = $1 AND inventory_id = $2
        `
		if err := tx.Get(&existingStorage, validationQuery, storage.Id, inventoryId); err != nil {
			return err
		}

		// Insert product-storage relationship
		relationQuery := `
            INSERT INTO product_storages (product_id, storage_id, created_at)
            VALUES ($1, $2, $3)
			ON CONFLICT DO NOTHING
		`

		if _, err := tx.Exec(relationQuery, productId, storage.Id, time.Now()); err != nil {
			return err
		}
	}

	return nil
}

func (repo *ProductRepository) handleProductImages(tx *sqlx.Tx, productId uuid.UUID, imageUrls []string) error {
	for i, imageUrl := range imageUrls {
		// Prepare the image object
		image := domain.Image{
			Id:        uuid.New(),
			URL:       imageUrl,
			ProductId: productId,
			IsPrimary: i == 0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Insert the image
		query := `
			INSERT INTO images 
			(id, url, product_id, is_primary, created_at, updated_at)
			VALUES (:id, :url, :product_id, :is_primary, :created_at, :updated_at)
		`
		if _, err := tx.NamedExec(query, &image); err != nil {
			return err
		}
	}
	return nil
}

func (repo *ProductRepository) clearProductRelationships(tx *sqlx.Tx, productId uuid.UUID) error {
	queries := []string{
		`DELETE FROM product_categories WHERE product_id = $1`,
		`DELETE FROM product_storages WHERE product_id = $1`,
		`DELETE FROM images WHERE product_id = $1`,
	}

	for _, query := range queries {
		if _, err := tx.Exec(query, productId); err != nil {
			return err
		}
	}

	return nil
}
