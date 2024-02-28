package handler

import (
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/Beigelman/nossas-despesas/internal/pkg/validator"
	"net/http"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type (
	SignInWithCredentialsRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}

	UserResponse struct {
		ID             int       `json:"id"`
		Name           string    `json:"name"`
		Email          string    `json:"email"`
		ProfilePicture *string   `json:"profile_picture,omitempty"`
		GroupID        *int      `json:"group_id,omitempty"`
		CreatedAt      time.Time `json:"created_at"`
		UpdatedAt      time.Time `json:"updated_at"`
	}

	UserLogIn struct {
		User         UserResponse `json:"user"`
		Token        string       `json:"token"`
		RefreshToken string       `json:"refresh_token"`
	}

	SignInWithCredentials func(ctx *fiber.Ctx) error
)

func NewSignInWithCredentials(signUpWithCredentials usecase.SignInWithCredentials) SignInWithCredentials {
	valid := validator.New()
	return func(ctx *fiber.Ctx) error {
		var req SignInWithCredentialsRequest
		if err := ctx.BodyParser(&req); err != nil {
			return fmt.Errorf("ctx.BodyParser: %w", err)
		}

		if err := valid.Validate(req); err != nil {
			return except.BadRequestError("invalid request body").SetInternal(err)
		}

		result, err := signUpWithCredentials(ctx.Context(), usecase.SignInWithCredentialsParams{
			Email:    req.Email,
			Password: req.Password,
		})
		if err != nil {
			return fmt.Errorf("signUpWithCredentials: %w", err)
		}

		return ctx.Status(http.StatusCreated).JSON(
			api.NewResponse(http.StatusCreated, UserLogIn{
				User: UserResponse{
					ID:             result.User.ID.Value,
					Name:           result.User.Name,
					Email:          result.User.Email,
					ProfilePicture: result.User.ProfilePicture,
					GroupID: func() *int {
						if result.User.GroupID == nil {
							return nil
						}
						return &result.User.GroupID.Value
					}(),
					CreatedAt: result.User.CreatedAt,
					UpdatedAt: result.User.UpdatedAt,
				},
				Token:        result.Token,
				RefreshToken: result.RefreshToken,
			}),
		)
	}
}
