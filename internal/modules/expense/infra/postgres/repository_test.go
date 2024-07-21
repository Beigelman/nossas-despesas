package postgres_test

import (
	"context"
	"github.com/Beigelman/nossas-despesas/internal/modules/category"
	postgres2 "github.com/Beigelman/nossas-despesas/internal/modules/category/infra/postgres"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense"
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/modules/group/infra/postgres"
	"testing"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/config"
	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	"github.com/Beigelman/nossas-despesas/internal/domain/repository"
	"github.com/Beigelman/nossas-despesas/internal/infra/postgres/userrepo"
	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/Beigelman/nossas-despesas/internal/tests"
	"github.com/stretchr/testify/suite"
)

type PgExpenseRepoTestSuite struct {
	suite.Suite
	ctx           context.Context
	err           error
	testContainer *tests.PostgresContainer

	expenseRepo       expense.Repository
	userRepo          repository.UserRepository
	categoryRepo      category.Repository
	categoryGroupRepo category.GroupRepository
	groupRepo         group.Repository

	payer         *entity.User
	receiver      *entity.User
	category      *category.Category
	categoryGroup *category.Group
	group         *group.Group

	db  db.Database
	cfg config.Config
}

func TestPgExpenseRepoTestSuite(t *testing.T) {
	suite.Run(t, new(PgExpenseRepoTestSuite))
}

func (s *PgExpenseRepoTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.testContainer, s.err = tests.StartPostgres(s.ctx)
	if s.err != nil {
		panic(s.err)
	}

	s.cfg = config.NewTestConfig(s.testContainer.Port, s.testContainer.Host, "postgres")

	s.db = db.New(&s.cfg)
	s.expenseRepo = NewPGRepository(s.db)
	s.userRepo = userrepo.NewPGRepository(s.db)
	s.categoryRepo = postgres2.NewPGRepository(s.db)
	s.categoryGroupRepo = postgres2.NewPGRepository(s.db)
	s.groupRepo = postgres.NewGroupRepository(s.db)

	s.err = s.db.MigrateUp()
	s.NoError(s.err)

	s.payer = entity.NewUser(entity.UserParams{
		ID:    s.userRepo.GetNextID(),
		Name:  "Payer",
		Email: "payer@email.com",
	})

	s.receiver = entity.NewUser(entity.UserParams{
		ID:    s.userRepo.GetNextID(),
		Name:  "Receiver",
		Email: "receiver@email.com",
	})

	s.categoryGroup = category.NewCategoryGroup(category.GroupAttributes{
		ID:   s.categoryGroupRepo.GetNextID(),
		Name: "Category",
		Icon: "test",
	})

	s.category = category.NewCategory(category.Attributes{
		ID:              s.categoryRepo.GetNextID(),
		Name:            "Category",
		Icon:            "test",
		CategoryGroupID: s.categoryGroup.ID,
	})

	s.group = group.NewGroup(group.Attributes{
		ID:   s.groupRepo.GetNextID(),
		Name: "Group",
	})

	s.NoError(s.userRepo.Store(s.ctx, s.payer))
	s.NoError(s.userRepo.Store(s.ctx, s.receiver))
	s.NoError(s.categoryGroupRepo.Store(s.ctx, s.categoryGroup))
	s.NoError(s.categoryRepo.Store(s.ctx, s.category))
	s.NoError(s.groupRepo.Store(s.ctx, s.group))
}

func (s *PgExpenseRepoTestSuite) TearDownSuite() {
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

func (s *PgExpenseRepoTestSuite) TearDownSubTest() {
	s.NoError(s.db.Clean())
}

func (s *PgExpenseRepoTestSuite) TestPgExpenseRepo_Store() {
	id := s.expenseRepo.GetNextID()
	expense, err := expense.New(expense.Attributes{
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

	s.NoError(s.expenseRepo.Store(s.ctx, expense))
}

func (s *PgExpenseRepoTestSuite) TestPgExpenseRepo_GetByID() {
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

func (s *PgExpenseRepoTestSuite) TestPgExpenseRepo_GetByGroupDate() {
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
