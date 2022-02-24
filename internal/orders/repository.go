package orders

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/magmel48/go-musthave-diploma/internal/logger"
	"strings"
)

//go:generate mockery --name=Repository
type Repository interface {
	Create(ctx context.Context, orderNumber string, userID int64) (*Order, error)
	FindByUser(ctx context.Context, orderNumber string, userID int64) (*Order, error)
	ListByUser(ctx context.Context, userID int64) ([]Order, error)
	FindUnprocessedOrders(ctx context.Context) ([]Order, error)
	Update(ctx context.Context, order Order) (int64, error)
}

var ErrConflict = errors.New("conflict")

// PostgreSQLRepository implements orders.Repository using PostgreSQL.
type PostgreSQLRepository struct {
	db *sql.DB
}

// NewPostgreSQLRepository creates new PostgreSQLRepository.
func NewPostgreSQLRepository(db *sql.DB) *PostgreSQLRepository {
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
		if strings.Contains(err.Error(), pgerrcode.UniqueViolation) {
			return nil, ErrConflict
		}

		return nil, err
	}

	return &result, nil
}

// FindByUser finds order belongs to specified user by order number.
func (repository *PostgreSQLRepository) FindByUser(ctx context.Context, orderNumber string, userID int64) (*Order, error) {
	var result Order

	if err := repository.db.QueryRowContext(
		ctx, `SELECT "id" FROM "orders" WHERE "number" = $1 AND "user_id" = $2`,
		orderNumber, userID).Scan(
		&result.ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &result, nil
}

// ListByUser finds all orders belong to specified user.
func (repository *PostgreSQLRepository) ListByUser(ctx context.Context, userID int64) ([]Order, error) {
	rows, err := repository.db.QueryContext(
		ctx,
		`SELECT "number", "status", "accrual", "uploaded_at" FROM "orders" WHERE "user_id" = $1 ORDER BY "uploaded_at" ASC`, userID)
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
		if err = rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadedAt); err != nil {
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

// FindUnprocessedOrders finds all orders with non-final status.
func (repository *PostgreSQLRepository) FindUnprocessedOrders(ctx context.Context) ([]Order, error) {
	rows, err := repository.db.QueryContext(
		ctx, `SELECT "id", "number", "status", "user_id" FROM "orders" WHERE "status" = ANY($1)`, UnprocessedStatuses)
	if err != nil {
		return nil, err
	}

	orders := make([]Order, 0)
	for rows.Next() {
		var order Order
		if err = rows.Scan(&order.ID, &order.Number, &order.Status, &order.UserID); err != nil {
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

func (repository *PostgreSQLRepository) Update(ctx context.Context, order Order) (int64, error) {
	result, err := repository.db.ExecContext(
		ctx,
		`UPDATE "orders" SET "status" = $1, "accrual" = $2 WHERE "id" = $3`, order.Status, order.Accrual, order.ID)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
