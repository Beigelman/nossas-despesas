package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Beigelman/nossas-despesas/internal/modules/auth"
	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/Beigelman/nossas-despesas/internal/pkg/dbtest"
)

type AuthRepositoryTestSuite struct {
	suite.Suite
	repository auth.Repository
	ctx        context.Context
	db         *db.Client
}

func TestAuthRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(AuthRepositoryTestSuite))
}

func (s *AuthRepositoryTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.db = dbtest.Setup(s.ctx, s.T())
	s.repository = NewAuthRepository(s.db)

	_, err := s.db.Conn().Exec(`
   		INSERT INTO users (id, name, email, created_at, updated_at, deleted_at, version)
			VALUES (1, 'john', 'john@email.com', NOW(), NOW(), NOW(), 0)
	`)
	s.NoError(err)
}

func (s *AuthRepositoryTestSuite) TearDownTest() {
	err := s.db.Clean("authentications")
	s.NoError(err)
}

func (s *AuthRepositoryTestSuite) TestPgUserRepo_Store() {
	id := s.repository.GetNextID()
	authentication, err := auth.NewCredentialAuth(auth.CredentialsAttributes{
		ID:       id,
		Email:    "john@email.com",
		Password: "test123",
	})
	s.NoError(err)

	s.NoError(s.repository.Store(s.ctx, authentication))
}

func (s *AuthRepositoryTestSuite) TestPgUserRepo_GetByID() {
	id := s.repository.GetNextID()
	expected, err := auth.NewCredentialAuth(auth.CredentialsAttributes{
		ID:       id,
		Email:    "john@email.com",
		Password: "test123",
	})
	s.NoError(err)
	s.NoError(s.repository.Store(s.ctx, expected))

	actual, err := s.repository.GetByID(s.ctx, id)
	s.NoError(err)

	s.Equal(expected.ID, actual.ID)
	s.Equal(expected.Email, actual.Email)
}

func (s *AuthRepositoryTestSuite) TestPgUserRepo_GetByEmail() {
	id := s.repository.GetNextID()
	expected, err := auth.NewCredentialAuth(auth.CredentialsAttributes{
		ID:       id,
		Email:    "john@email.com",
		Password: "test123",
	})
	s.NoError(err)

	s.NoError(s.repository.Store(s.ctx, expected))

	actual, err := s.repository.GetByEmail(s.ctx, expected.Email, auth.Types.Credentials)
	s.NoError(err)

	s.Equal(expected.ID, actual.ID)
	s.Equal(expected.Email, actual.Email)
	s.Equal(expected.Type, actual.Type)
}
