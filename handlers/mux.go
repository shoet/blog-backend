package handlers

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/shoet/blog/internal/config"
	"github.com/shoet/blog/internal/infrastracture"
	"github.com/shoet/blog/internal/interfaces"
	"github.com/shoet/blog/internal/interfaces/handler"
	"github.com/shoet/blog/internal/logging"
	"github.com/shoet/blog/services"
)

type MuxDependencies struct {
	Config         *config.Config
	DB             infrastracture.DB
	BlogService    *services.BlogService
	AuthService    *services.AuthService
	StorageService *services.AWSS3StorageService
	JWTer          *services.JWTManager
	Logger         *logging.Logger
	Validator      *validator.Validate
	Cookie         *interfaces.CookieManager
}

func NewMux(
	ctx context.Context, deps *MuxDependencies,
) (*chi.Mux, error) {
	router := chi.NewRouter()
	authMiddleWare := NewAuthorizationMiddleware(deps.JWTer)
	corsMiddleWare := NewCORSMiddleWare(deps.Config)
	router.Use(logging.WithLoggerMiddleware(deps.Logger), corsMiddleWare)

	setHealthRoute(router)
	setBlogsRoute(router, deps, authMiddleWare)
	setTagsRoute(router, deps)
	setFilesRoute(router, deps, authMiddleWare)
	setAuthRoute(router, deps)
	setAdminRoute(router, deps, authMiddleWare)
	return router, nil
}

func setHealthRoute(r chi.Router) {
	r.Route("/health", func(r chi.Router) {
		hh := &handler.HealthCheckHandler{}
		r.Get("/", hh.ServeHTTP)
	})
}

func setBlogsRoute(
	r chi.Router, deps *MuxDependencies, authMiddleWare *AuthorizationMiddleware,
) {
	r.Route("/blogs", func(r chi.Router) {
		blh := handler.NewBlogListHandler(deps.BlogService)
		r.Get("/", blh.ServeHTTP)

		bgh := handler.NewBlogGetHandler(deps.BlogService, deps.JWTer)
		r.Get("/{id}", bgh.ServeHTTP)

		bah := handler.NewBlogAddHandler(deps.BlogService, deps.Validator)
		r.With(authMiddleWare.Middleware).Post("/", bah.ServeHTTP)

		bdh := handler.NewBlogDeleteHandler(deps.BlogService, deps.Validator)
		r.With(authMiddleWare.Middleware).Delete("/", bdh.ServeHTTP)

		buh := handler.NewBlogPutHandler(deps.BlogService, deps.Validator)
		r.With(authMiddleWare.Middleware).Put("/", buh.ServeHTTP)
	})
}

func setTagsRoute(r chi.Router, deps *MuxDependencies) {
	r.Route("/tags", func(r chi.Router) {
		th := handler.NewTagListHandler(deps.BlogService)
		r.Get("/", th.ServeHTTP)
	})
}

func setFilesRoute(
	r chi.Router, deps *MuxDependencies, authMiddleWare *AuthorizationMiddleware,
) {
	r.Route("/files", func(r chi.Router) {
		gt := handler.NewGenerateThumbnailImageSignedURLHandler(deps.StorageService, deps.Validator)
		r.With(authMiddleWare.Middleware).Post("/thumbnail/new", gt.ServeHTTP)

		gc := handler.GenerateContentsImageSignedURLHandler{
			StorageService: deps.StorageService,
			Validator:      deps.Validator,
		}
		r.With(authMiddleWare.Middleware).Post("/content/new", gc.ServeHTTP)
	})
}

func setAuthRoute(r chi.Router, deps *MuxDependencies) {
	r.Route("/auth", func(r chi.Router) {
		ah := handler.NewAuthLoginHandler(deps.AuthService, deps.Validator, deps.Cookie)
		r.Post("/signin", ah.ServeHTTP)

		ash := handler.NewAuthSessionLoginHandler(deps.AuthService)
		r.Get("/login/me", ash.ServeHTTP)

		alh := handler.NewAuthLogoutHandler(deps.Cookie)
		r.Post("/admin/signout", alh.ServeHTTP)
	})
}

func setAdminRoute(
	r chi.Router, deps *MuxDependencies, authMiddleWare *AuthorizationMiddleware,
) {
	r.Route("/admin", func(r chi.Router) {
		bla := handler.NewBlogListAdminHandler(deps.BlogService)
		r.With(authMiddleWare.Middleware).Get("/blogs", bla.ServeHTTP)
	})
}
