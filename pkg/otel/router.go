package otel

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func Route(router chi.Router, method string, pattern string, fn http.HandlerFunc) {
	handler := otelhttp.WithRouteTag(pattern, fn)
	switch method {
	case http.MethodGet:
		router.Get(pattern, handler.ServeHTTP)
	case http.MethodPost:
		router.Post(pattern, handler.ServeHTTP)
	case http.MethodPut:
		router.Put(pattern, handler.ServeHTTP)
	case http.MethodPatch:
		router.Patch(pattern, handler.ServeHTTP)
	case http.MethodDelete:
		router.Delete(pattern, handler.ServeHTTP)
	}
}
