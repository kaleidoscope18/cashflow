package domain

import (
	"cashflow/domain/recurrency"
	"cashflow/models"
	"context"
)

type planService struct {
	repo *models.Repository
}

func NewPlanService(repo *models.Repository) models.PlanService {
	s := new(planService)
	s.repo = repo
	return s
}

func (s *planService) CreatePlan(ctx context.Context, name string, toAccount string, fromAccount *string, startDate string, rec string, totalAmount float64) (string, error) {
	toAccountTransaction, err := recurrency.GenerateTransactionFromPlan(name, toAccount, startDate, rec, totalAmount)
	if err != nil {
		return "", err
	}
	toAccountTransactionId, err := (*s.repo).CreateTransaction(ctx, toAccountTransaction)
	if err != nil {
		return "", err
	}

	var fromAccountValue string
	var fromAccountTransactionId *string
	if fromAccount != nil {
		fromAccountValue = *fromAccount
		fromAccountTransaction, err := recurrency.GenerateTransactionFromPlan(name, *fromAccount, startDate, rec, -totalAmount)
		if err != nil {
			return "", err
		}
		id, err := (*s.repo).CreateTransaction(ctx, fromAccountTransaction)
		if err != nil {
			return "", err
		}
		fromAccountTransactionId = &id
	}

	planId, err := (*s.repo).CreatePlan(ctx, models.Plan{
		Name:        name,
		ToAccount:   toAccount,
		FromAccount: fromAccountValue,
		StartDate:   startDate,
		Recurrency:  rec,
		TotalAmount: totalAmount,
	}, toAccountTransactionId, fromAccountTransactionId)
	if err != nil {
		return "", err
	}
	return planId, nil
}

func (s *planService) ListPlans(ctx context.Context) ([]models.Plan, error) {
	return (*s.repo).ReadPlans(ctx)
}
