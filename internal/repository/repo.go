package repository

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nerfthisdev/go-backend-test-task/internal/model"
	"github.com/pkg/errors"
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

func (r *Repository) Save(ctx context.Context, user model.User) error {
	query := `INSERT INTO auth
	(guid, user_agent, ip_address, token, created_at, expires_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.DB.Exec(ctx, query, user.Guid, user.UserAgent, user.IpAddress, user.Token, user.CreatedAt, user.ExpiresAt)
	if err != nil {
		return errors.Wrap(err, "failed to save token")
	}
	return nil
}

func (r *Repository) Get(ctx context.Context, guid string) (string, string, error) {
	var refreshToken string
	var ip string
	query := `SELECT token, ip_address FROM auth WHERE guid = $1`

	err := r.DB.QueryRow(ctx, query, guid).Scan(&refreshToken, &ip)
	if err != nil {
		return "", "", err
	}

	return refreshToken, ip, nil
}
