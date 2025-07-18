package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) UserExists(ctx context.Context, guid uuid.UUID) (bool, error) {
	const query = `SELECT 1 FROM users WHERE guid = $1 LIMIT 1`
	var dummy int
	err := r.db.QueryRow(ctx, query, guid).Scan(&dummy)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}
	return true, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, guid uuid.UUID) error {
	const query = `INSERT INTO users (guid) VALUES ($1)`
	_, err := r.db.Exec(ctx, query, guid)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}
