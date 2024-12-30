package repository

import (
	"context"
	"fmt"
	"os"
	"sync"

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

	pool, err := pgxpool.New(context.Background(), os.Getenv("SUPABASE_POSTGRES_CONNECTION_STRING"))
	if err != nil {
		return err
	}
	defer pool.Close()

	repo.db = pool

	fmt.Println("Connected to Supabase")
	return nil
}

func (repo *supabaseRepository) Close() error {
	fmt.Println("Close for supabase not implemented")
	return nil
}

func (repo *supabaseRepository) Health() error {
	fmt.Println("nothing")
	return nil
}
