package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nerfthisdev/go-backend-test-task/internal/config"
)

func InitDB(ctx context.Context, cfg config.Config) (*pgxpool.Pool, error) {
	dburi := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBAddress,
		cfg.DBPort,
		cfg.DBName,
	)

	pgCfg, err := pgxpool.ParseConfig(dburi)
	if err != nil {
		return nil, fmt.Errorf("failed to parse db config: %w", err)
	}

	pgCfg.MaxConns = 25
	pgCfg.MaxConnIdleTime = 5 * time.Minute
	pgCfg.MaxConnLifetime = 2 * time.Hour

	return pgxpool.NewWithConfig(ctx, pgCfg)
}
