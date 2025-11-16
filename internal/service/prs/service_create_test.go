package prs

import (
	"context"
	"errors"
	"testing"

	"github.com/NetPo4ki/pull-review/internal/repo/sqlc"
)

type prsRepoCreateStub struct {
	create func(ctx context.Context, prID, name, authorID string) (sqlc.PullRequest, error)
	cands  func(ctx context.Context, authorID string, limitCount int32) ([]string, error)
}

func (s prsRepoCreateStub) CreatePR(ctx context.Context, prID, name, authorID string) (sqlc.PullRequest, error) {
	return s.create(ctx, prID, name, authorID)
}
func (s prsRepoCreateStub) InsertPRReviewer(context.Context, string, string) error { return nil }
func (s prsRepoCreateStub) DeletePRReviewer(context.Context, string, string) error { return nil }
func (s prsRepoCreateStub) GetPR(context.Context, string) (sqlc.PullRequest, error) {
	return sqlc.PullRequest{}, nil
}
func (s prsRepoCreateStub) GetPRReviewers(context.Context, string) ([]string, error) { return nil, nil }
func (s prsRepoCreateStub) UpdatePRStatusIfOpen(context.Context, string) (sqlc.PullRequest, error) {
	return sqlc.PullRequest{}, nil
}
func (s prsRepoCreateStub) CandidatesForCreate(ctx context.Context, authorID string, limitCount int32) ([]string, error) {
	return s.cands(ctx, authorID, limitCount)
}
func (s prsRepoCreateStub) CandidateForReassign(context.Context, string, string) (string, error) {
	return "", nil
}

type usersRepoCreateStub struct {
	get func(ctx context.Context, userID string) (sqlc.User, error)
}

func (u usersRepoCreateStub) GetUserByID(ctx context.Context, userID string) (sqlc.User, error) {
	return u.get(ctx, userID)
}

type txNoopCreate struct{}

func (t *txNoopCreate) Do(ctx context.Context, fn func(context.Context, *sqlc.Queries) error) error {
	return fn(ctx, nil)
}

func TestCreatePR_AuthorNotFound(t *testing.T) {
	repo := prsRepoCreateStub{
		create: func(ctx context.Context, a, b, c string) (sqlc.PullRequest, error) { return sqlc.PullRequest{}, nil },
		cands:  func(ctx context.Context, a string, b int32) ([]string, error) { return nil, nil },
	}
	s := New(usersRepoCreateStub{
		get: func(ctx context.Context, id string) (sqlc.User, error) { return sqlc.User{}, errors.New("no rows") },
	}, repo, &txNoopCreate{})

	_, _, err := s.CreatePR(context.Background(), "x", "name", "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestCreatePR_AssignsUpToTwo(t *testing.T) {
	repo := prsRepoCreateStub{
		create: func(ctx context.Context, prID, name, authorID string) (sqlc.PullRequest, error) {
			return sqlc.PullRequest{PrID: prID, Name: name, AuthorID: authorID, Status: "OPEN"}, nil
		},
		cands: func(ctx context.Context, a string, b int32) ([]string, error) {
			return []string{"u1", "u2"}, nil
		},
	}
	s := New(usersRepoCreateStub{
		get: func(ctx context.Context, id string) (sqlc.User, error) { return sqlc.User{UserID: id}, nil },
	}, repo, &txNoopCreate{})

	_, reviewers, err := s.CreatePR(context.Background(), "p1", "x", "author")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(reviewers) != 2 {
		t.Fatalf("expected 2 reviewers, got %d", len(reviewers))
	}
}
