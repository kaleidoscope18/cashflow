package repository

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type mysqlRepository struct {
	db    *sql.DB
	mutex sync.RWMutex
}

var initDatabaseOnce sync.Once

func (repo *mysqlRepository) Init() error {
	var err error
	initDatabaseOnce.Do(func() {
		err = open(repo)
	})
	return err
}

func open(repo *mysqlRepository) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

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

	fmt.Println("Connected to MySQL database")
	return nil
}

func (repo *mysqlRepository) Close() error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	fmt.Println("Closing the MySQL DB connection")
	return repo.db.Close()
}

func (repo *mysqlRepository) Health() error {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	maxRetries := 3
	retryDelay := time.Second * 2

	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		err := repo.db.Ping()
		if err == nil {
			return nil
		}

		lastErr = err
		if attempt < maxRetries {
			time.Sleep(retryDelay)
		}
	}

	return fmt.Errorf("failed to ping database after %d attempts: %w", maxRetries, lastErr)
}
