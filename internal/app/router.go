package app

import (
	"log/slog"
	"net/http"

	prsapi "github.com/NetPo4ki/pull-review/internal/api/prs"
	statsapi "github.com/NetPo4ki/pull-review/internal/api/stats"
	teamsapi "github.com/NetPo4ki/pull-review/internal/api/teams"
	usersapi "github.com/NetPo4ki/pull-review/internal/api/users"
	mymw "github.com/NetPo4ki/pull-review/internal/middleware"
	prssvc "github.com/NetPo4ki/pull-review/internal/service/prs"
	statsvc "github.com/NetPo4ki/pull-review/internal/service/stats"
	teamssvc "github.com/NetPo4ki/pull-review/internal/service/teams"
	userssvc "github.com/NetPo4ki/pull-review/internal/service/users"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(l *slog.Logger, teamsSvc *teamssvc.Service, usersSvc *userssvc.Service, prsSvc *prssvc.Service, statsSvc *statsvc.Service, pool *pgxpool.Pool) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.RealIP)
	r.Use(mymw.RequestID)
	r.Use(mymw.Logging(l))
	r.Use(mymw.Recoverer(l))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	r.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
		if pool == nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("no db"))
			return
		}
		if err := pool.Ping(r.Context()); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("db not ready"))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready"))
	})

	r.Mount("/team", teamsapi.NewHandler(teamsSvc).Routes())
	r.Mount("/users", usersapi.NewHandler(usersSvc).Routes())
	r.Mount("/pullRequest", prsapi.NewHandler(prsSvc).Routes())
	r.Mount("/stats", statsapi.NewHandler(statsSvc).Routes())

	return r
}
