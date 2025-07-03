package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Beigelman/nossas-despesas/internal/modules/auth"
	"github.com/Beigelman/nossas-despesas/internal/modules/auth/usecase"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	mocks2 "github.com/Beigelman/nossas-despesas/internal/shared/mocks"
	"github.com/stretchr/testify/assert"
)

func TestRefreshToken(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	userRepo := mocks2.NewMockuserRepository(t)
	tokenProvider := mocks2.NewMockserviceTokenProvider(t)

	validToken := auth.Token{
		Raw: "validToken",
		Claims: auth.Claims{
			UserID:  0,
			GroupID: nil,
			Email:   "test@email",
		},
		IsValid: true,
	}

	usr := user.New(user.Attributes{
		ID:             user.ID{Value: 0},
		Name:           "test",
		Email:          "test@gmail.com",
		ProfilePicture: nil,
		GroupID:        nil,
	})

	refreshToken := usecase.NewRefreshAuthToken(userRepo, tokenProvider)

	t.Run("should return error if token is invalid", func(t *testing.T) {
		tokenProvider.EXPECT().ParseRefreshToken("invalidToken").Return(nil, errors.New("invalid token")).Once()

		resp, err := refreshToken(ctx, usecase.RefreshAuthTokenParams{RefreshToken: "invalidToken"})
		assert.Errorf(t, err, "tokenProvider.ParseRefreshToken(): invalid token")
		assert.Nil(t, resp)
	})

	t.Run("should return error if userRepo fails", func(t *testing.T) {
		tokenProvider.EXPECT().ParseRefreshToken("validToken").Return(&validToken, nil).Once()
		userRepo.EXPECT().GetByID(ctx, user.ID{Value: 0}).Return(nil, errors.New("test error")).Once()

		resp, err := refreshToken(ctx, usecase.RefreshAuthTokenParams{RefreshToken: "validToken"})
		assert.Errorf(t, err, "tokenProvider.ParseRefreshToken(): invalid token")
		assert.Nil(t, resp)
	})

	t.Run("should return error if user is not found fails", func(t *testing.T) {
		tokenProvider.EXPECT().ParseRefreshToken("validToken").Return(&validToken, nil).Once()
		userRepo.EXPECT().GetByID(ctx, user.ID{Value: 0}).Return(nil, nil).Once()

		resp, err := refreshToken(ctx, usecase.RefreshAuthTokenParams{RefreshToken: "validToken"})
		assert.Errorf(t, err, "user not found")
		assert.Nil(t, resp)
	})

	t.Run("should return error if new generated token fails", func(t *testing.T) {
		tokenProvider.EXPECT().ParseRefreshToken("validToken").Return(&validToken, nil).Once()
		userRepo.EXPECT().GetByID(ctx, user.ID{Value: 0}).Return(usr, nil).Once()
		tokenProvider.EXPECT().GenerateUserTokens(*usr).Return("", "", errors.New("test error")).Once()

		resp, err := refreshToken(ctx, usecase.RefreshAuthTokenParams{RefreshToken: "validToken"})
		assert.Errorf(t, err, "tokenProvider.GenerateUserTokens: test error")
		assert.Nil(t, resp)
	})

	t.Run("happy path", func(t *testing.T) {
		tokenProvider.EXPECT().ParseRefreshToken("validToken").Return(&validToken, nil).Once()
		userRepo.EXPECT().GetByID(ctx, user.ID{Value: 0}).Return(usr, nil).Once()
		tokenProvider.EXPECT().GenerateUserTokens(*usr).Return("new_token", "new_refresh_token", nil).Once()

		resp, err := refreshToken(ctx, usecase.RefreshAuthTokenParams{RefreshToken: "validToken"})
		assert.Nil(t, err)
		assert.Equal(t, resp.Token, "new_token")
		assert.Equal(t, resp.RefreshToken, "new_refresh_token")
		assert.Equal(t, resp.User.Name, usr.Name)
	})
}
