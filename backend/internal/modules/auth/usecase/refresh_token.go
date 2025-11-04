package usecase

import (
	"context"
	"fmt"

	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"github.com/Beigelman/nossas-despesas/internal/shared/service"
)

type RefreshAuthTokenParams struct {
	RefreshToken string
}

type RefreshAuthTokenResponse struct {
	User         *user.User
	Token        string
	RefreshToken string
}

type RefreshAuthToken func(ctx context.Context, p RefreshAuthTokenParams) (*RefreshAuthTokenResponse, error)

func NewRefreshAuthToken(userRepo user.Repository, tokenProvider service.TokenProvider) RefreshAuthToken {
	return func(ctx context.Context, p RefreshAuthTokenParams) (*RefreshAuthTokenResponse, error) {
		token, err := tokenProvider.ParseRefreshToken(p.RefreshToken)
		if err != nil {
			return nil, fmt.Errorf("tokenProvider.ParseToken: %w", err)
		}

		usr, err := userRepo.GetByID(ctx, user.ID{Value: token.Claims.UserID})
		if err != nil {
			return nil, fmt.Errorf("userRepo.GetByID: %w", err)
		}

		if usr == nil {
			return nil, fmt.Errorf("user not found")
		}

		authToken, refreshToken, err := tokenProvider.GenerateUserTokens(*usr)
		if err != nil {
			return nil, fmt.Errorf("tokenProvider.GenerateUserTokens: %w", err)
		}

		return &RefreshAuthTokenResponse{
			User:         usr,
			Token:        authToken,
			RefreshToken: refreshToken,
		}, nil
	}
}
