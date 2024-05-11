package transactions

import (
	"cashflow/models"
	"cashflow/utils"

	"github.com/dchest/uniuri"
)

type (
	service struct {
		repo models.Repository
	}
)

func New(db models.Repository) models.TransactionService {
	return service{
		repo: db,
	}
}

func (s service) ListTransactions(today *string) []*models.ComputedTransaction {
	return listTransactions(today, s.repo)
}

func (s service) WriteTransaction(date string, amount float64, description string) *models.Transaction {
	newTransaction := models.Transaction{
		Id:          uniuri.NewLen(5),
		Amount:      utils.RoundToTwoDigits(amount),
		Date:        utils.ParseDate(&date),
		Description: description,
	}
	s.repo.InsertTransaction(newTransaction)

	return &newTransaction
}

func (s service) WriteBalance(balance float64, date *string) *models.Balance {
	if date != nil {
		newBalance := s.repo.InsertBalance(utils.RoundToTwoDigits(balance), utils.ParseDate(date))
		return &newBalance
	}

	newBalance := s.repo.InsertBalance(balance, utils.GetTodayDate())
	return &newBalance
}

func (s service) ListBalances() []*models.Balance {
	return utils.ConvertStructToPointersArray(s.repo.ListBalances())
}
