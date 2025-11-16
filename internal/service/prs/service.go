package prs

import (
	"context"
	"errors"

	"github.com/NetPo4ki/pull-review/internal/repo/sqlc"
)

type UsersRepo interface {
	GetUserByID(ctx context.Context, userID string) (sqlc.User, error)
}

type PRsRepo interface {
	CreatePR(ctx context.Context, prID, name, authorID string) (sqlc.PullRequest, error)
	InsertPRReviewer(ctx context.Context, prID, userID string) error
	DeletePRReviewer(ctx context.Context, prID, userID string) error
	GetPR(ctx context.Context, prID string) (sqlc.PullRequest, error)
	GetPRReviewers(ctx context.Context, prID string) ([]string, error)
	UpdatePRStatusIfOpen(ctx context.Context, prID string) (sqlc.PullRequest, error)
	CandidatesForCreate(ctx context.Context, authorID string, limitCount int32) ([]string, error)
	CandidateForReassign(ctx context.Context, prID, oldUserID string) (string, error)
}

type TxRunner interface {
	Do(ctx context.Context, fn func(context.Context, *sqlc.Queries) error) error
}

type Service struct {
	users UsersRepo
	prs   PRsRepo
	tx    TxRunner
}

func New(users UsersRepo, prs PRsRepo, tx TxRunner) *Service {
	return &Service{users: users, prs: prs, tx: tx}
}

func (s *Service) CreatePR(ctx context.Context, prID, name, authorID string) (sqlc.PullRequest, []string, error) {
	if _, err := s.users.GetUserByID(ctx, authorID); err != nil {
		return sqlc.PullRequest{}, nil, ErrNotFound
	}
	cands, err := s.prs.CandidatesForCreate(ctx, authorID, 2)
	if err != nil {
		return sqlc.PullRequest{}, nil, err
	}

	var pr sqlc.PullRequest
	if err := s.tx.Do(ctx, func(ctx context.Context, q *sqlc.Queries) error {
		var e error
		pr, e = s.prs.CreatePR(ctx, prID, name, authorID)
		if e != nil {
			if isPGUnique(e) {
				return ErrPRExists
			}
			return e
		}
		for _, u := range cands {
			if e = s.prs.InsertPRReviewer(ctx, prID, u); e != nil {
				return e
			}
		}
		return nil
	}); err != nil {
		return sqlc.PullRequest{}, nil, err
	}

	return pr, cands, nil
}

func (s *Service) MergePR(ctx context.Context, prID string) (sqlc.PullRequest, error) {
	_, err := s.prs.GetPR(ctx, prID)
	if err != nil {
		return sqlc.PullRequest{}, ErrNotFound
	}
	pr, err := s.prs.UpdatePRStatusIfOpen(ctx, prID)
	if err == nil {
		return pr, nil
	}
	return s.prs.GetPR(ctx, prID)
}

func (s *Service) ReassignReviewer(ctx context.Context, prID, oldUserID string) (string, error) {
	pr, err := s.prs.GetPR(ctx, prID)
	if err != nil {
		return "", ErrNotFound
	}
	if string(pr.Status) == "MERGED" {
		return "", ErrPRMerged
	}

	revs, err := s.prs.GetPRReviewers(ctx, prID)
	if err != nil {
		return "", err
	}
	if !contains(revs, oldUserID) {
		return "", ErrNotAssigned
	}

	cand, err := s.prs.CandidateForReassign(ctx, prID, oldUserID)
	if err != nil || cand == "" {
		return "", ErrNoCandidate
	}

	if err := s.tx.Do(ctx, func(ctx context.Context, q *sqlc.Queries) error {
		if e := s.prs.DeletePRReviewer(ctx, prID, oldUserID); e != nil {
			return e
		}
		if e := s.prs.InsertPRReviewer(ctx, prID, cand); e != nil {
			return e
		}
		return nil
	}); err != nil {
		return "", err
	}
	return cand, nil
}

func (s *Service) GetPRByID(ctx context.Context, prID string) (sqlc.PullRequest, error) {
	return s.prs.GetPR(ctx, prID)
}

func (s *Service) GetPRReviewers(ctx context.Context, prID string) ([]string, error) {
	return s.prs.GetPRReviewers(ctx, prID)
}

func contains(s []string, t string) bool {
	for _, x := range s {
		if x == t {
			return true
		}
	}
	return false
}

func isPGUnique(err error) bool {
	var pgErr interface{ SQLState() string }
	if errors.As(err, &pgErr) && pgErr.SQLState() == "23505" {
		return true
	}
	return false
}
