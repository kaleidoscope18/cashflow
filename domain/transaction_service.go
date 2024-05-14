package domain

import (
	"cashflow/models"
	"cashflow/utils"

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

func (s *transactionService) ListTransactions(today *string) ([]*models.ComputedTransaction, error) {
	return listTransactions(today, s.repository, s.balanceService)
}

func (s *transactionService) WriteTransaction(date string, amount float64, description string) (*models.Transaction, error) {
	newTransaction := models.Transaction{
		Id:          uniuri.NewLen(5),
		Amount:      utils.RoundToTwoDigits(amount),
		Date:        utils.ParseDate(&date),
		Description: description,
	}
	(*s.repository).InsertTransaction(newTransaction)

	return &newTransaction, nil
}
