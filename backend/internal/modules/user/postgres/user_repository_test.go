package postgres_test

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"github.com/Beigelman/nossas-despesas/internal/modules/user/postgres"
	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/Beigelman/nossas-despesas/internal/pkg/dbtest"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	repository user.Repository
	ctx        context.Context
	db         *db.Client
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

func (s *UserRepositoryTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.db = dbtest.Setup(s.ctx, s.T())
	s.repository = postgres.NewUserRepository(s.db)
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
