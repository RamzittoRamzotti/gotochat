package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	models "github.com/RamzittoRamzotti/gotochat.git/internal/domain/models"
)

var ErrUserNotFound = errors.New("user not found")

type UserStorage struct {
	db *pgxpool.Pool
}

func NewUserStorage(db *pgxpool.Pool) *UserStorage {
	return &UserStorage{db: db}
}

func (s *UserStorage) GetUser(username string) (*models.User, error) {
	var u models.User

	row := s.db.QueryRow(context.Background(),
		`SELECT username, password_hash FROM users WHERE username = $1`,
		username,
	)
	err := row.Scan(&u.Username, &u.PasswordHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return &u, nil
}

func (s *UserStorage) CreateUser(username, passwordHash string) (*models.UserCreate, error) {
	_, err := s.db.Exec(context.Background(),
		`INSERT INTO users (username, password_hash) VALUES ($1, $2)`,
		username, []byte(passwordHash),
	)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &models.UserCreate{Username: username}, nil
}
