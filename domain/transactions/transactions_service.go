package transactions

import (
	"cashflow/models"
	"cashflow/utils"

	"github.com/dchest/uniuri"
)

type service struct {
	repository models.Repository
}

func New(injectedRepo models.Repository) models.TransactionService {
	service := new(service)
	service.repository = injectedRepo
	return service
}

func (s service) ListTransactions(today *string) []*models.ComputedTransaction {
	return listTransactions(today, s.repository)
}

func (s service) WriteTransaction(date string, amount float64, description string) *models.Transaction {
	newTransaction := models.Transaction{
		Id:          uniuri.NewLen(5),
		Amount:      utils.RoundToTwoDigits(amount),
		Date:        utils.ParseDate(&date),
		Description: description,
	}
	s.repository.InsertTransaction(newTransaction)

	return &newTransaction
}

func (service service) WriteBalance(balance float64, date *string) *models.Balance {
	if date != nil {
		newBalance := service.repository.InsertBalance(utils.RoundToTwoDigits(balance), utils.ParseDate(date))
		return &newBalance
	}

	newBalance := service.repository.InsertBalance(balance, utils.GetTodayDate())
	return &newBalance
}

func (service service) ListBalances() []*models.Balance {
	return utils.ConvertStructToPointersArray(service.repository.ListBalances())
}
