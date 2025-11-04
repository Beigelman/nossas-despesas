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

type CategoryRepositoryTestSuite struct {
	suite.Suite
	repository category.Repository
	ctx        context.Context
	db         *db.Client
}

func TestCategoryRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CategoryRepositoryTestSuite))
}

func (s *CategoryRepositoryTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.db = dbtest.Setup(s.ctx, s.T())
	s.repository = postgres.NewCategoryRepository(s.db)
}

func (s *CategoryRepositoryTestSuite) TearDownSubTest() {
	err := s.db.Clean()
	s.NoError(err)
}

func (s *CategoryRepositoryTestSuite) TestPgCategoryRepo_Store() {
	id := s.repository.GetNextID()
	catgry := category.New(category.Attributes{
		ID:              id,
		Name:            "shopping",
		Icon:            "test",
		CategoryGroupID: category.GroupID{Value: 1},
	})

	s.NoError(s.repository.Store(s.ctx, catgry))
}

func (s *CategoryRepositoryTestSuite) TestPgCategoryRepo_GetByID() {
	id := s.repository.GetNextID()
	expected := category.New(category.Attributes{
		ID:              id,
		Name:            "shopping",
		Icon:            "test",
		CategoryGroupID: category.GroupID{Value: 1},
	})

	err := s.repository.Store(s.ctx, expected)
	s.NoError(err)

	actual, err := s.repository.GetByID(s.ctx, id)
	s.NoError(err)

	s.Equal(expected.ID, actual.ID)
	s.Equal(expected.Name, actual.Name)
	s.Equal(expected.GroupCategoryID, actual.GroupCategoryID)
}

func (s *CategoryRepositoryTestSuite) TestPgCategoryRepo_GetByName() {
	id := s.repository.GetNextID()
	expected := category.New(category.Attributes{
		ID:              id,
		Name:            "shopping2",
		Icon:            "test",
		CategoryGroupID: category.GroupID{Value: 1},
	})

	err := s.repository.Store(s.ctx, expected)
	s.NoError(err)

	actual, err := s.repository.GetByName(s.ctx, "shopping2")
	s.NoError(err)

	s.Equal(expected.ID, actual.ID)
	s.Equal(expected.Name, actual.Name)
	s.Equal(expected.GroupCategoryID, actual.GroupCategoryID)
}
