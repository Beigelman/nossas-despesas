package usecase

import (
	"context"
	"fmt"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"github.com/Beigelman/ludaapi/internal/domain/repository"
	"github.com/Beigelman/ludaapi/internal/infra/token"
	"github.com/golang-jwt/jwt/v5"
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

func NewRefreshAuthToken(userRepo repository.UserRepository, tokenProvider *token.JWTProvider) RefreshAuthToken {
	return func(ctx context.Context, p RefreshAuthTokenParams) (*RefreshAuthTokenResponse, error) {
		parsedToken, err := tokenProvider.ParseToken(p.RefreshToken)
		if err != nil {
			return nil, fmt.Errorf("tokenProvider.ParseToken: %w", err)
		}

		claims := parsedToken.Claims.(jwt.MapClaims)
		userID := int(claims["user_id"].(float64))

		user, err := userRepo.GetByID(ctx, entity.UserID{Value: userID})
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
