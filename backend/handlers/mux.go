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
	"github.com/shoet/blog/util"
)

func NewMux(ctx context.Context, cfg *config.Config) (*chi.Mux, error) {
	db, err := store.NewDBMySQL(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create db: %w", err)
	}
	c := clocker.RealClocker{}
	blogRepo := store.BlogRepository{
		Clocker: &c,
	}
	blogService := services.NewBlogService(db, &blogRepo)
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

	userRepo := store.UserRepository{
		Clocker: &c,
	}
	kvs, err := store.NewRedisKVS(
		ctx,
		cfg.KVSHost,
		cfg.KVSPort,
		cfg.KVSUser,
		cfg.KVSPass,
		cfg.JWTExpiresInSec,
		cfg.KVSTlsEnabled,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create redis kvs: %w", err)
	}
	jwter := util.NewJWTer(kvs, cfg, &c)
	authService, err := services.NewAuthService(db, &userRepo, jwter)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth service: %w", err)
	}

	router := chi.NewRouter()
	authMiddleWare := NewAuthorizationMiddleware(jwter)
	corsMiddleWare := NewCORSMiddleWare(cfg)
	router.Route("/", func(r chi.Router) {
		r.Use(WithLoggerMiddleware(logger))
		r.Use(corsMiddleWare)
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

			// require login
			bah := &BlogAddHandler{
				Service:   blogService,
				Validator: validate,
			}
			r.With(authMiddleWare.Middleware).Post("/", bah.ServeHTTP)

			// require login
			bdh := &BlogDeleteHandler{
				Service:   blogService,
				Validator: validate,
			}
			r.With(authMiddleWare.Middleware).Delete("/", bdh.ServeHTTP)

			// require login
			buh := &BlogPutHandler{
				Service:   blogService,
				Validator: validate,
			}
			r.With(authMiddleWare.Middleware).Put("/", buh.ServeHTTP)
		})

		r.Route("/tags", func(r chi.Router) {
			// TODO: implement
			th := &TagListHandler{}
			r.Get("/", th.ServeHTTP)
		})

		r.Route("/files", func(r chi.Router) {
			// require login
			s := GenerateSignedURLHandler{
				StorageService: awsStorage,
				Validator:      validate,
			}
			r.With(authMiddleWare.Middleware).Post("/thumbnail/new", s.ServeHTTP)
		})

		r.Route("/auth", func(r chi.Router) {
			ah := NewAuthLoginHandler(authService, validate, cfg)
			r.Post("/signin", ah.ServeHTTP)

			ash := NewAuthSessionLoginHandler(authService, cfg)
			r.Get("/login/me", ash.ServeHTTP)

			alh := NewAuthLogoutHandler(cfg)
			r.Post("/admin/signout", alh.ServeHTTP)
		})

	})
	return router, nil
}
