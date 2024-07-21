package usecase

import (
	"context"
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/modules/auth"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/Beigelman/nossas-despesas/internal/shared/service"
	"google.golang.org/api/idtoken"
)

type SignInWithGoogleParams struct {
	IdToken string
}

type SignInWithGoogleResponse struct {
	User         *user.User
	Token        string
	RefreshToken string
}

type SignInWithGoogle func(ctx context.Context, p SignInWithGoogleParams) (*SignInWithGoogleResponse, error)

func NewSignInWithGoogle(userRepo user.Repository, authRepo auth.Repository, tokenProvider service.TokenProvider) SignInWithGoogle {
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

		var usr *user.User
		if existingUser != nil {
			if existingUser.ProfilePicture == nil && profilePicture != nil {
				existingUser.ProfilePicture = profilePicture
				if err := userRepo.Store(ctx, existingUser); err != nil {
					return nil, fmt.Errorf("userRepo.Store: %w", err)
				}
			}
			usr = existingUser
		} else {
			usr = user.New(user.Attributes{
				ID:             userRepo.GetNextID(),
				Name:           name,
				Email:          email,
				ProfilePicture: profilePicture,
			})

			if err := userRepo.Store(ctx, usr); err != nil {
				return nil, fmt.Errorf("userRepo.Store: %w", err)
			}
		}

		// Check se a autenticação já existe
		existingAuth, err := authRepo.GetByEmail(ctx, email, auth.Types.Google)
		if err != nil {
			return nil, fmt.Errorf("authRepo.GetByEmail: %w", err)
		}
		// TODO: isso aqui faz sentido?? Não deveria ser == nil?
		if existingAuth != nil {
			authentic := auth.NewGoogleAuth(auth.GoogleAuthAttributes{
				ID:         authRepo.GetNextID(),
				Email:      email,
				ProviderID: providerId,
			})

			if err := authRepo.Store(ctx, authentic); err != nil {
				return nil, fmt.Errorf("authRepo.Store: %w", err)
			}
		}

		// Geração do token de autenticação
		authToken, refreshToken, err := tokenProvider.GenerateUserTokens(*usr)
		if err != nil {
			return nil, fmt.Errorf("tokenProvider.GenerateUserTokens: %w", err)
		}

		return &SignInWithGoogleResponse{
			User:         usr,
			Token:        authToken,
			RefreshToken: refreshToken,
		}, nil
	}
}
