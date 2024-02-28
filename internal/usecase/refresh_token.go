package usecase

import (
	"context"
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	"github.com/Beigelman/nossas-despesas/internal/domain/repository"
	"github.com/Beigelman/nossas-despesas/internal/domain/service"
)

type RefreshAuthTokenParams struct {
	RefreshToken string
}

type RefreshAuthTokenResponse struct {
	User         *entity.User
	Token        string
	RefreshToken string
}

type RefreshAuthToken func(ctx context.Context, p RefreshAuthTokenParams) (*RefreshAuthTokenResponse, error)

func NewRefreshAuthToken(userRepo repository.UserRepository, tokenProvider service.TokenProvider) RefreshAuthToken {
	return func(ctx context.Context, p RefreshAuthTokenParams) (*RefreshAuthTokenResponse, error) {
		token, err := tokenProvider.ParseRefreshToken(p.RefreshToken)
		if err != nil {
			return nil, fmt.Errorf("tokenProvider.ParseToken: %w", err)
		}

		user, err := userRepo.GetByID(ctx, entity.UserID{Value: token.Claims.UserID})
		if err != nil {
			return nil, fmt.Errorf("userRepo.GetByID: %w", err)
		}

		if user == nil {
			return nil, fmt.Errorf("user not found")
		}

		authToken, refreshToken, err := tokenProvider.GenerateUserTokens(*user)
		if err != nil {
			return nil, fmt.Errorf("tokenProvider.GenerateUserTokens: %w", err)
		}

		return &RefreshAuthTokenResponse{
			User:         user,
			Token:        authToken,
			RefreshToken: refreshToken,
		}, nil
	}
}
