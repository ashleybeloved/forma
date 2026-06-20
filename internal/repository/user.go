package repository

import (
	"database/sql"
	"forma/internal/model"
	"log/slog"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) UsernameExists(username string) bool {
	var id int
	err := r.DB.QueryRow(`SELECT id FROM users WHERE username = ?`, username).Scan(&id)
	if err == sql.ErrNoRows {
		return false
	}

	if err != nil {
		slog.Error("database error", "error", err)
		return false
	}

	return true
}

func (r *UserRepository) AddUser(username string, hashedPassword string) error {
	exists := r.UsernameExists(username)
	if exists {
		return ErrUsernameAlreadyExists
	}

	_, err := r.DB.Exec(`INSERT INTO users (username, password) VALUES (?, ?)`, username, hashedPassword)
	if err != nil {
		slog.Warn("failed to insert user in database", "error", err)
		return err
	}
	return nil
}

func (r *UserRepository) GetUserByID(id int) (*model.User, error) {
	user := model.User{}

	err := r.DB.QueryRow(`SELECT id, username, password, created_at FROM users WHERE id = ?`, id).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}

		slog.Error("failed to execute query", "error", err)
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByUsername(username string) (*model.User, error) {
	user := model.User{}

	err := r.DB.QueryRow(`SELECT id, username, password, created_at FROM users WHERE username = ?`, username).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}

		slog.Error("failed to execute query", "error", err)
		return nil, err
	}

	return &user, nil
}
