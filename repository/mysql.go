package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type localDatabase struct {
	db *sql.DB
}

func (repo *localDatabase) Init() error {
	db, err := sql.Open("mysql", "root:new_password@tcp(127.0.0.1:3306)/cashflow")
	if err != nil {
		panic(err.Error())
	}
	repo.db = db

	// Ping the database to verify connection
	err = repo.db.Ping()
	if err != nil {
		return err
	}

	fmt.Println("Connected to MySQL database!")

	// Set connection pool options
	repo.db.SetMaxOpenConns(20)
	repo.db.SetMaxIdleConns(10)

	return nil
}

func (repo *localDatabase) Close() error {
	fmt.Println("Closing the mysql DB connection")
	return repo.db.Close()
}

func (repo *localDatabase) CheckConnection() error {
	return repo.db.Ping()
}
