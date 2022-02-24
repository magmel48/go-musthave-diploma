package auth

import (
	"context"
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/magmel48/go-musthave-diploma/internal/users"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockery --name=Auth
type Auth interface {
	CreateNew(ctx context.Context, user users.User) (*users.User, error)
	CheckUser(ctx context.Context, user users.User) (int64, error)
	StoreUser(ctx *gin.Context, user users.User) error
}

var ErrInvalidCredentials = errors.New("invalid credentials")

const UserIDKey = "user_id"

type Service struct {
	repository users.Repository
}

func NewService(repository users.Repository) *Service {
	return &Service{repository: repository}
}

func (service *Service) CreateNew(ctx context.Context, user users.User) (*users.User, error) {
	u, err := service.repository.Find(ctx, user.Login)
	if err != nil {
		return nil, err
	}

	if u == nil {
		hashedPassword, err := hashPassword(user.Password)
		if err != nil {
			return nil, err
		}

		u, err = service.repository.Create(ctx, users.User{Login: user.Login, Password: hashedPassword})
		if err != nil {
			return nil, err
		}

		return u, err
	}

	return nil, users.ErrConflict
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

func (service *Service) StoreUser(ctx *gin.Context, user users.User) error {
	session := sessions.Default(ctx)
	session.Set(UserIDKey, user.ID)
	err := session.Save()

	return err
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
