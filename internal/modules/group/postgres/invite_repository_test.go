package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/pkg/config"

	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/modules/group/postgres"
	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/Beigelman/nossas-despesas/internal/tests"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type GroupInviteRepositoryTestSuite struct {
	suite.Suite
	repository    group.InviteRepository
	ctx           context.Context
	db            *db.Client
	cfg           config.Config
	testContainer *tests.PostgresContainer
	err           error
}

func TestGroupInviteRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(GroupInviteRepositoryTestSuite))
}

func (s *GroupInviteRepositoryTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.testContainer, s.err = tests.StartPostgres(s.ctx)
	if s.err != nil {
		panic(s.err)
	}

	s.cfg = config.NewTestConfig(s.testContainer.Port, s.testContainer.Host)

	s.db, s.err = db.New(&s.cfg)
	s.NoError(s.err)

	s.repository = postgres.NewGroupInviteRepository(s.db)

	s.err = s.db.MigrateUp()
	s.NoError(s.err)
}

func (s *GroupInviteRepositoryTestSuite) TearDownSuite() {
	s.err = s.db.MigrateDown()
	s.NoError(s.err)

	s.err = s.db.Close()
	s.NoError(s.err)

	duration := 10 * time.Second
	s.err = s.testContainer.Stop(s.ctx, &duration)
	if s.err != nil {
		panic(s.err)
	}
}

func (s *GroupInviteRepositoryTestSuite) TearDownTest() {
	err := s.db.Clean("group_invites")
	s.NoError(err)
}

func (s *GroupInviteRepositoryTestSuite) TestPgGroupInviteRepo_Store() {
	id := s.repository.GetNextID()
	groupInvite := group.NewInvite(group.InviteAttributes{
		ID:        id,
		GroupID:   group.ID{Value: 1},
		Token:     uuid.NewString(),
		Email:     "john@email.com",
		ExpiresAt: time.Now().Add(time.Hour * 24),
	})

	s.NoError(s.repository.Store(s.ctx, groupInvite))
}

func (s *GroupInviteRepositoryTestSuite) TestPgGroupInviteRepo_GetByID() {
	id := s.repository.GetNextID()
	expected := group.NewInvite(group.InviteAttributes{
		ID:        id,
		GroupID:   group.ID{Value: 1},
		Token:     uuid.NewString(),
		Email:     "john@email.com",
		ExpiresAt: time.Now().Add(time.Hour * 24),
	})
	s.NoError(s.repository.Store(s.ctx, expected))

	actual, err := s.repository.GetByID(s.ctx, id)
	s.NoError(err)

	s.Equal(expected.ID, actual.ID)
	s.Equal(expected.Email, actual.Email)
	s.Equal(expected.Token, actual.Token)
}

func (s *GroupInviteRepositoryTestSuite) TestPgGroupInviteRepo_GetByToken() {
	id := s.repository.GetNextID()
	token := uuid.NewString()
	expected := group.NewInvite(group.InviteAttributes{
		ID:        id,
		GroupID:   group.ID{Value: 1},
		Token:     token,
		Email:     "john@email.com",
		ExpiresAt: time.Now().Add(time.Hour * 24),
	})
	s.NoError(s.repository.Store(s.ctx, expected))

	actual, err := s.repository.GetByToken(s.ctx, token)
	s.NoError(err)

	s.Equal(expected.ID, actual.ID)
	s.Equal(expected.Email, actual.Email)
	s.Equal(expected.Token, actual.Token)
}

func (s *GroupInviteRepositoryTestSuite) TestPgGroupInviteRepo_GetByEmail() {
	email := "john@email.com"
	s.NoError(s.repository.Store(s.ctx, group.NewInvite(group.InviteAttributes{
		ID:        s.repository.GetNextID(),
		GroupID:   group.ID{Value: 1},
		Token:     uuid.NewString(),
		Email:     email,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	})))

	s.NoError(s.repository.Store(s.ctx, group.NewInvite(group.InviteAttributes{
		ID:        s.repository.GetNextID(),
		GroupID:   group.ID{Value: 1},
		Token:     uuid.NewString(),
		Email:     email,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	})))

	actual, err := s.repository.GetGroupInvitesByEmail(s.ctx, group.ID{Value: 1}, email)
	s.NoError(err)

	s.Len(actual, 2)
	s.Equal(email, actual[0].Email)
	s.Equal(email, actual[1].Email)
}
