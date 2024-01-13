package jwt

import (
	"fmt"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"github.com/Beigelman/ludaapi/internal/domain/service"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Provider struct {
	secret string
}

func NewJWTProvider(secret string) service.TokenProvider {
	return &Provider{secret: secret}
}

func (p *Provider) GenerateUserTokens(user entity.User) (string, string, error) {
	tokenClaims := jwt.MapClaims{
		"user_id": user.ID.Value,
		"group_id": func() *int {
			if user.GroupID == nil {
				return nil
			}
			return &user.GroupID.Value
		}(),
		"email": user.Email,
		"exp":   jwt.NewNumericDate(time.Now().Add(time.Hour * 3)),
		"iat":   jwt.NewNumericDate(time.Now()),
	}

	refreshTokenClaims := jwt.MapClaims{
		"user_id": user.ID.Value,
		"exp":     jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 5)),
		"iat":     jwt.NewNumericDate(time.Now()),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims).SignedString([]byte(p.secret))
	if err != nil {
		return "", "", fmt.Errorf("new jwt: %w", err)
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims).SignedString([]byte(p.secret))
	if err != nil {
		return "", "", fmt.Errorf("new refresh jwt: %w", err)
	}

	return token, refreshToken, nil
}

func (p *Provider) ParseToken(stringToken string) (*entity.Token, error) {
	token, err := jwt.Parse(stringToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(p.secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse jwt: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid jwt")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid jwt claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid jwt claims")
	}

	var groupID *int
	groupNumber, ok := claims["group_id"].(float64)
	if !ok {
		groupID = nil
	} else {
		id := int(groupNumber)
		groupID = &id
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid jwt claims")
	}

	return &entity.Token{
		Raw:     token.Raw,
		IsValid: token.Valid,
		Claims: entity.Claims{
			UserID:  int(userID),
			GroupID: groupID,
			Email:   email,
		},
	}, nil
}

func (p *Provider) ParseRefreshToken(stringToken string) (*entity.Token, error) {
	token, err := jwt.Parse(stringToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(p.secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse jwt: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid jwt")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid jwt claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid jwt claims")
	}

	return &entity.Token{
		Raw:     token.Raw,
		IsValid: token.Valid,
		Claims: entity.Claims{
			UserID: int(userID),
		},
	}, nil
}
