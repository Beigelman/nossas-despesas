package group

import (
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
	"time"
)

type InviteStatus string

var InviteStatuses = struct {
	Pending  InviteStatus
	Sent     InviteStatus
	Accepted InviteStatus
}{
	Pending:  "pending",
	Sent:     "sent",
	Accepted: "accepted",
}

type InviteID struct{ Value int }

type Invite struct {
	ddd.Entity[InviteID]
	GroupID   ID
	Token     string
	Email     string
	ExpiresAt time.Time
	Status    InviteStatus
}

type InviteAttributes struct {
	ID        InviteID
	GroupID   ID
	Token     string
	Email     string
	ExpiresAt time.Time
}

func NewInvite(params InviteAttributes) *Invite {
	return &Invite{
		Entity: ddd.Entity[InviteID]{
			ID:        params.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Version:   0,
		},
		GroupID:   params.GroupID,
		Token:     params.Token,
		Email:     params.Email,
		ExpiresAt: params.ExpiresAt,
		Status:    InviteStatuses.Pending,
	}
}

func (g *Invite) Url(basePath string) string {
	return fmt.Sprintf("%s/group/%s/accept", basePath, g.Token)
}

func (g *Invite) CheckStatus() error {
	if g.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("group invite expired")
	}

	if g.Status != InviteStatuses.Sent {
		return fmt.Errorf("group invite status is %s and not sent", g.Status)
	}

	return nil
}

func (g *Invite) Accept() error {
	if g.Status != InviteStatuses.Sent {
		return fmt.Errorf("group invite was not sent")
	}

	g.Status = InviteStatuses.Accepted
	return nil
}

func (g *Invite) Sent() error {
	if g.Status != InviteStatuses.Pending && g.Status != InviteStatuses.Sent {
		return fmt.Errorf("group invite is not pending")
	}

	g.Status = InviteStatuses.Sent
	return nil
}
