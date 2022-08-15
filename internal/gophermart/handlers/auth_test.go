package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/mocks"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/models"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/services"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHandler_Register(t *testing.T) {
	user := models.User{Login: "login", HashedPassword: "hashed password"}
	validCredentials := models.Credentials{Login: "login", Password: "password"}
	takenCredentials := models.Credentials{Login: "taken", Password: "password"}
	type wantHeader struct {
		name  string
		value string
	}
	type want struct {
		statusCode int
		body       string
		header     wantHeader
	}
	tests := []struct {
		name string
		want want
		body string
	}{
		{
			name: "with valid credentials",
			want: want{
				statusCode: http.StatusOK,
				body:       "",
				header:     wantHeader{name: "Authorization", value: "Bearer token"},
			},
			body: "{\"login\": \"login\", \"password\":\"password\"}",
		},
		{
			name: "without login",
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "invalid credentials: login required",
			},
			body: "{\"password\":\"password\"}",
		},
		{
			name: "without password",
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "invalid credentials: password required",
			},
			body: "{\"login\": \"login\"}",
		},
		{
			name: "with invalid json",
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "cannot decode json",
			},
			body: "{login: login, password}",
		},
		{
			name: "with taken login",
			want: want{
				statusCode: http.StatusConflict,
				body:       "login is taken",
			},
			body: "{\"login\": \"taken\", \"password\":\"password\"}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuth := mocks.NewMockAuth(ctrl)
			mockAuth.EXPECT().AuthMiddleware().Return(emptyMiddleware).AnyTimes()

			mockAuth.EXPECT().Register(validCredentials).Return(user, nil).AnyTimes()
			mockAuth.EXPECT().Register(takenCredentials).Return(models.User{}, services.NewLoginTakenError(takenCredentials.Login, nil)).AnyTimes()
			mockAuth.EXPECT().GenerateToken(user).Return("token", nil).AnyTimes()

			mockOrders := mocks.NewMockOrdersManagerInterface(ctrl)

			mockBalance := mocks.NewMockBalanceProcessorInterface(ctrl)

			r := NewRouter(mockAuth, mockOrders, mockBalance)
			ts := httptest.NewServer(r)
			defer ts.Close()

			result, body := testRequest(t, ts, http.MethodPost, "/api/user/register", tt.body)
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.body, body)

			if tt.want.header.name != "" {
				assert.Equal(t, tt.want.header.value, result.Header.Get(tt.want.header.name))
			}
		})
	}
}

func TestHandler_Login(t *testing.T) {
	user := models.User{Login: "login", HashedPassword: "hashed password"}
	validCredentials := models.Credentials{Login: "login", Password: "password"}
	invalidCredentials := models.Credentials{Login: "invalid", Password: "password"}
	type wantHeader struct {
		name  string
		value string
	}
	type want struct {
		statusCode int
		body       string
		header     wantHeader
	}
	tests := []struct {
		name string
		want want
		body string
	}{
		{
			name: "with valid credentials",
			want: want{
				statusCode: http.StatusOK,
				body:       "",
				header:     wantHeader{name: "Authorization", value: "Bearer token"},
			},
			body: "{\"login\": \"login\", \"password\":\"password\"}",
		},
		{
			name: "without login",
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "invalid credentials: login required",
			},
			body: "{\"password\":\"password\"}",
		},
		{
			name: "without password",
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "invalid credentials: password required",
			},
			body: "{\"login\": \"login\"}",
		},
		{
			name: "with invalid json",
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "cannot decode json",
			},
			body: "{login: login, password}",
		},
		{
			name: "with invalid credentials",
			want: want{
				statusCode: http.StatusUnauthorized,
				body:       "incorrect login or password",
			},
			body: "{\"login\": \"invalid\", \"password\":\"password\"}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuth := mocks.NewMockAuth(ctrl)

			mockAuth.EXPECT().Login(validCredentials).Return(user, nil).AnyTimes()
			mockAuth.EXPECT().Login(invalidCredentials).Return(models.User{}, services.NewInvalidCredentialsError(invalidCredentials, nil)).AnyTimes()
			mockAuth.EXPECT().GenerateToken(user).Return("token", nil).AnyTimes()
			mockAuth.EXPECT().AuthMiddleware().Return(emptyMiddleware).AnyTimes()

			mockOrders := mocks.NewMockOrdersManagerInterface(ctrl)
			mockBalance := mocks.NewMockBalanceProcessorInterface(ctrl)

			r := NewRouter(mockAuth, mockOrders, mockBalance)
			ts := httptest.NewServer(r)
			defer ts.Close()

			result, body := testRequest(t, ts, http.MethodPost, "/api/user/login", tt.body)
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.body, body)

			if tt.want.header.name != "" {
				assert.Equal(t, tt.want.header.value, result.Header.Get(tt.want.header.name))
			}
		})
	}
}

func TestHandler_Middleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUsersStorage := mocks.NewMockUsersStorage(ctrl)
	password, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	require.NoError(t, err)
	user := models.User{ID: 100, HashedPassword: string(password)}
	mockUsersStorage.EXPECT().CreateNew(gomock.Any(), gomock.Any()).Return(user, nil)
	auth := services.NewAuth(mockUsersStorage, "secret")

	mockOrders := mocks.NewMockOrdersManagerInterface(ctrl)
	mockOrders.EXPECT().GetUsersOrders(gomock.Any()).Return(nil, nil).AnyTimes()

	mockBalance := mocks.NewMockBalanceProcessorInterface(ctrl)
	r := NewRouter(auth, mockOrders, mockBalance)
	ts := httptest.NewServer(r)
	defer ts.Close()

	result, _ := testRequest(t, ts, http.MethodPost, "/api/user/register", "{\"login\": \"login\",\"password\": \"password\"}")
	defer result.Body.Close()
	authHeader := result.Header.Get("Authorization")
	assert.NotEmpty(t, authHeader)

	result, _ = testRequestWithAuth(t, ts, http.MethodGet, "/api/user/orders", "", authHeader)
	defer result.Body.Close()

	assert.Equal(t, http.StatusNoContent, result.StatusCode)
}
