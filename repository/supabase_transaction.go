package repository

import (
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

	row, err := repo.db.Query(ctx, "SELECT id, description, amount, date::text, recurrency FROM transactions WHERE id = $1 AND author = $2", id, utils.GetContextValue(ctx, models.UserId))

	if err != nil {
		log.Printf("error querying transaction: %v", err)
		return models.Transaction{}, fmt.Errorf("error querying transaction: %w", err)
	}

	if !row.Next() {
		return models.Transaction{}, fmt.Errorf("transaction not found")
	}

	var transaction models.Transaction
	err = row.Scan(&transaction.Id, &transaction.Description, &transaction.Amount, &transaction.Date, &transaction.Recurrency)
	if err != nil {
		log.Printf("error scanning transaction: %v", err)
		return models.Transaction{}, fmt.Errorf("error scanning transaction: %w", err)
	}

	transaction.Date = utils.ParseDate(transaction.Date)
	return transaction, nil
}

func (repo *supabaseRepository) ListTransactions(ctx context.Context, from time.Time, to time.Time) ([]models.Transaction, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	rows, err := repo.db.Query(ctx, "SELECT id, description, amount, date::text, recurrency FROM transactions WHERE date >= $1 AND date <= $2 AND author = $3 ORDER BY date ASC", from.Format("2006-01-01"), to.Format("2006-01-01"), utils.GetContextValue(ctx, models.UserId))
	if err != nil {
		log.Printf("error listing transactions: %v", err)
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

	row := repo.db.QueryRow(ctx, "INSERT INTO transactions (description, amount, date, recurrency, author) VALUES ($1, $2, $3, $4, $5) RETURNING id", transaction.Description, transaction.Amount, transaction.Date, transaction.Recurrency, utils.GetContextValue(ctx, models.UserId))
	var id string
	err := row.Scan(&id)
	if err != nil {
		log.Printf("error inserting transaction: %v", err)
		return "", fmt.Errorf("error inserting transaction: %w", err)
	}

	return id, nil
}

func (repo *supabaseRepository) DeleteTransaction(ctx context.Context, id string) (string, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	_, err := repo.db.Exec(ctx, "DELETE FROM transactions WHERE id = $1 AND author = $2", id, utils.GetContextValue(ctx, models.UserId))

	if err != nil {
		log.Printf("error deleting transaction: %v", err)
		return id, err
	}

	return id, nil
}

func (repo *supabaseRepository) EditTransaction(ctx context.Context, transaction models.Transaction) (string, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	_, err := repo.db.Exec(ctx, "UPDATE transactions SET description = $1, amount = $2, recurrency = $3, date = $4 WHERE id = $5 AND author = $6", transaction.Description, transaction.Amount, transaction.Recurrency, transaction.Date, transaction.Id, utils.GetContextValue(ctx, models.UserId))

	if err != nil {
		log.Printf("error editing transaction: %v", err)
		return transaction.Id, err
	}

	return transaction.Id, nil
}
