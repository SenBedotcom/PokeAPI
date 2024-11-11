package models

import (
	"database/sql"
	"errors"
)

// User struct represents a user model in the database
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"` // Stored as a hashed password
}

// UserModel struct contains the database connection for executing queries
type UserModel struct {
	DB *sql.DB
}

// NewUserModel initializes a new UserModel with a database connection
func NewUserModel(db *sql.DB) *UserModel {
	return &UserModel{DB: db}
}

// CreateUser inserts a new user into the database
func (m *UserModel) CreateUser(user User) error {
	query := `INSERT INTO users (username, password) VALUES (?, ?)`
	_, err := m.DB.Exec(query, user.Username, user.Password)
	if err != nil {
		return err
	}
	return nil
}

// FindByUsername retrieves a user by username from the database
func (m *UserModel) FindByUsername(username string) (User, error) {
	var user User
	query := `SELECT id, username, password FROM users WHERE username = ?`
	row := m.DB.QueryRow(query, username)

	if err := row.Scan(&user.ID, &user.Username, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New("user not found")
		}
		return user, err
	}
	return user, nil
}
