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

	tr, br, err := repository.Init(models.Local)
	if err != nil {
		panic("Could not initiate the app : " + err.Error())
	}
	defer repository.Close()

	bs := domain.NewBalanceService(br)
	ts := domain.NewTransactionService(tr, &bs)

	app := &models.App{
		TransactionService: &ts,
		BalanceService:     &bs,
	}

	defer api.Run(app)
}
