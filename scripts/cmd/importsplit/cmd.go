package importsplit

import (
	"context"
	"fmt"
	"github.com/Beigelman/ludaapi/internal/config"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	vo "github.com/Beigelman/ludaapi/internal/domain/valueobject"
	"github.com/Beigelman/ludaapi/internal/infra/postgres/expenserepo"
	"github.com/Beigelman/ludaapi/internal/pkg/db"
	"github.com/Beigelman/ludaapi/internal/pkg/env"
	"github.com/Beigelman/ludaapi/scripts/utils"
	"github.com/spf13/cobra"
	"log"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var cmd = &cobra.Command{
	Use: "import-from-split-wize",
	Run: run,
}

var danId, luId, groupId int
var environment string

func run(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	cfg := config.New(env.Environment(environment))
	cfg.SetConfigPath("./internal/config/config.yml")
	if err := cfg.LoadConfig(); err != nil {
		panic(fmt.Errorf("cfg.LoadConfig: %w", err))
	}

	database := db.New(&cfg)
	expensesRepo := expenserepo.NewPGRepository(database)

	file, err := utils.ReadCSVFile("./scripts/data/luiel.csv")
	if err != nil {
		panic(fmt.Errorf("error reading csv file %w", err))
	}
	expensesCreated := 0
	for _, line := range file {
		//Data			Descrição		Categoria	Custo		Luíza Brito		Daniel Beigelman
		//2023-06-01	Ajuste maio		Geral		255.64		-255.64			255.64
		date, err := time.Parse("2006-01-02", line[0])
		if err != nil {
			panic(fmt.Errorf("error parsing date: %w", err))
		}
		name := line[1]
		category := SplitCategoryToCategory(line[2])
		amount, err := strconv.ParseFloat(line[3], 64)
		if err != nil {
			panic(fmt.Errorf("error parsing amount %w", err))
		}
		amountCents := int(100 * amount)
		danShare, err := strconv.ParseFloat(line[6], 64)
		if err != nil {
			log.Printf("error parsing dan share %v", err)
			panic(err)
		}

		ratio := danShare / amount
		var payerRatio, receiverRatio, payer, receiver int

		if ratio > 0 {
			payer = danId
			receiver = luId
			receiverRatio = int(math.Round(ratio * 100))
			payerRatio = 100 - receiverRatio
		} else {
			payer = luId
			receiver = danId
			receiverRatio = int(math.Round(ratio * -100))
			payerRatio = 100 - receiverRatio
		}

		splitRatio := vo.SplitRatio{
			Payer:    payerRatio,
			Receiver: receiverRatio,
		}

		regex, _ := regexp.Compile(`reembolso|cashback|ajuste`)
		createdAt := date.Add(time.Duration(int(rand.Float64()*86400)) * time.Millisecond)
		description := "Imported from splitwise"
		if regex.FindAllString(strings.ToLower(name), -1) != nil {
			createdAt = time.Time{}
			description = fmt.Sprintf("Imported from splitwise. Essa é uma transação legado que tem o objetivo de manter o balanço das contas. Data original: %s", date.Format("2006-01-02"))
		}

		expense, err := entity.NewExpense(entity.ExpenseParams{
			ID:          expensesRepo.GetNextID(),
			Name:        name,
			Amount:      amountCents,
			Description: description,
			GroupID:     entity.GroupID{Value: groupId},
			CategoryID:  entity.CategoryID{Value: category},
			SplitRatio:  splitRatio,
			PayerID:     entity.UserID{Value: payer},
			ReceiverID:  entity.UserID{Value: receiver},
			CreatedAt:   &createdAt,
		})

		if err := expensesRepo.Store(ctx, expense); err != nil {
			fmt.Println(fmt.Errorf("error storing expense %w", err))
		}
		expensesCreated++
	}

	fmt.Println("Created", expensesCreated, "expenses")

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
