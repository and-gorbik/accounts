package service

import (
	"context"
	"net/url"

	jsoniter "github.com/json-iterator/go"

	"accounts/domain"
)

const (
	TimeLayout = "2006-01-02 15:04:05"
)

type BusinessError struct {
	error
}

type AccountService struct {
	repo repo
}

func New(repo repo) *AccountService {
	return &AccountService{
		repo: repo,
	}
}

func (s *AccountService) FilterAccounts(ctx context.Context, params url.Values) ([]byte, error) {
	qps, err := ParseQueryParams(params, true)
	if err != nil {
		return nil, BusinessError{err}
	}

	accounts, err := s.repo.FilterAccounts(ctx, BuildFilter(qps))
	if err != nil {
		return nil, BusinessError{err}
	}

	return jsoniter.Marshal(accounts)
}

func (s *AccountService) AddAccount(ctx context.Context, body []byte) error {
	var account domain.AccountInput
	if err := jsoniter.Unmarshal(body, &account); err != nil {
		return BusinessError{err}
	}

	if err := s.repo.AddAccount(ctx, account); err != nil {
		return BusinessError{err}
	}

	return nil
}

func (s *AccountService) UpdateAccount(ctx context.Context, body []byte) error {
	var account domain.AccountUpdate
	if err := jsoniter.Unmarshal(body, &account); err != nil {
		return BusinessError{err}
	}

	if err := s.repo.UpdateAccount(ctx, account); err != nil {
		return BusinessError{err}
	}

	return nil
}

func (s *AccountService) AddLikes(ctx context.Context, body []byte) error {
	var likes []domain.LikeInput
	if err := jsoniter.Unmarshal(body, &likes); err != nil {
		return BusinessError{err}
	}

	if err := s.repo.AddLikes(ctx, &domain.LikesInput{Likes: likes}); err != nil {
		return BusinessError{err}
	}

	return nil
}
