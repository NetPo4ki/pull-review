package prs

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/NetPo4ki/pull-review/internal/repo/sqlc"
	svc "github.com/NetPo4ki/pull-review/internal/service/prs"
	"github.com/go-chi/chi/v5"
)

type Handler struct{ s *svc.Service }

func NewHandler(s *svc.Service) *Handler { return &Handler{s: s} }

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/create", h.create)
	r.Post("/merge", h.merge)
	r.Post("/reassign", h.reassign)
	return r
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var in CreatePRRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	pr, reviewers, err := h.s.CreatePR(r.Context(), in.PullRequestID, in.PullRequestName, in.AuthorID)
	if err != nil {
		switch err {
		case svc.ErrNotFound:
			writeError(w, http.StatusNotFound, "NOT_FOUND", "resource not found")
			return
		case svc.ErrPRExists:
			writeError(w, http.StatusConflict, "PR_EXISTS", "PR id already exists")
			return
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}
	resp := CreatePRResponse{PR: toDTO(pr, reviewers)}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) merge(w http.ResponseWriter, r *http.Request) {
	var in MergeRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	pr, err := h.s.MergePR(r.Context(), in.PullRequestID)
	if err != nil {
		if err == svc.ErrNotFound {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "resource not found")
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	revs, _ := h.s.GetPRReviewers(r.Context(), in.PullRequestID)
	resp := MergeResponse{PR: toDTO(pr, revs)}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) reassign(w http.ResponseWriter, r *http.Request) {
	var in ReassignRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	newID, err := h.s.ReassignReviewer(r.Context(), in.PullRequestID, in.OldUserID)
	if err != nil {
		switch err {
		case svc.ErrNotFound:
			writeError(w, http.StatusNotFound, "NOT_FOUND", "resource not found")
			return
		case svc.ErrPRMerged:
			writeError(w, http.StatusConflict, "PR_MERGED", "cannot reassign on merged PR")
			return
		case svc.ErrNotAssigned:
			writeError(w, http.StatusConflict, "NOT_ASSIGNED", "reviewer is not assigned to this PR")
			return
		case svc.ErrNoCandidate:
			writeError(w, http.StatusConflict, "NO_CANDIDATE", "no active replacement candidate in team")
			return
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}
	pr, _ := h.s.GetPRByID(r.Context(), in.PullRequestID)
	revs, _ := h.s.GetPRReviewers(r.Context(), in.PullRequestID)
	resp := ReassignResponse{PR: toDTO(pr, revs), ReplacedBy: newID}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
func toDTO(pr sqlc.PullRequest, reviewers []string) PRDTO {
	var merged *string
	if pr.MergedAt.Valid {
		s := pr.MergedAt.Time.UTC().Format(time.RFC3339)
		merged = &s
	}
	return PRDTO{
		PullRequestID:   pr.PrID,
		PullRequestName: pr.Name,
		AuthorID:        pr.AuthorID,
		Status:          string(pr.Status),
		Assigned:        reviewers,
		MergedAt:        merged,
	}
}

func writeError(w http.ResponseWriter, code int, errCode, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]any{"error": map[string]string{"code": errCode, "message": msg}})
}
