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
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
)

func TestSignUpWithCredentialsHandler(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name         string
		body         any
		usecase      usecase.SignUpWithCredentials
		expectedCode int
		assertBody   func(t *testing.T, resp *http.Response)
	}{
		{
			name: "success",
			body: controller.SignUpWithCredentialsRequest{Name: "John", Email: "john@example.com", Password: "secret123", ConfirmPassword: "secret123"},
			usecase: func(ctx context.Context, p usecase.SignUpWithCredentialsParams) (*usecase.SignUpWithCredentialsResponse, error) {
				gid := group.ID{Value: 1}
				usr := user.New(user.Attributes{ID: user.ID{Value: 2}, Name: p.Name, Email: p.Email, GroupID: &gid})
				return &usecase.SignUpWithCredentialsResponse{User: usr, Token: "token", RefreshToken: "refresh"}, nil
			},
			expectedCode: fiber.StatusCreated,
			assertBody: func(t *testing.T, resp *http.Response) {
				var res api.Response[controller.UserLogIn]
				assert.NoError(t, json.NewDecoder(resp.Body).Decode(&res))
				assert.Equal(t, 2, res.Data.User.ID)
				assert.Equal(t, "John", res.Data.User.Name)
				assert.Equal(t, "token", res.Data.Token)
				assert.Equal(t, "refresh", res.Data.RefreshToken)
			},
		},
		{
			name: "validation error",
			body: map[string]any{"name": "John"},
			usecase: func(ctx context.Context, p usecase.SignUpWithCredentialsParams) (*usecase.SignUpWithCredentialsResponse, error) {
				return nil, nil
			},
			expectedCode: fiber.StatusBadRequest,
			assertBody: func(t *testing.T, resp *http.Response) {
				var errRes api.ErrorResponse
				assert.NoError(t, json.NewDecoder(resp.Body).Decode(&errRes))
				assert.Equal(t, fiber.StatusBadRequest, errRes.StatusCode)
			},
		},
		{
			name: "body parser error",
			body: "invalid",
			usecase: func(ctx context.Context, p usecase.SignUpWithCredentialsParams) (*usecase.SignUpWithCredentialsResponse, error) {
				return nil, nil
			},
			expectedCode: fiber.StatusInternalServerError,
			assertBody: func(t *testing.T, resp *http.Response) {
				var errRes api.ErrorResponse
				assert.NoError(t, json.NewDecoder(resp.Body).Decode(&errRes))
				assert.Equal(t, fiber.StatusInternalServerError, errRes.StatusCode)
			},
		},
		{
			name: "usecase error",
			body: controller.SignUpWithCredentialsRequest{Name: "John", Email: "john@example.com", Password: "secret123", ConfirmPassword: "secret123"},
			usecase: func(ctx context.Context, p usecase.SignUpWithCredentialsParams) (*usecase.SignUpWithCredentialsResponse, error) {
				return nil, except.ConflictError()
			},
			expectedCode: fiber.StatusConflict,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			app := fiber.New(fiber.Config{ErrorHandler: api.ErrorHandler})
			handler := controller.NewSignUpWithCredentials(tt.usecase)
			app.Post("/signup", handler)

			var body []byte
			switch v := tt.body.(type) {
			case string:
				body = []byte(v)
			default:
				body, _ = json.Marshal(v)
			}

			req := httptest.NewRequest("POST", "/signup", bytes.NewBuffer(body))
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
