package repository

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/marwans200/cyberguard/internal/database"
	"github.com/marwans200/cyberguard/internal/models"
)

// UserRepo handles all database operations for the users table.
type UserRepo struct{}

// Create inserts a new user and returns the created record.
func (r *UserRepo) Create(name, email, passwordHash, role string) (*models.User, error) {
	user := &models.User{}
	err := database.DB.QueryRow(
		`INSERT INTO users (name, email, password_hash, role)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, name, email, password_hash, role, created_at, updated_at`,
		name, email, passwordHash, role,
	).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	return user, err
}

// FindByEmail retrieves a user by their email address.
func (r *UserRepo) FindByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := database.DB.QueryRow(
		`SELECT id, name, email, password_hash, role, created_at, updated_at
		 FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		log.Printf("[db error] FindByEmail: %v", err)
	}
	return user, err
}

// FindByID retrieves a user by their UUID.
func (r *UserRepo) FindByID(id string) (*models.User, error) {
	user := &models.User{}
	err := database.DB.QueryRow(
		`SELECT id, name, email, password_hash, role, created_at, updated_at
		 FROM users WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return user, err
}

// CountAll returns the total number of registered users.
func (r *UserRepo) CountAll() (int, error) {
	var count int
	err := database.DB.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count)
	return count, err
}

// ListAnalysts returns all users with the analyst or admin role (for assignment dropdowns).
func (r *UserRepo) ListAnalysts() ([]models.UserResponse, error) {
	rows, err := database.DB.Query(
		`SELECT id, name, email, role, created_at FROM users WHERE role IN ('analyst','admin') ORDER BY name`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var analysts []models.UserResponse
	for rows.Next() {
		var u models.UserResponse
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.CreatedAt); err != nil {
			return nil, err
		}
		analysts = append(analysts, u)
	}
	return analysts, rows.Err()
}

// UpdatedAt touches the updated_at column on a user record.
func (r *UserRepo) TouchUpdatedAt(id string) error {
	_, err := database.DB.Exec(`UPDATE users SET updated_at = $1 WHERE id = $2`, time.Now(), id)
	return err
}

// Delete removes a user by ID.
func (r *UserRepo) Delete(id string) error {
	_, err := database.DB.Exec(`DELETE FROM users WHERE id = $1`, id)
	return err
}
