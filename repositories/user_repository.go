// User related database operations.
package repositories

import (
	"database/sql"

	"github.com/Jaidenmagnan/waygates/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create a new user in the database.
func (r *UserRepository) Create(user models.CreateUser) (models.User, error) {
	query := "INSERT INTO users (email, username,password) VALUES (?, ?, ?)"

	result, err := r.db.Exec(query, user.Email, user.Username, user.Password)
	if err != nil {
		return models.User{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.User{}, err
	}

	return models.User{
		ID:       int(id),
		Email:    user.Email,
		Username: user.Username,
		Password: user.Password,
	}, nil
}

// Retrieves a user by email.
func (r *UserRepository) GetByEmail(email string) (*models.User, bool) {
	query := "SELECT id, email, username, password FROM users WHERE email = ?"

	row := r.db.QueryRow(query, email)

	var user models.User
	if err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password); err != nil {
		return nil, false
	}

	return &user, true
}

// Retrieves a user by ID.
func (r *UserRepository) GetByID(id int) (*models.User, bool) {
	query := "SELECT id, email, username, password FROM users WHERE id = ?"

	row := r.db.QueryRow(query, id)

	var user models.User
	if err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password); err != nil {
		return nil, false
	}

	return &user, true
}
