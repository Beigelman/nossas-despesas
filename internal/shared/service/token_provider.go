package service

import (
	"github.com/Beigelman/nossas-despesas/internal/modules/auth"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
)

type TokenProvider interface {
	GenerateUserTokens(user user.User) (string, string, error)
	ParseToken(token string) (*auth.Token, error)
	ParseRefreshToken(token string) (*auth.Token, error)
}
