package expenserepo

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type ExpenseModel struct {
	ID          int          `db:"id"`
	Name        string       `db:"name"`
	AmountCents int          `db:"amount_cents"`
	Description string       `db:"description"`
	GroupID     int          `db:"group_id"`
	CategoryID  int          `db:"category_id"`
	SplitRatio  SplitRatio   `db:"split_ratio"`
	PayerID     int          `db:"payer_id"`
	ReceiverID  int          `db:"receiver_id"`
	CreatedAt   time.Time    `db:"created_at"`
	UpdatedAt   time.Time    `db:"updated_at"`
	DeletedAt   sql.NullTime `db:"deleted_at"`
	Version     int          `db:"version"`
}

type SplitRatio struct {
	Payer    int `db:"payer" json:"payer"`
	Receiver int `db:"receiver" json:"receiver"`
}

func (sr SplitRatio) Value() (driver.Value, error) {
	return json.Marshal(sr)
}

func (sr *SplitRatio) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &sr)
}
