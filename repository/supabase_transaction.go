package repository

import (
	"cashflow/dev"
	"cashflow/models"
	"cashflow/utils"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

func (repo *supabaseRepository) GetTransactionById(ctx context.Context, id string) (models.Transaction, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	row, err := repo.db.Query(ctx, "SELECT id, description, amount, date, recurrency FROM transactions WHERE id = $1 AND author = $2", id, ctx.Value(models.UserId).(string))

	if err != nil {

		return models.Transaction{}, fmt.Errorf("error querying transaction: %w", err)
	}

	if !row.Next() {
		return models.Transaction{}, fmt.Errorf("transaction not found")
	}

	transaction, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.Transaction])
	if err != nil {
		return models.Transaction{}, fmt.Errorf("error scanning transaction: %w", err)
	}

	transaction.Date = utils.ParseDate(transaction.Date)
	return transaction, nil
}

func (repo *supabaseRepository) ListTransactions(ctx context.Context, from time.Time, to time.Time) ([]models.Transaction, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	rows, err := repo.db.Query(ctx, "SELECT id, description, amount, date, recurrency FROM transactions WHERE date >= $1 AND date <= $2 AND author = $3 ORDER BY date ASC", from.Format("2006-01-01"), to.Format("2006-01-01"), ctx.Value(models.UserId).(string))

	if err != nil {
		return nil, err
	}

	transactions, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Transaction])

	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (repo *supabaseRepository) InsertTransaction(ctx context.Context, transaction models.Transaction) (string, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	log.Printf("Inserting transaction for user %v", transaction)
	dev.PrintJson(map[string]interface{}{
		"description": transaction.Description,
		"amount":      transaction.Amount,
		"date":        transaction.Date,
		"recurrency":  transaction.Recurrency,
		"author":      ctx.Value(models.UserId).(string),
	})

	_, err := repo.db.Exec(ctx, "INSERT INTO transactions (description, amount, date, recurrency, author) VALUES ($1, $2, $3, $4, $5)", transaction.Description, transaction.Amount, transaction.Date, transaction.Recurrency, ctx.Value(models.UserId).(string))
	if err != nil {
		return "", fmt.Errorf("error inserting transaction: %w", err)
	}

	return transaction.Id, nil
}

func (repo *supabaseRepository) DeleteTransaction(ctx context.Context, id string) (string, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	_, err := repo.db.Exec(ctx, "DELETE FROM transactions WHERE id = $1 AND author = $2", id, ctx.Value(models.UserId).(string))

	if err != nil {
		return id, err
	}

	return id, nil
}

func (repo *supabaseRepository) EditTransaction(ctx context.Context, transaction models.Transaction) (string, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	_, err := repo.db.Exec(ctx, "UPDATE transactions SET description = $1, amount = $2, recurrency = $3, date = $4 WHERE id = $5 AND author = $6", transaction.Description, transaction.Amount, transaction.Recurrency, transaction.Date, transaction.Id, ctx.Value(models.UserId).(string))

	if err != nil {
		return transaction.Id, err
	}

	return transaction.Id, nil
}
