package repositories

import (
	"database/sql"

	"github.com/Jaidenmagnan/waygates/db"
	"github.com/Jaidenmagnan/waygates/models"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository() *UserRepository {
    return &UserRepository{db: db.DB}
}


func (r *UserRepository) Create(user models.User) (int64, error) {
    query := "INSERT INTO users (email, password) VALUES (?, ?)"
    
    result, err := r.db.Exec(query, user.Email, user.Password)
    if err != nil {
        return 0, err
    }
    
    return result.LastInsertId()
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := "SELECT id, email, password FROM users WHERE email = ?"

	rows, err := r.db.Query(query, email)
	if err != nil {
		return nil, err
	}
	
	defer rows.Close()

	var user models.User
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Email, &user.Password); err != nil {
			return nil, err
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &user, nil
}
