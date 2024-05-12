package repository

import (
	"cashflow/models"
	"fmt"
)

func (repo *localDatabase) ListTransactions() []models.Transaction {
	rows, err := repo.db.Query("SELECT * FROM transactions;")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		err = rows.Scan(&transaction.Id, &transaction.Description, &transaction.Amount, &transaction.Date)
		if err != nil {
			panic(err.Error())
		}
		transactions = append(transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		panic(err.Error())
	}

	return transactions
}

func (repo *localDatabase) InsertTransaction(transaction models.Transaction) models.Transaction {
	_, err := repo.db.Exec(fmt.Sprintf(`INSERT INTO transactions (amount, date, description) 
										VALUES (%.2f, "%s", "%s")`,
		transaction.Amount, transaction.Date, transaction.Description))
	if err != nil {
		panic(err)
	}

	return transaction
}
