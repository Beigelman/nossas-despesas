package postgres

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/civil"
)

type ExpenseModel struct {
	ID                int           `db:"id"`
	Name              string        `db:"name"`
	AmountCents       int           `db:"amount_cents"`
	RefundAmountCents sql.NullInt64 `db:"refund_amount_cents"`
	Description       string        `db:"description"`
	GroupID           int           `db:"group_id"`
	CategoryID        int           `db:"category_id"`
	SplitRatio        SplitRatio    `db:"split_ratio"`
	SplitType         string        `db:"split_type"`
	PayerID           int           `db:"payer_id"`
	ReceiverID        int           `db:"receiver_id"`
	CreatedAt         time.Time     `db:"created_at"`
	UpdatedAt         time.Time     `db:"updated_at"`
	DeletedAt         sql.NullTime  `db:"deleted_at"`
	Version           int           `db:"version"`
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

type ScheduledExpenseModel struct {
	ID              int                 `db:"id"`
	Name            string              `db:"name"`
	AmountCents     int                 `db:"amount_cents"`
	Description     string              `db:"description"`
	GroupID         int                 `db:"group_id"`
	CategoryID      int                 `db:"category_id"`
	SplitRatio      SplitRatio          `db:"split_ratio"`
	SplitType       string              `db:"split_type"`
	PayerID         int                 `db:"payer_id"`
	ReceiverID      int                 `db:"receiver_id"`
	FrequencyInDays int                 `db:"frequency_in_days"`
	LastGeneratedAt sql.Null[CivilDate] `db:"last_generated_at"`
	IsActive        bool                `db:"is_active"`
	CreatedAt       time.Time           `db:"created_at"`
	UpdatedAt       time.Time           `db:"updated_at"`
	Version         int                 `db:"version"`
}

type CivilDate civil.Date

func (d CivilDate) ToCivilDate() civil.Date {
	return civil.Date(d)
}

func (d *CivilDate) Scan(value interface{}) error {
	switch v := value.(type) {
	case time.Time:
		*d = CivilDate(civil.DateOf(v))
	case string:
		parsed, err := civil.ParseDate(v)
		if err != nil {
			return err
		}
		*d = CivilDate(parsed)
	default:
		return fmt.Errorf("cannot scan type %T into CivilDate", value)
	}
	return nil
}

func (d CivilDate) Value() (driver.Value, error) {
	date := civil.Date(d)
	return date.String(), nil // Return as string to be stored in the database
}
