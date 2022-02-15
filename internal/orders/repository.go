package orders

import (
	"context"
	"database/sql"
	"errors"
)

//go:generate mockery --name=Repository
type Repository interface {
	Create(ctx context.Context, order Order) (*Order, error)
	FindUserOrder(ctx context.Context, orderNumber string, userID int64) (*Order, error)
}

var ErrConflict = errors.New("conflict")

type PostgreSQLRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *PostgreSQLRepository {
	return &PostgreSQLRepository{db: db}
}

func (repository *PostgreSQLRepository) Create(ctx context.Context, order Order) (*Order, error) {
	var result Order

	if err := repository.db.QueryRowContext(
		ctx,
		`INSERT INTO "orders" ("number", "user_id") VALUES ($1, $2) RETURNING "id", "number", "status", "user_id", "uploaded_at"`,
		order.Number, order.UserID).Scan(
		&result.ID, &result.Number, &result.Status, &result.UserID, &result.UploadedAt); err != nil {
		return nil, err
	}

	return &result, nil
}

func (repository *PostgreSQLRepository) FindUserOrder(ctx context.Context, orderNumber string, userID int64) (*Order, error) {
	var result Order

	if err := repository.db.QueryRowContext(
		ctx, `SELECT "id" FROM "orders" WHERE "number" = $1 AND "user_id" = $2`,
		orderNumber, userID).Scan(
		&result.ID); err != nil {
		return nil, err
	}

	return &result, nil
}
