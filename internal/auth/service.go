package auth

import (
	"context"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/magmel48/go-musthave-diploma/internal/users"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

//go:generate mockery --name=Auth
type Auth interface {
	CreateNew(ctx context.Context, user users.User) (int64, error)
	CheckUser(ctx context.Context, user users.User) (int64, error)
}

var ErrConflict = errors.New("conflict")
var ErrInvalidCredentials = errors.New("invalid credentials")

type Service struct {
	repository users.Repository
}

func NewService(repository users.Repository) *Service {
	return &Service{repository: repository}
}

func (service *Service) CreateNew(ctx context.Context, user users.User) (int64, error) {
	u, err := service.repository.Find(ctx, user.Login)
	if err != nil {
		return 0, err
	}

	if u == nil {
		hashedPassword, err := hashPassword(user.Password)
		if err != nil {
			return 0, err
		}

		id, err := service.repository.Create(ctx, users.User{Login: user.Login, Password: hashedPassword})
		if err != nil {
			if strings.Contains(err.Error(), pgerrcode.UniqueViolation) {
				return 0, ErrConflict
			}

			return 0, err
		}

		return id, err
	}

	return 0, ErrConflict
}

func (service *Service) CheckUser(ctx context.Context, user users.User) (int64, error) {
	u, err := service.repository.Find(ctx, user.Login)
	if err != nil {
		return 0, err
	}

	if checkPasswordHash(user.Password, u.Password) {
		return u.ID, nil
	}

	return 0, ErrInvalidCredentials
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
