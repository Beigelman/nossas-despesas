package service

import "github.com/Beigelman/ludaapi/internal/domain/entity"

type TokenProvider interface {
	GenerateUserTokens(user entity.User) (string, string, error)
	ParseToken(token string) (*entity.Token, error)
}
