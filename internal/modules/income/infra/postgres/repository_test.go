package postgres

import (
	"context"
	"github.com/Beigelman/nossas-despesas/internal/config"
	"github.com/Beigelman/nossas-despesas/internal/modules/income"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"testing"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/Beigelman/nossas-despesas/internal/tests"
	"github.com/stretchr/testify/suite"
)

var userID = user.ID{Value: 1}

type IncomeRepositoryTestSuite struct {
	suite.Suite
	repository    income.Repository
	ctx           context.Context
	db            db.Database
	cfg           config.Config
	testContainer *tests.PostgresContainer
	err           error
}

func TestIncomeRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(IncomeRepositoryTestSuite))
}

func (s *IncomeRepositoryTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.testContainer, s.err = tests.StartPostgres(s.ctx)
	if s.err != nil {
		panic(s.err)
	}

	s.cfg = config.NewTestConfig(s.testContainer.Port, s.testContainer.Host)

	s.db = db.New(&s.cfg)
	s.repository = NewIncomeRepository(s.db)

	s.err = s.db.MigrateUp()
	s.NoError(s.err)

	_, err := s.db.Client().Exec(`
   		INSERT INTO users (id, name, email, created_at, updated_at, deleted_at, version)
		VALUES (1, 'john', 'john@email.com', now(), now(), now(), 0)`,
	)
	s.NoError(err)
}

func (s *IncomeRepositoryTestSuite) TearDownSuite() {
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

func (s *IncomeRepositoryTestSuite) TearDownTest() {
	err := s.db.Clean("incomes")
	s.NoError(err)
}

func (s *IncomeRepositoryTestSuite) TestPgUserRepo_Store() {
	id := s.repository.GetNextID()
	inc := income.New(income.Attributes{
		ID:     id,
		Amount: 100,
		UserID: userID,
		Type:   income.Types.Salary,
	})

	s.NoError(s.repository.Store(s.ctx, inc))
}

func (s *IncomeRepositoryTestSuite) TestPgUserRepo_GetByID() {
	id := s.repository.GetNextID()
	expected := income.New(income.Attributes{
		ID:     id,
		Amount: 100,
		UserID: userID,
		Type:   income.Types.Salary,
	})

	s.NoError(s.repository.Store(s.ctx, expected))

	actual, err := s.repository.GetByID(s.ctx, id)
	s.NoError(err)

	s.Equal(expected.ID, actual.ID)
	s.Equal(expected.Type, actual.Type)
}

func (s *IncomeRepositoryTestSuite) TestPgUserRepo_GetUserMonthlyIncomes() {
	thisMonth := time.Now()
	lasMonth := time.Now().AddDate(0, -1, 0)
	salary := income.New(income.Attributes{
		ID:        s.repository.GetNextID(),
		Amount:    100,
		UserID:    userID,
		Type:      income.Types.Salary,
		CreatedAt: &thisMonth,
	})

	benefit := income.New(income.Attributes{
		ID:        s.repository.GetNextID(),
		Amount:    200,
		UserID:    userID,
		Type:      income.Types.Benefit,
		CreatedAt: &thisMonth,
	})

	other := income.New(income.Attributes{
		ID:        s.repository.GetNextID(),
		Amount:    200,
		UserID:    userID,
		Type:      income.Types.Other,
		CreatedAt: &lasMonth,
	})

	s.NoError(s.repository.Store(s.ctx, salary))
	s.NoError(s.repository.Store(s.ctx, benefit))
	s.NoError(s.repository.Store(s.ctx, other))

	incomes, err := s.repository.GetUserMonthlyIncomes(s.ctx, userID, &thisMonth)
	s.NoError(err)
	s.Equal(2, len(incomes))
}
