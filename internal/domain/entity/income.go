package entity

import (
	"github.com/Beigelman/ludaapi/internal/pkg/ddd"
	"time"
)

type IncomeType string

func (i IncomeType) String() string {
	return string(i)
}

var IncomeTypes = struct {
	Salary           IncomeType
	Benefit          IncomeType
	ThirteenthSalary IncomeType
	Vacation         IncomeType
	Other            IncomeType
}{
	Salary:           "salary",
	Benefit:          "benefit",
	ThirteenthSalary: "thirteenth_salary",
	Vacation:         "vacation",
	Other:            "other",
}

type IncomeID struct{ Value int }

type Income struct {
	ddd.Entity[IncomeID]
	UserID UserID
	Amount int
	Type   IncomeType
}

type IncomeParams struct {
	ID        IncomeID
	UserID    UserID
	Amount    int
	Type      IncomeType
	CreatedAt *time.Time
}

func NewIncome(params IncomeParams) *Income {
	createdAt := time.Now()
	if params.CreatedAt != nil {
		createdAt = *params.CreatedAt
	}
	return &Income{
		Entity: ddd.Entity[IncomeID]{
			ID:        params.ID,
			CreatedAt: createdAt,
			UpdatedAt: time.Now(),
			Version:   0,
		},
		UserID: params.UserID,
		Amount: params.Amount,
		Type:   params.Type,
	}
}
