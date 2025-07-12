package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Beigelman/nossas-despesas/internal/modules/auth/controller"
	"github.com/Beigelman/nossas-despesas/internal/modules/auth/usecase"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestSignInWithGoogleHandler(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name         string
		body         any
		usecase      usecase.SignInWithGoogle
		expectedCode int
		assertBody   func(t *testing.T, resp *http.Response)
	}{
		{
			name: "success",
			body: controller.SignInWithGoogleRequest{Token: "idtoken"},
			usecase: func(ctx context.Context, p usecase.SignInWithGoogleParams) (*usecase.SignInWithGoogleResponse, error) {
				usr := user.New(user.Attributes{ID: user.ID{Value: 3}, Name: "John", Email: "john@example.com"})
				return &usecase.SignInWithGoogleResponse{User: usr, Token: "token", RefreshToken: "refresh"}, nil
			},
			expectedCode: fiber.StatusCreated,
			assertBody: func(t *testing.T, resp *http.Response) {
				var res api.Response[controller.UserLogIn]
				assert.NoError(t, json.NewDecoder(resp.Body).Decode(&res))
				assert.Equal(t, 3, res.Data.User.ID)
				assert.Equal(t, "token", res.Data.Token)
			},
		},
		{
			name: "validation error",
			body: map[string]string{},
			usecase: func(ctx context.Context, p usecase.SignInWithGoogleParams) (*usecase.SignInWithGoogleResponse, error) {
				return nil, nil
			},
			expectedCode: fiber.StatusBadRequest,
		},
		{
			name: "usecase error",
			body: controller.SignInWithGoogleRequest{Token: "idtoken"},
			usecase: func(ctx context.Context, p usecase.SignInWithGoogleParams) (*usecase.SignInWithGoogleResponse, error) {
				return nil, except.ForbiddenError()
			},
			expectedCode: fiber.StatusForbidden,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			app := fiber.New(fiber.Config{ErrorHandler: api.ErrorHandler})
			handler := controller.NewSignInWithGoogle(tt.usecase)
			app.Post("/google", handler)

			var body []byte
			switch v := tt.body.(type) {
			case string:
				body = []byte(v)
			default:
				body, _ = json.Marshal(v)
			}

			req := httptest.NewRequest("POST", "/google", bytes.NewBuffer(body))
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
