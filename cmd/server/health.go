package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"okj/pkg/composer"

	"github.com/alexliesenfeld/health"
	healthPsql "github.com/hellofresh/health-go/v5/checks/postgres"

	"github.com/go-chi/chi/v5"
)

type HealthServer struct {
	mux    *chi.Mux
	prefix string
}

func (s *HealthServer) Prefix() string {
	return s.prefix
}

func (s *HealthServer) Mux() http.Handler {
	return s.mux
}

func (s *HealthServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func SetupHealthCheck(cfg *Config) composer.Server {
	s := &HealthServer{
		prefix: "/healthz",
		mux:    chi.NewRouter(),
	}

	checker := health.NewChecker(
		health.WithCacheDuration(time.Duration(cfg.Server.Health.Cache)*time.Second),
		health.WithTimeout(time.Duration(cfg.Server.Health.Timeout)*time.Second),
		health.WithPeriodicCheck(
			time.Duration(cfg.Server.Health.Interval)*time.Second,
			time.Duration(cfg.Server.Health.Delay)*time.Second,
			health.Check{
				Name: "db",
				Check: healthPsql.New(healthPsql.Config{
					DSN: cfg.DB.DSN,
				}),
				MaxContiguousFails: uint(cfg.Server.Health.Retries),
			}),
		health.WithStatusListener(func(ctx context.Context, state health.CheckerState) {
			log.Printf("health status changed to %q", state.Status)
		}),
	)
	s.mux.Get("/readiness", health.NewHandler(checker))
	return s
}
