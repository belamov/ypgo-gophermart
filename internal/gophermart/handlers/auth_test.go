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
			name: "with valid validCredentials",
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

			mockAuth := mocks.NewMockAuthenticator(ctrl)

			mockAuth.EXPECT().Register(validCredentials).Return(user, nil).AnyTimes()
			mockAuth.EXPECT().Register(takenCredentials).Return(models.User{}, services.NewLoginTakenError(takenCredentials.Login, nil)).AnyTimes()
			mockAuth.EXPECT().GenerateToken(user).Return("token", nil).AnyTimes()

			r := NewRouter(mockAuth)
			ts := httptest.NewServer(r)
			defer ts.Close()

			result, body := testRequest(t, ts, http.MethodPost, "/api/user/register", tt.body, nil)
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.body, body)

			if tt.want.header.name != "" {
				assert.Equal(t, tt.want.header.value, result.Header.Get(tt.want.header.name))
			}
		})
	}
}
