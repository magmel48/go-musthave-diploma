package balances

import (
	"context"
	"database/sql"
)

//go:generate mockery --name=Repository
type Repository interface {
	FindByUser(ctx context.Context, userID int64) (*Balance, error)
}

type PostgreSQLRepository struct {
	db *sql.DB
}

func NewPostgreSQLRepository(db *sql.DB) *PostgreSQLRepository {
	return &PostgreSQLRepository{db: db}
}

func (repository *PostgreSQLRepository) FindByUser(ctx context.Context, userID int64) (*Balance, error) {
	var result Balance

	if err := repository.db.QueryRowContext(
		ctx,
		`INSERT INTO "balances" ("user_id", "current", "withdrawn") VALUES ($1, 0, 0)
			ON CONFLICT ("user_id") DO UPDATE SET "user_id" = $1 RETURNING "current", "withdrawn"`,
		userID).Scan(
		&result.Current, &result.Withdrawn); err != nil {
		return nil, err
	}

	return &result, nil
}
