package auth

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/ventry/internal/domain/models"
)

type AuthRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (repo *AuthRepository) CreateUser(user *models.User) error {
	query := `INSERT INTO users (id, username, email, password, role)
              VALUES ($1, $2, $3, $4, $5)
              RETURNING created_at, updated_at`

	user.ID = uuid.New()

	err := repo.db.QueryRow(
		query,
		user.ID,
		user.Username,
		user.Email,
		user.Password,
		user.Role,
	).Scan(&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	return nil
}

func (repo *AuthRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, email, password, role, created_at, updated_at 
			  FROM users WHERE email = $1`

	err := repo.db.Get(&user, query, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email %s: %v", email, err)
	}

	// Debug print
	fmt.Printf("Retrieved user: %+v\n", user)

	return &user, nil
}

func (repo *AuthRepository) GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, email, password, role, created_at, updated_at 
			  FROM users WHERE id = $1`

	err := repo.db.Get(&user, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID %s: %v", id, err)
	}

	return &user, nil
}
