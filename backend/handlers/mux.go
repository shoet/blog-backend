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

		r.Route("/tags", func(r chi.Router) {
			th := &TagListHandler{}
			r.Get("/", th.ServeHTTP)
		})

		r.Route("/blogs", func(r chi.Router) {
			blh := &BlogListHandler{}
			r.Get("/", blh.ServeHTTP)

			bgh := &BlogGetHandler{}
			r.Get("/{id}", bgh.ServeHTTP)
		})
	})
	return router
}
