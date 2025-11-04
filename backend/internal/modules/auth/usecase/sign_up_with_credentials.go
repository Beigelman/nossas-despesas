package usecase

import (
	"context"
	"fmt"

	"github.com/Beigelman/nossas-despesas/internal/modules/auth"
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/Beigelman/nossas-despesas/internal/shared/service"
)

type SignUpWithCredentialsParams struct {
	Name                 string
	Email                string
	Password             string
	ConfirmationPassword string
	ProfilePicture       *string
	GroupID              *group.ID
}

type SignUpWithCredentialsResponse struct {
	User         *user.User
	Token        string
	RefreshToken string
}

type SignUpWithCredentials func(ctx context.Context, p SignUpWithCredentialsParams) (*SignUpWithCredentialsResponse, error)

func NewSignUpWithCredentials(userRepo user.Repository, authRepo auth.Repository, tokenProvider service.TokenProvider) SignUpWithCredentials {
	return func(ctx context.Context, p SignUpWithCredentialsParams) (*SignUpWithCredentialsResponse, error) {
		existingAuth, err := authRepo.GetByEmail(ctx, p.Email, auth.Types.Credentials)
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

		var usr *user.User
		if existingUser != nil {
			usr = existingUser
		} else {
			usr = user.New(user.Attributes{
				ID:             userRepo.GetNextID(),
				Name:           p.Name,
				Email:          p.Email,
				ProfilePicture: p.ProfilePicture,
				GroupID:        p.GroupID,
			})

			if err := userRepo.Store(ctx, usr); err != nil {
				return nil, fmt.Errorf("userRepo.Store: %w", err)
			}
		}

		authentic, err := auth.NewCredentialAuth(auth.CredentialsAttributes{
			ID:       authRepo.GetNextID(),
			Email:    usr.Email,
			Password: p.Password,
		})
		if err != nil {
			return nil, fmt.Errorf("auth.NewCredentialAuth: %w", err)
		}

		if err := authRepo.Store(ctx, authentic); err != nil {
			return nil, fmt.Errorf("authRepo.Store: %w", err)
		}

		authToken, refreshToken, err := tokenProvider.GenerateUserTokens(*usr)
		if err != nil {
			return nil, fmt.Errorf("tokenProvider.GenerateUserTokens: %w", err)
		}

		return &SignUpWithCredentialsResponse{
			User:         usr,
			Token:        authToken,
			RefreshToken: refreshToken,
		}, nil
	}
}
