package repository

import (
	"cashflow/models"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func (repo *localDatabase) InsertBalance(amount float64, date string) models.Balance {
	_, err := repo.db.Exec(fmt.Sprintf(`INSERT INTO balances (amount, date) 
										VALUES (%.2f, "%s") 
										ON DUPLICATE KEY UPDATE
											date="%s", 
											amount=%.2f;`, amount, date, date, amount))
	if err != nil {
		panic(err)
	}

	balance := models.Balance{
		Amount: amount,
		Date:   date,
	}

	return balance
}

func (repo *localDatabase) ListBalances() []models.Balance {
	rows, err := repo.db.Query("SELECT * FROM balances;")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var balances []models.Balance
	for rows.Next() {
		var balance models.Balance
		err = rows.Scan(&balance.Amount, &balance.Date)
		if err != nil {
			panic(err.Error())
		}

		balances = append(balances, balance)
	}

	if err = rows.Err(); err != nil {
		panic(err.Error())
	}

	return balances
}
