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

var migrationPath string = "file:///Users/danielbeigelman/mydev/go-luda/api/database/migrations"

type PgUserRepoTestSuite struct {
	suite.Suite
	repository    repository.CategoryRepository
	ctx           context.Context
	db            db.Database
	cfg           config.Config
	testContainer *tests.PostgresContainer
	err           error
}

func TestPgUserRepoTestSuite(t *testing.T) {
	suite.Run(t, new(PgUserRepoTestSuite))
}

func (s *PgUserRepoTestSuite) SetupSuite() {
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
}

func (s *PgUserRepoTestSuite) TearDownSuite() {
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

func (s *PgUserRepoTestSuite) TearDownSubTest() {
	err := s.db.Clean()
	s.NoError(err)
}

func (s *PgUserRepoTestSuite) TestPgCategoryRepo_Store() {
	id := s.repository.GetNextID()
	category := entity.NewCategory(entity.CategoryParams{
		ID:   id,
		Name: "shopping",
	})

	s.NoError(s.repository.Store(s.ctx, category))
}

func (s *PgUserRepoTestSuite) TestPgCategoryRepo_GetByID() {
	id := s.repository.GetNextID()
	expected := entity.NewCategory(entity.CategoryParams{
		ID:   id,
		Name: "shopping",
	})

	err := s.repository.Store(s.ctx, expected)
	s.NoError(err)

	actual, err := s.repository.GetByID(s.ctx, id)
	s.NoError(err)

	s.Equal(expected.ID, actual.ID)
	s.Equal(expected.Name, actual.Name)
}
