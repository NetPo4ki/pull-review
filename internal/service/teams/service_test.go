package teams

import (
	"context"
	"errors"
	"testing"

	"github.com/NetPo4ki/pull-review/internal/repo/sqlc"
	"github.com/jackc/pgx/v5/pgconn"
)

type usersRepoStub struct {
	upsert func(ctx context.Context, userID, username, teamName string, isActive bool) (sqlc.User, error)
}

func (u usersRepoStub) UpsertUser(ctx context.Context, a, b, c string, d bool) (sqlc.User, error) {
	return u.upsert(ctx, a, b, c, d)
}

type teamsRepoStub struct {
	create func(ctx context.Context, teamName string) error
	get    func(ctx context.Context, teamName string) ([]sqlc.GetTeamWithMembersRow, error)
}

func (t teamsRepoStub) CreateTeam(ctx context.Context, teamName string) error {
	return t.create(ctx, teamName)
}
func (t teamsRepoStub) GetTeamWithMembers(ctx context.Context, teamName string) ([]sqlc.GetTeamWithMembersRow, error) {
	return t.get(ctx, teamName)
}

func TestAddTeam_Duplicate(t *testing.T) {
	pgdup := &pgconn.PgError{Code: "23505"}
	s := New(usersRepoStub{upsert: func(ctx context.Context, a, b, c string, d bool) (sqlc.User, error) {
		return sqlc.User{}, nil
	}}, teamsRepoStub{
		create: func(ctx context.Context, teamName string) error { return pgdup },
		get:    nil,
	})
	_, err := s.AddTeam(context.Background(), Team{TeamName: "backend"})
	if !errors.Is(err, ErrTeamExists) {
		t.Fatalf("expected ErrTeamExists, got %v", err)
	}
}

func TestGetTeam_NotFound(t *testing.T) {
	s := New(nil, teamsRepoStub{
		create: nil,
		get:    func(ctx context.Context, teamName string) ([]sqlc.GetTeamWithMembersRow, error) { return nil, nil },
	})
	_, err := s.GetTeam(context.Background(), "unknown")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
