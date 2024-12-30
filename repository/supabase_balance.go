package repository

import (
	"cashflow/models"
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

func (repo *supabaseRepository) InsertBalance(ctx context.Context, amount float64, date string) (models.Balance, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	balance := models.Balance{
		Amount: amount,
		Date:   date,
	}

	_, err := repo.db.Exec(ctx, "INSERT INTO balances (amount, date, author) VALUES ($1, $2, $3)", balance.Amount, balance.Date, ctx.Value(models.UserId).(string))

	if err != nil {
		return models.Balance{}, err
	}

	return balance, nil
}

func (repo *supabaseRepository) ListBalances(ctx context.Context, from time.Time, to time.Time) ([]models.Balance, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	rows, err := repo.db.Query(ctx, "SELECT amount, date FROM balances WHERE date >= $1 AND date <= $2 AND author = $3 ORDER BY date ASC", from.Format("2006-01-01"), to.Format("2006-01-01"), ctx.Value(models.UserId).(string))

	if err != nil {
		return make([]models.Balance, 0), err
	}

	var balances []models.Balance
	balances, err = pgx.CollectRows(rows, pgx.RowToStructByName[models.Balance])
	if err != nil {
		return make([]models.Balance, 0), err
	}

	return balances, nil
}

func (repo *supabaseRepository) DeleteBalance(ctx context.Context, date string) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	_, err := repo.db.Exec(ctx, "DELETE FROM balances WHERE date = $1 AND author = $2", date, ctx.Value(models.UserId).(string))

	if err != nil {
		return err
	}

	return nil
}
