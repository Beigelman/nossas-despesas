package usecase

import (
	"context"
	"fmt"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"github.com/Beigelman/ludaapi/internal/domain/repository"
	"github.com/Beigelman/ludaapi/internal/domain/service"
	"github.com/Beigelman/ludaapi/internal/pkg/except"
	"log/slog"
)

type SignInWithCredentialsParams struct {
	Email    string
	Password string
}

type SignInWithCredentialsResponse struct {
	User         *entity.User
	Token        string
	RefreshToken string
}

type SignInWithCredentials func(ctx context.Context, p SignInWithCredentialsParams) (*SignInWithCredentialsResponse, error)

func NewSignInWithCredentials(userRepo repository.UserRepository, authRepo repository.AuthRepository, tokenProvider service.TokenProvider) SignInWithCredentials {
	return func(ctx context.Context, p SignInWithCredentialsParams) (*SignInWithCredentialsResponse, error) {
		credentialAuth, err := authRepo.GetByEmail(ctx, p.Email, entity.AuthTypes.Credentials)
		if err != nil {
			return nil, fmt.Errorf("authRepo.GetByEmail: %w", err)
		}

		if credentialAuth == nil {
			return nil, except.BadRequestError("incorrect email or password")
		}

		if !credentialAuth.CheckPassword(p.Password) {
			return nil, except.BadRequestError("incorrect email or password")
		}

		user, err := userRepo.GetByEmail(ctx, credentialAuth.Email)
		if err != nil {
			return nil, fmt.Errorf("userRepo.GetByEmail: %w", err)
		}

		if user == nil {
			slog.Warn("token with user not found", slog.String("email", credentialAuth.Email))
			return nil, except.NotFoundError("user not found")
		}

		authToken, refreshToken, err := tokenProvider.GenerateUserTokens(*user)
		if err != nil {
			return nil, fmt.Errorf("tokenProvider.GenerateUserTokens: %w", err)
		}

		return &SignInWithCredentialsResponse{
			User:         user,
			Token:        authToken,
			RefreshToken: refreshToken,
		}, nil
	}
}
