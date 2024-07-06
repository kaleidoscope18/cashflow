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

var (
	initReposOnce  sync.Once
	singleInstance *singletons
)

func Init(storageType models.StorageStrategy) error {
	var err error
	initReposOnce.Do(func() {
		err = createRepos(storageType)
	})
	return err
}

func createRepos(storageType models.StorageStrategy) error {
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
