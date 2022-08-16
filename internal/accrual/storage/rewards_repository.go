package storage

import (
	"context"
	"errors"

	"github.com/belamov/ypgo-gophermart/internal/accrual/models"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
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

func (repo *RewardsRepository) Save(rewardCondition models.Reward) error {
	conn, err := repo.pool.Acquire(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("couldn't acquire connection from pool")
		return err
	}

	_, err = conn.Exec(
		context.Background(),
		"insert into rewards (match, reward, reward_type) values ($1, $2, $3)",
		rewardCondition.Match,
		rewardCondition.Reward,
		rewardCondition.RewardType,
	)

	conn.Release()

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.UniqueViolation {
			return NewNotUniqueError("match", err)
		}
	}

	return err
}

func (repo *RewardsRepository) GetMatchingReward(orderItem models.OrderItem) (models.Reward, error) {
	var reward models.Reward

	conn, err := repo.pool.Acquire(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("couldn't acquire connection from pool")
		return models.Reward{}, err
	}

	err = conn.QueryRow(
		context.Background(),
		"select reward, reward_type from rewards where $1 ILIKE '%' || match || '%' limit 1",
		orderItem.Description,
	).Scan(&reward.Reward, &reward.RewardType)

	conn.Release()

	if errors.Is(err, pgx.ErrNoRows) {
		return reward, nil
	}

	return reward, err
}
