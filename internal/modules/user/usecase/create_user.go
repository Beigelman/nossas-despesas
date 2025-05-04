package usecase

import (
	"context"
	"fmt"

	"github.com/Beigelman/nossas-despesas/internal/modules/group"

	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
)

type CreateUserParams struct {
	Name             string
	Email            string
	ProfilePicture   *string
	AuthenticationID *string
	GroupID          *group.ID
}

type CreateUser func(ctx context.Context, p CreateUserParams) (*user.User, error)

func NewCreateUser(userRepo user.Repository) CreateUser {
	return func(ctx context.Context, p CreateUserParams) (*user.User, error) {
		alreadyExists, err := userRepo.GetByEmail(ctx, p.Email)
		if err != nil {
			return nil, fmt.Errorf("userRepo.GetByName: %w", err)
		}

		if alreadyExists != nil {
			return nil, except.ConflictError("email already exists")
		}

		usr := user.New(user.Attributes{
			ID:             userRepo.GetNextID(),
			Name:           p.Name,
			Email:          p.Email,
			ProfilePicture: p.ProfilePicture,
			GroupID:        p.GroupID,
		})

		if err := userRepo.Store(ctx, usr); err != nil {
			return nil, fmt.Errorf("userRepo.Store: %w", err)
		}

		return usr, nil
	}
}
