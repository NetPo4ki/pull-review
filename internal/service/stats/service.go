package stats

import (
	"context"
)

type PRsStatsRepo interface {
	AssignmentsPerUser(ctx context.Context) ([]AssignmentsPerUserRow, error)
	AssignmentsPerPR(ctx context.Context) ([]AssignmentsPerPRRow, error)
}

type AssignmentsPerUserRow struct {
	UserID        string
	AssignedCount int32
}

type AssignmentsPerPRRow struct {
	PrID          string
	ReviewerCount int32
}

type Service struct {
	r PRsStatsRepo
}

func New(r PRsStatsRepo) *Service { return &Service{r: r} }

type Stats struct {
	Users        []AssignmentsPerUserRow
	PullRequests []AssignmentsPerPRRow
}

func (s *Service) Get(ctx context.Context) (Stats, error) {
	users, err := s.r.AssignmentsPerUser(ctx)
	if err != nil {
		return Stats{}, err
	}
	prs, err := s.r.AssignmentsPerPR(ctx)
	if err != nil {
		return Stats{}, err
	}
	return Stats{Users: users, PullRequests: prs}, nil
}
