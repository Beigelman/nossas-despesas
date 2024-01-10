package usecase

import (
	"context"
	"fmt"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"github.com/Beigelman/ludaapi/internal/domain/repository"
	"github.com/Beigelman/ludaapi/internal/domain/service"
	"github.com/Beigelman/ludaapi/internal/pkg/except"
)

type SignUpWithCredentialsParams struct {
	Name                 string
	Email                string
	Password             string
	ConfirmationPassword string
	ProfilePicture       *string
	GroupID              *entity.GroupID
}

type SignUpWithCredentialsResponse struct {
	User         *entity.User
	Token        string
	RefreshToken string
}

type SignUpWithCredentials func(ctx context.Context, p SignUpWithCredentialsParams) (*SignUpWithCredentialsResponse, error)

func NewSignUpWithCredentials(userRepo repository.UserRepository, authRepo repository.AuthRepository, tokenProvider service.TokenProvider) SignUpWithCredentials {
	return func(ctx context.Context, p SignUpWithCredentialsParams) (*SignUpWithCredentialsResponse, error) {
		existingAuth, err := authRepo.GetByEmail(ctx, p.Email, entity.AuthTypes.Credentials)
		if err != nil {
			return nil, fmt.Errorf("authRepo.GetByEmail: %w", err)
		}

		if existingAuth != nil {
			return nil, except.BadRequestError("email already registered")
		}

		if p.Password != p.ConfirmationPassword {
			return nil, except.UnprocessableEntityError("passwords do not match")
		}

		existingUser, err := userRepo.GetByEmail(ctx, p.Email)
		if err != nil {
			return nil, fmt.Errorf("userRepo.GetByEmail: %w", err)
		}

		var user *entity.User
		if existingUser != nil {
			user = existingUser
		} else {
			user = entity.NewUser(entity.UserParams{
				ID:             userRepo.GetNextID(),
				Name:           p.Name,
				Email:          p.Email,
				ProfilePicture: p.ProfilePicture,
				GroupID:        p.GroupID,
			})

			if err := userRepo.Store(ctx, user); err != nil {
				return nil, fmt.Errorf("userRepo.Store: %w", err)
			}
		}

		auth, err := entity.NewCredentialAuth(entity.CredentialsAuthParams{
			ID:       authRepo.GetNextID(),
			Email:    p.Email,
			Password: p.Password,
		})
		if err != nil {
			return nil, fmt.Errorf("entity.NewCredentialAuth: %w", err)
		}

		if err := authRepo.Store(ctx, auth); err != nil {
			return nil, fmt.Errorf("authRepo.Store: %w", err)
		}

		authToken, refreshToken, err := tokenProvider.GenerateUserTokens(*user)
		if err != nil {
			return nil, fmt.Errorf("tokenProvider.GenerateUserTokens: %w", err)
		}

		return &SignUpWithCredentialsResponse{
			User:         user,
			Token:        authToken,
			RefreshToken: refreshToken,
		}, nil
	}
}
