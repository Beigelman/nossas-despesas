package postgres_test

import (
	"context"
	postgres2 "github.com/Beigelman/nossas-despesas/internal/modules/category/postgres"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense/postgres"
	grouprepo "github.com/Beigelman/nossas-despesas/internal/modules/group/postgres"
	userrepo "github.com/Beigelman/nossas-despesas/internal/modules/user/postgres"
	"testing"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/config"
	"github.com/Beigelman/nossas-despesas/internal/modules/category"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense"
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/Beigelman/nossas-despesas/internal/tests"
	"github.com/stretchr/testify/suite"
)

type ExpenseRepositoryTestSuite struct {
	suite.Suite
	ctx           context.Context
	err           error
	testContainer *tests.PostgresContainer

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

	db  db.Database
	cfg config.Config
}

func TestExpenseRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ExpenseRepositoryTestSuite))
}

func (s *ExpenseRepositoryTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.testContainer, s.err = tests.StartPostgres(s.ctx)
	if s.err != nil {
		panic(s.err)
	}

	s.cfg = config.NewTestConfig(s.testContainer.Port, s.testContainer.Host)

	s.db, s.err = db.New(&s.cfg)
	s.NoError(s.err)
	s.expenseRepo = postgres.NewExpenseRepository(s.db)
	s.userRepo = userrepo.NewUserRepository(s.db)
	s.categoryRepo = postgres2.NewCategoryRepository(s.db)
	s.categoryGroupRepo = postgres2.NewCategoryGroupRepository(s.db)
	s.groupRepo = grouprepo.NewGroupRepository(s.db)

	s.err = s.db.MigrateUp()
	s.NoError(s.err)

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

func (s *ExpenseRepositoryTestSuite) TearDownSuite() {
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
