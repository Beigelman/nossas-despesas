package expenserepo_test

import (
	"context"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"github.com/Beigelman/ludaapi/internal/domain/repository"
	"github.com/Beigelman/ludaapi/internal/infra/postgres/categorygrouprepo"
	"testing"
	"time"

	"github.com/Beigelman/ludaapi/internal/config"
	"github.com/Beigelman/ludaapi/internal/infra/postgres/categoryrepo"
	"github.com/Beigelman/ludaapi/internal/infra/postgres/expenserepo"
	"github.com/Beigelman/ludaapi/internal/infra/postgres/grouprepo"
	"github.com/Beigelman/ludaapi/internal/infra/postgres/userrepo"
	"github.com/Beigelman/ludaapi/internal/pkg/db"
	"github.com/Beigelman/ludaapi/internal/tests"
	"github.com/stretchr/testify/suite"
)

var migrationPath = "file:///Users/danielbeigelman/mydev/go-luda/server/database/migrations"

type PgExpenseRepoTestSuite struct {
	suite.Suite
	ctx           context.Context
	err           error
	testContainer *tests.PostgresContainer

	expenseRepo       repository.ExpenseRepository
	userRepo          repository.UserRepository
	categoryRepo      repository.CategoryRepository
	categoryGroupRepo repository.CategoryGroupRepository
	groupRepo         repository.GroupRepository

	payer         *entity.User
	receiver      *entity.User
	category      *entity.Category
	categoryGroup *entity.CategoryGroup
	group         *entity.Group

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

	s.cfg = config.NewTestConfig(s.testContainer.Port, s.testContainer.Host)

	s.db = db.New(&s.cfg)
	s.expenseRepo = expenserepo.NewPGRepository(s.db)
	s.userRepo = userrepo.NewPGRepository(s.db)
	s.categoryRepo = categoryrepo.NewPGRepository(s.db)
	s.categoryGroupRepo = categorygrouprepo.NewPGRepository(s.db)
	s.groupRepo = grouprepo.NewPGRepository(s.db)

	s.err = s.db.MigrateUp(migrationPath)
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

	s.categoryGroup = entity.NewCategoryGroup(entity.CategoryGroupParams{
		ID:   s.categoryGroupRepo.GetNextID(),
		Name: "Category",
		Icon: "test",
	})

	s.category = entity.NewCategory(entity.CategoryParams{
		ID:              s.categoryRepo.GetNextID(),
		Name:            "Category",
		Icon:            "test",
		CategoryGroupID: s.categoryGroup.ID,
	})

	s.group = entity.NewGroup(entity.GroupParams{
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
	s.err = s.db.MigrateDown(migrationPath)
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
	expense, err := entity.NewExpense(entity.ExpenseParams{
		ID:          id,
		Name:        "my first expense",
		Amount:      100,
		Description: "My Description",
		PayerID:     s.payer.ID,
		ReceiverID:  s.receiver.ID,
		SplitRatio: entity.SplitRatio{
			Payer:    50,
			Receiver: 50,
		},
		CategoryID: s.category.ID,
		GroupID:    s.group.ID,
	})
	s.NoError(err)

	s.NoError(s.expenseRepo.Store(s.ctx, expense))
}

func (s *PgExpenseRepoTestSuite) TestPgExpenseRepo_GetByID() {
	id := s.expenseRepo.GetNextID()
	expected, err := entity.NewExpense(entity.ExpenseParams{
		ID:          id,
		Name:        "my first expense",
		Amount:      100,
		Description: "My Description",
		PayerID:     s.payer.ID,
		ReceiverID:  s.receiver.ID,
		SplitRatio: entity.SplitRatio{
			Payer:    50,
			Receiver: 50,
		},
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
