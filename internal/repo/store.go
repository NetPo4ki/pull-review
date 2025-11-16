package repo

import (
	"context"

	"github.com/NetPo4ki/pull-review/internal/repo/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	q *sqlc.Queries
}

func NewStore(pool *pgxpool.Pool) *Store { return &Store{q: sqlc.New(pool)} }

func (s *Store) CreateTeam(ctx context.Context, teamName string) error {
	return s.q.CreateTeam(ctx, teamName)
}
func (s *Store) GetTeamWithMembers(ctx context.Context, teamName string) ([]sqlc.GetTeamWithMembersRow, error) {
	return s.q.GetTeamWithMembers(ctx, teamName)
}

func (s *Store) SetIsActive(ctx context.Context, userID string, isActive bool) (sqlc.User, error) {
	return s.q.SetIsActive(ctx, sqlc.SetIsActiveParams{UserID: userID, IsActive: isActive})
}
func (s *Store) GetUserByID(ctx context.Context, userID string) (sqlc.User, error) {
	return s.q.GetUserByID(ctx, userID)
}
func (s *Store) GetPRsForReviewer(ctx context.Context, userID string) ([]sqlc.PullRequest, error) {
	return s.q.GetPRsForReviewer(ctx, userID)
}

func (s *Store) UpsertUser(ctx context.Context, userID, username, teamName string, isActive bool) (sqlc.User, error) {
	return s.q.UpsertUser(ctx, sqlc.UpsertUserParams{
		UserID:   userID,
		Username: username,
		TeamName: teamName,
		IsActive: isActive,
	})
}
