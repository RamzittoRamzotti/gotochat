package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/RamzittoRamzotti/gotochat.git/internal/domain"
)

var ErrUserNotFound = errors.New("user not found")

type UserStorage struct {
	db *pgxpool.Pool
}

func NewUserStorage(db *pgxpool.Pool) *UserStorage {
	return &UserStorage{db: db}
}

func (s *UserStorage) GetUser(ctx context.Context, username string) (*domain.User, error) {
	var u domain.User

	row := s.db.QueryRow(ctx,
		`SELECT id, username, password_hash FROM users WHERE username = $1`,
		username,
	)

	err := row.Scan(&u.ID, &u.Username, &u.PasswordHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	return &u, nil
}

func (s *UserStorage) CreateUser(ctx context.Context, u *domain.User) error {
	_, err := s.db.Exec(ctx,
		`INSERT INTO users (id, username, password_hash) VALUES ($1, $2, $3)`,
		u.ID, u.Username, u.PasswordHash,
	)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}
