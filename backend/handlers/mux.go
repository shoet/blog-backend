package handlers

import (
	"context"
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/shoet/blog/clocker"
	"github.com/shoet/blog/services"
	"github.com/shoet/blog/store"
)

func NewMux(ctx context.Context) (*chi.Mux, error) {
	db, err := store.NewDBSQLite3(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create db: %w", err)
	}
	repo := store.BlogRepository{
		Clocker: &clocker.RealClocker{},
	}
	blogService := services.NewBlogService(db, &repo)
	validate := validator.New()

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		r.Use(CORSMiddleWare)
		r.Route("/health", func(r chi.Router) {
			hh := &HealthCheckHandler{}
			r.Get("/", hh.ServeHTTP)
		})

		r.Route("/blogs", func(r chi.Router) {
			blh := &BlogListHandler{
				Service: blogService,
			}
			r.Get("/", blh.ServeHTTP)

			bgh := &BlogGetHandler{
				Service: blogService,
			}
			r.Get("/{id}", bgh.ServeHTTP)

			bah := &BlogAddHandler{
				Service:   blogService,
				Validator: validate,
			}
			r.Post("/", bah.ServeHTTP)
		})

		r.Route("/tags", func(r chi.Router) {
			th := &TagListHandler{}
			r.Get("/", th.ServeHTTP)
		})

	})
	return router, nil
}
