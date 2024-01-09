package token

import (
	"fmt"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JWTProvider struct {
	secret string
}

func NewJWTProvider(secret string) *JWTProvider {
	return &JWTProvider{secret: secret}
}

func (p *JWTProvider) GenerateUserTokens(user entity.User) (string, string, error) {
	tokenClaims := jwt.MapClaims{
		"user_id":  user.ID.Value,
		"group_id": user.GroupID,
		"email":    user.Email,
		"exp":      jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		"iat":      jwt.NewNumericDate(time.Now()),
	}

	refreshTokenClaims := jwt.MapClaims{
		"user_id": user.ID.Value,
		"exp":     jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		"iat":     jwt.NewNumericDate(time.Now()),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims).SignedString([]byte(p.secret))
	if err != nil {
		return "", "", fmt.Errorf("new token: %w", err)
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims).SignedString([]byte(p.secret))
	if err != nil {
		return "", "", fmt.Errorf("new refresh token: %w", err)
	}

	return token, refreshToken, nil
}

func (p *JWTProvider) ParseToken(stringToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(stringToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(p.secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}
