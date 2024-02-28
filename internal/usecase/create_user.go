package usecase

import (
	"context"
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	"github.com/Beigelman/nossas-despesas/internal/domain/repository"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
)

type CreateUserParams struct {
	Name             string
	Email            string
	ProfilePicture   *string
	AuthenticationID *string
	GroupID          *entity.GroupID
}

type CreateUser func(ctx context.Context, p CreateUserParams) (*entity.User, error)

func NewCreateUser(userRepo repository.UserRepository) CreateUser {
	return func(ctx context.Context, p CreateUserParams) (*entity.User, error) {
		alreadyExists, err := userRepo.GetByEmail(ctx, p.Email)
		if err != nil {
			return nil, fmt.Errorf("userRepo.GetByName: %w", err)
		}

		if alreadyExists != nil {
			return nil, except.ConflictError("email already exists")
		}

		userID := userRepo.GetNextID()

		user := entity.NewUser(entity.UserParams{
			ID:             userID,
			Name:           p.Name,
			Email:          p.Email,
			ProfilePicture: p.ProfilePicture,
			GroupID:        p.GroupID,
		})

		if err := userRepo.Store(ctx, user); err != nil {
			return nil, fmt.Errorf("userRepo.Store: %w", err)
		}

		return user, nil
	}
}
