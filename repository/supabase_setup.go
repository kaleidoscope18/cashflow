package repository

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
)

type supabaseRepository struct {
	mutex sync.RWMutex
	db    *pgx.Conn
}

var initSupabaseDatabaseOnce sync.Once

func (repo *supabaseRepository) Init() error {
	var err error
	initSupabaseDatabaseOnce.Do(func() {
		err = openSupabase(repo)
	})
	return err
}

func openSupabase(repo *supabaseRepository) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	connStr := os.Getenv("SUPABASE_CONNECTION_STRING")
	if connStr == "" {
		return fmt.Errorf("SUPABASE_CONNECTION_STRING environment variable is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to Supabase: %w", err)
	}
	repo.db = conn

	if err := repo.db.Ping(ctx); err != nil {
		repo.db.Close(context.Background())
		return fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("Connected to Supabase database")
	return nil
}

func (repo *supabaseRepository) Close() error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	if repo.db == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("Closing the Supabase DB connection")
	return repo.db.Close(ctx)
}

func (repo *supabaseRepository) Health() error {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	if repo.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	maxRetries := 3
	retryDelay := time.Second * 2

	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := repo.db.Ping(ctx)
		cancel()

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
