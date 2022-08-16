package models

import (
	"errors"
	"fmt"
)

const (
	AbsoluteType = "pt"
	PercentType  = "%"
)

type Reward struct {
	Match      string  `json:"match"`
	Reward     float64 `json:"reward"`
	RewardType string  `json:"reward_type"`
}

func (c *Reward) Validate() error {
	if c.Match == "" {
		return errors.New("match required")
	}

	if c.Reward == 0 {
		return errors.New("reward required")
	}

	if c.RewardType == "" {
		return errors.New("reward_type required")
	}

	if c.RewardType != PercentType && c.RewardType != AbsoluteType {
		return fmt.Errorf("reward_type must be %s or %s", PercentType, AbsoluteType)
	}

	return nil
}

func (c *Reward) CalculateReward(itemPrice float64) float64 {
	if c.RewardType == PercentType {
		return itemPrice / 100 * c.Reward
	}

	if c.RewardType == AbsoluteType {
		return c.Reward
	}

	return 0
}
