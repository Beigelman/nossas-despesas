package usecase

import (
	"context"
	"fmt"

	"github.com/Beigelman/nossas-despesas/internal/modules/auth"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/Beigelman/nossas-despesas/internal/shared/service"
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

func NewSignInWithGoogle(userRepo user.Repository, authRepo auth.Repository, tokenProvider service.TokenProvider, googleValidator service.GoogleTokenValidator) SignInWithGoogle {
	return func(ctx context.Context, p SignInWithGoogleParams) (*SignInWithGoogleResponse, error) {
		claims, err := googleValidator.ValidateToken(ctx, p.IdToken)
		if err != nil {
			return nil, fmt.Errorf("googleValidator.ValidateToken: %w", err)
		}

		if claims.Email == "" {
			return nil, except.UnprocessableEntityError("email not found in token")
		}

		if claims.Name == "" {
			return nil, except.UnprocessableEntityError("user name not found in token")
		}

		if claims.Sub == "" {
			return nil, except.UnprocessableEntityError("sub not found in token")
		}

		// Check se o usuário já existe
		existingUser, err := userRepo.GetByEmail(ctx, claims.Email)
		if err != nil {
			return nil, fmt.Errorf("userRepo.GetByEmail: %w", err)
		}

		var usr *user.User
		if existingUser != nil {
			if existingUser.ProfilePicture == nil && claims.Picture != nil {
				existingUser.ProfilePicture = claims.Picture
				if err := userRepo.Store(ctx, existingUser); err != nil {
					return nil, fmt.Errorf("userRepo.Store: %w", err)
				}
			}
			usr = existingUser
		} else {
			usr = user.New(user.Attributes{
				ID:             userRepo.GetNextID(),
				Name:           claims.Name,
				Email:          claims.Email,
				ProfilePicture: claims.Picture,
			})

			if err := userRepo.Store(ctx, usr); err != nil {
				return nil, fmt.Errorf("userRepo.Store: %w", err)
			}
		}

		// Check se a autenticação já existe
		existingAuth, err := authRepo.GetByEmail(ctx, claims.Email, auth.Types.Google)
		if err != nil {
			return nil, fmt.Errorf("authRepo.GetByEmail: %w", err)
		}
		// TODO: isso aqui faz sentido?? Não deveria ser == nil?
		if existingAuth == nil {
			authentic := auth.NewGoogleAuth(auth.GoogleAuthAttributes{
				ID:         authRepo.GetNextID(),
				Email:      claims.Email,
				ProviderID: claims.Sub,
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
