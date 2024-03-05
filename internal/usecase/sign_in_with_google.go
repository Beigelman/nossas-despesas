package usecase

import (
	"context"
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	"github.com/Beigelman/nossas-despesas/internal/domain/repository"
	"github.com/Beigelman/nossas-despesas/internal/domain/service"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"google.golang.org/api/idtoken"
)

type SignInWithGoogleParams struct {
	IdToken string
}

type SignInWithGoogleResponse struct {
	User         *entity.User
	Token        string
	RefreshToken string
}

type SignInWithGoogle func(ctx context.Context, p SignInWithGoogleParams) (*SignInWithGoogleResponse, error)

func NewSignInWithGoogle(userRepo repository.UserRepository, authRepo repository.AuthRepository, tokenProvider service.TokenProvider) SignInWithGoogle {
	return func(ctx context.Context, p SignInWithGoogleParams) (*SignInWithGoogleResponse, error) {
		token, err := idtoken.Validate(ctx, p.IdToken, "")
		if err != nil {
			return nil, fmt.Errorf("idtoken.Validate: %w", err)
		}

		email, ok := token.Claims["email"].(string)
		if !ok {
			return nil, except.UnprocessableEntityError("email not found in token")
		}

		name, ok := token.Claims["name"].(string)
		if !ok {
			return nil, except.UnprocessableEntityError("user name not found in token")
		}

		providerId, ok := token.Claims["sub"].(string)
		if !ok {
			return nil, except.UnprocessableEntityError("sub not found in token")
		}

		var profilePicture *string
		picture, ok := token.Claims["picture"].(string)
		if !ok {
			profilePicture = nil
		} else {
			profilePicture = &picture
		}

		// Check se o usuário já existe
		existingUser, err := userRepo.GetByEmail(ctx, email)
		if err != nil {
			return nil, fmt.Errorf("userRepo.GetByEmail: %w", err)
		}

		var user *entity.User
		if existingUser != nil {
			if existingUser.ProfilePicture == nil && profilePicture != nil {
				existingUser.ProfilePicture = profilePicture
				if err := userRepo.Store(ctx, existingUser); err != nil {
					return nil, fmt.Errorf("userRepo.Store: %w", err)
				}
			}
			user = existingUser
		} else {
			user = entity.NewUser(entity.UserParams{
				ID:             userRepo.GetNextID(),
				Name:           name,
				Email:          email,
				ProfilePicture: profilePicture,
			})

			if err := userRepo.Store(ctx, user); err != nil {
				return nil, fmt.Errorf("userRepo.Store: %w", err)
			}
		}

		// Check se a autenticação já existe
		existingAuth, err := authRepo.GetByEmail(ctx, email, entity.AuthTypes.Google)
		if err != nil {
			return nil, fmt.Errorf("authRepo.GetByEmail: %w", err)
		}

		if existingAuth != nil {
			auth := entity.NewGoogleAuth(entity.GoogleAuthParams{
				ID:         authRepo.GetNextID(),
				Email:      email,
				ProviderID: providerId,
			})

			if err := authRepo.Store(ctx, auth); err != nil {
				return nil, fmt.Errorf("authRepo.Store: %w", err)
			}
		}

		// Geração do token de autenticação
		authToken, refreshToken, err := tokenProvider.GenerateUserTokens(*user)
		if err != nil {
			return nil, fmt.Errorf("tokenProvider.GenerateUserTokens: %w", err)
		}

		return &SignInWithGoogleResponse{
			User:         user,
			Token:        authToken,
			RefreshToken: refreshToken,
		}, nil
	}
}
