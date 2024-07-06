package repository

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type mysqlDatabase struct {
	db *sql.DB
}

var (
	once  sync.Once
	mutex sync.Mutex
)

func open(repo *mysqlDatabase) error {
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_USER_PASSWORD")
	databaseName := os.Getenv("MYSQL_DATABASE_NAME")
	host := os.Getenv("DB_HOST")

	db, err := sql.Open("mysql", fmt.Sprintf(`%s:%s@tcp(%s)/%s`, user, password, host, databaseName))
	if err != nil {
		panic(err.Error())
	}
	repo.db = db

	// Set connection pool options
	repo.db.SetMaxOpenConns(20)
	repo.db.SetMaxIdleConns(10)
	repo.db.SetConnMaxLifetime(time.Minute * 5)

	// Verify the connection
	if err := repo.db.Ping(); err != nil {
		repo.db.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("Connected to MySQL database!")
	return nil
}

func (repo *mysqlDatabase) Init() error {
	var err error
	once.Do(func() {
		err = open(repo)
	})
	return err
}

func (repo *mysqlDatabase) Close() error {
	mutex.Lock()
	defer mutex.Unlock()
	fmt.Println("Closing the mysql DB connection")
	return repo.db.Close()
}
