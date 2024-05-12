package repository

import (
	"cashflow/models"
	"fmt"
	"sync"
)

type store struct {
	transactionsRepository models.TransactionRepository
	balancesRepository     models.BalanceRepository
}

var singleInstance *store
var lock = &sync.Mutex{}

func Init(storageType string) *store {
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
		singleInstance = &store{
			transactionsRepository: &inMemoryTransactionDatabase{},
			balancesRepository:     &inMemoryBalanceDatabase{},
		}
	case "Mocked":
		singleInstance = &store{transactionsRepository: &mockDb{}}
	case "Local":
		singleInstance = &store{
			transactionsRepository: &localDatabase{},
			balancesRepository:     &localDatabase{},
		}
	default:
		panic(fmt.Sprintln("Did not create an instance"))
	}

	singleInstance.transactionsRepository.Init()
	singleInstance.balancesRepository.Init()

	return singleInstance
}

func Close() {
	singleInstance.transactionsRepository.Close()
	singleInstance.balancesRepository.Close()
}

func GetBalanceRepo() *models.BalanceRepository {
	return &singleInstance.balancesRepository
}

func GetTransactionRepo() *models.TransactionRepository {
	return &singleInstance.transactionsRepository
}

// func (s *store) ListTransactions() []models.Transaction {
// 	lock.Lock()
// 	defer lock.Unlock()

// 	unorderedTransactions := s.transactionsRepository.ListTransactions()
// 	return utils.SortByDate(unorderedTransactions)
// }

// func (s *store) InsertTransaction(transaction models.Transaction) models.Transaction {
// 	lock.Lock()
// 	defer lock.Unlock()

// 	return s.transactionsRepository.InsertTransaction(transaction)
// }

// func (s *store) InsertBalance(amount float64, date string) models.Balance {
// 	lock.Lock()
// 	defer lock.Unlock()

// 	return s.transactionsRepository.InsertBalance(amount, date)
// }

// func (s *store) ListBalances() []models.Balance {
// 	lock.Lock()
// 	defer lock.Unlock()

// 	unorderedBalances := s.transactionsRepository.ListBalances()
// 	return utils.SortByDate(unorderedBalances)
// }
