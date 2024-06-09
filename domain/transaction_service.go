package domain

import (
	"cashflow/models"
	"cashflow/utils"
	"time"

	"github.com/dchest/uniuri"
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

func (s *transactionService) ListTransactions(from time.Time, to time.Time) ([]models.ComputedTransaction, error) {
	transactions, err := (*s.repository).ListTransactions(from, to)
	if err != nil {
		return make([]models.ComputedTransaction, 0), err
	}

	balances, err := (*s.balanceService).ListBalances(from, to)
	if err != nil {
		return make([]models.ComputedTransaction, 0), err
	}

	return listTransactions(utils.SortByDate(transactions), balances)
}

func (s *transactionService) WriteTransaction(date string, amount float64, description string) (*models.Transaction, error) {
	t, err := (*s.repository).InsertTransaction(models.Transaction{
		Id:          uniuri.NewLen(5),
		Amount:      utils.RoundToTwoDigits(amount),
		Date:        utils.ParseDate(&date),
		Description: description,
	})

	if err != nil {
		return &models.Transaction{}, err
	}

	return &t, nil
}
