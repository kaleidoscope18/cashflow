package repository

import (
	"cashflow/models"
	"cashflow/utils"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/supabase-community/postgrest-go"
)

func (repo *supabaseRepository) GetTransactionById(ctx context.Context, id string) (models.Transaction, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	res, count, err := repo.db.From("transactions").
		Select("id, description, amount, date, recurrency", "", false).
		Eq("id", id).
		Eq("author", "Insert author id here").
		Single().
		Execute()

	if err != nil {

		return models.Transaction{}, fmt.Errorf("error querying transaction: %w", err)
	}

	if count == 0 {
		return models.Transaction{}, fmt.Errorf("transaction not found")
	}

	var transaction models.Transaction
	err = json.Unmarshal(res, &transaction)
	if err != nil {
		return models.Transaction{}, fmt.Errorf("error scanning transaction: %w", err)
	}

	transaction.Date = utils.ParseDate(transaction.Date)
	return transaction, nil
}

func (repo *supabaseRepository) ListTransactions(ctx context.Context, from time.Time, to time.Time) ([]models.Transaction, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	rows, _, err := repo.db.From("transactions").
		Select("id, description, amount, date, recurrency", "", false).
		Filter("date", ">=", from.Format("2006-01-01")).
		Filter("date", "<=", to.Format("2006-01-01")).
		Filter("author", "=", "PUT AUTHOR ID HERE").
		Order("date", &postgrest.OrderOpts{
			Ascending: true,
		}).
		Execute()

	if err != nil {
		return nil, err
	}

	var transactions []models.Transaction
	err = json.Unmarshal(rows, &transactions)

	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (repo *supabaseRepository) InsertTransaction(ctx context.Context, transaction models.Transaction) (string, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	var id string
	_, _, err := repo.db.From("transactions").
		Insert(map[string]interface{}{
			"amount":      transaction.Amount,
			"date":        transaction.Date,
			"description": transaction.Description,
			"recurrency":  transaction.Recurrency,
			"author":      "PUT AUTHOR ID HERE",
		}, true, "", "", "").
		Execute()

	if err != nil {
		return "", fmt.Errorf("error inserting transaction: %w", err)
	}

	return id, nil
}

func (repo *supabaseRepository) DeleteTransaction(ctx context.Context, id string) (string, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	_, _, err := repo.db.From("transactions").
		Delete("", "").
		Eq("id", id).
		Eq("author", "PUT AUTHOR ID HERE").
		Execute()

	if err != nil {
		return id, err
	}

	return id, nil
}

func (repo *supabaseRepository) EditTransaction(ctx context.Context, transaction models.Transaction) (string, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	_, _, err := repo.db.From("transactions").
		Update(map[string]interface{}{
			"description": transaction.Description,
			"amount":      transaction.Amount,
			"recurrency":  transaction.Recurrency,
			"date":        transaction.Date,
		}, "", "").
		Eq("id", transaction.Id).
		Eq("author", "PUT AUTHOR ID HERE").
		Execute()

	if err != nil {
		return transaction.Id, err
	}

	return transaction.Id, nil
}
