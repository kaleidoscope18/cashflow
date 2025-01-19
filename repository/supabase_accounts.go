package repository

import (
	"cashflow/models"
	"cashflow/utils"
	"context"
	"log"
)

func (r *supabaseRepository) ListAccounts(ctx context.Context) ([]string, error) {
	rows, err := r.db.Query(ctx, "SELECT name FROM cashflow.accounts WHERE author = $1", utils.GetContextValue(ctx, models.UserId))
	if err != nil {
		log.Printf("error when querying SELECT on cashflow.accounts: %v", err)
		return nil, err
	}
	defer rows.Close()

	var accounts []string
	for rows.Next() {
		var account string
		rows.Scan(&account)
		accounts = append(accounts, account)
	}
	return accounts, nil
}
