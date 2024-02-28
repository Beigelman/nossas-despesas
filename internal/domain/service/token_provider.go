package service

import "github.com/Beigelman/nossas-despesas/internal/domain/entity"

type TokenProvider interface {
	GenerateUserTokens(user entity.User) (string, string, error)
	ParseToken(token string) (*entity.Token, error)
	ParseRefreshToken(token string) (*entity.Token, error)
}
