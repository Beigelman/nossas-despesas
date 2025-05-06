package importsplit

import (
	"context"
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/pkg/config"

	"github.com/Beigelman/nossas-despesas/internal/modules/expense/postgres"
	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/Beigelman/nossas-despesas/internal/pkg/env"
	"github.com/Beigelman/nossas-despesas/scripts/utils"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use: "import-from-split-wize",
	Run: run,
}

var (
	danId, luId, groupId int
	environment          string
)

func run(_ *cobra.Command, _ []string) {
	ctx := context.Background()

	cfg := config.New(env.Environment(environment))
	cfg.SetConfigPath("./internal/config/config.yml")
	if err := cfg.LoadConfig(); err != nil {
		panic(fmt.Errorf("cfg.LoadConfig: %w", err))
	}

	database, err := db.New(&cfg)
	if err != nil {
		panic(err)
	}
	expensesRepo := postgres.NewExpenseRepository(database)

	file, err := utils.ReadCSVFile("./scripts/data/luiel.csv")
	if err != nil {
		panic(fmt.Errorf("error reading csv file %w", err))
	}

	bar := progressbar.Default(int64(len(file)))
	for _, line := range file {
		expense, err := extractExpense(line, expensesRepo.GetNextID())
		if err != nil {
			fmt.Println(fmt.Errorf("error extracting expense %w", err))
		}

		if err := expensesRepo.Store(ctx, expense); err != nil {
			fmt.Println(fmt.Errorf("error storing expense %w", err))
		}

		if err := bar.Add(1); err != nil {
			fmt.Println(fmt.Errorf("error incrementing progress bar %w", err))
		}
	}

	if err := database.Close(); err != nil {
		fmt.Println(fmt.Errorf("error closing database %w", err))
	}
}

func Cmd() *cobra.Command {
	return cmd
}

func init() {
	cmd.Flags().IntVarP(&danId, "dan-id", "d", 1, "dan id")
	cmd.Flags().IntVarP(&luId, "lu-id", "l", 2, "lu id")
	cmd.Flags().IntVarP(&groupId, "group-id", "g", 1, "group id")
	cmd.Flags().StringVarP(&environment, "env", "e", "development", "environment to run the script (dev, stg, prd)")
}
