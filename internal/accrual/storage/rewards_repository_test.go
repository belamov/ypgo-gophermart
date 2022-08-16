package storage

import (
	"context"
	"os"
	"testing"

	"github.com/belamov/ypgo-gophermart/internal/accrual/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RewardsRepositoryTestSuite struct {
	suite.Suite
	rewardsRepo *RewardsRepository
}

func TestRewardsRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RewardsRepositoryTestSuite))
}

func (s *RewardsRepositoryTestSuite) SetupSuite() {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@db_accrual:5432/accrual?sslmode=disable"
	}
	rewardsRepo, err := NewRewardsRepository(dsn)
	require.NoError(s.T(), err)
	s.rewardsRepo = rewardsRepo

	err = RunMigrations(dsn)
	require.NoError(s.T(), err)
}

func (s *RewardsRepositoryTestSuite) TearDownTest() {
	conn, err := s.rewardsRepo.pool.Acquire(context.Background())
	require.NoError(s.T(), err)
	_, _ = conn.Exec(context.Background(), "truncate table rewards cascade")
	conn.Release()
}

func (s *RewardsRepositoryTestSuite) TestSave() {
	reward := models.RewardCondition{
		Match:      "match",
		Reward:     1.5,
		RewardType: "pt",
	}
	err := s.rewardsRepo.Save(reward)
	assert.NoError(s.T(), err)

	err = s.rewardsRepo.Save(reward)
	var nue *NotUniqueError
	assert.ErrorAs(s.T(), err, &nue)
}
