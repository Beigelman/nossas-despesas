package userrepo

import (
	"database/sql"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"time"

	"github.com/Beigelman/ludaapi/internal/pkg/ddd"
)

func toEntity(model UserModel) *entity.User {
	var profilePicture *string
	if model.ProfilePicture.Valid {
		profilePicture = &model.ProfilePicture.String
	}

	var deletedAt *time.Time
	if model.DeletedAt.Valid {
		deletedAt = &model.DeletedAt.Time
	}

	return &entity.User{
		Entity: ddd.Entity[entity.UserID]{
			ID:        entity.UserID{Value: model.ID},
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
			DeletedAt: deletedAt,
			Version:   model.Version,
		},
		Name:           model.Name,
		Email:          model.Email,
		ProfilePicture: profilePicture,
	}
}

func toModel(entity *entity.User) UserModel {
	var profilePicture sql.NullString
	if entity.ProfilePicture != nil {
		profilePicture = sql.NullString{String: *entity.ProfilePicture, Valid: true}
	} else {
		profilePicture = sql.NullString{String: "", Valid: false}
	}

	var deletedAt sql.NullTime
	if entity.DeletedAt != nil {
		deletedAt = sql.NullTime{Time: *entity.DeletedAt, Valid: true}
	} else {
		deletedAt = sql.NullTime{Time: time.Time{}, Valid: false}
	}

	return UserModel{
		ID:             entity.ID.Value,
		Name:           entity.Name,
		Email:          entity.Email,
		ProfilePicture: profilePicture,
		CreatedAt:      entity.CreatedAt,
		UpdatedAt:      entity.UpdatedAt,
		DeletedAt:      deletedAt,
		Version:        entity.Version,
	}
}
