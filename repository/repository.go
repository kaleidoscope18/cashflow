package repository

import (
	"cashflow/models"
	"fmt"
	"sync"
)

type singleton struct {
	transactionsRepository models.TransactionRepository
	balancesRepository     models.BalanceRepository
}

var singleInstance *singleton
var lock = &sync.Mutex{}

func Init(storageType models.StorageStrategy) (*models.TransactionRepository, *models.BalanceRepository, error) {
	lock.Lock()
	defer lock.Unlock()

	if singleInstance != nil {
		fmt.Println("Repositories already created, do not call repository.Init() again")
		tr, br := getRepos()
		return tr, br, nil
	}

	switch storageType {
	case models.InMemory:
		singleInstance = &singleton{
			transactionsRepository: &inMemoryTransactionDatabase{},
			balancesRepository:     &inMemoryBalanceDatabase{},
		}
	case models.Local:
		singleInstance = &singleton{
			transactionsRepository: &localDatabase{},
			balancesRepository:     &localDatabase{},
		}
	default:
		panic(fmt.Sprintln("Did not create repositories, can not run the app"))
	}

	err := singleInstance.transactionsRepository.Init()
	if err != nil {
		return nil, nil, err
	}
	err = singleInstance.balancesRepository.Init()
	if err != nil {
		return nil, nil, Close()
	}

	tr, br := getRepos()
	return tr, br, nil
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

func getRepos() (*models.TransactionRepository, *models.BalanceRepository) {
	return &singleInstance.transactionsRepository,
		&singleInstance.balancesRepository
}
