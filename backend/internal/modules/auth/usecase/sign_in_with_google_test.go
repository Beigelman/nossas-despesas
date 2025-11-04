package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Beigelman/nossas-despesas/internal/modules/auth"
	"github.com/Beigelman/nossas-despesas/internal/modules/auth/usecase"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"github.com/Beigelman/nossas-despesas/internal/shared/mocks"
	"github.com/Beigelman/nossas-despesas/internal/shared/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignInWithGoogle(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userRepo := mocks.NewMockuserRepository(t)
	authRepo := mocks.NewMockauthRepository(t)
	tokenProvider := mocks.NewMockserviceTokenProvider(t)
	googleValidator := mocks.NewMockserviceGoogleTokenValidator(t)

	usr := user.New(user.Attributes{
		ID:    user.ID{Value: 1},
		Name:  "Test User",
		Email: "test@email.com",
	})

	claims := &service.GoogleTokenClaims{
		Email:   "test@email.com",
		Name:    "Test User",
		Sub:     "google-sub-123",
		Picture: func() *string { s := "https://example.com/pic.jpg"; return &s }(),
	}

	signInWithGoogle := usecase.NewSignInWithGoogle(userRepo, authRepo, tokenProvider, googleValidator)

	t.Run("googleValidator.ValidateToken returns error", func(t *testing.T) {
		googleValidator.EXPECT().ValidateToken(ctx, "invalid-token").Return(nil, errors.New("invalid token")).Once()

		resp, err := signInWithGoogle(ctx, usecase.SignInWithGoogleParams{IdToken: "invalid-token"})
		assert.ErrorContains(t, err, "googleValidator.ValidateToken: invalid token")
		assert.Nil(t, resp)
	})

	t.Run("email not found in token", func(t *testing.T) {
		invalidClaims := &service.GoogleTokenClaims{
			Email: "",
			Name:  "Test User",
			Sub:   "google-sub-123",
		}
		googleValidator.EXPECT().ValidateToken(ctx, "token").Return(invalidClaims, nil).Once()

		resp, err := signInWithGoogle(ctx, usecase.SignInWithGoogleParams{IdToken: "token"})
		assert.ErrorContains(t, err, "email not found in token")
		assert.Nil(t, resp)
	})

	t.Run("user name not found in token", func(t *testing.T) {
		invalidClaims := &service.GoogleTokenClaims{
			Email: "test@email.com",
			Name:  "",
			Sub:   "google-sub-123",
		}
		googleValidator.EXPECT().ValidateToken(ctx, "token").Return(invalidClaims, nil).Once()

		resp, err := signInWithGoogle(ctx, usecase.SignInWithGoogleParams{IdToken: "token"})
		assert.ErrorContains(t, err, "user name not found in token")
		assert.Nil(t, resp)
	})

	t.Run("sub not found in token", func(t *testing.T) {
		invalidClaims := &service.GoogleTokenClaims{
			Email: "test@email.com",
			Name:  "Test User",
			Sub:   "",
		}
		googleValidator.EXPECT().ValidateToken(ctx, "token").Return(invalidClaims, nil).Once()

		resp, err := signInWithGoogle(ctx, usecase.SignInWithGoogleParams{IdToken: "token"})
		assert.ErrorContains(t, err, "sub not found in token")
		assert.Nil(t, resp)
	})

	t.Run("userRepo.GetByEmail returns error", func(t *testing.T) {
		googleValidator.EXPECT().ValidateToken(ctx, "token").Return(claims, nil).Once()
		userRepo.EXPECT().GetByEmail(ctx, "test@email.com").Return(nil, errors.New("db error")).Once()

		resp, err := signInWithGoogle(ctx, usecase.SignInWithGoogleParams{IdToken: "token"})
		assert.ErrorContains(t, err, "userRepo.GetByEmail: db error")
		assert.Nil(t, resp)
	})

	t.Run("userRepo.Store returns error for existing user with profile update", func(t *testing.T) {
		existingUser := user.New(user.Attributes{
			ID:             user.ID{Value: 1},
			Name:           "Test User",
			Email:          "test@email.com",
			ProfilePicture: nil, // No profile picture
		})

		googleValidator.EXPECT().ValidateToken(ctx, "token").Return(claims, nil).Once()
		userRepo.EXPECT().GetByEmail(ctx, "test@email.com").Return(existingUser, nil).Once()
		userRepo.EXPECT().Store(ctx, mock.Anything).Return(errors.New("store error")).Once()

		resp, err := signInWithGoogle(ctx, usecase.SignInWithGoogleParams{IdToken: "token"})
		assert.ErrorContains(t, err, "userRepo.Store: store error")
		assert.Nil(t, resp)
	})

	t.Run("userRepo.Store returns error for new user", func(t *testing.T) {
		googleValidator.EXPECT().ValidateToken(ctx, "token").Return(claims, nil).Once()
		userRepo.EXPECT().GetByEmail(ctx, "test@email.com").Return(nil, nil).Once()
		userRepo.EXPECT().GetNextID().Return(user.ID{Value: 1}).Once()
		userRepo.EXPECT().Store(ctx, mock.Anything).Return(errors.New("store error")).Once()

		resp, err := signInWithGoogle(ctx, usecase.SignInWithGoogleParams{IdToken: "token"})
		assert.ErrorContains(t, err, "userRepo.Store: store error")
		assert.Nil(t, resp)
	})

	t.Run("authRepo.GetByEmail returns error", func(t *testing.T) {
		userWithoutPic := user.New(user.Attributes{
			ID:             user.ID{Value: 1},
			Name:           "Test User",
			Email:          "test@email.com",
			ProfilePicture: nil,
		})

		googleValidator.EXPECT().ValidateToken(ctx, "token").Return(claims, nil).Once()
		userRepo.EXPECT().GetByEmail(ctx, "test@email.com").Return(userWithoutPic, nil).Once()
		userRepo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once() // Profile picture update
		authRepo.EXPECT().GetByEmail(ctx, "test@email.com", auth.Types.Google).Return(nil, errors.New("auth error")).Once()

		resp, err := signInWithGoogle(ctx, usecase.SignInWithGoogleParams{IdToken: "token"})
		assert.ErrorContains(t, err, "authRepo.GetByEmail: auth error")
		assert.Nil(t, resp)
	})

	t.Run("authRepo.Store returns error", func(t *testing.T) {
		userWithPic := user.New(user.Attributes{
			ID:             user.ID{Value: 1},
			Name:           "Test User",
			Email:          "test@email.com",
			ProfilePicture: func() *string { s := "existing-pic.jpg"; return &s }(),
		})

		googleValidator.EXPECT().ValidateToken(ctx, "token").Return(claims, nil).Once()
		userRepo.EXPECT().GetByEmail(ctx, "test@email.com").Return(userWithPic, nil).Once()
		authRepo.EXPECT().GetByEmail(ctx, "test@email.com", auth.Types.Google).Return(nil, nil).Once()
		authRepo.EXPECT().GetNextID().Return(auth.ID{Value: 1}).Once()
		authRepo.EXPECT().Store(ctx, mock.Anything).Return(errors.New("auth store error")).Once()

		resp, err := signInWithGoogle(ctx, usecase.SignInWithGoogleParams{IdToken: "token"})
		assert.ErrorContains(t, err, "authRepo.Store: auth store error")
		assert.Nil(t, resp)
	})

	t.Run("tokenProvider.GenerateUserTokens returns error", func(t *testing.T) {
		userWithPic := user.New(user.Attributes{
			ID:             user.ID{Value: 1},
			Name:           "Test User",
			Email:          "test@email.com",
			ProfilePicture: func() *string { s := "existing-pic.jpg"; return &s }(),
		})

		googleValidator.EXPECT().ValidateToken(ctx, "token").Return(claims, nil).Once()
		userRepo.EXPECT().GetByEmail(ctx, "test@email.com").Return(userWithPic, nil).Once()
		authRepo.EXPECT().GetByEmail(ctx, "test@email.com", auth.Types.Google).Return(nil, nil).Once()
		authRepo.EXPECT().GetNextID().Return(auth.ID{Value: 1}).Once()
		authRepo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()
		tokenProvider.EXPECT().GenerateUserTokens(*userWithPic).Return("", "", errors.New("token error")).Once()

		resp, err := signInWithGoogle(ctx, usecase.SignInWithGoogleParams{IdToken: "token"})
		assert.ErrorContains(t, err, "tokenProvider.GenerateUserTokens: token error")
		assert.Nil(t, resp)
	})

	t.Run("success - existing user with profile picture update", func(t *testing.T) {
		existingUser := user.New(user.Attributes{
			ID:             user.ID{Value: 1},
			Name:           "Test User",
			Email:          "test@email.com",
			ProfilePicture: nil, // No profile picture
		})

		googleValidator.EXPECT().ValidateToken(ctx, "token").Return(claims, nil).Once()
		userRepo.EXPECT().GetByEmail(ctx, "test@email.com").Return(existingUser, nil).Once()
		userRepo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()
		authRepo.EXPECT().GetByEmail(ctx, "test@email.com", auth.Types.Google).Return(nil, nil).Once()
		authRepo.EXPECT().GetNextID().Return(auth.ID{Value: 1}).Once()
		authRepo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()
		tokenProvider.EXPECT().GenerateUserTokens(mock.Anything).Return("auth-token", "refresh-token", nil).Once()

		resp, err := signInWithGoogle(ctx, usecase.SignInWithGoogleParams{IdToken: "token"})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "auth-token", resp.Token)
		assert.Equal(t, "refresh-token", resp.RefreshToken)
		assert.Equal(t, "Test User", resp.User.Name)
		assert.Equal(t, "test@email.com", resp.User.Email)
	})

	t.Run("success - existing user without profile picture update", func(t *testing.T) {
		existingUserWithPic := user.New(user.Attributes{
			ID:             user.ID{Value: 1},
			Name:           "Test User",
			Email:          "test@email.com",
			ProfilePicture: func() *string { s := "existing-pic.jpg"; return &s }(),
		})

		googleValidator.EXPECT().ValidateToken(ctx, "token").Return(claims, nil).Once()
		userRepo.EXPECT().GetByEmail(ctx, "test@email.com").Return(existingUserWithPic, nil).Once()
		authRepo.EXPECT().GetByEmail(ctx, "test@email.com", auth.Types.Google).Return(nil, nil).Once()
		authRepo.EXPECT().GetNextID().Return(auth.ID{Value: 1}).Once()
		authRepo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()
		tokenProvider.EXPECT().GenerateUserTokens(mock.Anything).Return("auth-token", "refresh-token", nil).Once()

		resp, err := signInWithGoogle(ctx, usecase.SignInWithGoogleParams{IdToken: "token"})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "auth-token", resp.Token)
		assert.Equal(t, "refresh-token", resp.RefreshToken)
	})

	t.Run("success - new user", func(t *testing.T) {
		googleValidator.EXPECT().ValidateToken(ctx, "token").Return(claims, nil).Once()
		userRepo.EXPECT().GetByEmail(ctx, "test@email.com").Return(nil, nil).Once()
		userRepo.EXPECT().GetNextID().Return(user.ID{Value: 1}).Once()
		userRepo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()
		authRepo.EXPECT().GetByEmail(ctx, "test@email.com", auth.Types.Google).Return(nil, nil).Once()
		authRepo.EXPECT().GetNextID().Return(auth.ID{Value: 1}).Once()
		authRepo.EXPECT().Store(ctx, mock.Anything).Return(nil).Once()
		tokenProvider.EXPECT().GenerateUserTokens(mock.Anything).Return("auth-token", "refresh-token", nil).Once()

		resp, err := signInWithGoogle(ctx, usecase.SignInWithGoogleParams{IdToken: "token"})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "auth-token", resp.Token)
		assert.Equal(t, "refresh-token", resp.RefreshToken)
		assert.Equal(t, "Test User", resp.User.Name)
		assert.Equal(t, "test@email.com", resp.User.Email)
	})

	t.Run("success - existing auth", func(t *testing.T) {
		existingAuth := auth.NewGoogleAuth(auth.GoogleAuthAttributes{
			ID:         auth.ID{Value: 1},
			Email:      "test@email.com",
			ProviderID: "google-sub-123",
		})

		googleValidator.EXPECT().ValidateToken(ctx, "token").Return(claims, nil).Once()
		userRepo.EXPECT().GetByEmail(ctx, "test@email.com").Return(usr, nil).Once()
		authRepo.EXPECT().GetByEmail(ctx, "test@email.com", auth.Types.Google).Return(existingAuth, nil).Once()
		tokenProvider.EXPECT().GenerateUserTokens(mock.Anything).Return("auth-token", "refresh-token", nil).Once()
		userRepo.EXPECT().Store(ctx, mock.Anything).Return(nil).Maybe()

		resp, err := signInWithGoogle(ctx, usecase.SignInWithGoogleParams{IdToken: "token"})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "auth-token", resp.Token)
		assert.Equal(t, "refresh-token", resp.RefreshToken)
	})
}
