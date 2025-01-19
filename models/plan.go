package models

import (
	"context"
)

type PlanRepository interface {
	CreatePlan(ctx context.Context, plan Plan, toAccountTransactionId *string, fromAccountTransactionId *string) (string, error)
	ReadPlans(ctx context.Context) ([]Plan, error)
	DeletePlan(ctx context.Context, id string) error
}

type PlanService interface {
	CreatePlan(ctx context.Context, name string, toAccount string, fromAccount string, startDate string, recurrency string, totalAmount float64) (string, error)
	ListPlans(ctx context.Context) ([]Plan, error)
	DeletePlan(ctx context.Context, id string) error
}

type Plan struct {
	Id          string
	Name        string
	ToAccount   string
	FromAccount string
	StartDate   string
	Recurrency  string
	TotalAmount float64
}
