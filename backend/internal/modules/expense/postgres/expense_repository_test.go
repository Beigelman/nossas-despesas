package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/Beigelman/nossas-despesas/internal/modules/category"
	postgres2 "github.com/Beigelman/nossas-despesas/internal/modules/category/postgres"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense/postgres"
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	grouprepo "github.com/Beigelman/nossas-despesas/internal/modules/group/postgres"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	userrepo "github.com/Beigelman/nossas-despesas/internal/modules/user/postgres"
	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/Beigelman/nossas-despesas/internal/pkg/dbtest"
)

type ExpenseRepositoryTestSuite struct {
	suite.Suite
	ctx context.Context

	expenseRepo       expense.Repository
	userRepo          user.Repository
	categoryRepo      category.Repository
	categoryGroupRepo category.GroupRepository
	groupRepo         group.Repository

	payer         *user.User
	receiver      *user.User
	category      *category.Category
	categoryGroup *category.Group
	group         *group.Group

	db *db.Client
}

func TestExpenseRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ExpenseRepositoryTestSuite))
}

func (s *ExpenseRepositoryTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.db = dbtest.Setup(s.ctx, s.T())

	s.expenseRepo = postgres.NewExpenseRepository(s.db)
	s.userRepo = userrepo.NewUserRepository(s.db)
	s.categoryRepo = postgres2.NewCategoryRepository(s.db)
	s.categoryGroupRepo = postgres2.NewCategoryGroupRepository(s.db)
	s.groupRepo = grouprepo.NewGroupRepository(s.db)

	s.payer = user.New(user.Attributes{
		ID:    s.userRepo.GetNextID(),
		Name:  "Payer",
		Email: "payer@email.com",
	})

	s.receiver = user.New(user.Attributes{
		ID:    s.userRepo.GetNextID(),
		Name:  "Receiver",
		Email: "receiver@email.com",
	})

	s.categoryGroup = category.NewGroup(category.GroupAttributes{
		ID:   s.categoryGroupRepo.GetNextID(),
		Name: "Category",
		Icon: "test",
	})

	s.category = category.New(category.Attributes{
		ID:              s.categoryRepo.GetNextID(),
		Name:            "Category",
		Icon:            "test",
		CategoryGroupID: s.categoryGroup.ID,
	})

	s.group = group.New(group.Attributes{
		ID:   s.groupRepo.GetNextID(),
		Name: "Group",
	})

	s.NoError(s.userRepo.Store(s.ctx, s.payer))
	s.NoError(s.userRepo.Store(s.ctx, s.receiver))
	s.NoError(s.categoryGroupRepo.Store(s.ctx, s.categoryGroup))
	s.NoError(s.categoryRepo.Store(s.ctx, s.category))
	s.NoError(s.groupRepo.Store(s.ctx, s.group))
}

func (s *ExpenseRepositoryTestSuite) TearDownSubTest() {
	s.NoError(s.db.Clean())
}

func (s *ExpenseRepositoryTestSuite) TestPgExpenseRepo_Store() {
	id := s.expenseRepo.GetNextID()
	expected, err := expense.New(expense.Attributes{
		ID:          id,
		Name:        "my first expense",
		Amount:      100,
		Description: "My Description",
		PayerID:     s.payer.ID,
		ReceiverID:  s.receiver.ID,
		SplitRatio: expense.SplitRatio{
			Payer:    50,
			Receiver: 50,
		},
		SplitType:  expense.SplitTypes.Equal,
		CategoryID: s.category.ID,
		GroupID:    s.group.ID,
	})
	s.NoError(err)

	s.NoError(s.expenseRepo.Store(s.ctx, expected))
}

func (s *ExpenseRepositoryTestSuite) TestPgExpenseRepo_GetByID() {
	id := s.expenseRepo.GetNextID()
	expected, err := expense.New(expense.Attributes{
		ID:          id,
		Name:        "my first expense",
		Amount:      100,
		Description: "My Description",
		PayerID:     s.payer.ID,
		ReceiverID:  s.receiver.ID,
		SplitRatio: expense.SplitRatio{
			Payer:    50,
			Receiver: 50,
		},
		SplitType:  expense.SplitTypes.Equal,
		CategoryID: s.category.ID,
		GroupID:    s.group.ID,
	})
	s.NoError(err)

	s.NoError(s.expenseRepo.Store(s.ctx, expected))

	actual, err := s.expenseRepo.GetByID(s.ctx, id)
	s.NoError(err)

	s.Equal(expected.ID, actual.ID)
	s.Equal(expected.Name, actual.Name)
	s.Equal(0, actual.Version)
}

func (s *ExpenseRepositoryTestSuite) TestPgExpenseRepo_GetByGroupDate() {
	var entities []expense.Expense
	for i := 0; i < 3; i++ {
		entity, err := expense.New(expense.Attributes{
			ID:          s.expenseRepo.GetNextID(),
			Name:        "my first expense",
			Amount:      100,
			Description: "My Description",
			PayerID:     s.payer.ID,
			ReceiverID:  s.receiver.ID,
			SplitRatio: expense.SplitRatio{
				Payer:    50,
				Receiver: 50,
			},
			SplitType:  expense.SplitTypes.Equal,
			CategoryID: s.category.ID,
			GroupID:    s.group.ID,
		})
		s.NoError(err)
		entities = append(entities, *entity)
	}
	s.NoError(s.expenseRepo.BulkStore(s.ctx, entities))

	expected, err := s.expenseRepo.GetByGroupDate(s.ctx, s.group.ID, time.Now())
	s.NoError(err)
	s.Len(expected, 3)
}
