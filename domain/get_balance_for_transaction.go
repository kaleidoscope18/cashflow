package domain

import (
	"cashflow/models"
	"cashflow/utils"
)

func GetBalanceForTransaction(transaction models.Transaction, previousTransaction models.ComputedTransaction, latestBalance models.Balance) float64 {
	if utils.IsDateBefore(previousTransaction.Date, latestBalance.Date) || previousTransaction.Date == latestBalance.Date {
		return latestBalance.Amount + transaction.Amount
	}
	return previousTransaction.Balance + transaction.Amount
}
