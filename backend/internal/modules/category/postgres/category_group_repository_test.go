package postgres_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Beigelman/nossas-despesas/internal/modules/category"
	"github.com/Beigelman/nossas-despesas/internal/modules/category/postgres"
	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/Beigelman/nossas-despesas/internal/pkg/dbtest"
)

type CategoryGroupTestSuite struct {
	suite.Suite
	repository category.GroupRepository
	ctx        context.Context
	db         *db.Client
}

func TestCategoryGroupTestSuite(t *testing.T) {
	suite.Run(t, new(CategoryGroupTestSuite))
}

func (s *CategoryGroupTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.db = dbtest.Setup(s.ctx, s.T())
	s.repository = postgres.NewCategoryGroupRepository(s.db)
}

func (s *CategoryGroupTestSuite) TearDownSubTest() {
	err := s.db.Clean()
	s.NoError(err)
}

func (s *CategoryGroupTestSuite) TestPgCategoryRepo_Store() {
	id := s.repository.GetNextID()
	groupCategory := category.NewGroup(category.GroupAttributes{
		ID:   id,
		Name: "shopping",
		Icon: "test",
	})

	s.NoError(s.repository.Store(s.ctx, groupCategory))
}

func (s *CategoryGroupTestSuite) TestPgCategoryRepo_GetByID() {
	id := s.repository.GetNextID()
	expected := category.NewGroup(category.GroupAttributes{
		ID:   id,
		Name: "shopping",
		Icon: "test",
	})

	err := s.repository.Store(s.ctx, expected)
	s.NoError(err)

	actual, err := s.repository.GetByID(s.ctx, id)
	s.NoError(err)

	s.Equal(expected.ID, actual.ID)
	s.Equal(expected.Name, actual.Name)
}

func (s *CategoryGroupTestSuite) TestPgCategoryRepo_GetByName() {
	id := s.repository.GetNextID()
	expected := category.NewGroup(category.GroupAttributes{
		ID:   id,
		Name: "shopping2",
		Icon: "test",
	})

	err := s.repository.Store(s.ctx, expected)
	s.NoError(err)

	actual, err := s.repository.GetByName(s.ctx, "shopping2")
	s.NoError(err)

	s.Equal(expected.ID, actual.ID)
	s.Equal(expected.Name, actual.Name)
}
