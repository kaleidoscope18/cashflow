package repository

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type supabaseRepository struct {
	mutex sync.RWMutex
	db    *pgxpool.Pool
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

	config, err := pgxpool.ParseConfig(os.Getenv("SUPABASE_POSTGRES_CONNECTION_STRING"))
	if err != nil {
		log.Fatalf("Unable to parse connection string: %v", err)
	}

	config.MaxConns = 10
	// Disable prepared statements
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return err
	}

	repo.db = pool

	err = pool.Ping(context.Background())
	if err != nil {
		return fmt.Errorf("could not connect to Supabase: %v", err)
	}

	return nil
}

func (repo *supabaseRepository) Close() error {
	repo.db.Close()
	fmt.Println("Closed connection to Supabase")
	return nil
}

func (repo *supabaseRepository) Health() error {
	err := repo.db.Ping(context.Background())
	if err != nil {
		return err
	}
	return nil
}
