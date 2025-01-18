package repository

import (
	"cashflow/models"
	"context"
)

func (r *supabaseRepository) CreatePlan(ctx context.Context, plan models.Plan) (string, error) {
	panic("not implemented")
}

func (r *supabaseRepository) ReadPlans(ctx context.Context) ([]models.Plan, error) {
	panic("not implemented")
}
