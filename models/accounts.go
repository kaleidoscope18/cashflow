package models

import "context"

type AccountRepository interface {
	ListAccounts(ctx context.Context) ([]string, error)
}
