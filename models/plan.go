package models

import (
	"context"
	"time"
)

type PlanRepository interface {
	CreatePlan(ctx context.Context, plan Plan) (string, error)
	ReadPlans(ctx context.Context) ([]Plan, error)
}

type PlanService interface {
	CreatePlan(ctx context.Context, name string, toAccount string, fromAccount string, startDate time.Time, recurrency string, totalAmount float64) (string, error)
	ListPlans(ctx context.Context) ([]Plan, error)
}

type Plan struct {
	Id          string
	Name        string
	ToAccount   string
	FromAccount string
	StartDate   time.Time
	Recurrency  string
	TotalAmount float64
}
