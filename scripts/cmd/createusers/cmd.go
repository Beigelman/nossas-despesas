package createusers

import (
	"context"
	"fmt"
	"github.com/Beigelman/ludaapi/internal/config"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"github.com/Beigelman/ludaapi/internal/infra/postgres/authrepo"
	"github.com/Beigelman/ludaapi/internal/infra/postgres/grouprepo"
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
