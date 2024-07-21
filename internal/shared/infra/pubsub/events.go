package pubsub

import (
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/modules/income"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"time"
)

type Event struct {
	Type    string
	GroupID group.ID
	UserID  user.ID
	SentAt  time.Time
}

type IncomeEvent struct {
	Event
	Income income.Income
}
