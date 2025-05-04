package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/modules/category/postgres"

	"github.com/Beigelman/nossas-despesas/internal/config"
	"github.com/Beigelman/nossas-despesas/internal/modules/category"
	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/Beigelman/nossas-despesas/internal/tests"
	"github.com/stretchr/testify/suite"
)

type CategoryRepositoryTestSuite struct {
	suite.Suite
	repository category.Repository

	ctx           context.Context
	db            db.Database
	cfg           config.Config
	testContainer *tests.PostgresContainer
	err           error
}

func TestCategoryRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CategoryRepositoryTestSuite))
}

func (s *CategoryRepositoryTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.testContainer, s.err = tests.StartPostgres(s.ctx)
	if s.err != nil {
		panic(s.err)
	}

	s.cfg = config.NewTestConfig(s.testContainer.Port, s.testContainer.Host)

	s.db, s.err = db.New(&s.cfg)
	s.NoError(s.err)
	s.repository = postgres.NewCategoryRepository(s.db)

	s.err = s.db.MigrateUp()
	s.NoError(s.err)
}

func (s *CategoryRepositoryTestSuite) TearDownSuite() {
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
