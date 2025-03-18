package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"sso/intenal/config"
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

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	querySelect := `SELECT 1 FROM users WHERE email = $1 LIMIT 1`

	row := tx.QueryRow(ctx, querySelect, email)

	query := `INSERT INTO users(email, passhash) VALUES($1, $2) RETURNING id`

	row := s.db.QueryRow(ctx, query, email, passHash)

	var userId int64
	err := row.Scan(&userId)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return userId, nil
}
