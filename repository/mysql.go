package repository

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type localDatabase struct {
	db *sql.DB
}

func (repo *localDatabase) Init() error {
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_USER_PASSWORD")
	databaseName := os.Getenv("MYSQL_DATABASE_NAME")
	host := os.Getenv("DB_HOST")

	db, err := sql.Open("mysql", fmt.Sprintf(`%s:%s@tcp(%s)/%s`, user, password, host, databaseName))
	if err != nil {
		panic(err.Error())
	}
	repo.db = db

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
