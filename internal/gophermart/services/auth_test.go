package services

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/mocks"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/models"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/storage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
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

			mockUsers := mocks.NewMockUsersStorage(ctrl)

			mockUsers.EXPECT().CreateNew("login", gomock.Not("password")).
				Return(models.User{Login: "login", HashedPassword: "some hash"}, nil).
				AnyTimes()
			mockUsers.EXPECT().CreateNew("existing login", gomock.Any()).
				Return(models.User{}, storage.NewNotUniqueError("login", errors.New(""))).
				AnyTimes()

			a := &JWTAuth{
				UserRepo: mockUsers,
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

func TestAuth_Login(t *testing.T) {
	validCredentials := models.Credentials{Login: "login", Password: "password"}
	invalidLogin := models.Credentials{Login: "invalid", Password: "password"}
	invalidPassword := models.Credentials{Login: "login", Password: "invalid"}
	tests := []struct {
		name        string
		credentials models.Credentials
		want        models.User
		wantErr     bool
	}{
		{
			name:        "it logins user",
			credentials: validCredentials,
			want: models.User{
				Login: "login",
			},
			wantErr: false,
		},
		{
			name:        "it doesnt login user with invalid login",
			credentials: invalidLogin,
			want:        models.User{},
			wantErr:     true,
		},
		{
			name:        "it doesnt login user with invalid password",
			credentials: invalidPassword,
			want:        models.User{},
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsers := mocks.NewMockUsersStorage(ctrl)

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(validCredentials.Password), bcrypt.DefaultCost)
			assert.NoError(t, err)
			mockUsers.EXPECT().FindByLogin(validCredentials.Login).
				Return(models.User{ID: 1, Login: validCredentials.Login, HashedPassword: string(hashedPassword)}, nil).
				AnyTimes()
			mockUsers.EXPECT().FindByLogin(invalidLogin.Login).
				Return(models.User{}, nil).
				AnyTimes()

			a := &JWTAuth{
				UserRepo: mockUsers,
			}
			loggedUser, err := a.Login(tt.credentials)
			if tt.wantErr {
				assert.Error(t, err)
				var invalidCredentialsError *InvalidCredentialsError
				assert.ErrorAs(t, err, &invalidCredentialsError)
				fmt.Println(invalidCredentialsError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, loggedUser.Login, tt.credentials.Login)
			}
		})
	}
}

func TestAuth_GenerateToken(t *testing.T) {
	tests := []struct {
		name    string
		user    models.User
		wantErr bool
	}{
		{
			name:    "it generates token with correct user_id",
			user:    models.User{ID: 1},
			wantErr: false,
		},
		{
			name:    "it does not generate token when user id is not set",
			user:    models.User{Login: "login"},
			wantErr: true,
		},
	}
	key := "secret"
	auth := NewAuth(nil, key)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenString, err := auth.GenerateToken(tt.user)
			if !tt.wantErr {
				assert.NoError(t, err)

				token, errDecode := auth.tokenAuth.Decode(tokenString)
				assert.NoError(t, errDecode)

				userClaim, ok := token.Get("user_id")
				assert.True(t, ok)

				parsedTokenString, errParse := strconv.Atoi(fmt.Sprintf("%v", userClaim))
				require.NoError(t, errParse)
				assert.Equal(t, tt.user.ID, parsedTokenString)

				assert.Greater(t, token.Expiration(), time.Now())
				assert.GreaterOrEqual(t, token.IssuedAt(), time.Unix(0, time.Now().Unix()/1e6*1e6))
			} else {
				assert.Error(t, err)
			}
		})
	}
}
