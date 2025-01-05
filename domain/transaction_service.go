package domain

import (
	"cashflow/domain/recurrency"
	"cashflow/models"
	"cashflow/utils"
	"context"
	"fmt"
	"log"
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

func (s *transactionService) DeleteRecurringTransaction(ctx context.Context, id string, date time.Time) (string, error) {
	existingT, err := (*s.repository).GetTransactionById(ctx, id)
	if err != nil {
		return "", fmt.Errorf("cannot delete transaction, transaction with id %s not found", id)
	}
	if existingT.Recurrency == "" {
		return "", fmt.Errorf("cannot delete transaction, transaction with id %s is not recurring", id)
	}

	/*
	     | ---------------- | ---------------- |
	   start             removed         end of recurrency
	                   on date only
	*/

	// TODO: validate that the given date is correct regarding the recurrency

	// Create new recurrent transaction to end until date - 1 day
	newRecurrency, err := recurrency.ChangeRecurrencyEnd(existingT.Recurrency, date.Format(utils.DateLayout))
	if err != nil {
		return "", err
	}
	newStartTId, err := (*s.repository).InsertTransaction(ctx, models.Transaction{
		Recurrency:  newRecurrency,
		Description: existingT.Description,
		Amount:      existingT.Amount,
		Date:        existingT.Date,
	})
	if err != nil {
		return "", err
	}

	// Create new recurrent transaction starting from date + 1 day
	newT := models.Transaction{
		Description: existingT.Description,
		Amount:      existingT.Amount,
		Date:        utils.ParseDateToTime(date.Format(utils.DateLayout)).AddDate(0, 0, 1).Format(utils.DateLayout),
		Recurrency:  existingT.Recurrency,
	}
	newEndTId, err := (*s.repository).InsertTransaction(ctx, newT)
	if err != nil {
		return (*s.repository).DeleteTransaction(ctx, newStartTId)
	}

	// Delete old transaction
	_, err = (*s.repository).DeleteTransaction(ctx, existingT.Id)
	if err != nil {
		_, errOne := (*s.repository).DeleteTransaction(ctx, newStartTId)
		_, errTwo := (*s.repository).DeleteTransaction(ctx, newEndTId)
		if errOne != nil || errTwo != nil {
			log.Printf("cannot rollback!! please fix manually, error: %s, %s", errOne, errTwo)
			return "", fmt.Errorf("cannot rollback, error: %s, %s", errOne, errTwo)
		}
		return "", err
	}
	return "", nil
}

func (s *transactionService) DeleteTransaction(ctx context.Context, id string) (string, error) {
	return (*s.repository).DeleteTransaction(ctx, id)
}

func (s *transactionService) EditTransaction(ctx context.Context, editType models.TransactionEditType, edited models.TransactionEdit) (string, error) {
	existingT, err := (*s.repository).GetTransactionById(ctx, edited.Id)
	if err != nil {
		return "", fmt.Errorf("cannot edit transaction, transaction with id %s not found", edited.Id)
	}

	switch editType {
	default:
	case models.All:
		editedTransaction := existingT
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
		newRecurrency, err := recurrency.ChangeRecurrencyEnd(existingT.Recurrency, *edited.Date)
		if err != nil {
			return "", err
		}

		newT := models.Transaction{
			Id:          existingT.Id,
			Amount:      existingT.Amount,
			Date:        *edited.Date,
			Description: existingT.Description,
			Recurrency:  existingT.Recurrency,
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
			Id:          existingT.Id,
			Recurrency:  newRecurrency,
			Description: existingT.Description,
			Amount:      existingT.Amount,
			Date:        existingT.Date,
		})
		if err != nil {
			return (*s.repository).DeleteTransaction(ctx, newTId)
		}
		return newTId, nil
	case models.OnDateOnly:
		/*
		       | ---------------- | ---------------- |
		   start date           edited         end of recurrency
		                     on date only
		*/

		// TODO: validate that the given date is correct regarding the recurrency

		// Create new recurrent transaction to end until new date - 1 day
		newRecurrency, err := recurrency.ChangeRecurrencyEnd(existingT.Recurrency, *edited.Date)
		if err != nil {
			return "", err
		}
		newStartTId, err := (*s.repository).InsertTransaction(ctx, models.Transaction{
			Recurrency:  newRecurrency,
			Description: existingT.Description,
			Amount:      existingT.Amount,
			Date:        existingT.Date,
		})
		if err != nil {
			return "", err
		}

		// Create new recurrent transaction starting from new date + 1 day
		newT := models.Transaction{
			Description: existingT.Description,
			Amount:      existingT.Amount,
			Date:        utils.ParseDateToTime(*edited.Date).AddDate(0, 0, 1).Format(utils.DateLayout),
			Recurrency:  existingT.Recurrency,
		}
		newEndTId, err := (*s.repository).InsertTransaction(ctx, newT)
		if err != nil {
			return (*s.repository).DeleteTransaction(ctx, newStartTId)
		}

		// Create single transaction on given date with edited data
		_, err = (*s.repository).InsertTransaction(ctx, models.Transaction{
			Description: *edited.Description,
			Amount:      *edited.Amount,
			Date:        *edited.Date,
			Recurrency:  "",
		})
		if err != nil {
			_, errOne := (*s.repository).DeleteTransaction(ctx, newStartTId)
			_, errTwo := (*s.repository).DeleteTransaction(ctx, newEndTId)
			if errOne != nil || errTwo != nil {
				log.Printf("cannot rollback!! please fix manually, error: %s, %s", errOne, errTwo)
				return "", fmt.Errorf("cannot rollback, error: %s, %s", errOne, errTwo)
			}
			return "", err
		}

		// Delete old transaction
		return (*s.repository).DeleteTransaction(ctx, existingT.Id)
	}

	return "", fmt.Errorf("error on edit transaction, edit type %s not found", editType)
}
