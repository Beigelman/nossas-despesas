package service

import (
	"context"

	"google.golang.org/api/idtoken"
)

type GoogleTokenClaims struct {
	Email   string
	Name    string
	Sub     string
	Picture *string
}

type GoogleTokenValidator interface {
	ValidateToken(ctx context.Context, token string) (*GoogleTokenClaims, error)
}

type GoogleTokenValidatorImpl struct{}

func NewGoogleTokenValidator() GoogleTokenValidator {
	return &GoogleTokenValidatorImpl{}
}

func (v *GoogleTokenValidatorImpl) ValidateToken(ctx context.Context, token string) (*GoogleTokenClaims, error) {
	payload, err := idtoken.Validate(ctx, token, "")
	if err != nil {
		return nil, err
	}

	email, ok := payload.Claims["email"].(string)
	if !ok {
		return nil, err
	}

	name, ok := payload.Claims["name"].(string)
	if !ok {
		return nil, err
	}

	sub, ok := payload.Claims["sub"].(string)
	if !ok {
		return nil, err
	}

	var picture *string
	if pic, ok := payload.Claims["picture"].(string); ok {
		picture = &pic
	}

	return &GoogleTokenClaims{
		Email:   email,
		Name:    name,
		Sub:     sub,
		Picture: picture,
	}, nil
}
