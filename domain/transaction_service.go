package domain

import (
	"cashflow/domain/recurrency"
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

	balances, err := (*s.balanceService).ListBalances(ctx, from, to)
	if err != nil {
		return nil, err
	}

	results, err := listTransactions(transactions, balances, from, to)
	return results, err
}

func (s *transactionService) WriteTransaction(ctx context.Context, date string, amount float64, description string, recurrency string) error {
	_, err := (*s.repository).InsertTransaction(ctx, models.Transaction{
		Amount:      utils.RoundToTwoDigits(amount),
		Date:        utils.ParseDate(date),
		Description: description,
		Recurrency:  recurrency,
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *transactionService) DeleteTransaction(ctx context.Context, id string) (string, error) {
	return (*s.repository).DeleteTransaction(ctx, id)
}

func (s *transactionService) EditTransaction(ctx context.Context, editType models.TransactionEditType, edited models.TransactionEdit) (string, error) {
	transaction, err := (*s.repository).GetTransactionById(ctx, edited.Id)
	if err != nil {
		return "", fmt.Errorf("cannot edit transaction, transaction with id %s not found", edited.Id)
	}

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
		return (*s.repository).EditTransaction(ctx, editedTransaction)
	case models.FromDate:
		newRecurrency, err := recurrency.ChangeRecurrencyEnd(transaction.Recurrency, *edited.Date)
		if err != nil {
			return "", err
		}

		newT := models.Transaction{
			Id:          transaction.Id,
			Amount:      transaction.Amount,
			Date:        *edited.Date,
			Description: transaction.Description,
			Recurrency:  transaction.Recurrency,
		}
		if edited.Amount != nil {
			newT.Amount = *edited.Amount
		}
		if edited.Description != nil && *edited.Description != "" {
			newT.Description = *edited.Description
		}
		if edited.Recurrency != nil && *edited.Recurrency != "" {
			newT.Recurrency = *edited.Recurrency
		}
		newTId, err := (*s.repository).InsertTransaction(ctx, newT)
		if err != nil {
			return "", err
		}

		_, err = (*s.repository).EditTransaction(ctx, models.Transaction{
			Id:          transaction.Id,
			Recurrency:  newRecurrency,
			Description: transaction.Description,
			Amount:      transaction.Amount,
			Date:        transaction.Date,
		})
		if err != nil {
			return (*s.repository).DeleteTransaction(ctx, newTId)
		}
		return newTId, nil
	case models.OnDateOnly:
	default:
	}

	return "", nil
}
