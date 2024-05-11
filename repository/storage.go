package repository

import (
	"cashflow/models"
	"cashflow/utils"
	"fmt"
	"sync"
)

type store struct {
	repository repo
}

type repo interface {
	init()
	close()
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
		singleInstance = &store{repository: &inMemoryDatabase{}}
	case "Mocked":
		singleInstance = &store{repository: &mockDb{}}
	case "Local":
		singleInstance = &store{repository: &localDatabase{}}
	default:
		panic(fmt.Sprintln("Did not create an instance"))
	}

	singleInstance.repository.init()
	return singleInstance
}

func Close() {
	singleInstance.repository.close()
}

func (s *store) ListTransactions() []models.Transaction {
	lock.Lock()
	defer lock.Unlock()

	unorderedTransactions := s.repository.ListTransactions()
	return utils.SortByDate(unorderedTransactions)
}

func (s *store) InsertTransaction(transaction models.Transaction) models.Transaction {
	lock.Lock()
	defer lock.Unlock()

	return s.repository.InsertTransaction(transaction)
}

func (s *store) InsertBalance(amount float64, date string) models.Balance {
	lock.Lock()
	defer lock.Unlock()

	return s.repository.InsertBalance(amount, date)
}

func (s *store) ListBalances() []models.Balance {
	lock.Lock()
	defer lock.Unlock()

	unorderedBalances := s.repository.ListBalances()
	return utils.SortByDate(unorderedBalances)
}
