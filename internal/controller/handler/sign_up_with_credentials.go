package handler

import (
	"fmt"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"github.com/Beigelman/ludaapi/internal/pkg/api"
	"github.com/Beigelman/ludaapi/internal/pkg/except"
	"github.com/Beigelman/ludaapi/internal/pkg/validator"
	"github.com/Beigelman/ludaapi/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type (
	SignUpWithCredentialsRequest struct {
		Name            string  `json:"name" validate:"required"`
		Email           string  `json:"email" validate:"required,email"`
		Password        string  `json:"password" validate:"required,min=8"`
		ConfirmPassword string  `json:"confirm_password" validate:"required,min=8"`
		ProfilePicture  *string `json:"profile_picture"`
		GroupID         *int    `json:"group_id"`
	}

	SignUpWithCredentials func(ctx *fiber.Ctx) error
)

func NewSignUpWithCredentials(signUpWithCredentials usecase.SignUpWithCredentials) SignUpWithCredentials {
	valid := validator.New()
	return func(ctx *fiber.Ctx) error {
		var req SignUpWithCredentialsRequest
		if err := ctx.BodyParser(&req); err != nil {
			return fmt.Errorf("ctx.BodyParser: %w", err)
		}

		if err := valid.Validate(req); err != nil {
			return except.BadRequestError("invalid request body").SetInternal(err)
		}

		result, err := signUpWithCredentials(ctx.Context(), usecase.SignUpWithCredentialsParams{
			Name:                 req.Name,
			Email:                req.Email,
			Password:             req.Password,
			ConfirmationPassword: req.ConfirmPassword,
			ProfilePicture:       req.ProfilePicture,
			GroupID: func() *entity.GroupID {
				if req.GroupID == nil {
					return nil
				}
				return &entity.GroupID{Value: *req.GroupID}
			}(),
		})
		if err != nil {
			return fmt.Errorf("signInWithCredentials: %w", err)
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
