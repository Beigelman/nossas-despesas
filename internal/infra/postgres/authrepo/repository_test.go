package authrepo

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

type PgUserRepoTestSuite struct {
	suite.Suite
	repository    repository.AuthRepository
	ctx           context.Context
	db            db.Database
	cfg           config.Config
	testContainer *tests.PostgresContainer
	err           error
}

func TestPgAuthRepoTestSuite(t *testing.T) {
	suite.Run(t, new(PgUserRepoTestSuite))
}

func (s *PgUserRepoTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.testContainer, s.err = tests.StartPostgres(s.ctx)
	if s.err != nil {
		panic(s.err)
	}

	s.cfg = config.NewTestConfig(s.testContainer.Port, s.testContainer.Host, "postgres")

	s.db = db.New(&s.cfg)
	s.repository = NewPGRepository(s.db)

	s.err = s.db.MigrateUp()
	s.NoError(s.err)

	_, err := s.db.Client().Exec(`
   		INSERT INTO users (id, name, email, created_at, updated_at, deleted_at, version)
		VALUES (1, 'john', 'john@email.com', now(), now(), now(), 0)`,
	)
	s.NoError(err)
}

func (s *PgUserRepoTestSuite) TearDownSuite() {
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

func (s *PgUserRepoTestSuite) TearDownTest() {
	err := s.db.Clean("authentications")
	s.NoError(err)
}

func (s *PgUserRepoTestSuite) TestPgUserRepo_Store() {
	id := s.repository.GetNextID()
	auth, err := entity.NewCredentialAuth(entity.CredentialsAuthParams{
		ID:       id,
		Email:    "john@email.com",
		Password: "test123",
	})
	s.NoError(err)

	s.NoError(s.repository.Store(s.ctx, auth))
}

func (s *PgUserRepoTestSuite) TestPgUserRepo_GetByID() {
	id := s.repository.GetNextID()
	expected, err := entity.NewCredentialAuth(entity.CredentialsAuthParams{
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

func (s *PgUserRepoTestSuite) TestPgUserRepo_GetByEmail() {
	id := s.repository.GetNextID()
	expected, err := entity.NewCredentialAuth(entity.CredentialsAuthParams{
		ID:       id,
		Email:    "john@email.com",
		Password: "test123",
	})
	s.NoError(err)

	s.NoError(s.repository.Store(s.ctx, expected))

	actual, err := s.repository.GetByEmail(s.ctx, expected.Email, entity.AuthTypes.Credentials)
	s.NoError(err)

	s.Equal(expected.ID, actual.ID)
	s.Equal(expected.Email, actual.Email)
	s.Equal(expected.Type, actual.Type)
}
