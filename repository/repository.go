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

func Init(storageType models.StorageStrategy) (*models.TransactionRepository, *models.BalanceRepository) {
	lock.Lock()
	defer lock.Unlock()

	if singleInstance != nil {
		fmt.Println("Repositories already created, do not call repository.Init() again")
		return getRepos()
	}

	switch storageType {
	case "InMemory":
		singleInstance = &store{
			transactionsRepository: &inMemoryTransactionDatabase{},
			balancesRepository:     &inMemoryBalanceDatabase{},
		}
	case "Mocked":
		singleInstance = &store{
			transactionsRepository: &mockTransactionDb{},
			balancesRepository:     &mockBalanceDb{},
		}
	case "Local":
		singleInstance = &store{
			transactionsRepository: &localDatabase{},
			balancesRepository:     &localDatabase{},
		}
	default:
		panic(fmt.Sprintln("Did not create repositories"))
	}

	singleInstance.transactionsRepository.Init()
	singleInstance.balancesRepository.Init()

	return getRepos()
}

func Close() {
	singleInstance.transactionsRepository.Close()
	singleInstance.balancesRepository.Close()
}

func getRepos() (*models.TransactionRepository, *models.BalanceRepository) {
	return &singleInstance.transactionsRepository,
		&singleInstance.balancesRepository
}
