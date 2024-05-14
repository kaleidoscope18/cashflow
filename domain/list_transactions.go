package domain

import (
	"cashflow/models"
	"cashflow/utils"
	"errors"
	"strconv"
)

func listTransactions(today *string, repo *models.TransactionRepository, balancesService *models.BalanceService) ([]*models.ComputedTransaction, error) {
	transactions := (*repo).ListTransactions()
	balances, err := (*balancesService).ListBalances()
	if err != nil {
		return make([]*models.ComputedTransaction, 0), err
	}

	enrichedTransactions := make([]*models.ComputedTransaction, len(transactions))

	for i, t := range transactions {
		latestBalance, latestBalanceError := getLatestBalanceBefore(t.Date, balances)
		previousTransaction, previousTransactionError := getPreviousTransaction(i, enrichedTransactions)

		enrichedTransactions[i] = &models.ComputedTransaction{Transaction: &(transactions)[i], Status: utils.GetStatusFromDate(&t.Date, &transactions[i].Date)}

		switch {
		case latestBalanceError != nil && previousTransactionError != nil:
			enrichedTransactions[i].Balance = utils.RoundToTwoDigits(t.Amount)
		case latestBalanceError != nil:
			enrichedTransactions[i].Balance = utils.RoundToTwoDigits(previousTransaction.Balance + t.Amount)
		case previousTransactionError != nil:
			enrichedTransactions[i].Balance = utils.RoundToTwoDigits(latestBalance.Amount + t.Amount)
		default:
			enrichedTransactions[i].Balance = utils.RoundToTwoDigits(getBalanceForTransaction(t, *previousTransaction, latestBalance))
		}
	}

	return enrichedTransactions, nil
}

func getBalanceForTransaction(transaction models.Transaction, previousTransaction models.ComputedTransaction, latestBalance models.Balance) float64 {
	if utils.IsDateBefore(previousTransaction.Date, latestBalance.Date) || previousTransaction.Date == latestBalance.Date {
		return latestBalance.Amount + transaction.Amount
	}
	return previousTransaction.Balance + transaction.Amount
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

func getPreviousTransaction(index int, transactionsWithBalances []*models.ComputedTransaction) (*models.ComputedTransaction, error) {
	if index > 0 {
		return transactionsWithBalances[index-1], nil
	}

	return nil, errors.New("no transaction precedes transaction #" + strconv.Itoa(index))
}
