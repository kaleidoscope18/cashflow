package repository

import (
	"cashflow/models"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type localDatabase struct {
	db *sql.DB
}

func (repo *localDatabase) Init() {
	db, err := sql.Open("mysql", "root:new_password@tcp(127.0.0.1:3306)/cashflow")
	if err != nil {
		panic(err.Error())
	}
	repo.db = db

	// Ping the database to verify connection
	err = repo.db.Ping()
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Connected to MySQL database!")

	// Set connection pool options
	repo.db.SetMaxOpenConns(20)
	repo.db.SetMaxIdleConns(10)
}

func (repo *localDatabase) Close() {
	fmt.Println("Closing the mysql DB connection")
	repo.db.Close()
}

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
