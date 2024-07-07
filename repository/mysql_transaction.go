package repository

import (
	"bytes"
	"cashflow/models"
	"cashflow/utils"
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

	var query bytes.Buffer
	data := map[string]interface{}{
		"from": from,
		"to":   to,
	}
	err := template.Must(template.New("ListTransactions").Parse(`
		SELECT * FROM transactions 
		WHERE date BETWEEN '{{.from}}' AND '{{.to}}
		ORDER BY date ASC';
	`)).Execute(&query, data)
	if err != nil {
		log.Printf("Failed to execute ListTransactions: %v", err)
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
		transaction.Date = utils.ParseDate(transaction.Date)
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
