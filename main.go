package main

import (
	"cashflow/api"
	"cashflow/domain"
	"cashflow/models"
	"cashflow/repository"
)

func main() {
	repository.Init("Local")
	defer repository.Close()

	bs := domain.NewBalanceService(repository.GetBalanceRepo())
	ts := domain.NewTransactionService(repository.GetTransactionRepo(), &bs)

	app := &models.App{
		TransactionService: &ts,
		BalanceService:     &bs,
	}

	defer api.Run(app)
}
