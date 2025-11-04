package controller

import (
	"fmt"
	"net/http"

	"github.com/Beigelman/nossas-despesas/internal/modules/auth/usecase"
	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/Beigelman/nossas-despesas/internal/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

type (
	SignInWithGoogleRequest struct {
		Token string `token:"email" validate:"required"`
	}

	SignInWithGoogle func(ctx *fiber.Ctx) error
)

func NewSignInWithGoogle(signInWithGoogle usecase.SignInWithGoogle) SignInWithGoogle {
	valid := validator.New()
	return func(ctx *fiber.Ctx) error {
		var req SignInWithGoogleRequest
		if err := ctx.BodyParser(&req); err != nil {
			return fmt.Errorf("ctx.BodyParser: %w", err)
		}

		if err := valid.Validate(req); err != nil {
			return except.BadRequestError("invalid request body").SetInternal(err)
		}

		result, err := signInWithGoogle(ctx.Context(), usecase.SignInWithGoogleParams{
			IdToken: req.Token,
		})
		if err != nil {
			return fmt.Errorf("signInWithGoogle: %w", err)
		}

		var groupID *int
		if result.User.GroupID != nil {
			groupID = &result.User.GroupID.Value
		}
		return ctx.Status(http.StatusCreated).JSON(
			api.NewResponse(http.StatusCreated, UserLogIn{
				User: UserResponse{
					ID:             result.User.ID.Value,
					Name:           result.User.Name,
					Email:          result.User.Email,
					ProfilePicture: result.User.ProfilePicture,
					GroupID:        groupID,
					Flags:          result.User.Flags,
					CreatedAt:      result.User.CreatedAt,
					UpdatedAt:      result.User.UpdatedAt,
				},
				Token:        result.Token,
				RefreshToken: result.RefreshToken,
			}),
		)
	}
}
