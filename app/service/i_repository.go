package service

import (
	"context"

	"accounts/app/repository"
	"accounts/domain"
)

type accountRepo interface {
	FilterAccounts(ctx context.Context, filter *repository.Filter) (*domain.AccountsOut, error)
	AddAccount(ctx context.Context, a domain.AccountInput) error
	UpdateAccount(ctx context.Context, a domain.AccountUpdate) error
	AddLikes(ctx context.Context, likes *domain.LikesInput) error
}
