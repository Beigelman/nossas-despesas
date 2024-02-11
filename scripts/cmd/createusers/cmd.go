package createusers

import (
	"context"
	"fmt"
	"github.com/Beigelman/ludaapi/internal/config"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"github.com/Beigelman/ludaapi/internal/infra/postgres/authrepo"
	"github.com/Beigelman/ludaapi/internal/infra/postgres/grouprepo"
	"github.com/Beigelman/ludaapi/internal/infra/postgres/incomerepo"
	"github.com/Beigelman/ludaapi/internal/infra/postgres/userrepo"
	"github.com/Beigelman/ludaapi/internal/pkg/db"
	"github.com/Beigelman/ludaapi/internal/pkg/env"
	"github.com/spf13/cobra"
)

var environment string

var cmd = &cobra.Command{
	Use: "create-users",
	Run: run,
}

func run(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	cfg := config.New(env.Environment(environment))
	cfg.SetConfigPath("./internal/config/config.yml")
	if err := cfg.LoadConfig(); err != nil {
		panic(fmt.Errorf("cfg.LoadConfig: %w", err))
	}

	database := db.New(&cfg)
	groupRepo := grouprepo.NewPGRepository(database)
	usersRepo := userrepo.NewPGRepository(database)
	authRepo := authrepo.NewPGRepository(database)
	incomeRepo := incomerepo.NewPGRepository(database)

	group := entity.NewGroup(entity.GroupParams{
		ID:   groupRepo.GetNextID(),
		Name: "Luiel",
	})

	if err := groupRepo.Store(ctx, group); err != nil {
		panic(fmt.Errorf("error saving group: %w", err))
	}

	dan := entity.NewUser(entity.UserParams{
		ID:      usersRepo.GetNextID(),
		Name:    "Daniel Beigelman",
		Email:   "daniel.b.beigelman@gmail.com",
		GroupID: &group.ID,
	})

	if err := usersRepo.Store(ctx, dan); err != nil {
		panic(fmt.Errorf("error saving user: %w", err))
	}

	danCreds, _ := entity.NewCredentialAuth(entity.CredentialsAuthParams{
		ID:       authRepo.GetNextID(),
		Email:    dan.Email,
		Password: "12345678",
	})

	if err := authRepo.Store(ctx, danCreds); err != nil {
		panic(fmt.Errorf("error saving user credentials: %w", err))
	}

	danIncome := entity.NewIncome(entity.IncomeParams{
		ID:     incomeRepo.GetNextID(),
		UserID: dan.ID,
		Amount: 1000000,
		Type:   entity.IncomeTypes.Salary,
	})

	if err := incomeRepo.Store(ctx, danIncome); err != nil {
		panic(fmt.Errorf("error saving income: %w", err))
	}

	lu := entity.NewUser(entity.UserParams{
		ID:      usersRepo.GetNextID(),
		Name:    "Luíza Brito",
		Email:   "brito.luiza27@gmail.com",
		GroupID: &group.ID,
	})

	if err := usersRepo.Store(ctx, lu); err != nil {
		panic(fmt.Errorf("error saving user: %w", err))
	}

	luCreds, _ := entity.NewCredentialAuth(entity.CredentialsAuthParams{
		ID:       authRepo.GetNextID(),
		Email:    lu.Email,
		Password: "12345678",
	})

	if err := authRepo.Store(ctx, luCreds); err != nil {
		panic(fmt.Errorf("error saving user credentials: %w", err))
	}

	luIncome := entity.NewIncome(entity.IncomeParams{
		ID:     incomeRepo.GetNextID(),
		UserID: lu.ID,
		Amount: 900000,
		Type:   entity.IncomeTypes.Salary,
	})

	if err := incomeRepo.Store(ctx, luIncome); err != nil {
		panic(fmt.Errorf("error saving income: %w", err))
	}

	if err := database.Close(); err != nil {
		fmt.Println(fmt.Errorf("error closing database %w", err))
	}
}

func init() {
	cmd.Flags().StringVarP(&environment, "env", "e", "development", "environment to run the script (local, dev, prod, etc)")
}

func Cmd() *cobra.Command {
	return cmd
}
