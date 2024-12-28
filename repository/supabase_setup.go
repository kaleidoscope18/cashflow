package repository

import (
	"fmt"
	"os"
	"sync"

	"github.com/supabase-community/supabase-go"
)

type supabaseRepository struct {
	mutex sync.RWMutex
	db    *supabase.Client
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

	client, err := supabase.NewClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_API_KEY"), &supabase.ClientOptions{})
	if err != nil {
		panic("cannot initalize client: " + err.Error())
	}

	repo.db = client

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
