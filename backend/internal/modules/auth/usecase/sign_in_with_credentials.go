package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Beigelman/nossas-despesas/internal/modules/auth"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/Beigelman/nossas-despesas/internal/shared/service"
)

type SignInWithCredentialsParams struct {
	Email    string
	Password string
}

type SignInWithCredentialsResponse struct {
	User         *user.User
	Token        string
	RefreshToken string
}

type SignInWithCredentials func(ctx context.Context, p SignInWithCredentialsParams) (*SignInWithCredentialsResponse, error)

func NewSignInWithCredentials(userRepo user.Repository, authRepo auth.Repository, tokenProvider service.TokenProvider) SignInWithCredentials {
	return func(ctx context.Context, p SignInWithCredentialsParams) (*SignInWithCredentialsResponse, error) {
		credentialAuth, err := authRepo.GetByEmail(ctx, p.Email, auth.Types.Credentials)
		if err != nil {
			return nil, fmt.Errorf("authRepo.GetByEmail: %w", err)
		}

		if credentialAuth == nil {
			return nil, except.BadRequestError("incorrect email or password")
		}

		if !credentialAuth.CheckPassword(p.Password) {
			return nil, except.BadRequestError("incorrect email or password")
		}

		usr, err := userRepo.GetByEmail(ctx, credentialAuth.Email)
		if err != nil {
			return nil, fmt.Errorf("userRepo.GetByEmail: %w", err)
		}

		if usr == nil {
			slog.Warn("token with user not found", slog.String("email", credentialAuth.Email))
			return nil, except.NotFoundError("user not found")
		}

		authToken, refreshToken, err := tokenProvider.GenerateUserTokens(*usr)
		if err != nil {
			return nil, fmt.Errorf("tokenProvider.GenerateUserTokens: %w", err)
		}

		return &SignInWithCredentialsResponse{
			User:         usr,
			Token:        authToken,
			RefreshToken: refreshToken,
		}, nil
	}
}
