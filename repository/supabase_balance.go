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

func (repo *supabaseRepository) InsertBalance(amount float64, date string) (models.Balance, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	_, err := repo.db.Exec(context.Background(), fmt.Sprintf(`INSERT INTO balances (amount, date) 
										VALUES (%.2f, '%s')
										ON CONFLICT (date) DO UPDATE SET
											amount = %.2f;`, amount, date, amount))
	if err != nil {
		return models.Balance{}, err
	}

	balance := models.Balance{
		Amount: amount,
		Date:   date,
	}

	return balance, nil
}

func (repo *supabaseRepository) ListBalances(from time.Time, to time.Time) ([]models.Balance, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	var query bytes.Buffer
	data := map[string]interface{}{
		"from": from,
		"to":   to,
	}
	err := template.Must(template.New("ListBalances").Parse(`
		SELECT * FROM balances 
		WHERE date BETWEEN '{{.from}}' AND '{{.to}}
		ORDER BY date ASC';
	`)).Execute(&query, data)

	if err != nil {
		log.Printf("Failed to execute ListBalances: %v", err)
		return nil, err
	}

	rows, err := repo.db.Query(context.Background(), query.String())
	if err != nil {
		return make([]models.Balance, 0), err
	}
	defer rows.Close()

	var balances []models.Balance
	for rows.Next() {
		var balance models.Balance
		err = rows.Scan(&balance.Amount, &balance.Date)
		if err != nil {
			return make([]models.Balance, 0), err
		}

		balance.Date = utils.ParseDate(balance.Date)
		balances = append(balances, balance)
	}

	if err = rows.Err(); err != nil {
		return make([]models.Balance, 0), err
	}

	return balances, nil
}

func (repo *supabaseRepository) DeleteBalance(ctx context.Context, date string) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	result, err := repo.db.Exec(context.Background(), fmt.Sprintf(`DELETE FROM balances 
										WHERE date = "%s";`, date))
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("balance not found")
	}

	return nil
}
