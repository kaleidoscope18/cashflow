package domain

import (
	"cashflow/domain/recurrency"
	"cashflow/models"
	"cashflow/utils"
	"errors"
	"strconv"
	"time"
)

func listTransactions(transactions []models.Transaction, balances []models.Balance, from time.Time, to time.Time) ([]models.ComputedTransaction, error) {
	withRecurrency, withoutRecurrency := recurrency.SplitTransactionsWithRecurrency(transactions)
	recurrencyOffsprings, _ := recurrency.GenerateTransactionsFromRecurrency(withRecurrency, from, to)
	transactionsToEnrich := append(withoutRecurrency, recurrencyOffsprings...)

	enrichedTransactions := make([]models.ComputedTransaction, len(transactionsToEnrich))
	for i, t := range transactionsToEnrich {
		balanceOnSameDay := getBalanceOnSameDay(t.Date, balances)
		latestBalance, latestBalanceError := getLatestBalanceBefore(t.Date, balances)
		previousTransaction, previousTransactionError := getPreviousTransaction(i, enrichedTransactions)

		item := t
		go func() {
			enrichedTransactions[i] = models.ComputedTransaction{
				Transaction: &item,
				Status:      utils.GetStatusFromDate(utils.GetTodayDate(), t.Date),
			}
		}()

		switch {
		case balanceOnSameDay != nil:
			enrichedTransactions[i].Balance = balanceOnSameDay.Amount
		case latestBalanceError != nil && previousTransactionError != nil:
			enrichedTransactions[i].Balance = t.Amount
		case latestBalanceError != nil:
			enrichedTransactions[i].Balance = previousTransaction.Balance + t.Amount
		case previousTransactionError != nil:
			enrichedTransactions[i].Balance = latestBalance.Amount + t.Amount
		default:
			enrichedTransactions[i].Balance = getBalanceForTransaction(t, *previousTransaction, latestBalance)
		}

		enrichedTransactions[i].Balance = utils.RoundToTwoDigits(enrichedTransactions[i].Balance)
	}

	return enrichedTransactions, nil
}

func getBalanceOnSameDay(date string, balances []models.Balance) *models.Balance {
	for _, b := range balances {
		if b.Date == date {
			return &b
		}
	}
	return nil
}

func getLatestBalanceBefore(date string, orderedBalances []models.Balance) (models.Balance, error) {
	if len(orderedBalances) == 0 || !utils.IsDateBefore(orderedBalances[0].Date, date) {
		return models.Balance{}, errors.New("no balance before " + date)
	}

	var currentBalance models.Balance
	for _, b := range orderedBalances {
		if !utils.IsDateBefore(b.Date, date) {
			break
		}
		currentBalance = b
	}

	return currentBalance, nil
}

func getPreviousTransaction(index int, transactionsWithBalances []models.ComputedTransaction) (*models.ComputedTransaction, error) {
	if index > 0 {
		return &transactionsWithBalances[index-1], nil
	}

	return nil, errors.New("no transaction precedes transaction #" + strconv.Itoa(index))
}

func getBalanceForTransaction(transaction models.Transaction, previousTransaction models.ComputedTransaction, latestBalance models.Balance) float64 {
	if transaction.Date == latestBalance.Date {
		return latestBalance.Amount
	}

	if utils.IsDateBefore(previousTransaction.Date, latestBalance.Date) || previousTransaction.Date == latestBalance.Date {
		return latestBalance.Amount + transaction.Amount
	}

	return previousTransaction.Balance + transaction.Amount
}
