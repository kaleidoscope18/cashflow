package main

import (
	"cashflow/api"
	"cashflow/domain"
	"cashflow/models"
	"cashflow/repository"
)

func main() {
	tr, br, err := repository.Init(models.InMemory)
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
