package importincomes

import (
	"context"
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/modules/income/postgres"
	"strconv"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/config"
	"github.com/Beigelman/nossas-despesas/internal/modules/income"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"

	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
	"github.com/Beigelman/nossas-despesas/internal/pkg/env"
	"github.com/Beigelman/nossas-despesas/scripts/utils"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use: "import-incomes",
	Run: run,
}

var (
	danId, luId int
	environment string
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
	incomesRepo := postgres.NewIncomeRepository(database)

	file, err := utils.ReadCSVFile("./scripts/data/incomes.csv")
	if err != nil {
		panic(fmt.Errorf("error reading csv file %w", err))
	}
	bar := progressbar.Default(int64(len(file)))
	for _, line := range file {
		// Data			Daniel Beigelman		Luíza Brito
		// 2023-06-01	100000					12000
		date, err := time.Parse("2006/01/02", line[0])
		if err != nil {
			panic(fmt.Errorf("error parsing date: %w", err))
		}

		offsetDate := date.Add(12 * time.Hour)

		danAmount, err := strconv.ParseFloat(line[1], 64)
		if err != nil {
			panic(fmt.Errorf("error parsing amount %w", err))
		}
		danIncomeCents := int(100 * danAmount)

		luAmount, err := strconv.ParseFloat(line[2], 64)
		if err != nil {
			panic(fmt.Errorf("error parsing amount %w", err))
		}
		luIncomeCents := int(100 * luAmount)

		danIncome := income.New(income.Attributes{
			ID:        incomesRepo.GetNextID(),
			UserID:    user.ID{Value: danId},
			Amount:    danIncomeCents,
			Type:      income.Types.Salary,
			CreatedAt: &offsetDate,
		})

		luIncome := income.New(income.Attributes{
			ID:        incomesRepo.GetNextID(),
			UserID:    user.ID{Value: luId},
			Amount:    luIncomeCents,
			Type:      income.Types.Salary,
			CreatedAt: &offsetDate,
		})

		if err := incomesRepo.Store(ctx, danIncome); err != nil {
			fmt.Println(fmt.Errorf("error storing expense %w", err))
		}

		if err := incomesRepo.Store(ctx, luIncome); err != nil {
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
	cmd.Flags().StringVarP(&environment, "env", "e", "development", "environment to run the script (dev, stg, prd)")
}
