package repository

import (
	"bytes"
	"cashflow/models"
	"context"
	"errors"
	"fmt"
	"log"
	"text/template"
	"time"
)

func (repo *mysqlRepository) ListTransactions(ctx context.Context, from time.Time, to time.Time) ([]models.Transaction, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	queryTemplate := template.Must(template.New("listTransactionsQueryTemplate").Parse(`
		SELECT * FROM transactions 
		WHERE date BETWEEN '{{.from}}' AND '{{.to}}
		ORDER BY date ASC';
	`))

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
		err = rows.Scan(&transaction.Id, &transaction.Description, &transaction.Amount, &transaction.Date, &transaction.Recurrency)
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

func (repo *mysqlRepository) InsertTransaction(transaction models.Transaction) (string, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	result, err := repo.db.Exec(fmt.Sprintf(`INSERT INTO transactions (amount, date, description, recurrency) 
										VALUES (%.2f, "%s", "%s", "%s")`,
		transaction.Amount, transaction.Date, transaction.Description, transaction.Recurrency))
	if err != nil {
		return "", err
	}

	id, err := result.LastInsertId()
	return fmt.Sprint(id), err
}

func (repo *mysqlRepository) DeleteTransaction(ctx context.Context, id string) (string, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	result, err := repo.db.Exec(fmt.Sprintf(`DELETE FROM transactions
											 WHERE id = "%s";`, id))
	if err != nil {
		return id, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return id, errors.New("transaction not found")
	}

	return id, nil
}
