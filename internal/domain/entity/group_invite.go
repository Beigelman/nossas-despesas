package entity

import (
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
	"time"
)

type GroupInviteStatus string

var GroupInviteStatuses = struct {
	Pending  GroupInviteStatus
	Sent     GroupInviteStatus
	Accepted GroupInviteStatus
}{
	Pending:  "pending",
	Sent:     "sent",
	Accepted: "accepted",
}

type GroupInviteID struct{ Value int }

type GroupInvite struct {
	ddd.Entity[GroupInviteID]
	GroupID   GroupID
	Token     string
	Email     string
	ExpiresAt time.Time
	Status    GroupInviteStatus
}

type GroupInviteParams struct {
	ID        GroupInviteID
	GroupID   GroupID
	Token     string
	Email     string
	ExpiresAt time.Time
}

func NewGroupInvite(params GroupInviteParams) *GroupInvite {
	return &GroupInvite{
		Entity: ddd.Entity[GroupInviteID]{
			ID:        params.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Version:   0,
		},
		GroupID:   params.GroupID,
		Token:     params.Token,
		Email:     params.Email,
		ExpiresAt: params.ExpiresAt,
		Status:    GroupInviteStatuses.Pending,
	}
}

func (g *GroupInvite) InviteURL(basePath string) string {
	return fmt.Sprintf("%s/group-invite/%s", basePath, g.Token)
}

func (g *GroupInvite) CheckStatus() error {
	if g.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("group invite expired")
	}

	if g.Status != GroupInviteStatuses.Sent {
		return fmt.Errorf("group invite status is %s and not sent", g.Status)
	}

	return nil
}

func (g *GroupInvite) Accept() error {
	if g.Status != GroupInviteStatuses.Sent {
		return fmt.Errorf("group invite was not sent")
	}

	g.Status = GroupInviteStatuses.Accepted
	return nil
}

func (g *GroupInvite) Sent() error {
	if g.Status != GroupInviteStatuses.Pending && g.Status != GroupInviteStatuses.Sent {
		return fmt.Errorf("group invite is not pending")
	}

	g.Status = GroupInviteStatuses.Sent
	return nil
}
