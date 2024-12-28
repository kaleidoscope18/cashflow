package repository

import (
	"cashflow/models"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/supabase-community/postgrest-go"
)

func (repo *supabaseRepository) InsertBalance(ctx context.Context, amount float64, date string) (models.Balance, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	balance := models.Balance{
		Amount: amount,
		Date:   date,
	}

	_, _, err := repo.db.From("balances").
		Insert(map[string]interface{}{
			"amount": balance.Amount,
			"date":   balance.Date,
			"author": "PUT AUTHOR ID HERE",
		}, true, "", "", "").
		Execute()

	if err != nil {
		return models.Balance{}, err
	}

	return balance, nil
}

func (repo *supabaseRepository) ListBalances(ctx context.Context, from time.Time, to time.Time) ([]models.Balance, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	rows, _, err := repo.db.From("balances").
		Select("amount, date", "", false).
		Filter("date", ">=", from.Format("2006-01-01")).
		Filter("date", "<=", to.Format("2006-01-01")).
		Filter("author", "=", "PUT AUTHOR ID HERE").
		Order("date", &postgrest.OrderOpts{
			Ascending: true,
		}).
		Execute()

	if err != nil {
		return make([]models.Balance, 0), err
	}

	var balances []models.Balance
	err = json.Unmarshal(rows, &balances)
	if err != nil {
		panic(err)
	}

	return balances, nil
}

func (repo *supabaseRepository) DeleteBalance(ctx context.Context, date string) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	_, count, err := repo.db.From("balances").
		Delete("", "").
		Eq("date", date).
		Eq("author", "PUT AUTHOR ID HERE").
		Execute()

	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("balance not found")
	}

	return nil
}
