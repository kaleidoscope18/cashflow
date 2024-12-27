package repository

import (
	"cashflow/models"
	"cashflow/utils"
	"context"
	"database/sql"
	"fmt"
	"time"
)

func (repo *supabaseRepository) GetTransactionById(ctx context.Context, id string) (models.Transaction, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	query := "SELECT id, description, amount, date, recurrency FROM transactions WHERE id = ?"

	var transaction models.Transaction
	err := repo.db.QueryRow(context.Background(), query, id).Scan(
		&transaction.Id,
		&transaction.Description,
		&transaction.Amount,
		&transaction.Date,
		&transaction.Recurrency,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Transaction{}, fmt.Errorf("transaction not found")
		}
		return models.Transaction{}, fmt.Errorf("error querying transaction: %w", err)
	}

	transaction.Date = utils.ParseDate(transaction.Date)
	return transaction, nil
}

func (repo *supabaseRepository) ListTransactions(ctx context.Context, from time.Time, to time.Time) ([]models.Transaction, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	query := `SELECT id, description, amount, date, recurrency FROM transactions WHERE date BETWEEN ? AND ? ORDER BY date ASC AND author = (SELECT auth.uid())`
	rows, err := repo.db.Query(context.Background(), query, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		err = rows.Scan(&transaction.Id, &transaction.Description, &transaction.Amount, &transaction.Date, &transaction.Recurrency)
		if err != nil {
			return nil, err
		}
		transaction.Date = utils.ParseDate(transaction.Date)
		transactions = append(transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (repo *supabaseRepository) InsertTransaction(ctx context.Context, transaction models.Transaction) (string, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	var id string
	err := repo.db.QueryRow(context.Background(),
		`INSERT INTO transactions (amount, date, description, recurrency, author) 
		VALUES ($1, $2, $3, $4, (SELECT auth.uid()))
		RETURNING id`,
		transaction.Amount, transaction.Date, transaction.Description, transaction.Recurrency).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("error inserting transaction: %w", err)
	}

	return id, nil
}

func (repo *supabaseRepository) DeleteTransaction(ctx context.Context, id string) (string, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx, err := repo.db.Begin(context.Background())
	if err != nil {
		return id, err
	}

	query := `DELETE FROM transactions WHERE id = ?`
	result, err := tx.Exec(context.Background(), query, id)
	if err != nil {
		tx.Rollback(context.Background())
		return id, err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		tx.Rollback(context.Background())
		return id, fmt.Errorf("no rows affected, transaction with id %s not found", id)
	}
	if err != nil {
		tx.Rollback(context.Background())
		return id, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		tx.Rollback(context.Background())
		return id, err
	}

	return id, nil
}

func (repo *supabaseRepository) EditTransaction(ctx context.Context, transaction models.Transaction) (string, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx, err := repo.db.Begin(context.Background())
	if err != nil {
		return transaction.Id, err
	}
	_, err = tx.Exec(context.Background(), `UPDATE transactions
					SET description=?, amount=?, recurrency=?, date=?
					WHERE id = ?;`,
		transaction.Description, transaction.Amount, transaction.Recurrency, transaction.Date, transaction.Id)

	if err != nil {
		tx.Rollback(context.Background())
		return transaction.Id, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		tx.Rollback(context.Background())
		return transaction.Id, err
	}

	return transaction.Id, nil
}
