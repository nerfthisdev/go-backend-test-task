package repository

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
}

func Init(ctx context.Context) (*Repository, error) {
	dburi := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	cfg, err := pgxpool.ParseConfig(dburi)

	cfg.MaxConns = 25
	cfg.MaxConnIdleTime = 5 * time.Minute
	cfg.MaxConnLifetime = 2 * time.Hour

	if err != nil {
		return nil, err
	}

	dbpool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &Repository{
		DB: dbpool,
	}, nil
}

func (r *Repository) InitSchema(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS refresh_tokens (
			guid UUID NOT NULL,
			token TEXT PRIMARY KEY,
			expires_at TIMESTAMP NOT NULL
		);`

	_, err := r.DB.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
