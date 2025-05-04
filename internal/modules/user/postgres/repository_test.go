package postgres_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/modules/user/postgres"

	"github.com/Beigelman/nossas-despesas/internal/config"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"

	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/Beigelman/nossas-despesas/internal/tests"
	"github.com/stretchr/testify/suite"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	repository    user.Repository
	ctx           context.Context
	db            db.Database
	cfg           config.Config
	testContainer *tests.PostgresContainer
	err           error
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

func (s *UserRepositoryTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.testContainer, s.err = tests.StartPostgres(s.ctx)
	if s.err != nil {
		panic(s.err)
	}

	s.cfg = config.NewTestConfig(s.testContainer.Port, s.testContainer.Host)

	s.db, s.err = db.New(&s.cfg)
	s.NoError(s.err)

	s.repository = postgres.NewUserRepository(s.db)

	s.err = s.db.MigrateUp()
	s.NoError(s.err)
}

func (s *UserRepositoryTestSuite) TearDownSuite() {
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

func (s *UserRepositoryTestSuite) TestPgUserRepo_Store() {
	id := s.repository.GetNextID()
	usr := user.New(user.Attributes{
		ID:    id,
		Name:  "John Doe",
		Email: "john@email.com",
	})

	log.Println("USER", usr.Flags)

	s.NoError(s.repository.Store(s.ctx, usr))

	usr.AddFlag("premium")

	s.NoError(s.repository.Store(s.ctx, usr))
}

func (s *UserRepositoryTestSuite) TestPgUserRepo_GetByID() {
	id := s.repository.GetNextID()
	expected := user.New(user.Attributes{
		ID:    id,
		Name:  "John Doe",
		Email: "john1@email.com",
	})
	expected.AddFlag(user.PREMIUM)
	expected.AddFlag(user.EDIT_PARTNER_INCOME)
	expected.AddFlag(user.Flag("test"))
	expected.RemoveFlag(user.Flag("test"))

	err := s.repository.Store(s.ctx, expected)
	s.NoError(err)

	actual, err := s.repository.GetByID(s.ctx, id)
	s.NoError(err)

	s.Equal(expected.ID, actual.ID)
	s.Equal(expected.Name, actual.Name)
	s.Equal(expected.Email, actual.Email)
	s.Len(actual.Flags, 2)
}

func (s *UserRepositoryTestSuite) TestPgUserRepo_GetByEmail() {
	id := s.repository.GetNextID()
	expected := user.New(user.Attributes{
		ID:    id,
		Name:  "John Doe",
		Email: "john2@email.com",
	})
	expected.AddFlag(user.PREMIUM)

	err := s.repository.Store(s.ctx, expected)
	s.NoError(err)

	actual, err := s.repository.GetByEmail(s.ctx, "john2@email.com")
	s.NoError(err)

	s.Equal(expected.ID, actual.ID)
	s.Equal(expected.Name, actual.Name)
	s.Equal(expected.Email, actual.Email)
	s.Len(actual.Flags, 1)
	s.Equal(user.PREMIUM, actual.Flags[0])
}
