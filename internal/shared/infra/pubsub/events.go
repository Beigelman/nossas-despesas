package pubsub

import (
	"time"

	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/modules/income"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
)

// TODO: mover para dentro de algum package de domínio
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
