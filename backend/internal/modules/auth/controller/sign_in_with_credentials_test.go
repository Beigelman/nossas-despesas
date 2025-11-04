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

func TestSignInWithCredentialsHandler(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name         string
		body         any
		usecase      usecase.SignInWithCredentials
		expectedCode int
		assertBody   func(t *testing.T, resp *http.Response)
	}{
		{
			name: "success",
			body: controller.SignInWithCredentialsRequest{Email: "john@example.com", Password: "secret"},
			usecase: func(ctx context.Context, p usecase.SignInWithCredentialsParams) (*usecase.SignInWithCredentialsResponse, error) {
				usr := user.New(user.Attributes{ID: user.ID{Value: 1}, Name: "John", Email: p.Email})
				return &usecase.SignInWithCredentialsResponse{User: usr, Token: "token", RefreshToken: "refresh"}, nil
			},
			expectedCode: fiber.StatusCreated,
			assertBody: func(t *testing.T, resp *http.Response) {
				var res api.Response[controller.UserLogIn]
				assert.NoError(t, json.NewDecoder(resp.Body).Decode(&res))
				assert.Equal(t, 1, res.Data.User.ID)
				assert.Equal(t, "John", res.Data.User.Name)
				assert.Equal(t, "token", res.Data.Token)
				assert.Equal(t, "refresh", res.Data.RefreshToken)
			},
		},
		{
			name: "validation error",
			body: map[string]string{"email": "john@example.com"},
			usecase: func(ctx context.Context, p usecase.SignInWithCredentialsParams) (*usecase.SignInWithCredentialsResponse, error) {
				return nil, nil
			},
			expectedCode: fiber.StatusBadRequest,
			assertBody: func(t *testing.T, resp *http.Response) {
				var errRes api.ErrorResponse
				assert.NoError(t, json.NewDecoder(resp.Body).Decode(&errRes))
				assert.Equal(t, fiber.StatusBadRequest, errRes.StatusCode)
				assert.Equal(t, "invalid request body", errRes.Message)
			},
		},
		{
			name: "usecase error",
			body: controller.SignInWithCredentialsRequest{Email: "john@example.com", Password: "secret"},
			usecase: func(ctx context.Context, p usecase.SignInWithCredentialsParams) (*usecase.SignInWithCredentialsResponse, error) {
				return nil, except.UnprocessableEntityError()
			},
			expectedCode: fiber.StatusUnprocessableEntity,
			assertBody: func(t *testing.T, resp *http.Response) {
				var errRes api.ErrorResponse
				assert.NoError(t, json.NewDecoder(resp.Body).Decode(&errRes))
				assert.Equal(t, fiber.StatusUnprocessableEntity, errRes.StatusCode)
			},
		},
		{
			name: "body parser error",
			body: "invalid",
			usecase: func(ctx context.Context, p usecase.SignInWithCredentialsParams) (*usecase.SignInWithCredentialsResponse, error) {
				return nil, nil
			},
			expectedCode: fiber.StatusInternalServerError,
			assertBody: func(t *testing.T, resp *http.Response) {
				var errRes api.ErrorResponse
				assert.NoError(t, json.NewDecoder(resp.Body).Decode(&errRes))
				assert.Equal(t, fiber.StatusInternalServerError, errRes.StatusCode)
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			app := fiber.New(fiber.Config{ErrorHandler: api.ErrorHandler})
			handler := controller.NewSignInWithCredentials(tt.usecase)
			app.Post("/login", handler)

			var body []byte
			switch v := tt.body.(type) {
			case string:
				body = []byte(v)
			default:
				body, _ = json.Marshal(v)
			}

			req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
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
