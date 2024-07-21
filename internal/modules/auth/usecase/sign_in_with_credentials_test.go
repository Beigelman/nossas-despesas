package usecase_test

import (
	"context"
	"errors"
	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	"github.com/Beigelman/nossas-despesas/internal/modules/auth"
	"github.com/Beigelman/nossas-despesas/internal/modules/auth/usecase"
	mockrepository "github.com/Beigelman/nossas-despesas/internal/tests/mocks/repository"
	mockservice "github.com/Beigelman/nossas-despesas/internal/tests/mocks/service"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSignInWithCredentials(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userRepo := mockrepository.NewMockUserRepository(t)
	authRepo := mockrepository.NewMockAuthRepository(t)
	tokenProvider := mockservice.NewMockTokenProvider(t)

	user := entity.NewUser(entity.UserParams{
		ID:             entity.UserID{Value: 3},
		Name:           "test",
		Email:          "test@email.com",
		ProfilePicture: nil,
		GroupID:        nil,
	})

	authorization, _ := auth.NewCredentialAuth(auth.CredentialsAttributes{
		ID:       auth.ID{Value: 3},
		Email:    "test@email.com",
		Password: "12345678",
	})

	signInWithCredentials := usecase.NewSignInWithCredentials(userRepo, authRepo, tokenProvider)

	t.Run("should return error with authRepo fails", func(t *testing.T) {
		authRepo.EXPECT().GetByEmail(ctx, "test@email.com", auth.Types.Credentials).Return(nil, errors.New("test error")).Once()

		resp, err := signInWithCredentials(ctx, usecase.SignInWithCredentialsParams{
			Email:    "test@email.com",
			Password: "12345678",
		})
		assert.Errorf(t, err, "authRepo.GetByEmail: test error")
		assert.Nil(t, resp)
	})

	t.Run("should return error if no auth method found", func(t *testing.T) {
		authRepo.EXPECT().GetByEmail(ctx, "test@email.com", auth.Types.Credentials).Return(nil, nil).Once()

		resp, err := signInWithCredentials(ctx, usecase.SignInWithCredentialsParams{
			Email:    "test@email.com",
			Password: "12345678",
		})
		assert.Errorf(t, err, "incorrect email or password")
		assert.Nil(t, resp)
	})

	t.Run("should return error if password is incorrect", func(t *testing.T) {
		authRepo.EXPECT().GetByEmail(ctx, "test@email.com", auth.Types.Credentials).Return(authorization, nil).Once()

		resp, err := signInWithCredentials(ctx, usecase.SignInWithCredentialsParams{
			Email:    "test@email.com",
			Password: "12345679",
		})
		assert.Errorf(t, err, "incorrect email or password")
		assert.Nil(t, resp)
	})

	t.Run("should return error if userRepo fails", func(t *testing.T) {
		authRepo.EXPECT().GetByEmail(ctx, "test@email.com", auth.Types.Credentials).Return(authorization, nil).Once()
		userRepo.EXPECT().GetByEmail(ctx, "test@email.com").Return(nil, errors.New("test error")).Once()
		resp, err := signInWithCredentials(ctx, usecase.SignInWithCredentialsParams{
			Email:    "test@email.com",
			Password: "12345678",
		})
		assert.Errorf(t, err, "userRepo.GetByEmail: test error")
		assert.Nil(t, resp)
	})

	t.Run("should return error if no user is found", func(t *testing.T) {
		authRepo.EXPECT().GetByEmail(ctx, "test@email.com", auth.Types.Credentials).Return(authorization, nil).Once()
		userRepo.EXPECT().GetByEmail(ctx, "test@email.com").Return(nil, nil).Once()
		resp, err := signInWithCredentials(ctx, usecase.SignInWithCredentialsParams{
			Email:    "test@email.com",
			Password: "12345678",
		})
		assert.Errorf(t, err, "user not found")
		assert.Nil(t, resp)
	})

	t.Run("should return error if tokenProvide fails", func(t *testing.T) {
		authRepo.EXPECT().GetByEmail(ctx, "test@email.com", auth.Types.Credentials).Return(authorization, nil).Once()
		userRepo.EXPECT().GetByEmail(ctx, "test@email.com").Return(user, nil).Once()
		tokenProvider.EXPECT().GenerateUserTokens(*user).Return("", "", errors.New("test error")).Once()
		resp, err := signInWithCredentials(ctx, usecase.SignInWithCredentialsParams{
			Email:    "test@email.com",
			Password: "12345678",
		})
		assert.Errorf(t, err, "tokenProvider.GenerateUserTokens: test error")
		assert.Nil(t, resp)
	})

	t.Run("happy path", func(t *testing.T) {
		authRepo.EXPECT().GetByEmail(ctx, "test@email.com", auth.Types.Credentials).Return(authorization, nil).Once()
		userRepo.EXPECT().GetByEmail(ctx, "test@email.com").Return(user, nil).Once()
		tokenProvider.EXPECT().GenerateUserTokens(*user).Return("new_token", "new_refresh_token", nil).Once()

		resp, err := signInWithCredentials(ctx, usecase.SignInWithCredentialsParams{
			Email:    "test@email.com",
			Password: "12345678",
		})
		assert.Nil(t, err)
		assert.Equal(t, resp.Token, "new_token")
		assert.Equal(t, resp.RefreshToken, "new_refresh_token")
		assert.Equal(t, resp.User.Name, user.Name)
	})
}
