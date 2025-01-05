package repository

import (
	"cashflow/models"
	"cashflow/utils"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

func (repo *supabaseRepository) InsertBalance(ctx context.Context, amount float64, date string) (models.Balance, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	_, err := repo.db.Exec(ctx, "INSERT INTO balances (amount, date, author) VALUES ($1, $2, $3) ON CONFLICT (date, author) DO UPDATE SET amount = $1", amount, date, utils.GetContextValue(ctx, models.UserId))
	if err != nil {
		return models.Balance{}, fmt.Errorf("error inserting balance: %w", err)
	}

	return models.Balance{
		Amount: amount,
		Date:   date,
	}, nil
}

func (repo *supabaseRepository) ListBalances(ctx context.Context, from time.Time, to time.Time) ([]models.Balance, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()

	rows, err := repo.db.Query(ctx, "SELECT amount, to_char(date, 'YYYY-MM-DD') as date FROM balances WHERE date >= $1 AND date <= $2 AND author = $3 ORDER BY date ASC", from.Format(utils.DateLayout), to.Format(utils.DateLayout), utils.GetContextValue(ctx, models.UserId))

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

	_, err := repo.db.Exec(ctx, "DELETE FROM balances WHERE date = $1 AND author = $2", date, utils.GetContextValue(ctx, models.UserId))

	if err != nil {
		return err
	}

	return nil
}
