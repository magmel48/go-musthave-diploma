package orders

import (
	"context"
	"database/sql"
	"errors"
	"github.com/magmel48/go-musthave-diploma/internal/logger"
)

//go:generate mockery --name=Repository
type Repository interface {
	Create(ctx context.Context, orderNumber string, userID int64) (*Order, error)
	FindUserOrder(ctx context.Context, orderNumber string, userID int64) (*Order, error)
	FindUserOrders(ctx context.Context, userID int64) ([]Order, error)
}

var ErrConflict = errors.New("conflict")

// PostgreSQLRepository implements orders.Repository using PostgreSQL.
type PostgreSQLRepository struct {
	db *sql.DB
}

// NewRepository creates new PostgreSQLRepository.
func NewRepository(db *sql.DB) *PostgreSQLRepository {
	return &PostgreSQLRepository{db: db}
}

// Create creates new order.
func (repository *PostgreSQLRepository) Create(ctx context.Context, orderNumber string, userID int64) (*Order, error) {
	var result Order

	if err := repository.db.QueryRowContext(
		ctx,
		`INSERT INTO "orders" ("number", "user_id") VALUES ($1, $2) RETURNING "id", "number", "status", "user_id", "uploaded_at"`,
		orderNumber, userID).Scan(
		&result.ID, &result.Number, &result.Status, &result.UserID, &result.UploadedAt); err != nil {
		return nil, err
	}

	return &result, nil
}

// FindUserOrder finds order belongs to specified user by order number.
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

// FindUserOrders finds all orders belong to specified user.
func (repository *PostgreSQLRepository) FindUserOrders(ctx context.Context, userID int64) ([]Order, error) {
	rows, err := repository.db.QueryContext(
		ctx, `SELECT "number", "status", "accrual", "uploaded_at" WHERE "user_id" = $1`, userID)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}()

	orders := make([]Order, 0)
	for rows.Next() {
		var order Order
		err := rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadedAt)
		if err != nil {
			return orders, nil
		}

		orders = append(orders, order)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return orders, nil
}
