package storage

import (
	"context"

	"github.com/belamov/ypgo-gophermart/internal/accrual/models"
	"github.com/jackc/pgx/v4"
)

type RewardsRepository struct {
	conn *pgx.Conn
}

func NewRewardsRepository(dsn string) (*RewardsRepository, error) {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	return &RewardsRepository{
		conn: conn,
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
