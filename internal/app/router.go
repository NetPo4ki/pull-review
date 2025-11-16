package app

import (
	"log/slog"
	"net/http"

	teamsapi "github.com/NetPo4ki/pull-review/internal/api/teams"
	mymw "github.com/NetPo4ki/pull-review/internal/middleware"
	teamssvc "github.com/NetPo4ki/pull-review/internal/service/teams"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(l *slog.Logger, teamsSvc *teamssvc.Service) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.RealIP)
	r.Use(mymw.RequestID)
	r.Use(mymw.Logging(l))
	r.Use(mymw.Recoverer(l))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	r.Mount("/team", teamsapi.NewHandler(teamsSvc).Routes())

	return r
}
