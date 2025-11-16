package stats

import (
	"encoding/json"
	"net/http"

	"github.com/NetPo4ki/pull-review/internal/service/stats"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	s *stats.Service
}

func NewHandler(s *stats.Service) *Handler { return &Handler{s: s} }

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.getStats)
	return r
}

func (h *Handler) getStats(w http.ResponseWriter, r *http.Request) {
	res, err := h.s.Get(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	out := StatsResponse{
		Users:        make([]UserAssignmentDTO, 0, len(res.Users)),
		PullRequests: make([]PRAssignmentDTO, 0, len(res.PullRequests)),
	}
	for _, u := range res.Users {
		out.Users = append(out.Users, UserAssignmentDTO{
			UserID:        u.UserID,
			AssignedCount: u.AssignedCount,
		})
	}
	for _, p := range res.PullRequests {
		out.PullRequests = append(out.PullRequests, PRAssignmentDTO{
			PullRequestID: p.PrID,
			ReviewerCount: p.ReviewerCount,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}
