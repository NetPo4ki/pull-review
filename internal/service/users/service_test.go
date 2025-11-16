package users

import (
	"context"
	"errors"
	"testing"

	"github.com/NetPo4ki/pull-review/internal/repo/sqlc"
)

type usersRepoStub struct {
	set func(ctx context.Context, userID string, isActive bool) (sqlc.User, error)
	get func(ctx context.Context, userID string) (sqlc.User, error)
}

func (u usersRepoStub) SetIsActive(ctx context.Context, userID string, isActive bool) (sqlc.User, error) {
	return u.set(ctx, userID, isActive)
}

func (u usersRepoStub) GetUserByID(ctx context.Context, userID string) (sqlc.User, error) {
	return u.get(ctx, userID)
}

type prsRepoStub struct {
	get func(ctx context.Context, userID string) ([]sqlc.PullRequest, error)
}

func (p prsRepoStub) GetPRsForReviewer(ctx context.Context, userID string) ([]sqlc.PullRequest, error) {
	return p.get(ctx, userID)
}

func TestSetIsActive_NotFound(t *testing.T) {
	s := New(usersRepoStub{
		set: func(ctx context.Context, id string, a bool) (sqlc.User, error) {
			return sqlc.User{}, errors.New("no rows")
		},
		get: nil,
	}, prsRepoStub{get: nil})

	_, err := s.SetIsActive(context.Background(), "missing", false)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestGetReview_OK(t *testing.T) {
	s := New(usersRepoStub{}, prsRepoStub{
		get: func(ctx context.Context, id string) ([]sqlc.PullRequest, error) {
			return []sqlc.PullRequest{{PrID: "pr-1", Name: "X", AuthorID: "u1", Status: "OPEN"}}, nil
		},
	})
	out, err := s.GetReview(context.Background(), "u2")
	if err != nil || len(out) != 1 || out[0].PrID != "pr-1" {
		t.Fatalf("unexpected: %+v, err=%v", out, err)
	}
}
