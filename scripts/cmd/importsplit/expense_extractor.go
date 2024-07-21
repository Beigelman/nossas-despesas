package importsplit

import (
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense"
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
)

func extractExpense(line []string, id expense.ID) (*expense.Expense, error) {
	date, err := time.Parse("2006-01-02", line[0])
	if err != nil {
		panic(fmt.Errorf("error parsing date: %w", err))
	}

	name := line[1]
	category := SplitCategoryToCategory(line[2])
	amount, err := strconv.ParseFloat(line[3], 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing amount %w", err)
	}

	amountCents := int(100 * amount)
	danShare, err := strconv.ParseFloat(line[6], 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing dan share %v", err)
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

	splitRatio := expense.SplitRatio{
		Payer:    payerRatio,
		Receiver: receiverRatio,
	}

	var splitType expense.SplitType
	if payerRatio == 50 || receiverRatio == 50 {
		splitType = expense.SplitTypes.Equal
	} else if receiverRatio == 100 {
		splitType = expense.SplitTypes.Transfer
	} else {
		splitType = expense.SplitTypes.Proportional
	}

	regex, _ := regexp.Compile(`reembolso|cashback|ajuste`)
	createdAt := date.Add(time.Hour*4 + time.Duration(int(rand.Float64()*86400))*time.Millisecond)
	description := "Imported from splitwise"
	if regex.FindAllString(strings.ToLower(name), -1) != nil {
		createdAt = time.Time{}
		description = fmt.Sprintf("Imported from splitwise. Essa é uma transação legado que tem o objetivo de manter o balanço das contas. Data original: %s", date.Format("2006-01-02"))
	}

	return expense.New(expense.Attributes{
		ID:          id,
		Name:        name,
		Amount:      amountCents,
		Description: description,
		GroupID:     group.ID{Value: groupId},
		CategoryID:  category.CategoryID{Value: category},
		SplitRatio:  splitRatio,
		SplitType:   splitType,
		PayerID:     entity.UserID{Value: payer},
		ReceiverID:  entity.UserID{Value: receiver},
		CreatedAt:   &createdAt,
	})
}
