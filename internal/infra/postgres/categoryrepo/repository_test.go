package categoryrepo

import (
	"context"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"github.com/Beigelman/ludaapi/internal/domain/repository"
	"testing"
	"time"

	"github.com/Beigelman/ludaapi/internal/config"
	"github.com/Beigelman/ludaapi/internal/pkg/db"
	"github.com/Beigelman/ludaapi/internal/tests"
	"github.com/stretchr/testify/suite"
)

var migrationPath string = "file:///Users/danielbeigelman/mydev/go-luda/server/database/migrations"

type PgCategoryRepoTestSuite struct {
	suite.Suite
	repository repository.CategoryRepository

	ctx           context.Context
	db            db.Database
	cfg           config.Config
	testContainer *tests.PostgresContainer
	err           error
}

func TestPgUserRepoTestSuite(t *testing.T) {
	suite.Run(t, new(PgCategoryRepoTestSuite))
}

func (s *PgCategoryRepoTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.testContainer, s.err = tests.StartPostgres(s.ctx)
	if s.err != nil {
		panic(s.err)
	}

	s.cfg = config.NewTestConfig(s.testContainer.Port, s.testContainer.Host)

	s.db = db.New(&s.cfg)
	s.repository = NewPGRepository(s.db)

	s.err = s.db.MigrateUp(migrationPath)
	s.NoError(s.err)

	_, err := s.db.Client().Exec(`
   		INSERT INTO category_groups (id, name, icon, created_at, updated_at, deleted_at, version)
		VALUES (1, 'category group', 'test', now(), now(), now(), 0)`,
	)
	s.NoError(err)
}

func (s *PgCategoryRepoTestSuite) TearDownSuite() {
	s.err = s.db.MigrateDown(migrationPath)
	s.NoError(s.err)

	s.err = s.db.Close()
	s.NoError(s.err)

	duration := 10 * time.Second
	s.err = s.testContainer.Stop(s.ctx, &duration)
	if s.err != nil {
		panic(s.err)
	}
}

func (s *PgCategoryRepoTestSuite) TearDownSubTest() {
	err := s.db.Clean()
	s.NoError(err)
}

func (s *PgCategoryRepoTestSuite) TestPgCategoryRepo_Store() {
	id := s.repository.GetNextID()
	category := entity.NewCategory(entity.CategoryParams{
		ID:              id,
		Name:            "shopping",
		Icon:            "test",
		CategoryGroupID: entity.CategoryGroupID{Value: 1},
	})

	s.NoError(s.repository.Store(s.ctx, category))
}

func (s *PgCategoryRepoTestSuite) TestPgCategoryRepo_GetByID() {
	id := s.repository.GetNextID()
	expected := entity.NewCategory(entity.CategoryParams{
		ID:              id,
		Name:            "shopping",
		Icon:            "test",
		CategoryGroupID: entity.CategoryGroupID{Value: 1},
	})

	err := s.repository.Store(s.ctx, expected)
	s.NoError(err)

	actual, err := s.repository.GetByID(s.ctx, id)
	s.NoError(err)

	s.Equal(expected.ID, actual.ID)
	s.Equal(expected.Name, actual.Name)
	s.Equal(expected.GroupCategoryID, actual.GroupCategoryID)
}

func (s *PgCategoryRepoTestSuite) TestPgCategoryRepo_GetByName() {
	id := s.repository.GetNextID()
	expected := entity.NewCategory(entity.CategoryParams{
		ID:              id,
		Name:            "shopping2",
		Icon:            "test",
		CategoryGroupID: entity.CategoryGroupID{Value: 1},
	})

	err := s.repository.Store(s.ctx, expected)
	s.NoError(err)

	actual, err := s.repository.GetByName(s.ctx, "shopping2")
	s.NoError(err)

	s.Equal(expected.ID, actual.ID)
	s.Equal(expected.Name, actual.Name)
	s.Equal(expected.GroupCategoryID, actual.GroupCategoryID)
}
