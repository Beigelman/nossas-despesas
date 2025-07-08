package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/Beigelman/nossas-despesas/internal/pkg/dbtest"
	"github.com/Beigelman/nossas-despesas/internal/shared/fixture"
	"github.com/stretchr/testify/suite"
)

type GetExpensesTestSuite struct {
	suite.Suite
	db          *db.Client
	ctx         context.Context
	getExpenses GetExpenses
}

func TestGetExpensesTestSuite(t *testing.T) {
	suite.Run(t, new(GetExpensesTestSuite))
}

func (s *GetExpensesTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.db = dbtest.Setup(s.ctx, s.T())

	err := fixture.ExecuteSQLFiles(s.db, []string{
		"./fixtures/basic_setup.sql",
		"./fixtures/get_expenses.sql",
	})
	s.NoError(err)

	s.getExpenses = NewGetExpenses(s.db)
}

func (s *GetExpensesTestSuite) TearDownSubTest() {
	s.NoError(s.db.Clean())
}

func (s *GetExpensesTestSuite) TestGetExpenses_WithoutSearch() {
	input := GetExpensesInput{
		GroupID:         100,
		LastExpenseDate: time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
		LastExpenseID:   999999,
		Limit:           10,
		Search:          "",
	}

	result, err := s.getExpenses(s.ctx, input)
	s.NoError(err)
	s.NotEmpty(result)

	// Verify all results belong to group 100
	for _, expense := range result {
		s.Equal(100, expense.GroupID)
	}
}

func (s *GetExpensesTestSuite) TestGetExpenses_WithSearch() {
	input := GetExpensesInput{
		GroupID:         100,
		LastExpenseDate: time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
		LastExpenseID:   999999,
		Limit:           10,
		Search:          "McDonald",
	}

	result, err := s.getExpenses(s.ctx, input)
	s.NoError(err)
	s.Len(result, 1)

	expense := result[0]
	s.Equal("Almoço McDonald", expense.Name)
	s.Contains(expense.Description, "McDonald")
}

func (s *GetExpensesTestSuite) TestGetExpenses_EmptyResult() {
	input := GetExpensesInput{
		GroupID:         999, // Non-existent group
		LastExpenseDate: time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
		LastExpenseID:   999999,
		Limit:           10,
		Search:          "",
	}

	result, err := s.getExpenses(s.ctx, input)
	s.NoError(err)
	s.Empty(result)
}
