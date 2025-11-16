package users

import (
	"encoding/json"
	"net/http"
	"strings"

	svc "github.com/NetPo4ki/pull-review/internal/service/users"
	"github.com/go-chi/chi/v5"
)

type Handler struct{ s *svc.Service }

func NewHandler(s *svc.Service) *Handler { return &Handler{s: s} }

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/setIsActive", h.setIsActive)
	r.Get("/getReview", h.getReview)
	return r
}

func (h *Handler) setIsActive(w http.ResponseWriter, r *http.Request) {
	var in SetIsActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(in.UserID) == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}
	u, err := h.s.SetIsActive(r.Context(), in.UserID, in.IsActive)
	if err != nil {
		if err == svc.ErrNotFound {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "resource not found")
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	resp := SetIsActiveResponse{
		User: UserDTO{
			UserID:   u.UserID,
			Username: u.Username,
			TeamName: u.TeamName,
			IsActive: u.IsActive,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) getReview(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}
	list, err := h.s.GetReview(r.Context(), userID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	resp := GetReviewResponse{
		UserID:       userID,
		PullRequests: make([]PullRequestShortDTO, 0, len(list)),
	}
	for _, p := range list {
		resp.PullRequests = append(resp.PullRequests, PullRequestShortDTO{
			PullRequestID:   p.PrID,
			PullRequestName: p.Name,
			AuthorID:        p.AuthorID,
			Status:          p.Status,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func writeError(w http.ResponseWriter, code int, errCode, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]string{
			"code":    errCode,
			"message": msg,
		},
	})
}
