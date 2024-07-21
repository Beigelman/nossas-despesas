package usecase

import (
	"context"
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	vo "github.com/Beigelman/nossas-despesas/internal/shared/infra/email"
	"html/template"
	"log/slog"
	"strings"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/domain/repository"
	"github.com/Beigelman/nossas-despesas/internal/domain/service"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/google/uuid"
)

type (
	InviteUserToGroupInput struct {
		GroupID group.ID
		Email   string
		BaseURL string
	}

	InviteUserToGroup func(ctx context.Context, input InviteUserToGroupInput) (*group.Invite, error)
)

func NewInviteUserToGroup(
	userRepo repository.UserRepository,
	groupRepo group.Repository,
	groupInviteRepo group.InviteRepository,
	emailProvider service.EmailProvider,
) InviteUserToGroup {
	return func(ctx context.Context, input InviteUserToGroupInput) (*group.Invite, error) {
		group, err := groupRepo.GetByID(ctx, input.GroupID)
		if err != nil {
			return nil, fmt.Errorf("groupRepo.GetByID: %w", err)
		}

		if group == nil {
			return nil, except.NotFoundError("group not found")
		}

		invitee, err := userRepo.GetByEmail(ctx, input.Email)
		if err != nil {
			return nil, fmt.Errorf("userRepo.GetByEmail: %w", err)
		}

		if invitee != nil && invitee.GroupID != nil {
			return nil, except.UnprocessableEntityError("user already in a group")
		}

		invites, err := groupInviteRepo.GetGroupInvitesByEmail(ctx, group.ID, input.Email)
		if err != nil {
			return nil, fmt.Errorf("groupInviteRepo.GetByEmail: %w", err)
		}

		if len(invites) > 3 && invites[2].CreatedAt.After(time.Now().Add(-time.Hour*48)) {
			return nil, except.NewHTTPError(429, "too many invites sent to this email recently")
		}

		groupInvite := group.NewGroupInvite(group.GroupInviteParams{
			ID:        groupInviteRepo.GetNextID(),
			GroupID:   group.ID,
			Token:     uuid.NewString(),
			Email:     input.Email,
			ExpiresAt: time.Now().Add(time.Hour * 48),
		})

		if err := groupInviteRepo.Store(ctx, groupInvite); err != nil {
			return nil, fmt.Errorf("groupInviteRepo.Store: %w", err)
		}

		go func() {
			tmpl, err := template.ParseFiles("./templates/group_invite.html")
			if err != nil {
				slog.Error("failed to parse template file", "error", err)
				return
			}

			html := strings.Builder{}
			if err := tmpl.Execute(&html, map[string]any{
				"GroupName": group.Name,
				"Link":      groupInvite.InviteURL(input.BaseURL),
			}); err != nil {
				slog.Error("failed to execute template", "error", err)
				return
			}

			if err := emailProvider.Send(ctx, vo.Email{
				From:    "noreplay@nossasdespesas.com.br",
				To:      []string{input.Email},
				Html:    html.String(),
				Subject: "Convite para compartilhar despesas",
			}); err != nil {
				slog.Error("failed to send group invite email", "error", err)
				return
			}

			if err := groupInvite.Sent(); err != nil {
				slog.Error("failed to update group invite status", "error", err)
				return
			}

			if err := groupInviteRepo.Store(ctx, groupInvite); err != nil {
				slog.Error("failed to update group invite status", "error", err)
				return
			}

			slog.Info("group invite email sent", "email", input.Email)
		}()

		return groupInvite, nil
	}
}
