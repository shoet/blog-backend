package handlers

import (
	"github.com/go-chi/chi/v5"
)

func NewMux() *chi.Mux {
	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		r.Use(CORSMiddleWare)
		r.Route("/health", func(r chi.Router) {
			hh := &HealthCheckHandler{}
			r.Get("/", hh.ServeHTTP)
		})
	})
	return router
}
