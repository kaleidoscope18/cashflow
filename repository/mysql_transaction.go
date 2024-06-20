package repository

import (
	"bytes"
	"cashflow/models"
	"fmt"
	"log"
	"text/template"
	"time"
)

const listTransactionsQuery = `SELECT * FROM transactions 
    							WHERE date 
    							BETWEEN '{{.from}}' AND '{{.to}}';`

func (repo *localDatabase) ListTransactions(from time.Time, to time.Time) ([]models.Transaction, error) {
	queryTemplate := template.Must(template.New("listTransactionsQueryTemplate").Parse(listTransactionsQuery))

	var query bytes.Buffer
	data := map[string]interface{}{
		"from": from,
		"to":   to,
	}
	err := queryTemplate.Execute(&query, data)
	if err != nil {
		log.Printf("Failed to execute listTransactionsQueryTemplate: %v", err)
		return nil, err
	}

	rows, err := repo.db.Query(query.String())
	if err != nil {
		return make([]models.Transaction, 0), err
	}

	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		err = rows.Scan(&transaction.Id, &transaction.Description, &transaction.Amount, &transaction.Date)
		if err != nil {
			return make([]models.Transaction, 0), err
		}
		transactions = append(transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		return make([]models.Transaction, 0), err
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
