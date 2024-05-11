package repository

import (
	"cashflow/models"
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type localDatabase struct {
	db *sql.DB
}

func (repo *localDatabase) init() {
	// Open a connection to the MySQL database
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

func (repo *localDatabase) close() {
	fmt.Println("Closing the mysql DB connection")
	repo.db.Close()
}

func (repo *localDatabase) ListTransactions() []models.Transaction {
	fmt.Println("Querying the DB with SELECT * FROM transactions;")
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

		js, _ := json.MarshalIndent(transaction, "", " ")
		fmt.Println(string(js))

		transactions = append(transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		panic(err.Error())
	}

	return transactions
}

func (repo *localDatabase) InsertTransaction(transaction models.Transaction) models.Transaction {
	panic(fmt.Errorf("not implemented"))

}

func (repo *localDatabase) InsertBalance(amount float64, date string) models.Balance {
	panic(fmt.Errorf("not implemented"))

}

func (repo *localDatabase) ListBalances() []models.Balance {
	panic(fmt.Errorf("not implemented"))
}
