package repository

import (
	"cashflow/models"
	"cashflow/utils"
	"context"
	"database/sql"
	"errors"
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

	const insertQuery = `INSERT INTO transactions (amount, date, description, recurrency) VALUES (?, ?, ?, ?)`
	result, err := repo.db.ExecContext(ctx, insertQuery, transaction.Amount, transaction.Date, transaction.Description, transaction.Recurrency)
	if err != nil {
		return "", err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(id, 10), nil
}

func (repo *mysqlRepository) DeleteTransaction(ctx context.Context, id string) (string, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	query := `DELETE FROM transactions WHERE id = ?`
	result, err := repo.db.ExecContext(ctx, query, id)

	if err != nil {
		return id, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return id, fmt.Errorf("error deleting transaction: %w", err)
	}
	if rowsAffected == 0 {
		return id, errors.New("transaction not found")
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
