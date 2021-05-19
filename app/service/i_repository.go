package service

import (
	"context"

	"accounts/domain"
)

type repo interface {
	FilterAccounts(ctx context.Context, filter string) ([]domain.AccountOut, error)
	AddAccount(ctx context.Context, a domain.AccountInput) error
	UpdateAccount(ctx context.Context, a domain.AccountUpdate) error
	AddLikes(ctx context.Context, likes []domain.LikeInput) error
}
