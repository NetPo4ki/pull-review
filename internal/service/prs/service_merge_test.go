package prs

import (
	"context"
	"errors"
	"testing"

	"github.com/NetPo4ki/pull-review/internal/repo/sqlc"
)

type prsRepoStub struct {
	get     func(ctx context.Context, prID string) (sqlc.PullRequest, error)
	updOpen func(ctx context.Context, prID string) (sqlc.PullRequest, error)
}

func (s prsRepoStub) CreatePR(context.Context, string, string, string) (sqlc.PullRequest, error) {
	return sqlc.PullRequest{}, nil
}
func (s prsRepoStub) InsertPRReviewer(context.Context, string, string) error { return nil }
func (s prsRepoStub) DeletePRReviewer(context.Context, string, string) error { return nil }
func (s prsRepoStub) GetPR(ctx context.Context, prID string) (sqlc.PullRequest, error) {
	return s.get(ctx, prID)
}
func (s prsRepoStub) GetPRReviewers(context.Context, string) ([]string, error) { return nil, nil }
func (s prsRepoStub) UpdatePRStatusIfOpen(ctx context.Context, prID string) (sqlc.PullRequest, error) {
	return s.updOpen(ctx, prID)
}
func (s prsRepoStub) CandidatesForCreate(context.Context, string, int32) ([]string, error) {
	return nil, nil
}
func (s prsRepoStub) CandidateForReassign(context.Context, string, string) (string, error) {
	return "", nil
}

type usersRepoStub struct{}

func (u usersRepoStub) GetUserByID(ctx context.Context, userID string) (sqlc.User, error) {
	return sqlc.User{}, nil
}

type txNoop struct{}

func (t *txNoop) Do(ctx context.Context, fn func(context.Context, *sqlc.Queries) error) error {
	return fn(ctx, nil)
}

func TestMerge_Idempotent(t *testing.T) {
	open := sqlc.PullRequest{PrID: "p1", Name: "x", AuthorID: "a", Status: "OPEN"}
	merged := sqlc.PullRequest{PrID: "p1", Name: "x", AuthorID: "a", Status: "MERGED"}

	repo := prsRepoStub{
		get: func(ctx context.Context, prID string) (sqlc.PullRequest, error) {
			return merged, nil
		},
		updOpen: func(ctx context.Context, prID string) (sqlc.PullRequest, error) {
			return merged, nil
		},
	}
	s := New(usersRepoStub{}, repo, &txNoop{})

	repo.get = func(ctx context.Context, prID string) (sqlc.PullRequest, error) { return open, nil }
	got, err := s.MergePR(context.Background(), "p1")
	if err != nil || string(got.Status) != "MERGED" {
		t.Fatalf("expected MERGED, got %+v, err=%v", got, err)
	}

	repo.get = func(ctx context.Context, prID string) (sqlc.PullRequest, error) { return merged, nil }
	_, err = s.MergePR(context.Background(), "p1")
	if err != nil {
		t.Fatalf("expected nil on idempotent merge, got %v", err)
	}
}

func TestMerge_NotFound(t *testing.T) {
	repo := prsRepoStub{
		get: func(ctx context.Context, prID string) (sqlc.PullRequest, error) {
			return sqlc.PullRequest{}, errors.New("no rows")
		},
		updOpen: func(ctx context.Context, prID string) (sqlc.PullRequest, error) {
			return sqlc.PullRequest{}, nil
		},
	}
	s := New(usersRepoStub{}, repo, &txNoop{})
	_, err := s.MergePR(context.Background(), "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
