package domain

import (
	"cashflow/models"
	"cashflow/repository"
)

func setup() models.TransactionService {
	tr, br := repository.Init(models.Local)
	defer repository.Close()

	bs := NewBalanceService(br)
	ts := NewTransactionService(tr, &bs)
	return ts
}
