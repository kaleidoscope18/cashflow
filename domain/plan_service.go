package domain

import (
	"cashflow/domain/recurrency"
	"cashflow/models"
	"context"
	"fmt"
	"log"
	"slices"
)

type planService struct {
	repo *models.Repository
}

func NewPlanService(repo *models.Repository) models.PlanService {
	s := new(planService)
	s.repo = repo
	return s
}

func (s *planService) CreatePlan(ctx context.Context, name string, toAccount string, fromAccount string, startDate string, rec string, totalAmount float64) (string, error) {
	accounts, err := (*s.repo).ListAccounts(ctx)
	if err != nil {
		log.Printf("error on creating a plan: %v", err)
		return "", err
	}

	var toAccountTransactionId *string
	if slices.Contains(accounts, toAccount) {
		toAccountTransaction, err := recurrency.GenerateTransactionFromPlan(name, toAccount, startDate, rec, totalAmount)
		if err != nil {
			log.Printf("error on creating a plan: %v", err)
			return "", err
		}
		id, err := (*s.repo).CreateTransaction(ctx, toAccountTransaction)
		if err != nil {
			log.Printf("error on creating a plan: %v", err)
			return "", err
		}
		toAccountTransactionId = &id
	}

	var fromAccountTransactionId *string
	if slices.Contains(accounts, fromAccount) {
		fromAccountTransaction, err := recurrency.GenerateTransactionFromPlan(name, fromAccount, startDate, rec, -totalAmount)
		if err != nil {
			log.Printf("error on creating a plan: %v", err)
			return "", err
		}
		id, err := (*s.repo).CreateTransaction(ctx, fromAccountTransaction)
		if err != nil {
			log.Printf("error on creating a plan: %v", err)
			return "", err
		}
		fromAccountTransactionId = &id
	}

	if toAccountTransactionId == nil && fromAccountTransactionId == nil {
		return "", fmt.Errorf("from and to accounts cannot be both external")
	}

	planId, err := (*s.repo).CreatePlan(ctx, models.Plan{
		Name:        name,
		ToAccount:   toAccount,
		FromAccount: fromAccount,
		StartDate:   startDate,
		Recurrency:  rec,
		TotalAmount: totalAmount,
	}, toAccountTransactionId, fromAccountTransactionId)
	if err != nil {
		log.Printf("error on creating a plan: %v", err)
		return "", err
	}
	return planId, nil
}

func (s *planService) ListPlans(ctx context.Context) ([]models.Plan, error) {
	plans, err := (*s.repo).ReadPlans(ctx)
	if err != nil {
		log.Printf("error on listing plans: %v", err)
		return nil, err
	}
	return plans, nil
}

func (s *planService) DeletePlan(ctx context.Context, id string) error {
	err := (*s.repo).DeletePlan(ctx, id)
	if err != nil {
		log.Printf("error on deleting a plan: %v", err)
		return err
	}
	return nil
}
