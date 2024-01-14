package interfaces

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/shoet/blog/internal/config"
	"github.com/shoet/blog/internal/infrastracture"
	"github.com/shoet/blog/internal/infrastracture/repository"
	"github.com/shoet/blog/internal/infrastracture/services"
	"github.com/shoet/blog/internal/interfaces/cookie"
	"github.com/shoet/blog/internal/interfaces/handler"
	"github.com/shoet/blog/internal/interfaces/middleware"
	"github.com/shoet/blog/internal/logging"
	"github.com/shoet/blog/internal/usecase/create_blog"
	"github.com/shoet/blog/internal/usecase/get_blog_detail"
	"github.com/shoet/blog/internal/usecase/get_blogs"
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
	Cookie         *cookie.CookieController
	BlogRepository *repository.BlogRepository
}

func NewMux(
	ctx context.Context, deps *MuxDependencies,
) (*chi.Mux, error) {
	router := chi.NewRouter()
	authMiddleWare := middleware.NewAuthorizationMiddleware(deps.JWTer)
	corsMiddleWare := middleware.NewCORSMiddleWare(deps.Config)
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
	r chi.Router, deps *MuxDependencies, authMiddleWare *middleware.AuthorizationMiddleware,
) {
	r.Route("/blogs", func(r chi.Router) {
		blh := handler.NewBlogListHandler(get_blogs.NewUsecase(deps.DB, deps.BlogRepository))
		r.Get("/", blh.ServeHTTP)

		bgh := handler.NewBlogGetHandler(
			get_blog_detail.NewUsecase(deps.DB, deps.BlogRepository), deps.JWTer)
		r.Get("/{id}", bgh.ServeHTTP)

		bah := handler.NewBlogAddHandler(
			create_blog.NewUsecase(deps.DB, deps.BlogRepository, deps.BlogService),
			deps.Validator)
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
	r chi.Router, deps *MuxDependencies, authMiddleWare *middleware.AuthorizationMiddleware,
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
	r chi.Router, deps *MuxDependencies, authMiddleWare *middleware.AuthorizationMiddleware,
) {
	r.Route("/admin", func(r chi.Router) {
		bla := handler.NewBlogListAdminHandler(get_blogs.NewUsecase(deps.DB, deps.BlogRepository))
		r.With(authMiddleWare.Middleware).Get("/blogs", bla.ServeHTTP)
	})
}
