package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
	"golang.org/x/crypto/bcrypt"
)

type Type string

var Types = struct {
	Credentials Type
	Google      Type
	MagicLink   Type
}{
	Credentials: "credentials",
	Google:      "google",
	MagicLink:   "magic_link",
}

type ID struct{ Value int }

type Auth struct {
	ddd.Entity[ID]
	Email      string
	Password   *string
	ProviderID *string
	Type       Type
}

type CredentialsAttributes struct {
	ID       ID
	Email    string
	Password string
}

func NewCredentialAuth(attr CredentialsAttributes) (*Auth, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(attr.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("bcrypt.GenerateFromPassword: %w", err)
	}

	sHashPassword := string(hashPassword)

	return &Auth{
		Entity: ddd.Entity[ID]{
			ID:        attr.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Version:   0,
		},
		Email:    attr.Email,
		Password: &sHashPassword,
		Type:     Types.Credentials,
	}, nil
}

type GoogleAuthAttributes struct {
	ID         ID
	Email      string
	ProviderID string
}

func NewGoogleAuth(attr GoogleAuthAttributes) *Auth {
	return &Auth{
		Entity: ddd.Entity[ID]{
			ID:        attr.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Version:   0,
		},
		Email:      attr.Email,
		ProviderID: &attr.ProviderID,
		Type:       Types.Google,
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

type Repository interface {
	ddd.Repository[ID, Auth]
	GetByEmail(ctx context.Context, email string, authType Type) (*Auth, error)
}
