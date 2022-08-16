package storage

import (
	"context"

	"github.com/belamov/ypgo-gophermart/internal/accrual/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type RewardsRepository struct {
	pool *pgxpool.Pool
}

func NewRewardsRepository(dsn string) (*RewardsRepository, error) {
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	return &RewardsRepository{
		pool: pool,
	}, nil
}

func (r RewardsRepository) Exists(match string) (bool, error) {
	// TODO implement me
	panic("implement me")
}

func (r RewardsRepository) CreateNew(rewardCondition models.RewardCondition) error {
	// TODO implement me
	panic("implement me")
}
