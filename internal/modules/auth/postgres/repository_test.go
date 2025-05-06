package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/pkg/config"

	"github.com/Beigelman/nossas-despesas/internal/modules/auth"
	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/Beigelman/nossas-despesas/internal/tests"
	"github.com/stretchr/testify/suite"
)

type AuthRepositoryTestSuite struct {
	suite.Suite
	repository    auth.Repository
	ctx           context.Context
	db            db.Database
	cfg           config.Config
	testContainer *tests.PostgresContainer
	err           error
}

func TestAuthRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(AuthRepositoryTestSuite))
}

func (s *AuthRepositoryTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.testContainer, s.err = tests.StartPostgres(s.ctx)
	if s.err != nil {
		panic(s.err)
	}

	s.cfg = config.NewTestConfig(s.testContainer.Port, s.testContainer.Host)

	s.db, s.err = db.New(&s.cfg)
	s.NoError(s.err)
	s.repository = NewAuthRepository(s.db)

	s.err = s.db.MigrateUp()
	s.NoError(s.err)

	_, err := s.db.Client().Exec(`
   		INSERT INTO users (id, name, email, created_at, updated_at, deleted_at, version)
		VALUES (1, 'john', 'john@email.com', now(), now(), now(), 0)`,
	)
	s.NoError(err)
}

func (s *AuthRepositoryTestSuite) TearDownSuite() {
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
