package repository

import (
	"cashflow/models"
	"fmt"
	"sync"
)

var (
	initReposOnce sync.Once
	repository    models.Repository
)

func Init(storageType models.StorageStrategy) error {
	var err error
	initReposOnce.Do(func() {
		err = createRepository(storageType)
	})
	return err
}

func createRepository(storageType models.StorageStrategy) error {
	switch storageType {
	case models.InMemory:
		repository = &inMemoryRepository{}
	case models.Local:
		repository = &mysqlRepository{}
	default:
		panic(fmt.Sprintln("Did not create repository instance, can not run the app"))
	}

	err := repository.Init()
	if err != nil {
		return err
	}

	return nil
}

func Close() error {
	if repository == nil {
		return nil
	}
	return repository.Close()
}

func Get() *models.Repository {
	return &repository
}

func Health() error {
	return repository.Health()
}
