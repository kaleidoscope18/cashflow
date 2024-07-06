package repository

import (
	"cashflow/models"
	"fmt"
	"sync"
)

type singletons struct {
	transactionsRepository models.TransactionRepository
	balancesRepository     models.BalanceRepository
}

var singleInstance *singletons
var lock = &sync.Mutex{}

func Init(storageType models.StorageStrategy) error {
	lock.Lock()
	defer lock.Unlock()

	if singleInstance != nil {
		fmt.Println("Repositories already created, do not call repository.Init() again")
		return nil
	}

	switch storageType {
	case models.InMemory:
		singleInstance = &singletons{
			transactionsRepository: &inMemoryTransactionDatabase{},
			balancesRepository:     &inMemoryBalanceDatabase{},
		}
	case models.Local:
		singleInstance = &singletons{
			transactionsRepository: &mysqlDatabase{},
			balancesRepository:     &mysqlDatabase{},
		}
	default:
		panic(fmt.Sprintln("Did not create repositories, can not run the app"))
	}

	err := singleInstance.transactionsRepository.Init()
	if err != nil {
		return err
	}
	err = singleInstance.balancesRepository.Init()
	if err != nil {
		return Close()
	}

	return nil
}

func Close() error {
	err := singleInstance.transactionsRepository.Close()
	if err != nil {
		return err
	}
	err = singleInstance.balancesRepository.Close()
	if err != nil {
		return err
	}
	return nil
}

func GetRepos() (*models.TransactionRepository, *models.BalanceRepository) {
	return &singleInstance.transactionsRepository,
		&singleInstance.balancesRepository
}
