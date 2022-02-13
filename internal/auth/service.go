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
	CheckUser(ctx context.Context, login string, password string) error
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
			if strings.Index(err.Error(), pgerrcode.UniqueViolation) != -1 {
				return 0, ErrConflict
			}

			return 0, err
		}

		return id, err
	}

	return 0, ErrConflict
}

func (service *Service) CheckUser(ctx context.Context, login string, password string) error {
	user, err := service.repository.Find(ctx, login)
	if err != nil {
		return err
	}

	if checkPasswordHash(password, user.Password) {
		return nil
	}

	return ErrInvalidCredentials
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
