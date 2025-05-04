package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/config"
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/Beigelman/nossas-despesas/internal/tests"
	"github.com/stretchr/testify/suite"
)

type GroupRepositoryTestSuite struct {
	suite.Suite
	repository    group.Repository
	ctx           context.Context
	db            db.Database
	cfg           config.Config
	testContainer *tests.PostgresContainer
	err           error
}

func TestGroupRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(GroupRepositoryTestSuite))
}

func (s *GroupRepositoryTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.testContainer, s.err = tests.StartPostgres(s.ctx)
	if s.err != nil {
		panic(s.err)
	}

	s.cfg = config.NewTestConfig(s.testContainer.Port, s.testContainer.Host)

	s.db, s.err = db.New(&s.cfg)
	s.NoError(s.err)
	s.repository = NewGroupRepository(s.db)

	s.err = s.db.MigrateUp()
	s.NoError(s.err)
}

func (s *GroupRepositoryTestSuite) TearDownSuite() {
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

func (s *GroupRepositoryTestSuite) TearDownSubTest() {
	err := s.db.Clean()
	s.NoError(err)
}

func (s *GroupRepositoryTestSuite) TestPgGroupRepo_Store() {
	id := s.repository.GetNextID()
	grp := group.New(group.Attributes{
		ID:   id,
		Name: "My Group",
	})

	err := s.repository.Store(s.ctx, grp)
	s.NoError(err)
}

func (s *GroupRepositoryTestSuite) TestPgGroupRepo_GetByID() {
	id := s.repository.GetNextID()
	expected := group.New(group.Attributes{
		ID:   id,
		Name: "My Group",
	})
	err := s.repository.Store(s.ctx, expected)
	s.NoError(err)

	actual, err := s.repository.GetByID(s.ctx, id)
	s.NoError(err)

	s.Equal(expected.ID, actual.ID)
	s.Equal(expected.Name, actual.Name)
}

func (s *GroupRepositoryTestSuite) TestPgGroupRepo_GetByName() {
	id := s.repository.GetNextID()
	expected := group.New(group.Attributes{
		ID:   id,
		Name: "My Group 2",
	})

	err := s.repository.Store(s.ctx, expected)
	s.NoError(err)

	actual, err := s.repository.GetByName(s.ctx, expected.Name)
	s.NoError(err)

	s.Equal(expected.ID, actual.ID)
	s.Equal(expected.Name, actual.Name)
}
