package incomerepo

import (
	"context"
	"testing"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	"github.com/Beigelman/nossas-despesas/internal/domain/repository"

	"github.com/Beigelman/nossas-despesas/internal/config"
	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/Beigelman/nossas-despesas/internal/tests"
	"github.com/stretchr/testify/suite"
)

var userID = entity.UserID{Value: 1}

type PgIncomeRepoTestSuite struct {
	suite.Suite
	repository    repository.IncomeRepository
	ctx           context.Context
	db            db.Database
	cfg           config.Config
	testContainer *tests.PostgresContainer
	err           error
}

func TestPgAuthRepoTestSuite(t *testing.T) {
	suite.Run(t, new(PgIncomeRepoTestSuite))
}

func (s *PgIncomeRepoTestSuite) SetupSuite() {
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

func (s *PgIncomeRepoTestSuite) TearDownSuite() {
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

func (s *PgIncomeRepoTestSuite) TearDownTest() {
	err := s.db.Clean("incomes")
	s.NoError(err)
}

func (s *PgIncomeRepoTestSuite) TestPgUserRepo_Store() {
	id := s.repository.GetNextID()
	income := entity.NewIncome(entity.IncomeParams{
		ID:     id,
		Amount: 100,
		UserID: userID,
		Type:   entity.IncomeTypes.Salary,
	})

	s.NoError(s.repository.Store(s.ctx, income))
}

func (s *PgIncomeRepoTestSuite) TestPgUserRepo_GetByID() {
	id := s.repository.GetNextID()
	expected := entity.NewIncome(entity.IncomeParams{
		ID:     id,
		Amount: 100,
		UserID: userID,
		Type:   entity.IncomeTypes.Salary,
	})

	s.NoError(s.repository.Store(s.ctx, expected))

	actual, err := s.repository.GetByID(s.ctx, id)
	s.NoError(err)

	s.Equal(expected.ID, actual.ID)
	s.Equal(expected.Type, actual.Type)
}

func (s *PgIncomeRepoTestSuite) TestPgUserRepo_GetUserMonthlyIncomes() {
	thisMonth := time.Now()
	lasMonth := time.Now().AddDate(0, -1, 0)
	salary := entity.NewIncome(entity.IncomeParams{
		ID:        s.repository.GetNextID(),
		Amount:    100,
		UserID:    userID,
		Type:      entity.IncomeTypes.Salary,
		CreatedAt: &thisMonth,
	})

	benefit := entity.NewIncome(entity.IncomeParams{
		ID:        s.repository.GetNextID(),
		Amount:    200,
		UserID:    userID,
		Type:      entity.IncomeTypes.Benefit,
		CreatedAt: &thisMonth,
	})

	other := entity.NewIncome(entity.IncomeParams{
		ID:        s.repository.GetNextID(),
		Amount:    200,
		UserID:    userID,
		Type:      entity.IncomeTypes.Benefit,
		CreatedAt: &lasMonth,
	})

	s.NoError(s.repository.Store(s.ctx, salary))
	s.NoError(s.repository.Store(s.ctx, benefit))
	s.NoError(s.repository.Store(s.ctx, other))

	incomes, err := s.repository.GetUserMonthlyIncomes(s.ctx, userID, &thisMonth)
	s.NoError(err)
	s.Equal(2, len(incomes))
}
