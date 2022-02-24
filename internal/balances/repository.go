package balances

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgerrcode"
	"strings"
)

//go:generate mockery --name=Repository
type Repository interface {
	FindByUser(ctx context.Context, userID int64) (*Balance, error)
	Change(ctx context.Context, userID int64, amount float64) error
}

var ErrInsufficientFunds = errors.New("insufficient funds")

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
		`INSERT INTO "balances" ("user_id", "current") VALUES ($1, 0)
			ON CONFLICT ("user_id") DO UPDATE SET "user_id" = $1 RETURNING "current", "withdrawn"`,
		userID).Scan(
		&result.Current, &result.Withdrawn); err != nil {
		return nil, err
	}

	return &result, nil
}

func (repository *PostgreSQLRepository) Change(ctx context.Context, userID int64, amount float64) error {
	if amount < 0 {
		if result, err := repository.db.ExecContext(
			ctx,
			`UPDATE "balances" SET "current" = "current" + $1, "withdrawn" = "withdrawn" - $1 WHERE "user_id" = $2`,
			amount, userID); err != nil {
			if strings.Contains(err.Error(), pgerrcode.CheckViolation) {
				return ErrInsufficientFunds
			}

			return err
		} else {
			rowsAffected, err := result.RowsAffected()
			if err != nil {
				return err
			}

			// if no rows affected - no funds then on user balance
			if rowsAffected != 1 {
				return ErrInsufficientFunds
			}
		}

		return nil
	}

	// in case of positive amount to add on top of user balance - just safely upsert
	if _, err := repository.db.ExecContext(
		ctx,
		`INSERT INTO "balances" ("user_id", "current") VALUES ($1, $2)
			ON CONFLICT ("user_id") DO UPDATE SET "current" = "balances"."current" + $2`,
		userID, amount); err != nil {
		if strings.Contains(err.Error(), pgerrcode.CheckViolation) {
			return ErrInsufficientFunds
		}
	}

	return nil
}
