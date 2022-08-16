package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReward_CalculateReward(t *testing.T) {
	tests := []struct {
		name   string
		reward Reward
		price  float64
		want   float64
	}{
		{name: "empty reward", reward: Reward{}, price: 100, want: 0},
		{name: "percent reward", reward: Reward{RewardType: PercentType, Reward: 10}, price: 100, want: 10},
		{name: "percent reward, zero reward", reward: Reward{RewardType: PercentType, Reward: 0}, price: 100, want: 0},
		{name: "absolute reward", reward: Reward{RewardType: AbsoluteType, Reward: 10}, price: 100, want: 10},
		{name: "absolute reward, zero reward", reward: Reward{RewardType: AbsoluteType, Reward: 0}, price: 100, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.reward.CalculateReward(tt.price))
		})
	}
}
