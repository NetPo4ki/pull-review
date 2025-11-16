package teams

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	srv "github.com/NetPo4ki/pull-review/internal/service/teams"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	s *srv.Service
}

func NewHandler(s *srv.Service) *Handler { return &Handler{s: s} }

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/add", h.addTeam)
	r.Get("/get", h.getTeam)
	return r
}

func (h *Handler) addTeam(w http.ResponseWriter, r *http.Request) {
	var in AddTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(in.TeamName) == "" {
		http.Error(w, "team_name is required", http.StatusBadRequest)
		return
	}
	for _, m := range in.Members {
		if strings.TrimSpace(m.UserID) == "" || strings.TrimSpace(m.Username) == "" {
			http.Error(w, "member user_id and username are required", http.StatusBadRequest)
			return
		}
	}
	team := srv.Team{TeamName: in.TeamName}
	for _, m := range in.Members {
		team.Members = append(team.Members, srv.MemberInput{
			UserID:   m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
		})
	}

	out, err := h.s.AddTeam(r.Context(), team)
	if err != nil {
		switch err {
		case srv.ErrTeamExists:
			writeError(w, http.StatusBadRequest, "TEAM_EXISTS", "team_name already exists")
			return
		default:
			slog.Error("team_add_error", "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}

	resp := AddTeamResponse{
		Team: TeamDTO{
			TeamName: out.TeamName,
			Members:  make([]TeamMemberDTO, 0, len(out.Members)),
		},
	}
	for _, m := range out.Members {
		resp.Team.Members = append(resp.Team.Members, TeamMemberDTO{
			UserID:   m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) getTeam(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		http.Error(w, "team_name is required", http.StatusBadRequest)
		return
	}

	out, err := h.s.GetTeam(r.Context(), teamName)
	if err != nil {
		switch err {
		case srv.ErrNotFound:
			writeError(w, http.StatusNotFound, "NOT_FOUND", "resource not found")
			return
		default:
			slog.Error("team_get_error", "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}

	resp := TeamDTO{
		TeamName: out.TeamName,
		Members:  make([]TeamMemberDTO, 0, len(out.Members)),
	}
	for _, m := range out.Members {
		resp.Members = append(resp.Members, TeamMemberDTO{
			UserID:   m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
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
