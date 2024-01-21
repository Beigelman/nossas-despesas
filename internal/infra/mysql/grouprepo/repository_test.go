package grouprepo

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

var migrationPath = "file:///Users/danielbeigelman/mydev/go-luda/server/database/migrations"

type PgGroupRepoTestSuite struct {
	suite.Suite
	repository    repository.GroupRepository
	ctx           context.Context
	db            db.Database
	cfg           config.Config
	testContainer *tests.MySqlContainer
	err           error
}

func TestPgGroupRepoTestSuite(t *testing.T) {
	suite.Run(t, new(PgGroupRepoTestSuite))
}

func (s *PgGroupRepoTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.testContainer, s.err = tests.StartMySql(s.ctx)
	if s.err != nil {
		panic(s.err)
	}

	s.cfg = config.NewTestConfig(s.testContainer.Port, s.testContainer.Host)

	s.db = db.New(&s.cfg)
	s.repository = NewPGRepository(s.db)

	s.err = s.db.MigrateUp(migrationPath)
	s.NoError(s.err)
}

func (s *PgGroupRepoTestSuite) TearDownSuite() {
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

func (s *PgGroupRepoTestSuite) TearDownSubTest() {
	err := s.db.Clean()
	s.NoError(err)
}

func (s *PgGroupRepoTestSuite) TestPgGroupRepo_Store() {
	id := s.repository.GetNextID()
	group := entity.NewGroup(entity.GroupParams{
		ID:   id,
		Name: "My Group",
	})

	err := s.repository.Store(s.ctx, group)
	s.NoError(err)
}

func (s *PgGroupRepoTestSuite) TestPgGroupRepo_GetByID() {
	id := s.repository.GetNextID()
	expected := entity.NewGroup(entity.GroupParams{
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

func (s *PgGroupRepoTestSuite) TestPgGroupRepo_GetByName() {
	id := s.repository.GetNextID()
	expected := entity.NewGroup(entity.GroupParams{
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
