package services

import (
	"errors"
	"fmt"
	"testing"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/mocks"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/models"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/storage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAuth_Register(t *testing.T) {
	tests := []struct {
		name        string
		credentials models.Credentials
		want        models.User
		wantErr     bool
	}{
		{
			name: "it registers user",
			credentials: models.Credentials{
				Login:    "login",
				Password: "password",
			},
			want: models.User{
				Login: "login",
			},
			wantErr: false,
		},
		{
			name: "it doesnt register user with not unique login",
			credentials: models.Credentials{
				Login:    "existing login",
				Password: "password",
			},
			want:    models.User{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsers := mocks.NewMockUsers(ctrl)

			mockUsers.EXPECT().CreateNew("login", gomock.Not("password")).
				Return(models.User{Login: "login", HashedPassword: "some hash"}, nil).
				AnyTimes()
			mockUsers.EXPECT().CreateNew("existing login", gomock.Any()).
				Return(models.User{}, storage.NewNotUniqueError("login", errors.New(""))).
				AnyTimes()

			a := &Auth{
				userRepo: mockUsers,
			}
			registeredUser, err := a.Register(tt.credentials)
			if tt.wantErr {
				assert.Error(t, err)
				var loginTakenError *LoginTakenError
				assert.ErrorAs(t, err, &loginTakenError)
				fmt.Println(loginTakenError.Error())
				assert.Contains(t, loginTakenError.Error(), "login is already taken: existing login")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, registeredUser.Login, tt.credentials.Login)
			}
		})
	}
}
