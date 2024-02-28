package entity

import (
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthType string

var AuthTypes = struct {
	Credentials AuthType
	Google      AuthType
}{
	Credentials: "credentials",
	Google:      "google",
}

type AuthID struct{ Value int }

type Auth struct {
	ddd.Entity[AuthID]
	Email      string
	Password   *string
	ProviderID *string
	Type       AuthType
}

type CredentialsAuthParams struct {
	ID       AuthID
	Email    string
	Password string
}

func NewCredentialAuth(params CredentialsAuthParams) (*Auth, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("bcrypt.GenerateFromPassword: %w", err)
	}

	sHashPassword := string(hashPassword)

	return &Auth{
		Entity: ddd.Entity[AuthID]{
			ID:        params.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Version:   0,
		},
		Email:    params.Email,
		Password: &sHashPassword,
		Type:     AuthTypes.Credentials,
	}, nil
}

type GoogleAuthParams struct {
	ID         AuthID
	Email      string
	ProviderID string
}

func NewGoogleAuth(params GoogleAuthParams) *Auth {
	return &Auth{
		Entity: ddd.Entity[AuthID]{
			ID:        params.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Version:   0,
		},
		Email:      params.Email,
		ProviderID: &params.ProviderID,
		Type:       AuthTypes.Google,
	}
}

func (a *Auth) CheckPassword(password string) bool {
	if a.Password == nil {
		return true
	}

	err := bcrypt.CompareHashAndPassword([]byte(*a.Password), []byte(password))
	return err == nil
}

type Claims struct {
	UserID  int
	GroupID *int
	Email   string
}

type Token struct {
	Raw     string
	Claims  Claims
	IsValid bool
}
