package products

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/ventry/internal/domain"
)

type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(data *sqlx.DB) *ProductRepository {
	return &ProductRepository{db: data}
}

func (repo *ProductRepository) ListProducts(inventoryId uuid.UUID) (*[]domain.Product, error) {
	// First, get all products for the inventory
	products := []domain.Product{}
	productsQuery := `SELECT * FROM products WHERE inventory_id = $1`

	if err := repo.db.Select(&products, productsQuery, inventoryId); err != nil {
		return nil, err
	}

	// For each product, fetch related data
	for i := range products {
		// Fetch categories
		categoriesQuery := `
			SELECT c.* FROM categories c
			INNER JOIN product_categories pc ON c.id = pc.category_id
			WHERE pc.product_id = $1
		`
		if err := repo.db.Select(&products[i].Categories, categoriesQuery, products[i].Id); err != nil {
			return nil, err
		}

		// Fetch storages
		storagesQuery := `
			SELECT s.* FROM storages s
			INNER JOIN product_storages ps ON s.id = ps.storage_id
			WHERE ps.product_id = $1
		`
		if err := repo.db.Select(&products[i].Storages, storagesQuery, products[i].Id); err != nil {
			return nil, err
		}

		// Fetch images
		imagesQuery := `
			SELECT * FROM images 
			WHERE product_id = $1 
			ORDER BY is_primary DESC, created_at ASC
		`
		if err := repo.db.Select(&products[i].Images, imagesQuery, products[i].Id); err != nil {
			return nil, err
		}
	}

	return &products, nil
}

func (repo *ProductRepository) GetProduct(productId uuid.UUID) (*domain.Product, error) {
	var product domain.Product
	query := `SELECT * FROM products WHERE id = $1`

	err := repo.db.Get(&product, query, productId)
	return &product, err
}

func (repo *ProductRepository) CreateProduct(product *domain.Product, categoryNames []string, storages []domain.Storage, imageUrls []string) error {
	// Start transaction
	tx, err := repo.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
		}
	}()

	// Insert the product into the database
	query := `INSERT INTO products (
				id, name, description, sku, code, quantity, restock_level, optimal_level, 
				cost, price, inventory_id, created_at, updated_at
			  ) VALUES (
			  	:id, :name, :description, :sku, :code, :quantity, :restock_level, :optimal_level, 
				:cost, :price, :inventory_id, :created_at, :updated_at
			  )`

	if _, err := tx.NamedExec(query, product); err != nil {
		_ = tx.Rollback()
		return err
	}

	// Handle category relationships
	if err := repo.handleProductCategories(tx, product.Id, product.InventoryId, categoryNames); err != nil {
		_ = tx.Rollback()
		return err
	}

	// Handle storage relationships
	if err := repo.handleProductStorages(tx, product.Id, product.InventoryId, storages); err != nil {
		_ = tx.Rollback()
		return err
	}

	// Handle image relationships
	if len(imageUrls) > 0 {
		if err := repo.handleProductImages(tx, product.Id, imageUrls); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	return tx.Commit()
}

func (repo *ProductRepository) EditProduct(product *domain.Product, categoryNames []string, storages []domain.Storage, imageUrls []string) error {
	// Start transaction
	tx, err := repo.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
		}
	}()

	// Update the product
	query := `UPDATE products SET 
				name = :name,
				description = :description,
				sku = :sku,
				code = :code,
				quantity = :quantity,
				restock_level = :restock_level,
				optimal_level = :optimal_level,
				cost = :cost,
				price = :price,
				inventory_id = :inventory_id,
				updated_at = :updated_at
			 WHERE id = :id`

	if _, err := tx.NamedExec(query, product); err != nil {
		_ = tx.Rollback()
		return err
	}

	// Clear existing relationships
	if err := repo.clearProductRelationships(tx, product.Id); err != nil {
		_ = tx.Rollback()
		return err
	}

	// Handle category relationships
	if err := repo.handleProductCategories(tx, product.Id, product.InventoryId, categoryNames); err != nil {
		_ = tx.Rollback()
		return err
	}

	// Handle storage relationships
	if err := repo.handleProductStorages(tx, product.Id, product.InventoryId, storages); err != nil {
		_ = tx.Rollback()
		return err
	}

	// Handle image relationships
	if len(imageUrls) > 0 {
		if err := repo.handleProductImages(tx, product.Id, imageUrls); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (repo *ProductRepository) DeleteProduct(productId uuid.UUID) error {
	query := `DELETE FROM products WHERE id = $1`

	_, err := repo.db.Exec(query, productId)
	return err
}
