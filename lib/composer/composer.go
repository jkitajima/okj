package composer

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server interface {
	Mux() http.Handler
	Prefix() string
}

type Composer struct {
	servers []Server
	Mux     *chi.Mux
}

func NewComposer(middlewares ...func(http.Handler) http.Handler) *Composer {
	c := &Composer{Mux: chi.NewRouter()}
	c.Mux.Use(middlewares...)
	return c
}

func (c *Composer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.Mux.ServeHTTP(w, r)
}

func (c *Composer) Compose(servers ...Server) error {
	if len(c.servers) > 0 {
		return errors.New("composer: composer is already filled with servers")
	}

	for _, s := range servers {
		prefix := s.Prefix()
		if prefix == "" {
			return errors.New("composer: server prefix is empty")
		}

		mux := s.Mux()
		if mux == nil {
			return errors.New("composer: server Mux is nil")
		}

		c.Mux.Mount(prefix, mux)
	}

	return nil
}
