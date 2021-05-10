package service

import (
	"context"
	"net/url"

	"accounts/domain"
)

const (
	TimeLayout = "2006-01-02 15:04:05"
)

type AccountService struct {
	repo repo
}

func (s *AccountService) FilterAccounts(ctx context.Context, params url.Values) ([]domain.AccountOut, error) {
	qps, err := ParseQueryParams(params, true)
	if err != nil {
		return nil, err
	}

	return s.repo.FilterAccounts(ctx, BuildFilter(qps))
}
