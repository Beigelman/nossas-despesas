package usecase_test

import (
	"context"
	"errors"
	"github.com/Beigelman/nossas-despesas/internal/modules/auth"
	"github.com/Beigelman/nossas-despesas/internal/modules/auth/usecase"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	mockrepository "github.com/Beigelman/nossas-despesas/internal/tests/mocks/repository"
	mockservice "github.com/Beigelman/nossas-despesas/internal/tests/mocks/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestSignUpWithCredentials(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userRepo := mockrepository.NewMockUserRepository(t)
	authRepo := mockrepository.NewMockAuthRepository(t)
	tokenProvider := mockservice.NewMockTokenProvider(t)

	usr := user.New(user.Attributes{
		ID:             user.ID{Value: 1},
		Name:           "test",
		Email:          "test@email.com",
		ProfilePicture: nil,
		GroupID:        nil,
	})

	authorization, _ := auth.NewCredentialAuth(auth.CredentialsAttributes{
		ID:       auth.ID{Value: 1},
		Email:    "test@email.com",
		Password: "12345678",
	})

	signUpWithCredentials := usecase.NewSignUpWithCredentials(userRepo, authRepo, tokenProvider)

	t.Run("should return error with authRepo fails", func(t *testing.T) {
		authRepo.EXPECT().GetByEmail(ctx, "test@email.com", auth.Types.Credentials).Return(nil, errors.New("test error")).Once()

		resp, err := signUpWithCredentials(ctx, usecase.SignUpWithCredentialsParams{
			Name:                 "test",
			Email:                "test@email.com",
			Password:             "12345678",
			ConfirmationPassword: "12345678",
		})
		assert.EqualError(t, err, "authRepo.GetByEmail: test error")
		assert.Nil(t, resp)
	})

	t.Run("should return error if auth method already exists", func(t *testing.T) {
		authRepo.EXPECT().GetByEmail(ctx, "test@email.com", auth.Types.Credentials).Return(authorization, nil).Once()

		resp, err := signUpWithCredentials(ctx, usecase.SignUpWithCredentialsParams{
			Name:                 "test",
			Email:                "test@email.com",
			Password:             "12345678",
			ConfirmationPassword: "12345678",
		})
		assert.EqualError(t, err, "email already registered")
		assert.Nil(t, resp)
	})

	t.Run("should return error if password and confirm password does not match", func(t *testing.T) {
		authRepo.EXPECT().GetByEmail(ctx, "test@email.com", auth.Types.Credentials).Return(nil, nil).Once()

		resp, err := signUpWithCredentials(ctx, usecase.SignUpWithCredentialsParams{
			Name:                 "test",
			Email:                "test@email.com",
			Password:             "12345678",
			ConfirmationPassword: "12345679",
		})
		assert.EqualError(t, err, "passwords do not match")
		assert.Nil(t, resp)
	})

	t.Run("should return error if userRepo fails", func(t *testing.T) {
		authRepo.EXPECT().GetByEmail(ctx, "test@email.com", auth.Types.Credentials).Return(nil, nil).Once()
		userRepo.EXPECT().GetByEmail(ctx, "test@email.com").Return(nil, errors.New("test error")).Once()
		resp, err := signUpWithCredentials(ctx, usecase.SignUpWithCredentialsParams{
			Name:                 "test",
			Email:                "test@email.com",
			Password:             "12345678",
			ConfirmationPassword: "12345678",
		})
		assert.EqualError(t, err, "userRepo.GetByEmail: test error")
		assert.Nil(t, resp)
	})

	t.Run("should return error if authRepo fails", func(t *testing.T) {
		authRepo.EXPECT().GetByEmail(ctx, "test@email.com", auth.Types.Credentials).Return(nil, nil).Once()
		userRepo.EXPECT().GetByEmail(ctx, usr.Email).Return(usr, nil).Once()
		authRepo.EXPECT().GetNextID().Return(auth.ID{Value: 1}).Once()
		authRepo.EXPECT().Store(ctx, mock.Anything).Return(errors.New("test error")).Once()

		resp, err := signUpWithCredentials(ctx, usecase.SignUpWithCredentialsParams{
			Name:                 "test",
			Email:                "test@email.com",
			Password:             "12345678",
			ConfirmationPassword: "12345678",
		})
		assert.EqualError(t, err, "authRepo.Store: test error")
		assert.Nil(t, resp)
	})

	t.Run("should return error if tokenProvider fails", func(t *testing.T) {
		authRepo.EXPECT().GetByEmail(ctx, "test@email.com", auth.Types.Credentials).Return(nil, nil).Once()
		userRepo.EXPECT().GetByEmail(ctx, usr.Email).Return(usr, nil).Once()
		authRepo.EXPECT().GetNextID().Return(auth.ID{Value: 1}).Once()
		authRepo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()
		tokenProvider.EXPECT().GenerateUserTokens(*usr).Return("", "", errors.New("test error")).Once()

		resp, err := signUpWithCredentials(ctx, usecase.SignUpWithCredentialsParams{
			Name:                 "test",
			Email:                "test@email.com",
			Password:             "12345678",
			ConfirmationPassword: "12345678",
		})
		assert.EqualError(t, err, "tokenProvider.GenerateUserTokens: test error")
		assert.Nil(t, resp)
	})

	t.Run("happy path", func(t *testing.T) {
		authRepo.EXPECT().GetByEmail(ctx, "test@email.com", auth.Types.Credentials).Return(nil, nil).Once()
		userRepo.EXPECT().GetByEmail(ctx, usr.Email).Return(usr, nil).Once()
		authRepo.EXPECT().GetNextID().Return(auth.ID{Value: 1}).Once()
		authRepo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()
		tokenProvider.EXPECT().GenerateUserTokens(*usr).Return("token", "refresh_token", nil).Once()

		resp, err := signUpWithCredentials(ctx, usecase.SignUpWithCredentialsParams{
			Name:                 "test",
			Email:                "test@email.com",
			Password:             "12345678",
			ConfirmationPassword: "12345678",
		})
		assert.Nil(t, err)
		assert.Equal(t, resp.Token, "token")
		assert.Equal(t, resp.RefreshToken, "refresh_token")
		assert.Equal(t, resp.User.Name, usr.Name)
	})

}
