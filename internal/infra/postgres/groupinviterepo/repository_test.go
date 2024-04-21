package groupinviterepo

import (
	"context"
	"testing"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	"github.com/Beigelman/nossas-despesas/internal/domain/repository"
	"github.com/google/uuid"

	"github.com/Beigelman/nossas-despesas/internal/config"
	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/Beigelman/nossas-despesas/internal/tests"
	"github.com/stretchr/testify/suite"
)

type PgGroupInviteRepoTestSuite struct {
	suite.Suite
	repository    repository.GroupInviteRepository
	ctx           context.Context
	db            db.Database
	cfg           config.Config
	testContainer *tests.PostgresContainer
	err           error
}

func TestPgGroupInviteRepoTestSuite(t *testing.T) {
	suite.Run(t, new(PgGroupInviteRepoTestSuite))
}

func (s *PgGroupInviteRepoTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.testContainer, s.err = tests.StartPostgres(s.ctx)
	if s.err != nil {
		panic(s.err)
	}

	s.cfg = config.NewTestConfig(s.testContainer.Port, s.testContainer.Host, "postgres")

	s.db = db.New(&s.cfg)
	s.repository = NewPGRepository(s.db)

	s.err = s.db.MigrateUp()
	s.NoError(s.err)
}

func (s *PgGroupInviteRepoTestSuite) TearDownSuite() {
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

func (s *PgGroupInviteRepoTestSuite) TearDownTest() {
	err := s.db.Clean("group_invites")
	s.NoError(err)
}

func (s *PgGroupInviteRepoTestSuite) TestPgGroupInviteRepo_Store() {
	id := s.repository.GetNextID()
	groupInvite := entity.NewGroupInvite(entity.GroupInviteParams{
		ID:        id,
		GroupID:   entity.GroupID{Value: 1},
		Token:     uuid.NewString(),
		Email:     "john@email.com",
		ExpiresAt: time.Now().Add(time.Hour * 24),
	})

	s.NoError(s.repository.Store(s.ctx, groupInvite))
}

func (s *PgGroupInviteRepoTestSuite) TestPgGroupInviteRepo_GetByID() {
	id := s.repository.GetNextID()
	expected := entity.NewGroupInvite(entity.GroupInviteParams{
		ID:        id,
		GroupID:   entity.GroupID{Value: 1},
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

func (s *PgGroupInviteRepoTestSuite) TestPgGroupInviteRepo_GetByToken() {
	id := s.repository.GetNextID()
	token := uuid.NewString()
	expected := entity.NewGroupInvite(entity.GroupInviteParams{
		ID:        id,
		GroupID:   entity.GroupID{Value: 1},
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

func (s *PgGroupInviteRepoTestSuite) TestPgGroupInviteRepo_GetByEmail() {
	email := "john@email.com"
	s.NoError(s.repository.Store(s.ctx, entity.NewGroupInvite(entity.GroupInviteParams{
		ID:        s.repository.GetNextID(),
		GroupID:   entity.GroupID{Value: 1},
		Token:     uuid.NewString(),
		Email:     email,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	})))

	s.NoError(s.repository.Store(s.ctx, entity.NewGroupInvite(entity.GroupInviteParams{
		ID:        s.repository.GetNextID(),
		GroupID:   entity.GroupID{Value: 1},
		Token:     uuid.NewString(),
		Email:     email,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	})))

	actual, err := s.repository.GetGroupInvitesByEmail(s.ctx, entity.GroupID{Value: 1}, email)
	s.NoError(err)

	s.Len(actual, 2)
	s.Equal(email, actual[0].Email)
	s.Equal(email, actual[1].Email)
}
