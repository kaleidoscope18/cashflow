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

func (s *transactionService) ListTransactions(ctx context.Context, from time.Time, to time.Time, account string) ([]models.ComputedTransaction, error) {
	transactions, err := (*s.repository).ReadTransactions(ctx, to, account)
	if err != nil {
		return nil, err
	}

	balances, err := (*s.balanceService).ListBalances(ctx, from, to, account)
	if err != nil {
		return nil, err
	}

	results, err := listTransactions(transactions, balances, from, to)
	return results, err
}

func (s *transactionService) CreateTransaction(ctx context.Context, date time.Time, amount float64, description string, recurrency string, isManual bool, account string) (string, error) {
	id, err := (*s.repository).CreateTransaction(ctx, models.Transaction{
		Amount:      utils.RoundToTwoDigits(amount),
		Date:        date.Format(utils.DateLayout),
		Description: description,
		Recurrency:  recurrency,
		Account:     account,
		IsManual:    isManual,
	})

	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *transactionService) DeleteRecurringTransaction(ctx context.Context, id string, date time.Time, account string) error {
	existingT, err := (*s.repository).ReadTransactionById(ctx, id, account)
	if err != nil {
		return fmt.Errorf("cannot delete transaction, transaction with id %s not found", id)
	}
	if existingT.Recurrency == "" {
		return fmt.Errorf("cannot delete transaction, transaction with id %s is not recurring", id)
	}

	/*
	     | ---------------- | ---------------- |
	   start             removed         end of recurrency
	                   on date only
	*/

	// TODO: validate that the given date is correct regarding the recurrency

	// Create new recurrent transaction to end until date - 1 day
	newRecurrency, err := recurrency.ChangeRecurrencyEnd(existingT.Recurrency, date)
	if err != nil {
		return err
	}
	newStartTId, err := (*s.repository).CreateTransaction(ctx, models.Transaction{
		Recurrency:  newRecurrency,
		Description: existingT.Description,
		Amount:      existingT.Amount,
		Date:        existingT.Date,
		Account:     account,
		IsManual:    existingT.IsManual,
	})
	if err != nil {
		return err
	}

	// Create new recurrent transaction starting from date + 1 day
	newT := models.Transaction{
		Description: existingT.Description,
		Amount:      existingT.Amount,
		Date:        utils.ParseDateToTime(date.Format(utils.DateLayout)).AddDate(0, 0, 1).Format(utils.DateLayout),
		Recurrency:  existingT.Recurrency,
		Account:     account,
		IsManual:    existingT.IsManual,
	}
	newEndTId, err := (*s.repository).CreateTransaction(ctx, newT)
	if err != nil {
		return (*s.repository).DeleteTransaction(ctx, string(newStartTId), account)
	}

	// Delete old transaction
	err = (*s.repository).DeleteTransaction(ctx, existingT.Id, account)
	if err != nil {
		errOne := (*s.repository).DeleteTransaction(ctx, string(newStartTId), account)
		errTwo := (*s.repository).DeleteTransaction(ctx, string(newEndTId), account)
		if errOne != nil || errTwo != nil {
			log.Printf("cannot rollback!! please fix manually, error: %s, %s", errOne, errTwo)
			return fmt.Errorf("cannot rollback, error: %s, %s", errOne, errTwo)
		}
		return err
	}
	return nil
}

func (s *transactionService) DeleteTransaction(ctx context.Context, id string, account string) error {
	return (*s.repository).DeleteTransaction(ctx, id, account)
}

