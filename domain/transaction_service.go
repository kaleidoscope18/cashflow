package domain

import (
	"cashflow/dev"
	"cashflow/models"
	"cashflow/utils"
	"context"
	"fmt"
	"time"
)

type transactionService struct {
	repository     *models.Repository
	balanceService *models.BalanceService
}

func NewTransactionService(repo *models.Repository, balanceService *models.BalanceService) models.TransactionService {
	s := new(transactionService)
	s.repository = repo
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
	return results, err
}

func (s *transactionService) WriteTransaction(ctx context.Context, date string, amount float64, description string, recurrency string) (string, error) {
	id, err := (*s.repository).InsertTransaction(ctx, models.Transaction{
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

func (s *transactionService) EditTransaction(ctx context.Context, editType models.TransactionEditType, edited models.TransactionEdit) (string, error) {
	transaction, err := (*s.repository).GetTransactionById(ctx, edited.Id)
	if err != nil {
		return "", fmt.Errorf("cannot edit transaction, transaction with id %s not found", edited.Id)
	}

	dev.PrintJson(transaction)

	switch editType {
	case models.All:
		editedTransaction := transaction
		if edited.Amount != nil {
			editedTransaction.Amount = *edited.Amount
		}
		if edited.Date != nil {
			editedTransaction.Date = *edited.Date
		}
		if edited.Recurrency != nil {
			editedTransaction.Recurrency = *edited.Recurrency
		}
		if edited.Description != nil {
			editedTransaction.Description = *edited.Description
		}
		dev.PrintJson(editedTransaction)
		return (*s.repository).EditTransaction(ctx, editedTransaction)
	case models.FromDate:
	case models.OnDateOnly:
	default:
	}

	return "", nil
}
