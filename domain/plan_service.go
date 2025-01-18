package domain

import (
	"cashflow/models"
	"context"
	"time"
)

type planService struct {
	repo *models.Repository
}

func NewPlanService(repo *models.Repository) models.PlanService {
	s := new(planService)
	s.repo = repo
	return s
}

func (s *planService) CreatePlan(ctx context.Context, name string, toAccount string, fromAccount string, startDate time.Time, recurrency string, totalAmount float64) (string, error) {
	panic("not implemented")
}

func (s *planService) ListPlans(ctx context.Context) ([]models.Plan, error) {
	panic("not implemented")
}
