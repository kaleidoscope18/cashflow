package domain

import (
	"cashflow/models"
	"cashflow/utils"
	"context"
	"time"
)

type balanceService struct {
	repository *models.Repository
}

func NewBalanceService(repo *models.Repository) models.BalanceService {
	s := new(balanceService)
	s.repository = repo
	return s
}

func (s *balanceService) WriteBalance(ctx context.Context, balance float64, date *time.Time, account string) (models.Balance, error) {
	if date != nil {
		newBalance, err := (*s.repository).InsertBalance(ctx, utils.RoundToTwoDigits(balance), *date, account)
		if err != nil {
			return models.Balance{}, err
		}
		return newBalance, nil
	}

	newBalance, err := (*s.repository).InsertBalance(ctx, utils.RoundToTwoDigits(balance), utils.GetTodayDateAsTime(), account)
	if err != nil {
		return models.Balance{}, err
	}
	return newBalance, nil
}

func (s *balanceService) ListBalances(ctx context.Context, from time.Time, to time.Time, account string) ([]models.Balance, error) {
	result, err := (*s.repository).ListBalances(ctx, from, to, account)
	if err != nil {
		return make([]models.Balance, 0), err
	}

	return result, nil
}

func (s *balanceService) DeleteBalance(ctx context.Context, date time.Time, account string) error {
	err := (*s.repository).DeleteBalance(ctx, date, account)
	return err
}
