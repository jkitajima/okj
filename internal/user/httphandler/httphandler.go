package httphandler

import (
	"net/http"

	"okj/internal/user"
	repo "okj/internal/user/repo/gorm"
	"okj/pkg/composer"

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

func NewServer(auth *jwtauth.JWTAuth, db *gorm.DB, validtr *validator.Validate) composer.Server {
	s := &UserServer{
		entity:         "users",
		prefix:         "/users",
		mux:            chi.NewRouter(),
		auth:           auth,
		db:             repo.NewRepo(db),
		inputValidator: validtr,
	}

	s.service = &user.Service{Repo: s.db}
	s.addRoutes()
	return s
}
