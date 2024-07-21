package createusers

import (
	"context"
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/config"
	"github.com/Beigelman/nossas-despesas/internal/modules/auth"
	authrepo "github.com/Beigelman/nossas-despesas/internal/modules/auth/infra/postgres"
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	grouprepo "github.com/Beigelman/nossas-despesas/internal/modules/group/infra/postgres"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	userrepo "github.com/Beigelman/nossas-despesas/internal/modules/user/infra/postgres"
	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/Beigelman/nossas-despesas/internal/pkg/env"
	"github.com/spf13/cobra"
)

var (
	environment string
	password    string
)

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
	groupRepo := grouprepo.NewGroupRepository(database)
	usersRepo := userrepo.NewUserRepository(database)
	authRepo := authrepo.NewAuthRepository(database)

	grp := group.New(group.Attributes{
		ID:   groupRepo.GetNextID(),
		Name: "Luiel",
	})

	if err := groupRepo.Store(ctx, grp); err != nil {
		panic(fmt.Errorf("error saving group: %w", err))
	}

	dan := user.New(user.Attributes{
		ID:      usersRepo.GetNextID(),
		Name:    "Daniel Beigelman",
		Email:   "daniel.b.beigelman@gmail.com",
		GroupID: &grp.ID,
	})

	if err := usersRepo.Store(ctx, dan); err != nil {
		panic(fmt.Errorf("error saving user: %w", err))
	}

	danCreds, _ := auth.NewCredentialAuth(auth.CredentialsAttributes{
		ID:       authRepo.GetNextID(),
		Email:    dan.Email,
		Password: password,
	})

	if err := authRepo.Store(ctx, danCreds); err != nil {
		panic(fmt.Errorf("error saving user credentials: %w", err))
	}

	lu := user.New(user.Attributes{
		ID:      usersRepo.GetNextID(),
		Name:    "Luíza Brito",
		Email:   "brito.luiza27@gmail.com",
		GroupID: &grp.ID,
	})

	if err := usersRepo.Store(ctx, lu); err != nil {
		panic(fmt.Errorf("error saving user: %w", err))
	}

	luCreds, _ := auth.NewCredentialAuth(auth.CredentialsAttributes{
		ID:       authRepo.GetNextID(),
		Email:    lu.Email,
		Password: password,
	})

	if err := authRepo.Store(ctx, luCreds); err != nil {
		panic(fmt.Errorf("error saving user credentials: %w", err))
	}

	if err := database.Close(); err != nil {
		fmt.Println(fmt.Errorf("error closing database %w", err))
	}
}

func init() {
	cmd.Flags().StringVarP(&environment, "env", "e", "development", "environment to run the script (local, dev, prod, etc)")
	cmd.Flags().StringVarP(&password, "password", "p", "12345678", "password for the users created")
}

func Cmd() *cobra.Command {
	return cmd
}
