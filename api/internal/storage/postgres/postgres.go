package postgres

import (
	"context"
	"database/sql"
	"errors"
	"gRPC/api/internal/domain/models"
	"gRPC/api/internal/storage"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(dsn string) (*Storage, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	var id int64
	err := s.db.QueryRowContext(ctx,
		"INSERT INTO users (email, pass_hash) VALUES ($1, $2) RETURNING id",
		email, passHash,
	).Scan(&id)

	if err != nil {
		// Проверка на unique violation
		return 0, storage.ErrUserExists
	}
	return id, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	var user models.User
	err := s.db.QueryRowContext(ctx,
		"SELECT id, email, pass_hash, is_admin FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Email, &user.PassHash, &user.IsAdmin)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, storage.ErrUserNotFound
		}
		return models.User{}, err
	}
	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	var isAdmin bool
	err := s.db.QueryRowContext(ctx,
		"SELECT is_admin FROM users WHERE id = $1",
		userID,
	).Scan(&isAdmin)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, storage.ErrUserNotFound
		}
		return false, err
	}
	return isAdmin, nil
}

func (s *Storage) App(ctx context.Context, appID int) (models.App, error) {
	var app models.App
	err := s.db.QueryRowContext(ctx,
		"SELECT id, name, secret FROM apps WHERE id = $1",
		appID,
	).Scan(&app.ID, &app.Name, &app.Secret)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, storage.ErrAppNotFound
		}
		return models.App{}, err
	}
	return app, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
