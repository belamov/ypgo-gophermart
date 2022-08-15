package models

import "errors"

type RewardCondition struct {
	Match      string  `json:"match"`
	Reward     float64 `json:"reward"`
	RewardType string  `json:"reward_type"`
}

func (c *RewardCondition) Validate() error {
	if c.Match == "" {
		return errors.New("match required")
	}

	if c.Reward == 0 {
		return errors.New("reward required")
	}

	if c.RewardType == "" {
		return errors.New("reward_type required")
	}

	if c.RewardType != "%" && c.RewardType != "pt" {
		return errors.New("reward_type must be % or pt")
	}

	return nil
}
