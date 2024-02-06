package cmd

import (
	"context"
	"fmt"
	"github.com/Beigelman/ludaapi/internal/config"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	vo "github.com/Beigelman/ludaapi/internal/domain/valueobject"
	"github.com/Beigelman/ludaapi/internal/infra/postgres/expenserepo"
	"github.com/Beigelman/ludaapi/internal/pkg/db"
	"github.com/Beigelman/ludaapi/scripts/utils"
	"github.com/spf13/cobra"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"
)

var importFromSplitwiseCmd = &cobra.Command{
	Use: "import-from-split-wize",
	Run: run,
}

var danId, luId, groupId int

func run(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	cfg := config.Config{
		Env:         "local",
		ServiceName: "import-script",
		Port:        "8080",
		LogLevel:    "INFO",
		Db: config.Db{
			Host:         "localhost",
			Port:         "5432",
			Name:         "app",
			User:         "root",
			Password:     "root",
			Type:         "postgres",
			MaxIdleConns: 1,
			MaxOpenConns: 1,
		},
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

func init() {
	importFromSplitwiseCmd.Flags().IntVarP(&danId, "dan-id", "d", 1, "dan id")
	importFromSplitwiseCmd.Flags().IntVarP(&luId, "lu-id", "l", 2, "lu id")
	importFromSplitwiseCmd.Flags().IntVarP(&groupId, "group-id", "g", 1, "group id")
}

func SplitCategoryToCategory(category string) int {
	switch category {
	case "Geral":
		return 50
	case "Mercado":
		return 14
	case "Jantar fora":
		return 15
	case "Táxi":
		return 31
	case "Seguro":
		return 34
	case "Hotel":
		return 52
	case "Filmes":
		return 19
	case "Bebidas alcoólicas":
		return 17
	case "Entretenimento - Outros":
		return 26
	case "Combustível":
		return 27
	case "Presentes":
		return 39
	case "Transporte - Outros":
		return 35
	case "Estacionamento":
		return 28
	case "Casa - Outros":
		return 13
	case "Vida - Outros":
		return 41
	case "Produtos de limpeza":
		return 12
	case "Aluguel":
		return 1
	case "TV/Telefone/Internet":
		return 6
	case "Manutenção":
		return 7
	case "Eletricidade":
		return 4
	case "Aquecimento/gás":
		return 5
	case "Vestuário":
		return 40
	case "Despesas médicas":
		return 48
	case "Animais de estimação":
		return 8
	case "Móveis":
		return 10
	case "Carro":
		return 35
	case "Eletrônicos":
		return 13
	case "Comidas e bebidas - Outros":
		return 18
	case "Esportes":
		return 36
	case "Serviços":
		return 11
	case "Serviços públicos - Outros":
		return 13
	case "Ônibus/trem":
		return 54
	case "Avião":
		return 51
	case "Educação":
		return 38
	case "Limpeza":
		return 11
	case "Música":
		return 23
	case "Jogos":
		return 21
	default:
		return 50
	}
}