func (s *transactionService) EditTransaction(ctx context.Context, editType models.TransactionEditType, edited models.TransactionEdit) error {
	existingT, err := (*s.repository).ReadTransactionById(ctx, edited.Id, edited.Account)
	if err != nil {
		return fmt.Errorf("transaction with id %s not found", edited.Id)
	}

	switch editType {
	default:
	case models.All:
		editedTransaction := existingT
		if edited.Amount != nil {
			editedTransaction.Amount = *edited.Amount
		}
		if edited.Date != nil {
			editedTransaction.Date = (*edited.Date).Format(utils.DateLayout)
		}
		if edited.Description != nil {
			editedTransaction.Description = *edited.Description
		}
		if edited.Recurrency != nil {
			editedTransaction.Recurrency = *edited.Recurrency
		}
		if edited.IsManual != nil {
			editedTransaction.IsManual = *edited.IsManual
		}
		return (*s.repository).UpdateTransaction(ctx, editedTransaction)
	case models.FromDate:
		/*
		     | ---------------- | ---- edited ---- |
		   start            from date             end
		*/
		recurrencyWithStop, err := recurrency.ChangeRecurrencyEnd(existingT.Recurrency, *edited.Date)
		if err != nil {
			return err
		}

		newT := models.Transaction{
			Id:          existingT.Id,
			Amount:      existingT.Amount,
			Date:        (*edited.Date).Format(utils.DateLayout),
			Description: existingT.Description,
			Recurrency:  existingT.Recurrency,
			Account:     existingT.Account,
			IsManual:    existingT.IsManual,
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
		if edited.IsManual != nil {
			newT.IsManual = *edited.IsManual
		}
		newTId, err := (*s.repository).CreateTransaction(ctx, newT)
		if err != nil {
			return err
		}

		err = (*s.repository).UpdateTransaction(ctx, models.Transaction{
			Id:          existingT.Id,
			Recurrency:  recurrencyWithStop,
			Description: existingT.Description,
			Amount:      existingT.Amount,
			Date:        existingT.Date,
			Account:     existingT.Account,
			IsManual:    existingT.IsManual,
		})
		if err != nil {
			return (*s.repository).DeleteTransaction(ctx, string(newTId), edited.Account)
		}
		return nil
	case models.OnDateOnly:
		/*
		       | ---------------- X ---------------- |
		   start date           edited         end of recurrency
		                     on date only
		*/

		// Create new recurrent transaction to end until new date - 1 day
		recurrencyWithStop, err := recurrency.ChangeRecurrencyEnd(existingT.Recurrency, *edited.Date)
		if err != nil {
			return err
		}
		newStartTId, err := (*s.repository).CreateTransaction(ctx, models.Transaction{
			Recurrency:  recurrencyWithStop,
			Description: existingT.Description,
			Amount:      existingT.Amount,
			Date:        existingT.Date,
			Account:     existingT.Account,
			IsManual:    existingT.IsManual,
		})
		if err != nil {
			return err
		}

		// Create new recurrent transaction starting from new date + 1 day
		newT := models.Transaction{
			Description: existingT.Description,
			Amount:      existingT.Amount,
			Date:        (*edited.Date).AddDate(0, 0, 1).Format(utils.DateLayout),
			Recurrency:  existingT.Recurrency,
			Account:     existingT.Account,
			IsManual:    existingT.IsManual,
		}
		newEndTId, err := (*s.repository).CreateTransaction(ctx, newT)
		if err != nil {
			return (*s.repository).DeleteTransaction(ctx, string(newStartTId), edited.Account)
		}

		// Create single transaction on given date with edited data
		description := existingT.Description
		if edited.Description != nil {
			description = *edited.Description
		}
		isManual := existingT.IsManual
		if edited.IsManual != nil {
			isManual = *edited.IsManual
		}
		_, err = (*s.repository).CreateTransaction(ctx, models.Transaction{
			Description: description,
			Amount:      *edited.Amount,
			Date:        (*edited.Date).Format(utils.DateLayout),
			Recurrency:  "",
			Account:     edited.Account,
			IsManual:    isManual,
		})
		if err != nil {
			errOne := (*s.repository).DeleteTransaction(ctx, string(newStartTId), edited.Account)
			errTwo := (*s.repository).DeleteTransaction(ctx, string(newEndTId), edited.Account)
			if errOne != nil || errTwo != nil {
				log.Printf("cannot rollback!! please fix manually, error: %s, %s", errOne, errTwo)
				return fmt.Errorf("cannot rollback, error: %s, %s", errOne, errTwo)
			}
			return err
		}

		// Delete old transaction
		return (*s.repository).DeleteTransaction(ctx, existingT.Id, edited.Account)
	}

	return fmt.Errorf("error on edit transaction, edit type %s not found", editType)
}
