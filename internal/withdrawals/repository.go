package withdrawals

import (
	"context"
	"database/sql"
	"github.com/magmel48/go-musthave-diploma/internal/logger"
)

//go:generate mockery --name=Repository
type Repository interface {
	Create(ctx context.Context, userID int64, orderNumber string, amount float64) error
	FindByUser(ctx context.Context, userID int64) ([]Withdrawal, error)
}

type PostgreSQLRepository struct {
	db *sql.DB
}

func NewPostgreSQLRepository(db *sql.DB) *PostgreSQLRepository {
	return &PostgreSQLRepository{db: db}
}

func (repository *PostgreSQLRepository) Create(ctx context.Context, userID int64, orderNumber string, amount float64) error {
	_, err := repository.db.ExecContext(
		ctx, `INSERT INTO "withdrawals" ("user_id", "order", "sum") VALUES ($1, $2, $3)`, userID, orderNumber, amount)

	return err
}

func (repository *PostgreSQLRepository) FindByUser(ctx context.Context, userID int64) ([]Withdrawal, error) {
	rows, err := repository.db.QueryContext(
		ctx,
		`SELECT "order", "sum", "processed_at" FROM "withdrawals" WHERE "user_id" = $1 ORDER BY "processed_at" ASC`, userID)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}()

	withdrawals := make([]Withdrawal, 0)
	for rows.Next() {
		var withdrawal Withdrawal
		if err = rows.Scan(&withdrawal.Order, &withdrawal.Sum, &withdrawal.ProcessedAt); err != nil {
			return withdrawals, nil
		}

		withdrawals = append(withdrawals, withdrawal)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return withdrawals, nil
}
