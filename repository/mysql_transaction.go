package repository

import (
	"cashflow/models"
	"fmt"
	"time"
)

func (repo *localDatabase) ListTransactions(from time.Time, to time.Time) ([]models.Transaction, error) {
	rows, err := repo.db.Query("SELECT * FROM transactions;") // TODO filter with from and to + get the latest transaction before from date
	if err != nil {
		return make([]models.Transaction, 0), nil
	}

	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		err = rows.Scan(&transaction.Id, &transaction.Description, &transaction.Amount, &transaction.Date)
		if err != nil {
			return make([]models.Transaction, 0), nil
		}
		transactions = append(transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		return make([]models.Transaction, 0), nil
	}

	return transactions, nil
}

func (repo *localDatabase) InsertTransaction(transaction models.Transaction) (models.Transaction, error) {
	_, err := repo.db.Exec(fmt.Sprintf(`INSERT INTO transactions (amount, date, description) 
										VALUES (%.2f, "%s", "%s")`,
		transaction.Amount, transaction.Date, transaction.Description))
	if err != nil {
		return models.Transaction{}, nil
	}

	return transaction, nil
}
