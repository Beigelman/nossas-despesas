package pubsub

import (
	"time"

	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
)

type Event struct {
	Type    string
	GroupID entity.GroupID
	UserID  entity.UserID
	SentAt  time.Time
}

type IncomeEvent struct {
	Event
	Income entity.Income
}
