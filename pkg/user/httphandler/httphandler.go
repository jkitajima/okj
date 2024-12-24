package httphandler

import (
	"log/slog"
	"net/http"

	"okj/lib/composer"
	"okj/pkg/user"
	repo "okj/pkg/user/repo/gorm"

	"github.com/go-playground/validator/v10"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"gorm.io/gorm"
)

type UserServer struct {
	entity         string
	mux            *chi.Mux
	prefix         string
	service        *user.Service
	auth           *jwtauth.JWTAuth
	db             user.Repoer
	inputValidator *validator.Validate
	logger         *slog.Logger
}

func (s *UserServer) Prefix() string {
	return s.prefix
}

func (s *UserServer) Mux() http.Handler {
	return s.mux
}

func (s *UserServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func NewServer(auth *jwtauth.JWTAuth, db *gorm.DB, validtr *validator.Validate, logger *slog.Logger) composer.Server {
	s := &UserServer{
		entity:         "users",
		prefix:         "/users",
		mux:            chi.NewRouter(),
		auth:           auth,
		db:             repo.NewRepo(db, logger),
		inputValidator: validtr,
		logger:         logger,
	}

	s.service = &user.Service{Repo: s.db}
	s.addRoutes()
	return s
}
