package usecase

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"github.com/Beigelman/ludaapi/internal/domain/repository"
	"github.com/Beigelman/ludaapi/internal/pkg/except"
)

type CreateUserParams struct {
	Name             string
	Email            string
	ProfilePicture   *string
	AuthenticationID *string
	GroupID          *entity.GroupID
}

type CreateUser func(ctx context.Context, p CreateUserParams) (*entity.User, error)

func NewCreateUser(repo repository.UserRepository, auth *auth.Client) CreateUser {
	return func(ctx context.Context, p CreateUserParams) (*entity.User, error) {
		alreadyExists, err := repo.GetByEmail(ctx, p.Email)
		if err != nil {
			return nil, fmt.Errorf("repo.GetByName: %w", err)
		}

		if alreadyExists != nil {
			return nil, except.ConflictError("email already exists")
		}

		userID := repo.GetNextID()

		user := entity.NewUser(entity.UserParams{
			ID:               userID,
			Name:             p.Name,
			Email:            p.Email,
			ProfilePicture:   p.ProfilePicture,
			GroupID:          p.GroupID,
			AuthenticationID: p.AuthenticationID,
		})

		if err := repo.Store(ctx, user); err != nil {
			return nil, fmt.Errorf("repo.Store: %w", err)
		}

		return user, nil
	}
}
