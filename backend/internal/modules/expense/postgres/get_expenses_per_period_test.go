package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/Beigelman/nossas-despesas/internal/pkg/dbtest"
	"github.com/Beigelman/nossas-despesas/internal/shared/fixture"
)

type GetExpensesPerPeriodTestSuite struct {
	suite.Suite
	db                   *db.Client
	ctx                  context.Context
	getExpensesPerPeriod GetExpensesPerPeriod
}

func TestGetExpensesPerPeriodTestSuite(t *testing.T) {
	suite.Run(t, new(GetExpensesPerPeriodTestSuite))
}

func (s *GetExpensesPerPeriodTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.db = dbtest.Setup(s.ctx, s.T())

	err := fixture.ExecuteSQLFiles(s.db, []string{
		"./fixtures/basic_setup.sql",
		"./fixtures/get_expenses_per_period.sql",
	})
	s.NoError(err)

	s.getExpensesPerPeriod = NewGetExpensesPerPeriod(s.db)
}

func (s *GetExpensesPerPeriodTestSuite) TearDownSubTest() {
	s.NoError(s.db.Clean())
}

func (s *GetExpensesPerPeriodTestSuite) TestGetExpensesPerPeriod_DailyAggregation() {
	input := GetExpensesPerPeriodInput{
		GroupID:   100,
		Aggregate: "day",
		StartDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 6, 30, 23, 59, 59, 0, time.UTC),
	}

	result, err := s.getExpensesPerPeriod(s.ctx, input)
	s.NoError(err)
	s.Len(result, 4) // 4 different days in June

	// Create a map for easier verification
	dailyData := make(map[string]ExpensesPerPeriod)
	for _, item := range result {
		dailyData[item.Date] = item
	}

	// Verify June 1st (should have 2 expenses totaling 2500)
	june1 := dailyData["2024-06-01"]
	s.Equal("2024-06-01", june1.Date)
	s.Equal(2500, june1.Amount) // 1000 + 1500
	s.Equal(2, june1.Count)

	// Verify June 2nd
	june2 := dailyData["2024-06-02"]
	s.Equal("2024-06-02", june2.Date)
	s.Equal(2000, june2.Amount)
	s.Equal(1, june2.Count)

	// Verify June 15th
	june15 := dailyData["2024-06-15"]
	s.Equal("2024-06-15", june15.Date)
	s.Equal(3000, june15.Amount)
	s.Equal(1, june15.Count)

	// Verify June 30th
	june30 := dailyData["2024-06-30"]
	s.Equal("2024-06-30", june30.Date)
	s.Equal(2500, june30.Amount)
	s.Equal(1, june30.Count)
}

func (s *GetExpensesPerPeriodTestSuite) TestGetExpensesPerPeriod_MonthlyAggregation() {
	input := GetExpensesPerPeriodInput{
		GroupID:   100,
		Aggregate: "month",
		StartDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 8, 31, 23, 59, 59, 0, time.UTC),
	}

	result, err := s.getExpensesPerPeriod(s.ctx, input)
	s.NoError(err)
	s.Len(result, 3) // June, July, August

	// Create a map for easier verification
	monthlyData := make(map[string]ExpensesPerPeriod)
	for _, item := range result {
		monthlyData[item.Date] = item
	}

	// Verify June 2024
	june := monthlyData["2024-06"]
	s.Equal("2024-06", june.Date)
	s.Equal(10000, june.Amount) // 1000 + 1500 + 2000 + 3000 + 2500
	s.Equal(5, june.Count)

	// Verify July 2024
	july := monthlyData["2024-07"]
	s.Equal("2024-07", july.Date)
	s.Equal(12500, july.Amount) // 4000 + 3500 + 5000
	s.Equal(3, july.Count)

	// Verify August 2024
	august := monthlyData["2024-08"]
	s.Equal("2024-08", august.Date)
	s.Equal(6000, august.Amount)
	s.Equal(1, august.Count)
}

func (s *GetExpensesPerPeriodTestSuite) TestGetExpensesPerPeriod_EmptyResult() {
	input := GetExpensesPerPeriodInput{
		GroupID:   999, // Non-existent group
		Aggregate: "day",
		StartDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 6, 30, 23, 59, 59, 0, time.UTC),
	}

	result, err := s.getExpensesPerPeriod(s.ctx, input)
	s.NoError(err)
	s.Empty(result)
}

func (s *GetExpensesPerPeriodTestSuite) TestGetExpensesPerPeriod_NoExpensesInDateRange() {
	input := GetExpensesPerPeriodInput{
		GroupID:   100,
		Aggregate: "day",
		StartDate: time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 9, 30, 23, 59, 59, 0, time.UTC),
	}

	result, err := s.getExpensesPerPeriod(s.ctx, input)
	s.NoError(err)
	s.Empty(result)
}

func (s *GetExpensesPerPeriodTestSuite) TestGetExpensesPerPeriod_DateRange() {
	// Test with a specific date range that includes only July
	input := GetExpensesPerPeriodInput{
		GroupID:   100,
		Aggregate: "day",
		StartDate: time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 7, 31, 23, 59, 59, 0, time.UTC),
	}

	result, err := s.getExpensesPerPeriod(s.ctx, input)
	s.NoError(err)
	s.Len(result, 3) // 3 different days in July

	// Verify all results are from July 2024
	for _, item := range result {
		s.Contains(item.Date, "2024-07")
	}

	// Calculate total amount for July
	totalAmount := 0
	totalCount := 0
	for _, item := range result {
		totalAmount += item.Amount
		totalCount += item.Count
	}
	s.Equal(12500, totalAmount) // 4000 + 3500 + 5000
	s.Equal(3, totalCount)
}

func (s *GetExpensesPerPeriodTestSuite) TestGetExpensesPerPeriod_OrderedResults() {
	input := GetExpensesPerPeriodInput{
		GroupID:   100,
		Aggregate: "day",
		StartDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 7, 31, 23, 59, 59, 0, time.UTC),
	}

	result, err := s.getExpensesPerPeriod(s.ctx, input)
	s.NoError(err)
	s.NotEmpty(result)

	// Verify results are ordered by date
	for i := 1; i < len(result); i++ {
		s.True(result[i-1].Date <= result[i].Date, "Results should be ordered by date")
	}
}
