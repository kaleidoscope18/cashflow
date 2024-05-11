package main

import (
	"cashflow/api"
	"cashflow/domain/transactions"
	"cashflow/models"
	"cashflow/repository"
)

func main() {
	storage := repository.New("Local")
	defer repository.Close()

	ts := transactions.New(storage)

	app := &models.App{
		TransactionService: ts,
	}

	defer api.Run(app)
}
