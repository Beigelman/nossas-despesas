package postgres

import (
	"database/sql"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"

	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
)

func toEntity(model UserModel) *user.User {
	var profilePicture *string
	if model.ProfilePicture.Valid {
		profilePicture = &model.ProfilePicture.String
	}

	var deletedAt *time.Time
	if model.DeletedAt.Valid {
		deletedAt = &model.DeletedAt.Time
	}

	var groupID *entity.GroupID
	if model.GroupID.Valid {
		groupID = &entity.GroupID{Value: int(model.GroupID.Int64)}
	}

	return &user.User{
		Entity: ddd.Entity[user.ID]{
			ID:        user.ID{Value: model.ID},
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
			DeletedAt: deletedAt,
			Version:   model.Version,
		},
		GroupID:        groupID,
		Name:           model.Name,
		Email:          model.Email,
		ProfilePicture: profilePicture,
	}
}

func toModel(entity *user.User) UserModel {
	profilePicture := sql.NullString{String: "", Valid: false}
	if entity.ProfilePicture != nil {
		profilePicture = sql.NullString{String: *entity.ProfilePicture, Valid: true}
	}
	deletedAt := sql.NullTime{Time: time.Time{}, Valid: false}
	if entity.DeletedAt != nil {
		deletedAt = sql.NullTime{Time: *entity.DeletedAt, Valid: true}
	}

	groupID := sql.NullInt64{Int64: 0, Valid: false}
	if entity.GroupID != nil {
		groupID = sql.NullInt64{Int64: int64(entity.GroupID.Value), Valid: true}
	}

	return UserModel{
		ID:             entity.ID.Value,
		Name:           entity.Name,
		Email:          entity.Email,
		GroupID:        groupID,
		ProfilePicture: profilePicture,
		CreatedAt:      entity.CreatedAt,
		UpdatedAt:      entity.UpdatedAt,
		DeletedAt:      deletedAt,
		Version:        entity.Version,
	}
}
