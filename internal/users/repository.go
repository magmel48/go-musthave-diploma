package users

import (
	"context"
	"database/sql"
	"github.com/magmel48/go-musthave-diploma/internal/logger"
)

//go:generate mockery --name=Repository
type Repository interface {
	Find(ctx context.Context, login string) (*User, error)
	Create(ctx context.Context, user User) (int64, error)
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repository *UserRepository) Find(ctx context.Context, login string) (*User, error) {
	rows, err := repository.db.QueryContext(ctx, `SELECT "id", "login", "password" FROM "users"`)
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

		if err = rows.Scan(user.ID, user.Login, user.Password); err != nil {
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

func (repository *UserRepository) Create(ctx context.Context, user User) (int64, error) {
	var userID int64

	if err := repository.db.QueryRowContext(
		ctx,
		`INSERT INTO "users" (login, password) VALUES ($1, $2) RETURNING "id"`,
		user.Login,
		user.Password).Scan(&userID); err != nil {
		return 0, err
	}

	return userID, nil
}
