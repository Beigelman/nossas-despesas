package controller

import (
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/modules/auth/usecase"
	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/Beigelman/nossas-despesas/internal/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type (
	RefreshAuthTokenRequest struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	RefreshAuthToken func(ctx *fiber.Ctx) error
)

func NewRefreshAuthToken(refreshAuthToken usecase.RefreshAuthToken) RefreshAuthToken {
	valid := validator.New()
	return func(ctx *fiber.Ctx) error {
		var req RefreshAuthTokenRequest
		if err := ctx.BodyParser(&req); err != nil {
			return fmt.Errorf("ctx.BodyParser: %w", err)
		}

		if err := valid.Validate(req); err != nil {
			return except.BadRequestError("invalid request body").SetInternal(err)
		}

		result, err := refreshAuthToken(ctx.Context(), usecase.RefreshAuthTokenParams{
			RefreshToken: req.RefreshToken,
		})
		if err != nil {
			return fmt.Errorf("refreshAuthToken: %w", err)
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
