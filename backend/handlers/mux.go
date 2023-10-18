package handlers

import (
	"context"
	"fmt"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"github.com/shoet/blog/clocker"
	"github.com/shoet/blog/config"
	"github.com/shoet/blog/services"
	"github.com/shoet/blog/store"
)

func NewMux(ctx context.Context, cfg *config.Config) (*chi.Mux, error) {
	db, err := store.NewDBSQLite3(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create db: %w", err)
	}
	repo := store.BlogRepository{
		Clocker: &clocker.RealClocker{},
	}
	blogService := services.NewBlogService(db, &repo)
	validate := validator.New()

	logger := zerolog.
		New(os.Stdout).
		With().
		Timestamp().
		Logger()

	awsStorage, err := services.NewAWSS3StorageService(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create aws storage: %w", err)
	}

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		r.Use(WithLoggerMiddleware(logger))
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

			bdh := &BlogDeleteHandler{
				Service:   blogService,
				Validator: validate,
			}
			r.Post("/delete", bdh.ServeHTTP)

		})

		r.Route("/tags", func(r chi.Router) {
			th := &TagListHandler{}
			r.Get("/", th.ServeHTTP)
		})

		r.Route("/files", func(r chi.Router) {
			s := GenerateSignedURLHandler{
				StorageService: awsStorage,
				Validator:      validate,
			}
			r.Post("/thumbnail/new", s.ServeHTTP)
		})

	})
	return router, nil
}
