package importincomes

import (
	"context"
	"fmt"
	"github.com/Beigelman/ludaapi/internal/config"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"github.com/Beigelman/ludaapi/internal/infra/postgres/incomerepo"
	"github.com/Beigelman/ludaapi/internal/pkg/db"
	"github.com/Beigelman/ludaapi/internal/pkg/env"
	"github.com/Beigelman/ludaapi/scripts/utils"
	"github.com/spf13/cobra"
	"strconv"
	"time"
)

var cmd = &cobra.Command{
	Use: "import-incomes",
	Run: run,
}

var danId, luId int
var environment string

func run(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	cfg := config.New(env.Environment(environment))
	cfg.SetConfigPath("./internal/config/config.yml")
	if err := cfg.LoadConfig(); err != nil {
		panic(fmt.Errorf("cfg.LoadConfig: %w", err))
	}

	database := db.New(&cfg)
	incomesRepo := incomerepo.NewPGRepository(database)

	file, err := utils.ReadCSVFile("./scripts/data/incomes.csv")
	if err != nil {
		panic(fmt.Errorf("error reading csv file %w", err))
	}
	incomesCreated := 0
	for _, line := range file {
		//Data			Daniel Beigelman		Luíza Brito
		//2023-06-01	100000					12000
		date, err := time.Parse("2006/01/02", line[0])
		if err != nil {
			panic(fmt.Errorf("error parsing date: %w", err))
		}

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

		danIncome := entity.NewIncome(entity.IncomeParams{
			ID:        incomesRepo.GetNextID(),
			UserID:    entity.UserID{Value: danId},
			Amount:    danIncomeCents,
			Type:      entity.IncomeTypes.Salary,
			CreatedAt: &date,
		})

		luIncome := entity.NewIncome(entity.IncomeParams{
			ID:        incomesRepo.GetNextID(),
			UserID:    entity.UserID{Value: luId},
			Amount:    luIncomeCents,
			Type:      entity.IncomeTypes.Salary,
			CreatedAt: &date,
		})

		if err := incomesRepo.Store(ctx, danIncome); err != nil {
			fmt.Println(fmt.Errorf("error storing expense %w", err))
		}

		if err := incomesRepo.Store(ctx, luIncome); err != nil {
			fmt.Println(fmt.Errorf("error storing expense %w", err))
		}

		incomesCreated += 2
	}

	fmt.Println("Created", incomesCreated, "incomes")

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
