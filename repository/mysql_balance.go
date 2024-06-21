package repository

import (
	"cashflow/models"
	"context"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func (repo *localDatabase) InsertBalance(amount float64, date string) (models.Balance, error) {
	_, err := repo.db.Exec(fmt.Sprintf(`INSERT INTO balances (amount, date) 
										VALUES (%.2f, "%s") 
										ON DUPLICATE KEY UPDATE
											date="%s", 
											amount=%.2f;`, amount, date, date, amount))
	if err != nil {
		return models.Balance{}, err
	}

	balance := models.Balance{
		Amount: amount,
		Date:   date,
	}

	return balance, nil
}

func (repo *localDatabase) ListBalances(from time.Time, to time.Time) ([]models.Balance, error) {
	rows, err := repo.db.Query("SELECT * FROM balances;") // TODO filter with from and to + get the latest balance before from
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

		balances = append(balances, balance)
	}

	if err = rows.Err(); err != nil {
		return make([]models.Balance, 0), err
	}

	return balances, nil
}

func (repo *localDatabase) DeleteBalance(ctx context.Context, date string) error {
	_, err := repo.db.Exec(fmt.Sprintf(`DELETE FROM balances 
										WHERE date = "%s";`, date))
	return err
}
