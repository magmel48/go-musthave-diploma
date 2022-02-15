package users

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
	Find(ctx context.Context, login string) (*User, error)
	Create(ctx context.Context, user User) (*User, error)
}

var ErrConflict = errors.New("conflict")

type PostgreSQLRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *PostgreSQLRepository {
	return &PostgreSQLRepository{db: db}
}

func (repository *PostgreSQLRepository) Find(ctx context.Context, login string) (*User, error) {
	rows, err := repository.db.QueryContext(
		ctx, `SELECT "id", "login", "password" FROM "users" WHERE "login" = $1`, login)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}()

	if rows.Next() {
		user := new(User)

		if err = rows.Scan(&user.ID, &user.Login, &user.Password); err != nil {
			return nil, err
		}

		return user, nil
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (repository *PostgreSQLRepository) Create(ctx context.Context, user User) (*User, error) {
	var result User

	if err := repository.db.QueryRowContext(
		ctx,
		`INSERT INTO "users" ("login", "password") VALUES ($1, $2) RETURNING "id", "login", "password"`,
		user.Login,
		user.Password).Scan(&result.ID, &result.Login, &result.Password); err != nil {
		if strings.Contains(err.Error(), pgerrcode.UniqueViolation) {
			return nil, ErrConflict
		}

		return &result, err
	}

	return &result, nil
}
