package userrepo

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

type PgUserRepoTestSuite struct {
	suite.Suite
	repository    repository.UserRepository
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

func (s *PgUserRepoTestSuite) TestPgUserRepo_Store() {
	id := s.repository.GetNextID()
	user := entity.NewUser(entity.UserParams{
		ID:    id,
		Name:  "John Doe",
		Email: "john@email.com",
	})

	s.NoError(s.repository.Store(s.ctx, user))
}

func (s *PgUserRepoTestSuite) TestPgUserRepo_GetByID() {
	id := s.repository.GetNextID()
	expected := entity.NewUser(entity.UserParams{
		ID:    id,
		Name:  "John Doe",
		Email: "john@email.com",
	})

	err := s.repository.Store(s.ctx, expected)
	s.NoError(err)

	actual, err := s.repository.GetByID(s.ctx, id)
	s.NoError(err)

	s.Equal(expected.ID, actual.ID)
	s.Equal(expected.Name, actual.Name)
	s.Equal(expected.Email, actual.Email)
}

func (s *PgUserRepoTestSuite) TestPgUserRepo_GetByEmail() {
	id := s.repository.GetNextID()
	expected := entity.NewUser(entity.UserParams{
		ID:    id,
		Name:  "John Doe",
		Email: "john@email.com",
	})

	err := s.repository.Store(s.ctx, expected)
	s.NoError(err)

	actual, err := s.repository.GetByEmail(s.ctx, "john@email.com")
	s.NoError(err)

	s.Equal(expected.ID, actual.ID)
	s.Equal(expected.Name, actual.Name)
	s.Equal(expected.Email, actual.Email)
}
