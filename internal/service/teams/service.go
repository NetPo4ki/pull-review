package teams

import (
	"context"
	"errors"

	"github.com/NetPo4ki/pull-review/internal/repo/sqlc"
	"github.com/jackc/pgx/v5/pgconn"
)

type UsersRepo interface {
	UpsertUser(ctx context.Context, userID, username, teamName string, isActive bool) (sqlc.User, error)
}

type TeamsRepo interface {
	CreateTeam(ctx context.Context, teamName string) error
	GetTeamWithMembers(ctx context.Context, teamName string) ([]sqlc.GetTeamWithMembersRow, error)
}

type Service struct {
	users UsersRepo
	teams TeamsRepo
}

func New(users UsersRepo, teams TeamsRepo) *Service {
	return &Service{users: users, teams: teams}
}

type MemberInput struct {
	UserID   string
	Username string
	IsActive bool
}

type Team struct {
	TeamName string
	Members  []MemberInput
}

type TeamDTO struct {
	TeamName string
	Members  []struct {
		UserID   string
		Username string
		IsActive bool
	}
}

func (s *Service) AddTeam(ctx context.Context, team Team) (TeamDTO, error) {
	if err := s.teams.CreateTeam(ctx, team.TeamName); err != nil {
		var pgErr *pgconn.PgError
		if errorsAs(err, &pgErr) && pgErr.Code == "23505" {
			return TeamDTO{}, ErrTeamExists
		}
		return TeamDTO{}, err
	}

	out := TeamDTO{TeamName: team.TeamName}
	for _, m := range team.Members {
		_, err := s.users.UpsertUser(ctx, m.UserID, m.Username, team.TeamName, m.IsActive)
		if err != nil {
			return TeamDTO{}, err
		}
		out.Members = append(out.Members, struct {
			UserID   string
			Username string
			IsActive bool
		}{UserID: m.UserID, Username: m.Username, IsActive: m.IsActive})
	}
	return out, nil
}

func (s *Service) GetTeam(ctx context.Context, teamName string) (TeamDTO, error) {
	rows, err := s.teams.GetTeamWithMembers(ctx, teamName)
	if err != nil {
		return TeamDTO{}, err
	}
	if len(rows) == 0 {
		return TeamDTO{}, ErrNotFound
	}
	out := TeamDTO{TeamName: teamName}
	for _, r := range rows {
		if !r.UserID.Valid {
			continue
		}
		out.Members = append(out.Members, struct {
			UserID   string
			Username string
			IsActive bool
		}{UserID: r.UserID.String, Username: r.Username.String, IsActive: r.IsActive.Bool})
	}
	return out, nil
}

func errorsAs(err error, target any) bool { return errors.As(err, target) }
