package authrepo

import (
	"database/sql"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"time"

	"github.com/Beigelman/ludaapi/internal/pkg/ddd"
)

func toEntity(model AuthModel) *entity.Auth {
	var deletedAt *time.Time
	if model.DeletedAt.Valid {
		deletedAt = &model.DeletedAt.Time
	}

	var providerID *string
	if model.ProviderID.Valid {
		providerID = &model.ProviderID.String
	}

	var password *string
	if model.Password.Valid {
		password = &model.Password.String
	}

	return &entity.Auth{
		Entity: ddd.Entity[entity.AuthID]{
			ID:        entity.AuthID{Value: model.ID},
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
			DeletedAt: deletedAt,
			Version:   model.Version,
		},
		Email:      model.Email,
		Password:   password,
		ProviderID: providerID,
		Type: func() entity.AuthType {
			if model.Type == "credentials" {
				return entity.AuthTypes.Credentials
			} else {
				return entity.AuthTypes.Google
			}
		}(),
	}
}

func toModel(entity *entity.Auth) AuthModel {
	deletedAt := sql.NullTime{Time: time.Time{}, Valid: false}
	if entity.DeletedAt != nil {
		deletedAt = sql.NullTime{Time: *entity.DeletedAt, Valid: true}
	}

	providerID := sql.NullString{String: "", Valid: false}
	if entity.ProviderID != nil {
		providerID = sql.NullString{String: *entity.ProviderID, Valid: true}
	}

	password := sql.NullString{String: "", Valid: false}
	if entity.Password != nil {
		password = sql.NullString{String: *entity.Password, Valid: true}
	}

	return AuthModel{
		ID:         entity.ID.Value,
		Email:      entity.Email,
		Password:   password,
		ProviderID: providerID,
		Type: func() string {
			if entity.Type == "credentials" {
				return "credentials"
			} else {
				return "google"
			}
		}(),
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: deletedAt,
		Version:   entity.Version,
	}
}
