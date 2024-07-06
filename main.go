package main

import (
	"cashflow/api"
	"cashflow/domain"
	"cashflow/models"
	"cashflow/repository"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = repository.Init(models.Local)
	if err != nil {
		panic("Could not initiate the app : " + err.Error())
	}
	defer repository.Close()

	repo := repository.Get()
	bs := domain.NewBalanceService(repo)
	ts := domain.NewTransactionService(repo, &bs)

	app := &models.App{
		TransactionService: &ts,
		BalanceService:     &bs,
	}

	defer api.Run(app)
}
