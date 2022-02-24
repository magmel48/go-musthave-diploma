package auth

import (
	"context"
	"github.com/magmel48/go-musthave-diploma/internal/users"
	"github.com/magmel48/go-musthave-diploma/internal/users/mocks"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

func TestService_CreateNew(t *testing.T) {
	type fields struct {
		repository users.Repository
	}
	type args struct {
		ctx  context.Context
		user users.User
	}

	user := users.User{Login: "login", Password: "password"}

	conflictingUserRepository := mocks.Repository{}
	conflictingUserRepository.On(
		"Find", mock.Anything, mock.Anything).Return(&users.User{Login: "login"}, nil)

	emptyUserRepository := mocks.Repository{}
	emptyUserRepository.On(
		"Find", mock.Anything, mock.Anything).Return(nil, nil)
	emptyUserRepository.On("Create", mock.Anything, mock.Anything).Return(&user, nil)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *users.User
		wantErr bool
	}{
		{
			name:    "should not create a new user if login already taken",
			fields:  fields{repository: &conflictingUserRepository},
			args:    args{ctx: context.TODO(), user: user},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "should create new user if login was not taken before",
			fields:  fields{repository: &emptyUserRepository},
			args:    args{ctx: context.TODO(), user: user},
			want:    &user,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &Service{
				repository: tt.fields.repository,
			}

			got, err := service.CreateNew(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateNew() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateNew() got = %v, want %v", got, tt.want)
			}
		})
	}
}
