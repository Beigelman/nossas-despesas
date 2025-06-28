package postgres_test

import (
	"context"
	"testing"
	"time"

	"cloud.google.com/go/civil"

	"github.com/Beigelman/nossas-despesas/internal/modules/category"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense/postgres"
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/Beigelman/nossas-despesas/internal/pkg/dbtest"
	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
	"github.com/stretchr/testify/suite"
)

type ScheduledExpenseRepositoryTestSuite struct {
	suite.Suite
	ctx context.Context

	scheduledExpenseRepo expense.ScheduledExpenseRepository

	db *db.Client
}

func TestScheduledExpenseRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ScheduledExpenseRepositoryTestSuite))
}

func (s *ScheduledExpenseRepositoryTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.db = dbtest.Setup(s.ctx, s.T())
	s.scheduledExpenseRepo = postgres.NewScheduledExpenseRepository(s.db)
}

func (s *ScheduledExpenseRepositoryTestSuite) TearDownSubTest() {
	s.NoError(s.db.Clean("scheduled_expenses"))
}

func (s *ScheduledExpenseRepositoryTestSuite) TestPgScheduledExpenseRepo_Store() {
	scheduledExpense, err := expense.NewScheduledExpense(expense.ScheduledExpenseAttributes{
		ID:              s.scheduledExpenseRepo.GetNextID(),
		Name:            "Test Scheduled Expense",
		Amount:          1000,
		Description:     "Test Description",
		GroupID:         group.ID{Value: 1},
		CategoryID:      category.ID{Value: 1},
		SplitType:       expense.SplitTypes.Equal,
		PayerID:         user.ID{Value: 1},
		ReceiverID:      user.ID{Value: 2},
		FrequencyInDays: 7,
		LastGeneratedAt: &[]civil.Date{civil.DateOf(time.Now())}[0],
	})
	s.NoError(err)

	s.NoError(s.scheduledExpenseRepo.Store(s.ctx, scheduledExpense))
}

func (s *ScheduledExpenseRepositoryTestSuite) TestPgScheduledExpenseRepo_BulkStore() {
	scheduledExpenses := []expense.ScheduledExpense{
		{
			Entity: ddd.Entity[expense.ScheduledExpenseID]{
				ID:        s.scheduledExpenseRepo.GetNextID(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Version:   1,
			},
			Name:            "Test Scheduled Expense",
			Amount:          1000,
			Description:     "Test Description",
			GroupID:         group.ID{Value: 1},
			CategoryID:      category.ID{Value: 1},
			SplitType:       expense.SplitTypes.Equal,
			PayerID:         user.ID{Value: 1},
			ReceiverID:      user.ID{Value: 2},
			FrequencyInDays: 7,
			LastGeneratedAt: &[]civil.Date{civil.DateOf(time.Now())}[0],
		},
		{
			Entity: ddd.Entity[expense.ScheduledExpenseID]{
				ID:        s.scheduledExpenseRepo.GetNextID(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Version:   1,
			},
			Name:            "Test Scheduled Expense 2",
			Amount:          2000,
			Description:     "Test Description 2",
			GroupID:         group.ID{Value: 1},
			CategoryID:      category.ID{Value: 1},
			SplitType:       expense.SplitTypes.Equal,
			PayerID:         user.ID{Value: 1},
			ReceiverID:      user.ID{Value: 2},
			FrequencyInDays: 7,
			LastGeneratedAt: &[]civil.Date{civil.DateOf(time.Now())}[0],
		},
	}

	s.NoError(s.scheduledExpenseRepo.BulkStore(s.ctx, scheduledExpenses))
}

func (s *ScheduledExpenseRepositoryTestSuite) TestPgScheduledExpenseRepo_GetByID() {
	scheduledExpense, err := expense.NewScheduledExpense(expense.ScheduledExpenseAttributes{
		ID:              s.scheduledExpenseRepo.GetNextID(),
		Name:            "Test Scheduled Expense",
		Amount:          1000,
		Description:     "Test Description",
		GroupID:         group.ID{Value: 1},
		CategoryID:      category.ID{Value: 1},
		SplitType:       expense.SplitTypes.Equal,
		PayerID:         user.ID{Value: 1},
		ReceiverID:      user.ID{Value: 2},
		FrequencyInDays: 7,
	})
	s.NoError(err)

	s.NoError(s.scheduledExpenseRepo.Store(s.ctx, scheduledExpense))

	retrieved, err := s.scheduledExpenseRepo.GetByID(s.ctx, scheduledExpense.ID)
	s.NoError(err)
	s.NotNil(retrieved)
	s.Equal(scheduledExpense.ID, retrieved.ID)
	s.Equal(scheduledExpense.Name, retrieved.Name)
	s.Equal(scheduledExpense.Amount, retrieved.Amount)
	s.Equal(scheduledExpense.Description, retrieved.Description)
	s.Equal(scheduledExpense.GroupID, retrieved.GroupID)
	s.Equal(scheduledExpense.CategoryID, retrieved.CategoryID)
	s.Equal(scheduledExpense.SplitType, retrieved.SplitType)
	s.Equal(scheduledExpense.PayerID, retrieved.PayerID)
	s.Equal(scheduledExpense.ReceiverID, retrieved.ReceiverID)
	s.Equal(scheduledExpense.FrequencyInDays, retrieved.FrequencyInDays)
}

func (s *ScheduledExpenseRepositoryTestSuite) TestPgScheduledExpenseRepo_GetActiveScheduledExpenses() {
	scheduledExpense, err := expense.NewScheduledExpense(expense.ScheduledExpenseAttributes{
		ID:              s.scheduledExpenseRepo.GetNextID(),
		Name:            "Test Scheduled Expense",
		Amount:          1000,
		Description:     "Test Description",
		GroupID:         group.ID{Value: 1},
		CategoryID:      category.ID{Value: 1},
		SplitType:       expense.SplitTypes.Equal,
		PayerID:         user.ID{Value: 1},
		ReceiverID:      user.ID{Value: 2},
		FrequencyInDays: 7,
	})
	s.NoError(err)

	s.NoError(s.scheduledExpenseRepo.Store(s.ctx, scheduledExpense))

	activeExpenses, err := s.scheduledExpenseRepo.GetActiveScheduledExpenses(s.ctx)
	s.NoError(err)
	s.NotNil(activeExpenses)
	s.Len(activeExpenses, 1)
	s.Equal(scheduledExpense.ID, activeExpenses[0].ID)
	s.Equal(scheduledExpense.Name, activeExpenses[0].Name)
	s.Equal(scheduledExpense.Amount, activeExpenses[0].Amount)
	s.Equal(scheduledExpense.Description, activeExpenses[0].Description)
	s.Equal(scheduledExpense.GroupID, activeExpenses[0].GroupID)
	s.Equal(scheduledExpense.CategoryID, activeExpenses[0].CategoryID)
	s.Equal(scheduledExpense.SplitType, activeExpenses[0].SplitType)
	s.Equal(scheduledExpense.PayerID, activeExpenses[0].PayerID)
	s.Equal(scheduledExpense.ReceiverID, activeExpenses[0].ReceiverID)
	s.Equal(scheduledExpense.FrequencyInDays, activeExpenses[0].FrequencyInDays)
}
