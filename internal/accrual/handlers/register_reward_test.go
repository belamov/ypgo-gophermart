package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/belamov/ypgo-gophermart/internal/accrual/mocks"
	"github.com/belamov/ypgo-gophermart/internal/accrual/models"
	"github.com/belamov/ypgo-gophermart/internal/accrual/storage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_RegisterReward(t *testing.T) {
	type want struct {
		statusCode int
	}
	validReward := models.Reward{
		Match:      "match",
		RewardType: "%",
		Reward:     10.0,
	}
	alreadyRegisteredReward := models.Reward{
		Match:      "registered match",
		RewardType: "%",
		Reward:     10.0,
	}
	tests := []struct {
		name   string
		want   want
		reward models.Reward
	}{
		{
			name: "it accepts new reward",
			want: want{
				statusCode: http.StatusOK,
			},
			reward: validReward,
		},
		{
			name: "it responds with 409 when reward is already registered",
			want: want{
				statusCode: http.StatusConflict,
			},
			reward: alreadyRegisteredReward,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockOrderManager := mocks.NewMockOrderManagementInterface(ctrl)

			mockRewards := mocks.NewMockRewardsStorage(ctrl)
			mockRewards.EXPECT().Save(validReward).Return(nil).AnyTimes()
			mockRewards.EXPECT().Save(alreadyRegisteredReward).Return(storage.NewNotUniqueError("match", nil)).AnyTimes()

			r := NewRouter(mockOrderManager, mockRewards)
			ts := httptest.NewServer(r)
			defer ts.Close()

			requestJSON, err := json.Marshal(tt.reward)
			require.NoError(t, err)
			result, _ := testRequest(
				t,
				ts,
				http.MethodPost,
				"/api/goods",
				string(requestJSON),
			)
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}
