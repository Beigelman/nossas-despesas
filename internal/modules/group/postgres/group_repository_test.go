package postgres

import (
	"context"
	"testing"

	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/Beigelman/nossas-despesas/internal/pkg/dbtest"
	"github.com/stretchr/testify/suite"
)

type GroupRepositoryTestSuite struct {
	suite.Suite
	repository group.Repository
	ctx        context.Context
	db         *db.Client
}

func TestGroupRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(GroupRepositoryTestSuite))
}

func (s *GroupRepositoryTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.db = dbtest.Setup(s.ctx, s.T())
	s.repository = NewGroupRepository(s.db)
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
