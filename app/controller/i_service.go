package controller

import (
	"context"
	"net/url"
)

type accountService interface {
	FilterAccounts(ctx context.Context, params url.Values) ([]byte, error)
	AddAccount(ctx context.Context, body []byte) error
	UpdateAccount(ctx context.Context, body []byte) error
	AddLikes(ctx context.Context, body []byte) error
}
