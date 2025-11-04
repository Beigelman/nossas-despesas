package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"

	"github.com/Beigelman/nossas-despesas/internal/modules/auth/controller"
	"github.com/Beigelman/nossas-despesas/internal/modules/auth/usecase"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
)

func TestRefreshTokenHandler(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name         string
		body         any
		usecase      usecase.RefreshAuthToken
		expectedCode int
		assertBody   func(t *testing.T, resp *http.Response)
	}{
		{
			name: "success",
			body: controller.RefreshAuthTokenRequest{RefreshToken: "tok"},
			usecase: func(ctx context.Context, p usecase.RefreshAuthTokenParams) (*usecase.RefreshAuthTokenResponse, error) {
				usr := user.New(user.Attributes{ID: user.ID{Value: 4}, Name: "John", Email: "john@example.com"})
				return &usecase.RefreshAuthTokenResponse{User: usr, Token: "token", RefreshToken: "refresh"}, nil
			},
			expectedCode: fiber.StatusCreated,
			assertBody: func(t *testing.T, resp *http.Response) {
				var res api.Response[controller.UserLogIn]
				assert.NoError(t, json.NewDecoder(resp.Body).Decode(&res))
				assert.Equal(t, 4, res.Data.User.ID)
				assert.Equal(t, "refresh", res.Data.RefreshToken)
			},
		},
		{
			name: "validation error",
			body: map[string]string{},
			usecase: func(ctx context.Context, p usecase.RefreshAuthTokenParams) (*usecase.RefreshAuthTokenResponse, error) {
				return nil, nil
			},
			expectedCode: fiber.StatusBadRequest,
		},
		{
			name: "usecase error",
			body: controller.RefreshAuthTokenRequest{RefreshToken: "tok"},
			usecase: func(ctx context.Context, p usecase.RefreshAuthTokenParams) (*usecase.RefreshAuthTokenResponse, error) {
				return nil, except.UnauthorizedError()
			},
			expectedCode: fiber.StatusUnauthorized,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			app := fiber.New(fiber.Config{ErrorHandler: api.ErrorHandler})
			handler := controller.NewRefreshAuthToken(tt.usecase)
			app.Post("/refresh", handler)

			var body []byte
			switch v := tt.body.(type) {
			case string:
				body = []byte(v)
			default:
				body, _ = json.Marshal(v)
			}

			req := httptest.NewRequest("POST", "/refresh", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, resp.StatusCode)

			if tt.assertBody != nil {
				tt.assertBody(t, resp)
			}
			assert.NoError(t, resp.Body.Close())
		})
	}
}
