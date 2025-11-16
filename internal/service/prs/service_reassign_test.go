package prs

import (
	"context"
	"errors"
	"testing"

	"github.com/NetPo4ki/pull-review/internal/repo/sqlc"
)

type prsRepoReassignStub struct {
	get       func(ctx context.Context, prID string) (sqlc.PullRequest, error)
	reviewers func(ctx context.Context, prID string) ([]string, error)
	cand      func(ctx context.Context, prID, old string) (string, error)
	del       func(ctx context.Context, prID, old string) error
	ins       func(ctx context.Context, prID, uid string) error
}

func (s prsRepoReassignStub) CreatePR(context.Context, string, string, string) (sqlc.PullRequest, error) {
	return sqlc.PullRequest{}, nil
}
func (s prsRepoReassignStub) InsertPRReviewer(ctx context.Context, prID, userID string) error {
	if s.ins != nil {
		return s.ins(ctx, prID, userID)
	}
	return nil
}
func (s prsRepoReassignStub) DeletePRReviewer(ctx context.Context, prID, userID string) error {
	if s.del != nil {
		return s.del(ctx, prID, userID)
	}
	return nil
}
func (s prsRepoReassignStub) GetPR(ctx context.Context, prID string) (sqlc.PullRequest, error) {
	return s.get(ctx, prID)
}
func (s prsRepoReassignStub) GetPRReviewers(ctx context.Context, prID string) ([]string, error) {
	return s.reviewers(ctx, prID)
}
func (s prsRepoReassignStub) UpdatePRStatusIfOpen(context.Context, string) (sqlc.PullRequest, error) {
	return sqlc.PullRequest{}, nil
}
func (s prsRepoReassignStub) CandidatesForCreate(context.Context, string, int32) ([]string, error) {
	return nil, nil
}
func (s prsRepoReassignStub) CandidateForReassign(ctx context.Context, prID, oldUserID string) (string, error) {
	return s.cand(ctx, prID, oldUserID)
}

type usersRepoReassignStub struct{}

func (u usersRepoReassignStub) GetUserByID(ctx context.Context, userID string) (sqlc.User, error) {
	return sqlc.User{UserID: userID}, nil
}

type txNoopReassign struct{}

func (t *txNoopReassign) Do(ctx context.Context, fn func(context.Context, *sqlc.Queries) error) error {
	return fn(ctx, nil)
}

func TestReassign_NotAssigned(t *testing.T) {
	repo := prsRepoReassignStub{
		get: func(ctx context.Context, prID string) (sqlc.PullRequest, error) {
			return sqlc.PullRequest{PrID: prID, Status: "OPEN"}, nil
		},
		reviewers: func(ctx context.Context, prID string) ([]string, error) { return []string{"u1", "u2"}, nil },
		cand:      func(ctx context.Context, prID, old string) (string, error) { return "", errors.New("no") },
	}
	s := New(usersRepoReassignStub{}, repo, &txNoopReassign{})
	_, err := s.ReassignReviewer(context.Background(), "pr1", "u9")
	if !errors.Is(err, ErrNotAssigned) {
		t.Fatalf("expected ErrNotAssigned, got %v", err)
	}
}

func TestReassign_Success(t *testing.T) {
	repo := prsRepoReassignStub{
		get: func(ctx context.Context, prID string) (sqlc.PullRequest, error) {
			return sqlc.PullRequest{PrID: prID, Status: "OPEN"}, nil
		},
		reviewers: func(ctx context.Context, prID string) ([]string, error) { return []string{"u1", "u2"}, nil },
		cand:      func(ctx context.Context, prID, old string) (string, error) { return "u3", nil },
		del:       func(ctx context.Context, prID, old string) error { return nil },
		ins:       func(ctx context.Context, prID, uid string) error { return nil },
	}
	s := New(usersRepoReassignStub{}, repo, &txNoopReassign{})
	newID, err := s.ReassignReviewer(context.Background(), "pr1", "u1")
	if err != nil || newID != "u3" {
		t.Fatalf("expected u3, got %q err=%v", newID, err)
	}
}

func TestReassign_NoCandidate(t *testing.T) {
	repo := prsRepoReassignStub{
		get: func(ctx context.Context, prID string) (sqlc.PullRequest, error) {
			return sqlc.PullRequest{PrID: prID, Status: "OPEN"}, nil
		},
		reviewers: func(ctx context.Context, prID string) ([]string, error) { return []string{"u1"}, nil },
		cand:      func(ctx context.Context, prID, old string) (string, error) { return "", nil },
	}
	s := New(usersRepoReassignStub{}, repo, &txNoopReassign{})
	_, err := s.ReassignReviewer(context.Background(), "pr1", "u1")
	if !errors.Is(err, ErrNoCandidate) {
		t.Fatalf("expected ErrNoCandidate, got %v", err)
	}
}

func TestReassign_PRIsMerged(t *testing.T) {
	repo := prsRepoReassignStub{
		get: func(ctx context.Context, prID string) (sqlc.PullRequest, error) {
			return sqlc.PullRequest{PrID: prID, Status: "MERGED"}, nil
		},
		reviewers: func(ctx context.Context, prID string) ([]string, error) { return []string{"u1", "u2"}, nil },
		cand:      func(ctx context.Context, prID, old string) (string, error) { return "u3", nil },
	}
	s := New(usersRepoReassignStub{}, repo, &txNoopReassign{})
	_, err := s.ReassignReviewer(context.Background(), "pr1", "u1")
	if !errors.Is(err, ErrPRMerged) {
		t.Fatalf("expected ErrPRMerged, got %v", err)
	}
}
