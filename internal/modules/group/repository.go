package group

import (
	"context"
	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
)

type Repository interface {
	ddd.Repository[ID, Group]
	GetByName(ctx context.Context, name string) (*Group, error)
}

type InviteRepository interface {
	ddd.Repository[InviteID, Invite]
	GetGroupInvitesByEmail(ctx context.Context, groupID ID, email string) ([]Invite, error)
	GetByToken(ctx context.Context, token string) (*Invite, error)
}
