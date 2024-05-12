package domain

import (
	"cashflow/models"
	"cashflow/utils"
)

type balanceService struct {
	repository *models.BalanceRepository
}

func NewBalanceService(repo *models.BalanceRepository) models.BalanceService {
	s := new(balanceService)
	s.repository = repo
	return s
}

func (s balanceService) WriteBalance(balance float64, date *string) (*models.Balance, error) {
	if date != nil {
		newBalance := (*s.repository).InsertBalance(utils.RoundToTwoDigits(balance), utils.ParseDate(date))
		return &newBalance, nil
	}

	newBalance := (*s.repository).InsertBalance(balance, utils.GetTodayDate())
	return &newBalance, nil
}

func (s balanceService) ListBalances() ([]*models.Balance, error) {
	return utils.ConvertStructToPointersArray((*s.repository).ListBalances()), nil
}
