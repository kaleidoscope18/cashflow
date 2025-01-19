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

func (repo *supabaseRepository) ReadTransactionById(ctx context.Context, id string, account string) (models.Transaction, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	author := utils.GetContextValue(ctx, models.UserId)
	row := repo.db.QueryRow(ctx, "SELECT id, description, amount, date::text, recurrency, account FROM cashflow.transactions WHERE id = $1 AND author = $2 AND account = $3", id, author, account)
	var transaction models.Transaction
	err := row.Scan(&transaction.Id, &transaction.Description, &transaction.Amount, &transaction.Date, &transaction.Recurrency, &transaction.Account)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.Transaction{}, fmt.Errorf("transaction not found")
		}
		log.Printf("error SELECT on cashflow.transactions: %v", err)
		return models.Transaction{}, err
	}

	transaction.Date = utils.ParseDate(transaction.Date)
	return transaction, nil
}

func (repo *supabaseRepository) ReadTransactions(ctx context.Context, from time.Time, to time.Time, account string) ([]models.Transaction, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	rows, err := repo.db.Query(ctx, "SELECT id, description, amount, date::text, recurrency, account FROM cashflow.transactions WHERE date >= $1 AND date <= $2 AND author = $3 AND account = $4 ORDER BY date ASC", from.Format(utils.DateLayout), to.Format(utils.DateLayout), utils.GetContextValue(ctx, models.UserId), account)
	if err != nil {
		log.Printf("error SELECT on cashflow.transactions: %v", err)
		return nil, err
	}
	defer rows.Close()

	transactions, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Transaction])
	if err != nil {
		log.Printf("error collecting rows: %v", err)
		return nil, err
	}

	return transactions, nil
}

func (repo *supabaseRepository) CreateTransaction(ctx context.Context, transaction models.Transaction) (string, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	row := repo.db.QueryRow(ctx, "INSERT INTO cashflow.transactions (description, amount, date, recurrency, author, account) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", transaction.Description, transaction.Amount, transaction.Date, transaction.Recurrency, utils.GetContextValue(ctx, models.UserId), transaction.Account)
	var id string
	err := row.Scan(&id)
	if err != nil {
		log.Printf("error INSERT on cashflow.transactions: %v", err)
		return "", err
	}

	return id, nil
}

func (repo *supabaseRepository) DeleteTransaction(ctx context.Context, id string, account string) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	_, err := repo.db.Exec(ctx, "DELETE FROM cashflow.transactions WHERE id = $1 AND author = $2 AND account = $3", id, utils.GetContextValue(ctx, models.UserId), account)

	if err != nil {
		log.Printf("error DELETE on cashflow.transactions: %v", err)
		return err
	}

	return nil
}

func (repo *supabaseRepository) UpdateTransaction(ctx context.Context, transaction models.Transaction) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	_, err := repo.db.Exec(ctx, "UPDATE cashflow.transactions SET description = $1, amount = $2, recurrency = $3, date = $4 WHERE id = $5 AND author = $6 AND account = $7", transaction.Description, transaction.Amount, transaction.Recurrency, transaction.Date, transaction.Id, utils.GetContextValue(ctx, models.UserId), transaction.Account)

	if err != nil {
		log.Printf("error UPDATE on cashflow.transactions: %v", err)
		return err
	}

	return nil
}
