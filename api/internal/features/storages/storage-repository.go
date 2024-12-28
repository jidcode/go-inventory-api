package storages

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/ventry/internal/domain"
)

type StorageRepository struct {
	db *sqlx.DB
}

func NewStorageRepository(db *sqlx.DB) *StorageRepository {
	return &StorageRepository{db: db}
}

func (repo *StorageRepository) ListStorages(inventoryId uuid.UUID) ([]domain.Storage, error) {
	storages := []domain.Storage{}
	query := `
		SELECT * FROM storages s
		WHERE s.inventory_id = $1
		GROUP BY s.id
		ORDER BY s.created_at DESC
	`

	if err := repo.db.Select(&storages, query, inventoryId); err != nil {
		return nil, err
	}

	// For each storage, fetch its units and products
	for i := range storages {
		units, err := repo.ListUnits(storages[i].Id)
		if err != nil {
			return nil, err
		}
		storages[i].Units = units
	}

	return storages, nil
}

func (repo *StorageRepository) GetStorage(storageId uuid.UUID) (*domain.Storage, error) {
	var storage domain.Storage
	query := `SELECT * FROM storages WHERE id = $1`

	if err := repo.db.Get(&storage, query, storageId); err != nil {
		return nil, err
	}

	units, err := repo.ListUnits(storageId)
	if err != nil {
		return nil, err
	}
	storage.Units = units

	return &storage, nil
}

func (repo *StorageRepository) CreateStorage(storage *domain.Storage) error {
	tx, err := repo.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
		}
	}()

	storageQuery := `
		INSERT INTO storages (
			id, name, location, capacity, inventory_id, created_at, updated_at
		) VALUES (
			:id, :name, :location, :capacity, :inventory_id, :created_at, :updated_at
		)
	`
	if _, err := tx.NamedExec(storageQuery, storage); err != nil {
		_ = tx.Rollback()
		return err
	}

	unitCount := 10
	if storage.Capacity <= 20 {
		unitCount = storage.Capacity
	}

	for i := 1; i <= unitCount; i++ {
		unit := domain.StorageUnit{
			Id:         uuid.New(),
			Name:       GenerateUnitName(i),
			StorageId:  storage.Id,
			IsOccupied: false,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		unitQuery := `
			INSERT INTO storage_units (
				id, name, storage_id, is_occupied, created_at, updated_at
			) VALUES (
				:id, :name, :storage_id, :is_occupied, :created_at, :updated_at
			)
		`
		if _, err := tx.NamedExec(unitQuery, unit); err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (repo *StorageRepository) EditStorage(storage *domain.Storage) error {
	query := `UPDATE storages
				SET name = :name, location = :location, capacity = :capacity,
				 	inventory_id = :inventory_id, updated_at = :updated_at
				WHERE id = :id`

	_, err := repo.db.NamedExec(query, storage)
	return err
}

func (repo *StorageRepository) DeleteStorage(storageId uuid.UUID) error {
	tx, err := repo.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
		}
	}()

	// Check for occupied units
	var occupiedCount int
	checkQuery := `SELECT COUNT(*) FROM storage_units WHERE storage_id = $1 AND is_occupied = true`
	if err := tx.Get(&occupiedCount, checkQuery, storageId); err != nil {
		_ = tx.Rollback()
		return err
	}

	if occupiedCount > 0 {
		_ = tx.Rollback()
		return errors.New("cannot delete storage with occupied units")
	}

	// Delete storage (cascade will handle units)
	if _, err := tx.Exec(`DELETE FROM storages WHERE id = $1`, storageId); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
