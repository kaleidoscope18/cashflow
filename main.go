package main

import (
	"cashflow/api"
	"cashflow/domain/transactions"
	"cashflow/models"
	"cashflow/reporting"
	"cashflow/repository"
)

func main() {
	storage := repository.New("Mocked")

	ts := transactions.New(storage)
	reporting.PrintCommandLine(ts)

	app := &models.App{
		TransactionService: ts,
	}

	defer api.Run(app)
}
