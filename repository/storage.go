package repository

import (
	"cashflow/models"
	"cashflow/utils"
	"fmt"
	"sync"
)

type store struct {
	db database
}

type database interface {
	init()
	models.Repository
}

var singleInstance *store
var lock = &sync.Mutex{}

func New(storageType string) models.Repository {
	if singleInstance != nil {
		fmt.Println("Store instance already created")
		return singleInstance
	}

	lock.Lock()
	defer lock.Unlock()

	if singleInstance != nil {
		fmt.Println("Database already created")
		return singleInstance
	}

	switch storageType {
	case "InMemory":
		singleInstance = &store{db: &inMemoryDatabase{}}
		singleInstance.db.init()
	case "Mocked":
		singleInstance = &store{db: &mockDb{}}
	default:
		fmt.Println("Did not create an instance")
	}

	return singleInstance
}

func (s *store) ListTransactions() []models.Transaction {
	lock.Lock()
	defer lock.Unlock()

	unorderedTransactions := s.db.ListTransactions()
	return utils.SortByDate(unorderedTransactions)
}

func (s *store) InsertTransaction(transaction models.Transaction) models.Transaction {
	lock.Lock()
	defer lock.Unlock()

	return s.db.InsertTransaction(transaction)
}

func (s *store) InsertBalance(amount float64, date string) models.Balance {
	lock.Lock()
	defer lock.Unlock()

	return s.db.InsertBalance(amount, date)
}

func (s *store) ListBalances() []models.Balance {
	lock.Lock()
	defer lock.Unlock()

	unorderedBalances := s.db.ListBalances()
	return utils.SortByDate(unorderedBalances)
}
