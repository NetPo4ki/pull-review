package users

import (
	"context"

	"github.com/NetPo4ki/pull-review/internal/repo/sqlc"
)

type UsersRepo interface {
	SetIsActive(ctx context.Context, userID string, isActive bool) (sqlc.User, error)
	GetUserByID(ctx context.Context, userID string) (sqlc.User, error)
}

type PRsRepo interface {
	GetPRsForReviewer(ctx context.Context, userID string) ([]sqlc.PullRequest, error)
}

type Service struct {
	users UsersRepo
	prs   PRsRepo
}

func New(users UsersRepo, prs PRsRepo) *Service { return &Service{users: users, prs: prs} }

func (s *Service) SetIsActive(ctx context.Context, userID string, isActive bool) (sqlc.User, error) {
	u, err := s.users.SetIsActive(ctx, userID, isActive)
	if err != nil {
		return sqlc.User{}, ErrNotFound
	}
	return u, nil
}

type PRShort struct {
	PrID     string
	Name     string
	AuthorID string
	Status   string
}

func (s *Service) GetReview(ctx context.Context, userID string) ([]PRShort, error) {
	list, err := s.prs.GetPRsForReviewer(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]PRShort, 0, len(list))
	for _, p := range list {
		out = append(out, PRShort{
			PrID:     p.PrID,
			Name:     p.Name,
			AuthorID: p.AuthorID,
			Status:   string(p.Status),
		})
	}
	return out, nil
}
