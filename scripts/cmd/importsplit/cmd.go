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
	"strconv"
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
		danShare, err := strconv.ParseFloat(line[5], 64)
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

		spiceDate := date.Add(time.Duration(int(rand.Float64()*86400)) * time.Millisecond)
		expense, err := entity.NewExpense(entity.ExpenseParams{
			ID:          expensesRepo.GetNextID(),
			Name:        name,
			Amount:      amountCents,
			Description: "imported from splitwise",
			GroupID:     entity.GroupID{Value: groupId},
			CategoryID:  entity.CategoryID{Value: category},
			SplitRatio:  splitRatio,
			PayerID:     entity.UserID{Value: payer},
			ReceiverID:  entity.UserID{Value: receiver},
			CreatedAt:   &spiceDate,
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

func SplitCategoryToCategory(category string) int {
	switch category {
	case "Geral":
		return 61
	case "Mercado":
		return 15
	case "Jantar fora":
		return 16
	case "Táxi":
		return 32
	case "Seguro":
		return 35
	case "Hotel":
		return 54
	case "Filmes":
		return 20
	case "Bebidas alcoólicas":
		return 18
	case "Entretenimento - Outros":
		return 27
	case "Combustível":
		return 28
	case "Presentes":
		return 41
	case "Transporte - Outros":
		return 37
	case "Estacionamento":
		return 29
	case "Casa - Outros":
		return 14
	case "Vida - Outros":
		return 43
	case "Produtos de limpeza":
		return 13
	case "Aluguel":
		return 1
	case "TV/Telefone/Internet":
		return 6
	case "Manutenção":
		return 8
	case "Eletricidade":
		return 4
	case "Aquecimento/gás":
		return 5
	case "Vestuário":
		return 42
	case "Despesas médicas":
		return 48
	case "Animais de estimação":
		return 9
	case "Móveis":
		return 11
	case "Carro":
		return 37
	case "Eletrônicos":
		return 14
	case "Comidas e bebidas - Outros":
		return 19
	case "Esportes":
		return 38
	case "Serviços":
		return 12
	case "Serviços públicos - Outros":
		return 14
	case "Ônibus/trem":
		return 31
	case "Avião":
		return 53
	case "Educação":
		return 40
	case "Limpeza":
		return 13
	case "Música":
		return 24
	case "Jogos":
		return 22
	case "Impostos":
		return 60
	default:
		return 61
	}
}
