package domain

import (
	"cashflow/models"
	"cashflow/utils"
	"time"
)

type balanceService struct {
	repository *models.BalanceRepository
}

func NewBalanceService(repo *models.BalanceRepository) models.BalanceService {
	s := new(balanceService)
	s.repository = repo
	return s
}

func (s *balanceService) WriteBalance(balance float64, date *string) (models.Balance, error) {
	if date != nil {
		newBalance, err := (*s.repository).InsertBalance(utils.RoundToTwoDigits(balance), utils.ParseDate(date))
		if err != nil {
			return models.Balance{}, err
		}
		return newBalance, nil
	}

	newBalance, err := (*s.repository).InsertBalance(balance, utils.GetTodayDate())
	if err != nil {
		return models.Balance{}, err
	}
	return newBalance, nil
}

func (s *balanceService) ListBalances(from time.Time, to time.Time) ([]models.Balance, error) {
	result, err := (*s.repository).ListBalances(from, to)
	if err != nil {
		return make([]models.Balance, 0), err
	}

	return utils.SortByDate(result), nil
}
