package repository

import (
	"cashflow/models"
	"cashflow/utils"
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

func (r *supabaseRepository) CreatePlan(ctx context.Context, plan models.Plan, toAccountTransactionId *string, fromAccountTransactionId *string) (string, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	var toAccountField string
	if toAccountTransactionId != nil {
		toAccountField = "to_account"
	} else {
		toAccountField = "to_external_account"
	}

	var fromAccountField string
	if fromAccountTransactionId != nil {
		fromAccountField = "from_account"
	} else {
		fromAccountField = "from_external_account"
	}

	row := r.db.QueryRow(ctx,
		"INSERT INTO cashflow.plans "+
			"(name, "+toAccountField+", "+fromAccountField+", "+
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

	var id string
	err := row.Scan(&id)
	if err != nil {
		log.Printf("error on INSERT into cashflow.plans: %v", err)
		return "", err
	}
	return id, nil
}

func (r *supabaseRepository) ReadPlans(ctx context.Context) ([]models.Plan, error) {
	r.mutex.Lock()
	rows, err := r.db.Query(ctx,
		"SELECT id, name, to_account, from_account, "+
			"start_date, recurrency, total_amount, "+
			"from_external_account, to_external_account "+
			"FROM cashflow.plans "+
			"WHERE author = $1 "+
			"ORDER BY start_date ASC",
		utils.GetContextValue(ctx, models.UserId))
	r.mutex.Unlock()

	if err != nil {
		log.Printf("error when querying SELECT on cashflow.plans: %v", err)
		return nil, err
	}
	defer rows.Close()

	type plan struct {
		Id                  string
		Name                string
		StartDate           string
		Recurrency          string
		TotalAmount         float64
		FromAccount         *string
		FromExternalAccount *string
		ToAccount           *string
		ToExternalAccount   *string
	}
	plans, err := pgx.CollectRows(rows, pgx.RowToStructByName[plan])
	if err != nil {
		log.Printf("error when collecting rows: %v", err)
		return nil, err
	}

	convertedPlans := []models.Plan{}
	for _, plan := range plans {
		var toAccount, fromAccount string
		if plan.ToAccount != nil {
			toAccount = *plan.ToAccount
		} else {
			toAccount = *plan.ToExternalAccount
		}
		if plan.FromAccount != nil {
			fromAccount = *plan.FromAccount
		} else {
			fromAccount = *plan.FromExternalAccount
		}
		convertedPlans = append(convertedPlans, models.Plan{
			Id:          plan.Id,
			Name:        plan.Name,
			ToAccount:   toAccount,
			FromAccount: fromAccount,
			StartDate:   plan.StartDate,
			Recurrency:  plan.Recurrency,
			TotalAmount: plan.TotalAmount,
		})
	}

	return convertedPlans, nil
}

func (r *supabaseRepository) DeletePlan(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, err := r.db.Exec(ctx,
		"DELETE FROM cashflow.plans WHERE id = $1 AND author = $2",
		id,
		utils.GetContextValue(ctx, models.UserId))

	if err != nil {
		log.Printf("error during DELETE on cashflow.plans: %v", err)
		return err
	}

	return nil
}
