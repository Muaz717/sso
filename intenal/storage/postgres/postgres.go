package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"sso/intenal/config"
	"sso/intenal/domain/models"
	"sso/intenal/storage"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(ctx context.Context, cfg config.DBConfig) (*Storage, error) {
	const op = "storage.postgres.New"

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.Username,
		cfg.DBPassword,
		cfg.Host,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect db: %w", err)
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(
	ctx context.Context,
	email string,
	passHash []byte,
) (int64, error) {
	const op = "postgres.SaveUser"

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback(ctx)

	querySelect := `SELECT COUNT(id) as count FROM users WHERE email = $1`
	row := tx.QueryRow(ctx, querySelect, email)

	var count int
	err = row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	if count > 0 {
		return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
	}

	queryInsert := `INSERT INTO users(email, passhash) VALUES($1, $2) RETURNING id`
	row = tx.QueryRow(ctx, queryInsert, email, passHash)

	var userId int64
	err = row.Scan(&userId)
	if err != nil {
		_ = tx.Rollback(ctx)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("%s: %w", err)
	}

	return userId, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "postgres.User"

	query := `SELECT id, email, passhash FROM users WHERE email = $1`
	row := s.db.QueryRow(ctx, query, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "postgres.IsAdmin"

	query := `SELECT is_admin FROM users WHERE id = $1`
	row := s.db.QueryRow(ctx, query, userID)

	var isAdmin bool
	err := row.Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}

func (s *Storage) App(ctx context.Context, appID int) (models.App, error) {
	const op = "postgres.App"

	query := `SELECT id, name, secret FROM apps WHERE id = $1`
	row := s.db.QueryRow(ctx, query, appID)

	var app models.App
	err := row.Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}
