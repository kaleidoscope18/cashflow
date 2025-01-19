package domain

import (
	"cashflow/domain/recurrency"
	"cashflow/domain/status"
	"cashflow/models"
	"cashflow/utils"
	"errors"
	"sort"
	"strconv"
	"time"
)

func listTransactions(transactions []models.Transaction, balances []models.Balance, from time.Time, to time.Time) ([]models.ComputedTransaction, error) {
	withRecurrency, withoutRecurrency := recurrency.SplitTransactionsWithRecurrency(transactions)
	recurrencyOffsprings, _ := recurrency.GenerateTransactionsFromRecurrency(withRecurrency, from, to)
	transactionsToEnrich := append(withoutRecurrency, recurrencyOffsprings...)

	enrichedTransactions := make([]models.ComputedTransaction, len(transactionsToEnrich))
	sort.Slice(transactionsToEnrich, func(i, j int) bool {
		dI := utils.ParseDateToTime(transactionsToEnrich[i].Date)
		dJ := utils.ParseDateToTime(transactionsToEnrich[j].Date)

		if dI.Equal(dJ) {
			if transactionsToEnrich[i].IsManual == transactionsToEnrich[j].IsManual {
				return transactionsToEnrich[i].Amount > transactionsToEnrich[j].Amount
			}
			return !transactionsToEnrich[i].IsManual
		}

		return dI.Before(dJ)
	})
	today := utils.GetTodayDate()

	for i, t := range transactionsToEnrich {
		balanceOneDayBefore := getBalanceOneDayBefore(t.Date, balances)
		latestBalance, latestBalanceError := getLatestBalanceBefore(t.Date, balances)
		previousTransaction, previousTransactionError := getPreviousTransaction(i, enrichedTransactions)

		tCopy := t
		enrichedTransactions[i] = models.ComputedTransaction{
			Transaction: &tCopy,
			Status:      status.GetStatusFromDate(today, t.Date),
		}

		switch {
		case balanceOneDayBefore != nil:
			enrichedTransactions[i].Balance = balanceOneDayBefore.Amount + t.Amount
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

func getBalanceOneDayBefore(date string, balances []models.Balance) *models.Balance {
	for _, b := range balances {
		if b.Date == utils.ParseDateToTime(date).AddDate(0, 0, -1).Format(utils.DateLayout) {
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
