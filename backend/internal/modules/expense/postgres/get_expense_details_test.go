package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/Beigelman/nossas-despesas/internal/pkg/dbtest"
	"github.com/Beigelman/nossas-despesas/internal/shared/fixture"
)

type GetExpenseDetailsTestSuite struct {
	suite.Suite
	db                *db.Client
	ctx               context.Context
	getExpenseDetails GetExpenseDetails
}

func TestGetExpenseDetailsTestSuite(t *testing.T) {
	suite.Run(t, new(GetExpenseDetailsTestSuite))
}

func (s *GetExpenseDetailsTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.db = dbtest.Setup(s.ctx, s.T())

	err := fixture.ExecuteSQLFiles(s.db, []string{
		"./fixtures/basic_setup.sql",
		"./fixtures/get_expense_details.sql",
	})
	s.NoError(err)

	s.getExpenseDetails = NewGetExpenseDetails(s.db)
}

func (s *GetExpenseDetailsTestSuite) TearDownSubTest() {
	s.NoError(s.db.Clean())
}

func (s *GetExpenseDetailsTestSuite) TestGetExpenseDetails_Success() {
	result, err := s.getExpenseDetails(s.ctx, 1)
	s.NoError(err)
	s.Len(result, 1)

	expense := result[0]
	s.Equal(1, expense.ID)
	s.Equal("Almoço Básico", expense.Name)
	s.Equal(float32(2500), expense.Amount)
	s.Nil(expense.RefundAmount)
	s.Equal("Almoço no restaurante", expense.Description)
	s.Equal(100, expense.CategoryID)
	s.Equal(100, expense.PayerID)
	s.Equal(101, expense.ReceiverID)
	s.Equal(100, expense.GroupID)
	s.Equal("equal", expense.SplitType)
	s.Equal(50, expense.SplitRatio.Payer)
	s.Equal(50, expense.SplitRatio.Receiver)
	s.Nil(expense.DeletedAt)
}

func (s *GetExpenseDetailsTestSuite) TestGetExpenseDetails_WithRefundAmount() {
	result, err := s.getExpenseDetails(s.ctx, 3)
	s.NoError(err)
	s.Len(result, 1)

	expense := result[0]
	s.Equal(3, expense.ID)
	s.Equal("Compra com Reembolso", expense.Name)
	s.Equal(float32(10000), expense.Amount)
	s.NotNil(expense.RefundAmount)
	s.Equal(float32(2000), *expense.RefundAmount)
	s.Equal("Compra que teve reembolso parcial", expense.Description)
}

func (s *GetExpenseDetailsTestSuite) TestGetExpenseDetails_NotFound() {
	result, err := s.getExpenseDetails(s.ctx, 999)
	s.NoError(err)
	s.Len(result, 0)
}

func (s *GetExpenseDetailsTestSuite) TestGetExpenseDetails_AllFields() {
	result, err := s.getExpenseDetails(s.ctx, 2)
	s.NoError(err)
	s.Len(result, 1)

	expense := result[0]
	s.Equal(2, expense.ID)
	s.Equal("Jantar Completo", expense.Name)
	s.Equal(float32(7500), expense.Amount)
	s.Equal("Jantar com todos os campos preenchidos", expense.Description)
	s.Equal(100, expense.CategoryID)
	s.Equal(100, expense.PayerID)
	s.Equal(101, expense.ReceiverID)
	s.Equal(100, expense.GroupID)
	s.Equal("proportional", expense.SplitType)
	s.Equal(60, expense.SplitRatio.Payer)
	s.Equal(40, expense.SplitRatio.Receiver)
	s.NotNil(expense.CreatedAt)
	s.NotNil(expense.UpdatedAt)
	s.Nil(expense.DeletedAt)
}
