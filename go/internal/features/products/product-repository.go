package products

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/ventry/internal/domain/models"
)

type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// BASIC CRUD FUNCTIONALITY

func (repo *ProductRepository) ListProducts(inventoryID uuid.UUID) ([]models.Product, error) {
	products := []models.Product{}
	query := `SELECT * FROM products WHERE inventory_id = $1`
	err := repo.db.Select(&products, query, inventoryID)
	return products, err
}

func (repo *ProductRepository) GetProduct(id uuid.UUID) (models.Product, error) {
	var product models.Product
	query := `SELECT * FROM products WHERE id = $1`
	err := repo.db.Get(&product, query, id)
	return product, err
}

func (repo *ProductRepository) CreateProduct(product *models.Product, categoryNames []string, storages []models.Storage, imageUrls []string) error {
	// Start transaction
	tx, err := repo.db.Beginx()
	if err != nil {
		return err
	}

	// Ensure transaction rollback on panic
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
		}
	}()

	// Initialize product fields
	product.ID = uuid.New()
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	// Insert the product into the database
	query := `
        INSERT INTO products 
        (id, name, description, code, cost, price, quantity, inventory_id, created_at, updated_at)
        VALUES (:id, :name, :description, :code, :cost, :price, :quantity, :inventory_id, :created_at, :updated_at)
    `
	if _, err := tx.NamedExec(query, product); err != nil {
		_ = tx.Rollback()
		return err
	}

	// Handle category relationships
	if err := repo.handleProductCategories(tx, product.ID, product.InventoryID, categoryNames); err != nil {
		_ = tx.Rollback()
		return err
	}

	// Handle storage relationships
	if err := repo.handleProductStorages(tx, product.ID, product.InventoryID, storages); err != nil {
		_ = tx.Rollback()
		return err
	}

	// Handle image relationships
	if len(imageUrls) > 0 {
		if err := repo.handleProductImages(tx, product.ID, imageUrls); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	return tx.Commit()
}

func (repo *ProductRepository) UpdateProduct(product *models.Product) error {
	product.UpdatedAt = time.Now()

	query := `
        UPDATE products
        SET name = :name, description = :description, code = :code, cost = :cost, 
            price = :price, quantity = :quantity, inventory_id = :inventory_id, updated_at = :updated_at
        WHERE id = :id
    `
	_, err := repo.db.NamedExec(query, product)
	return err
}

func (repo *ProductRepository) DeleteProduct(id uuid.UUID) error {
	tx, err := repo.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	deleteStoragesQuery := `DELETE FROM product_storages WHERE product_id = $1`
	_, err = tx.Exec(deleteStoragesQuery, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	deleteCategoriesQuery := `DELETE FROM product_categories WHERE product_id = $1`
	_, err = tx.Exec(deleteCategoriesQuery, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	deleteProductQuery := `DELETE FROM products WHERE id = $1`
	_, err = tx.Exec(deleteProductQuery, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// PRODUCT RELATION HELPERS

func (repo *ProductRepository) GetProductWithRelations(productID uuid.UUID) (models.Product, error) {
	var product models.Product

	productQuery := `
        SELECT * FROM products WHERE id = $1
    `
	if err := repo.db.Get(&product, productQuery, productID); err != nil {
		return product, err
	}

	categoriesQuery := `
        SELECT c.* FROM categories c
        INNER JOIN product_categories pc ON c.id = pc.category_id
        WHERE pc.product_id = $1
    `
	if err := repo.db.Select(&product.Categories, categoriesQuery, productID); err != nil {
		return product, err
	}

	storagesQuery := `
        SELECT s.* FROM storages s
        INNER JOIN product_storages ps ON s.id = ps.storage_id
        WHERE ps.product_id = $1
    `
	if err := repo.db.Select(&product.Storages, storagesQuery, productID); err != nil {
		return product, err
	}

	imagesQuery := `
        SELECT * FROM images 
        WHERE product_id = $1 
        ORDER BY is_primary DESC, created_at ASC
    `
	if err := repo.db.Select(&product.Images, imagesQuery, productID); err != nil {
		return product, err
	}

	return product, nil
}

func (repo *ProductRepository) handleProductCategories(tx *sqlx.Tx, productID, inventoryID uuid.UUID, categoryNames []string) error {
	for _, categoryName := range categoryNames {
		var category models.Category
		query := `SELECT * FROM categories WHERE name = $1 AND inventory_id = $2`

		// Try to get the category
		if err := tx.Get(&category, query, categoryName, inventoryID); err != nil {
			if err == sql.ErrNoRows {
				// Create a new category if not found
				category = models.Category{
					ID:          uuid.New(),
					Name:        categoryName,
					InventoryID: inventoryID,
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
		if _, err := tx.Exec(relationQuery, productID, category.ID, time.Now()); err != nil {
			return err
		}
	}
	return nil
}

func (repo *ProductRepository) handleProductStorages(tx *sqlx.Tx, productID, inventoryID uuid.UUID, storages []models.Storage) error {
	for _, storage := range storages {
		// Validate that the storage exists and belongs to the inventory
		var existingStorage models.Storage
		validationQuery := `
            SELECT * FROM storages WHERE id = $1 AND inventory_id = $2
        `
		if err := tx.Get(&existingStorage, validationQuery, storage.ID, inventoryID); err != nil {
			return err
		}

		// Insert product-storage relationship
		relationQuery := `
            INSERT INTO product_storages (product_id, storage_id, quantity, created_at, updated_at)
            VALUES ($1, $2, $3, $4, $5)
        `
		if _, err := tx.Exec(relationQuery, productID, storage.ID, 0, time.Now(), time.Now()); err != nil {
			return err
		}
	}

	return nil
}

func (repo *ProductRepository) handleProductImages(tx *sqlx.Tx, productID uuid.UUID, imageUrls []string) error {
	for i, imageUrl := range imageUrls {
		// Prepare the image object
		image := models.Image{
			ID:        uuid.New(),
			URL:       imageUrl,
			ProductID: productID,
			IsPrimary: i == 0, // Set the first image as primary
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
