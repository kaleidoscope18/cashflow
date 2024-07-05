package domain

import (
	"cashflow/models"
	"cashflow/utils"
	"context"
	"time"
)

type transactionService struct {
	repository     *models.TransactionRepository
	balanceService *models.BalanceService
}

func NewTransactionService(transactionRepository *models.TransactionRepository, balanceService *models.BalanceService) models.TransactionService {
	s := new(transactionService)
	s.repository = transactionRepository
	s.balanceService = balanceService
	return s
}

func (s *transactionService) ListTransactions(ctx context.Context, from time.Time, to time.Time) ([]models.ComputedTransaction, error) {
	transactions, err := (*s.repository).ListTransactions(ctx, from, to)
	if err != nil {
		return nil, err
	}

	balances, err := (*s.balanceService).ListBalances(from, to)
	if err != nil {
		return nil, err
	}

	results, err := listTransactions(transactions, balances, from, to)
	return utils.SortByDate(results), err
}

func (s *transactionService) WriteTransaction(date string, amount float64, description string, recurrency string) (string, error) {
	id, err := (*s.repository).InsertTransaction(models.Transaction{
		Amount:      utils.RoundToTwoDigits(amount),
		Date:        utils.ParseDate(date),
		Description: description,
		Recurrency:  recurrency,
	})

	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *transactionService) DeleteTransaction(ctx context.Context, id string) (string, error) {
	return (*s.repository).DeleteTransaction(ctx, id)
}
