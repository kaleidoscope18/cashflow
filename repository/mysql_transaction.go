package repository

import (
	"cashflow/models"
	"cashflow/utils"
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

func (repo *mysqlRepository) GetTransactionById(ctx context.Context, id string) (models.Transaction, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	query := "SELECT id, description, amount, date, recurrency FROM transactions WHERE id = ?"

	var transaction models.Transaction
	err := repo.db.QueryRowContext(ctx, query, id).Scan(
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

func (repo *mysqlRepository) ListTransactions(ctx context.Context, from time.Time, to time.Time) ([]models.Transaction, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	query := `SELECT id, description, amount, date, recurrency FROM transactions WHERE date BETWEEN ? AND ? ORDER BY date ASC`
	rows, err := repo.db.QueryContext(ctx, query, from, to)
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

func (repo *mysqlRepository) InsertTransaction(ctx context.Context, transaction models.Transaction) (string, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	stmt, err := repo.db.PrepareContext(ctx, `INSERT INTO transactions (amount, date, description, recurrency) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return "", err
	}
	result, err := stmt.ExecContext(ctx, transaction.Amount, transaction.Date, transaction.Description, transaction.Recurrency)
	if err != nil {
		return "", fmt.Errorf("error inserting transaction: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return "", fmt.Errorf("error getting last insert ID: %w", err)
	}
	return strconv.FormatInt(id, 10), nil
}

func (repo *mysqlRepository) DeleteTransaction(ctx context.Context, id string) (string, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return id, err
	}

	query := `DELETE FROM transactions WHERE id = ?`
	result, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		tx.Rollback()
		return id, err
	}

	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		tx.Rollback()
		return id, fmt.Errorf("no rows affected, transaction with id %s not found", id)
	}
	if err != nil {
		tx.Rollback()
		return id, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return id, err
	}

	return id, nil
}

func (repo *mysqlRepository) EditTransaction(ctx context.Context, transaction models.Transaction) (string, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx, err := repo.db.Begin()
	if err != nil {
		return transaction.Id, err
	}
	_, err = tx.Exec(`UPDATE transactions
					SET description=?, amount=?, recurrency=?, date=?
					WHERE id = ?;`,
		transaction.Description, transaction.Amount, transaction.Recurrency, transaction.Date, transaction.Id)

	if err != nil {
		tx.Rollback()
		return transaction.Id, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return transaction.Id, err
	}

	return transaction.Id, nil
}
