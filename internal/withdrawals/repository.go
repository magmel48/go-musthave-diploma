package withdrawals

import (
	"context"
	"database/sql"
)

//go:generate mockery --name=Repository
type Repository interface {
	Create(ctx context.Context, userID int64, orderNumber string, amount float64) error
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
