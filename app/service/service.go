package service

import (
	"accounts/domain"
)

type AccountService struct {
	repo repository
}

type repository interface {
}

func (s *AccountService) FilterAccounts(params map[string]QueryParamWithOp) ([]domain.AccountOut, error) {
	return nil, nil
}
