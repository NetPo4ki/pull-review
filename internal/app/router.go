package app

import (
	"log/slog"
	"net/http"

	prsapi "github.com/NetPo4ki/pull-review/internal/api/prs"
	teamsapi "github.com/NetPo4ki/pull-review/internal/api/teams"
	usersapi "github.com/NetPo4ki/pull-review/internal/api/users"
	mymw "github.com/NetPo4ki/pull-review/internal/middleware"
	prssvc "github.com/NetPo4ki/pull-review/internal/service/prs"
	teamssvc "github.com/NetPo4ki/pull-review/internal/service/teams"
	userssvc "github.com/NetPo4ki/pull-review/internal/service/users"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(l *slog.Logger, teamsSvc *teamssvc.Service, usersSvc *userssvc.Service, prsSvc *prssvc.Service) http.Handler {
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
	r.Mount("/users", usersapi.NewHandler(usersSvc).Routes())
	r.Mount("/pullRequest", prsapi.NewHandler(prsSvc).Routes())

	return r
}
