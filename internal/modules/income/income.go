package income

import (
	"context"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
)

type Type string

func (i Type) String() string {
	return string(i)
}

var Types = struct {
	Salary           Type
	Benefit          Type
	ThirteenthSalary Type
	Vacation         Type
	Other            Type
}{
	Salary:           "salary",
	Benefit:          "benefit",
	ThirteenthSalary: "thirteenth_salary",
	Vacation:         "vacation",
	Other:            "other",
}

type ID struct{ Value int }

type Income struct {
	ddd.Entity[ID]
	UserID user.ID
	Amount int
	Type   Type
}

type Attributes struct {
	ID        ID
	UserID    user.ID
	Amount    int
	Type      Type
	CreatedAt *time.Time
}

func New(params Attributes) *Income {
	createdAt := time.Now()
	if params.CreatedAt != nil {
		createdAt = *params.CreatedAt
	}
	return &Income{
		Entity: ddd.Entity[ID]{
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

type UpdateAttributes struct {
	Amount    *int
	Type      *Type
	CreatedAt *time.Time
}

func (i *Income) Update(attr UpdateAttributes) {
	if attr.Type != nil {
		i.Type = *attr.Type
	}

	if attr.Amount != nil {
		i.Amount = *attr.Amount
	}

	if attr.CreatedAt != nil {
		i.CreatedAt = *attr.CreatedAt
	}

	i.UpdatedAt = time.Now()
}

func (i *Income) Delete() {
	now := time.Now()
	i.UpdatedAt = now
	i.DeletedAt = &now
}

type Repository interface {
	ddd.Repository[ID, Income]
	GetUserMonthlyIncomes(ctx context.Context, userID user.ID, date *time.Time) ([]Income, error)
}
