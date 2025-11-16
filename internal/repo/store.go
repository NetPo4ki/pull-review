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

func (s *Store) CreatePR(ctx context.Context, prID, name, authorID string) (sqlc.PullRequest, error) {
	return s.q.CreatePR(ctx, sqlc.CreatePRParams{PrID: prID, Name: name, AuthorID: authorID})
}
func (s *Store) InsertPRReviewer(ctx context.Context, prID, userID string) error {
	return s.q.InsertPRReviewer(ctx, sqlc.InsertPRReviewerParams{PrID: prID, UserID: userID})
}
func (s *Store) DeletePRReviewer(ctx context.Context, prID, userID string) error {
	return s.q.DeletePRReviewer(ctx, sqlc.DeletePRReviewerParams{PrID: prID, UserID: userID})
}
func (s *Store) GetPR(ctx context.Context, prID string) (sqlc.PullRequest, error) {
	return s.q.GetPR(ctx, prID)
}
func (s *Store) GetPRReviewers(ctx context.Context, prID string) ([]string, error) {
	return s.q.GetPRReviewers(ctx, prID)
}
func (s *Store) UpdatePRStatusIfOpen(ctx context.Context, prID string) (sqlc.PullRequest, error) {
	return s.q.UpdatePRStatusIfOpen(ctx, prID)
}
func (s *Store) CandidatesForCreate(ctx context.Context, authorID string, limitCount int32) ([]string, error) {
	return s.q.CandidatesForCreate(ctx, sqlc.CandidatesForCreateParams{AuthorID: authorID, LimitCount: limitCount})
}
func (s *Store) CandidateForReassign(ctx context.Context, prID, oldUserID string) (string, error) {
	return s.q.CandidateForReassign(ctx, sqlc.CandidateForReassignParams{OldUserID: oldUserID, PrID: prID})
}
