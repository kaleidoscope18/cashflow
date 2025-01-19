package repository

import (
	"cashflow/models"
	"cashflow/utils"
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

func (r *supabaseRepository) CreatePlan(ctx context.Context, plan models.Plan, toAccountTransactionId string, fromAccountTransactionId *string) (string, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	id, err := r.db.Exec(ctx,
		"INSERT INTO cashflow.plans "+
			"(name, to_account, from_account, "+
			"start_date, recurrency, total_amount, "+
			"from_transaction_id, to_transaction_id, author) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) "+
			"RETURNING id",
		plan.Name,
		plan.ToAccount,
		plan.FromAccount,
		plan.StartDate,
		plan.Recurrency,
		plan.TotalAmount,
		fromAccountTransactionId,
		toAccountTransactionId,
		utils.GetContextValue(ctx, models.UserId),
	)

	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func (r *supabaseRepository) ReadPlans(ctx context.Context) ([]models.Plan, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	rows, err := r.db.Query(ctx,
		"SELECT id, name, to_account, from_account, "+
			"start_date, recurrency, total_amount "+
			"FROM cashflow.plans "+
			"WHERE author = $1 "+
			"ORDER BY start_date ASC",
		utils.GetContextValue(ctx, models.UserId))

	if err != nil {
		log.Printf("error listing plans: %v", err)
		return nil, err
	}

	plans, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Plan])

	if err != nil {
		return nil, err
	}

	return plans, nil
}
