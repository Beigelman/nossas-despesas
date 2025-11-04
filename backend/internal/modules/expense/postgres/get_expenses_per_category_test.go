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

type GetExpensesPerCategoryTestSuite struct {
	suite.Suite
	db                     *db.Client
	ctx                    context.Context
	getExpensesPerCategory GetExpensesPerCategory
}

func TestGetExpensesPerCategoryTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(GetExpensesPerCategoryTestSuite))
}

func (s *GetExpensesPerCategoryTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.db = dbtest.Setup(s.ctx, s.T())

	err := fixture.ExecuteSQLFiles(s.db, []string{
		"./fixtures/basic_setup.sql",
		"./fixtures/get_expenses_per_category.sql",
	})
	s.NoError(err)

	s.getExpensesPerCategory = NewGetExpensesPerCategory(s.db)
}

func (s *GetExpensesPerCategoryTestSuite) TearDownSubTest() {
	s.NoError(s.db.Clean())
}

func (s *GetExpensesPerCategoryTestSuite) TestGetExpensesPerCategory_Success() {
	input := GetExpensesPerCategoryInput{
		GroupID:   100,
		StartDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 6, 30, 23, 59, 59, 0, time.UTC),
	}

	result, err := s.getExpensesPerCategory(s.ctx, input)
	s.NoError(err)
	s.Len(result, 3) // Alimentação, Transporte, Lazer

	// Find each category group
	var alimentacao, transporte, lazer *ExpensesPerCategory
	for i := range result {
		switch result[i].CategoryGroup {
		case "Alimentação":
			alimentacao = &result[i]
		case "Transporte":
			transporte = &result[i]
		case "Lazer":
			lazer = &result[i]
		}
	}

	// Verify Alimentação category group
	s.NotNil(alimentacao)
	s.Equal("Alimentação", alimentacao.CategoryGroup)
	s.Equal(25500, alimentacao.Amount) // 2500 + 3000 + 8000 + 12000
	s.Len(alimentacao.Categories, 2)   // Restaurante, Supermercado

	// Verify Transporte category group
	s.NotNil(transporte)
	s.Equal("Transporte", transporte.CategoryGroup)
	s.Equal(3500, transporte.Amount) // 1500 + 2000
	s.Len(transporte.Categories, 1)  // Uber

	// Verify Lazer category group
	s.NotNil(lazer)
	s.Equal("Lazer", lazer.CategoryGroup)
	s.Equal(3000, lazer.Amount) // 3000
	s.Len(lazer.Categories, 1)  // Cinema
}

func (s *GetExpensesPerCategoryTestSuite) TestGetExpensesPerCategory_DateFilter() {
	// Test with a narrower date range that excludes some expenses
	input := GetExpensesPerCategoryInput{
		GroupID:   100,
		StartDate: time.Date(2024, 6, 3, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 6, 5, 23, 59, 59, 0, time.UTC),
	}

	result, err := s.getExpensesPerCategory(s.ctx, input)
	s.NoError(err)
	s.NotEmpty(result)

	// Should only include expenses from 2024-06-03 to 2024-06-05
	// This should include: Supermercado 1 (8000), Supermercado 2 (12000), Uber 1 (1500)
	totalAmount := 0
	for _, category := range result {
		totalAmount += category.Amount
	}
	s.Equal(21500, totalAmount) // 8000 + 12000 + 1500
}

func (s *GetExpensesPerCategoryTestSuite) TestGetExpensesPerCategory_EmptyResult() {
	input := GetExpensesPerCategoryInput{
		GroupID:   999, // Non-existent group
		StartDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 6, 30, 23, 59, 59, 0, time.UTC),
	}

	result, err := s.getExpensesPerCategory(s.ctx, input)
	s.NoError(err)
	s.Empty(result)
}

func (s *GetExpensesPerCategoryTestSuite) TestGetExpensesPerCategory_NoExpensesInDateRange() {
	input := GetExpensesPerCategoryInput{
		GroupID:   100,
		StartDate: time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 7, 31, 23, 59, 59, 0, time.UTC),
	}

	result, err := s.getExpensesPerCategory(s.ctx, input)
	s.NoError(err)
	s.Empty(result)
}

func (s *GetExpensesPerCategoryTestSuite) TestGetExpensesPerCategory_CategoryDetails() {
	input := GetExpensesPerCategoryInput{
		GroupID:   100,
		StartDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 6, 30, 23, 59, 59, 0, time.UTC),
	}

	result, err := s.getExpensesPerCategory(s.ctx, input)
	s.NoError(err)
	s.NotEmpty(result)

	// Find Alimentação category group and verify its categories
	var alimentacao *ExpensesPerCategory
	for i := range result {
		if result[i].CategoryGroup == "Alimentação" {
			alimentacao = &result[i]
			break
		}
	}

	s.NotNil(alimentacao)
	s.Len(alimentacao.Categories, 2)

	// Verify individual categories within Alimentação
	categoryAmounts := make(map[string]int)
	for _, cat := range alimentacao.Categories {
		categoryAmounts[cat.Category] = cat.Amount
	}

	s.Equal(5500, categoryAmounts["Restaurante"])   // 2500 + 3000
	s.Equal(20000, categoryAmounts["Supermercado"]) // 8000 + 12000
}
