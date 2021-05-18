package repository

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"

	"accounts/domain"
)

type Repository struct {
	conn *pgxpool.Conn
}

func (r *Repository) FilterAccounts(ctx context.Context, filter string) ([]domain.AccountOut, error) {
	return nil, nil
}
