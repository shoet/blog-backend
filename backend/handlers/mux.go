package handlers

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/shoet/blog/config"
	"github.com/shoet/blog/logging"
	"github.com/shoet/blog/services"
	"github.com/shoet/blog/store"
	"github.com/shoet/blog/util"
)

type MuxDependencies struct {
	Config         *config.Config
	DB             store.DB
	BlogService    *services.BlogService
	AuthService    *services.AuthService
	StorageService *services.AWSS3StorageService
	JWTer          *util.JWTer
	Logger         *logging.Logger
	Validator      *validator.Validate
	Cookie         *CookieManager
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
		hh := &HealthCheckHandler{}
		r.Get("/", hh.ServeHTTP)
	})
}

func setBlogsRoute(
	r chi.Router, deps *MuxDependencies, authMiddleWare *AuthorizationMiddleware,
) {
	r.Route("/blogs", func(r chi.Router) {
		blh := NewBlogListHandler(deps.BlogService)
		r.Get("/", blh.ServeHTTP)

		bgh := NewBlogGetHandler(deps.BlogService, deps.JWTer)
		r.Get("/{id}", bgh.ServeHTTP)

		bah := NewBlogAddHandler(deps.BlogService, deps.Validator)
		r.With(authMiddleWare.Middleware).Post("/", bah.ServeHTTP)

		bdh := NewBlogDeleteHandler(deps.BlogService, deps.Validator)
		r.With(authMiddleWare.Middleware).Delete("/", bdh.ServeHTTP)

		buh := NewBlogPutHandler(deps.BlogService, deps.Validator)
		r.With(authMiddleWare.Middleware).Put("/", buh.ServeHTTP)
	})
}

func setTagsRoute(r chi.Router, deps *MuxDependencies) {
	r.Route("/tags", func(r chi.Router) {
		th := NewTagListHandler(deps.BlogService)
		r.Get("/", th.ServeHTTP)
	})
}

func setFilesRoute(
	r chi.Router, deps *MuxDependencies, authMiddleWare *AuthorizationMiddleware,
) {
	r.Route("/files", func(r chi.Router) {
		gt := NewGenerateThumbnailImageSignedURLHandler(deps.StorageService, deps.Validator)
		r.With(authMiddleWare.Middleware).Post("/thumbnail/new", gt.ServeHTTP)

		gc := GenerateContentsImageSignedURLHandler{
			StorageService: deps.StorageService,
			Validator:      deps.Validator,
		}
		r.With(authMiddleWare.Middleware).Post("/content/new", gc.ServeHTTP)
	})
}

func setAuthRoute(r chi.Router, deps *MuxDependencies) {
	r.Route("/auth", func(r chi.Router) {
		ah := NewAuthLoginHandler(deps.AuthService, deps.Validator, deps.Cookie)
		r.Post("/signin", ah.ServeHTTP)

		ash := NewAuthSessionLoginHandler(deps.AuthService)
		r.Get("/login/me", ash.ServeHTTP)

		alh := NewAuthLogoutHandler(deps.Cookie)
		r.Post("/admin/signout", alh.ServeHTTP)
	})
}

func setAdminRoute(
	r chi.Router, deps *MuxDependencies, authMiddleWare *AuthorizationMiddleware,
) {
	r.Route("/admin", func(r chi.Router) {
		bla := &BlogListAdminHandler{
			Service: deps.BlogService,
		}
		r.With(authMiddleWare.Middleware).Get("/blogs", bla.ServeHTTP)
	})
}
