package repository

import (
	"context"
	"github.com/Beigelman/nossas-despesas/internal/domain/entity"

	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
)

type GroupInviteRepository interface {
	ddd.Repository[entity.GroupInviteID, entity.GroupInvite]
	GetGroupInvitesByEmail(ctx context.Context, groupID entity.GroupID, email string) ([]entity.GroupInvite, error)
	GetByToken(ctx context.Context, token string) (*entity.GroupInvite, error)
}
