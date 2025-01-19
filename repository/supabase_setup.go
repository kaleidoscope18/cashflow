package repository

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

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
		if err == nil {
			// Start pool monitoring
			go repo.monitorConnectionPool()
		}
	})
	return err
}

func (repo *supabaseRepository) Close() error {
	// Wait for active queries to finish (optional)
	stats := repo.db.Stat()
	if stats.AcquiredConns() > 0 {
		log.Printf("Waiting for %d active connections to finish...", stats.AcquiredConns())
	}

	repo.db.Close()
	log.Println("Closed connection to Supabase")
	return nil
}

func (repo *supabaseRepository) Health() error {
	err := repo.db.Ping(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func openSupabase(repo *supabaseRepository) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	config, err := pgxpool.ParseConfig(os.Getenv("SUPABASE_POSTGRES_CONNECTION_STRING"))
	if err != nil {
		log.Fatalf("Unable to parse connection string: %v", err)
	}

	config.MaxConns = 20 // Increased since prepared statements aren't cached
	config.MinConns = 5  // Keep some connections ready
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 15 * time.Minute
	config.HealthCheckPeriod = 1 * time.Minute

	// Required setting for supabase, for http limitations on Render.com
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return err
	}

	repo.db = pool

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = pool.Ping(ctx)
	if err != nil {
		return fmt.Errorf("could not connect to Supabase: %v", err)
	}

	return nil
}

func (repo *supabaseRepository) monitorConnectionPool() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		stats := repo.db.Stat()

		// Log if we're close to max connections
		if float64(stats.AcquiredConns())/float64(stats.MaxConns()) > 0.7 {
			log.Printf("WARNING: High connection usage - %.2f%% of max connections in use",
				float64(stats.AcquiredConns())/float64(stats.MaxConns())*100)
		}

		// Regular stats logging
		log.Printf("Pool Stats - Total: %d, Acquired: %d, Idle: %d, Max Conns: %d",
			stats.TotalConns(),
			stats.AcquiredConns(),
			stats.IdleConns(),
			stats.MaxConns(),
		)
	}
}
