package interfaces

import (
	"context"
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/shoet/blog/internal/clocker"
	"github.com/shoet/blog/internal/config"
	"github.com/shoet/blog/internal/infrastracture"
	"github.com/shoet/blog/internal/infrastracture/adapter"
	"github.com/shoet/blog/internal/infrastracture/repository"
	"github.com/shoet/blog/internal/infrastracture/services/auth_service"
	"github.com/shoet/blog/internal/infrastracture/services/blog_service"
	"github.com/shoet/blog/internal/infrastracture/services/contents_service"
	"github.com/shoet/blog/internal/infrastracture/services/jwt_service"
	"github.com/shoet/blog/internal/interfaces/cookie"
	"github.com/shoet/blog/internal/interfaces/handler"
	"github.com/shoet/blog/internal/interfaces/middleware"
	"github.com/shoet/blog/internal/logging"
	"github.com/shoet/blog/internal/usecase/create_blog"
	"github.com/shoet/blog/internal/usecase/delete_blog"
	"github.com/shoet/blog/internal/usecase/get_blog_detail"
	"github.com/shoet/blog/internal/usecase/get_blogs"
	"github.com/shoet/blog/internal/usecase/get_blogs_offset_paging"
	"github.com/shoet/blog/internal/usecase/get_github_contributions"
	"github.com/shoet/blog/internal/usecase/get_github_contributions_latest_week"
	"github.com/shoet/blog/internal/usecase/get_tags"
	"github.com/shoet/blog/internal/usecase/login_user"
	"github.com/shoet/blog/internal/usecase/login_user_session"
	"github.com/shoet/blog/internal/usecase/put_blog"
	"github.com/shoet/blog/internal/usecase/storage_presigned_content"
	"github.com/shoet/blog/internal/usecase/storage_presigned_thumbnail"
)

type MuxDependencies struct {
	Config               *config.Config
	DB                   infrastracture.DB
	BlogRepository       *repository.BlogRepository
	BlogRepositoryOffset *repository.BlogRepositoryOffset
	BlogService          *blog_service.BlogService
	AuthService          *auth_service.AuthService
	ContentsService      *contents_service.ContentsService
	JWTer                *jwt_service.JWTService
	Logger               *logging.Logger
	Validator            *validator.Validate
	Cookie               *cookie.CookieController
	GitHubAPIAdapter     *adapter.GitHubV4APIClient
	Clocker              clocker.Clocker
}

func NewMux(
	ctx context.Context, deps *MuxDependencies,
) (*chi.Mux, error) {
	log.Printf("set middleware")
	router := chi.NewRouter()
	authMiddleWare := middleware.NewAuthorizationMiddleware(deps.JWTer)
	corsMiddleWare := middleware.NewCORSMiddleWare(deps.Config)
	router.Use(logging.WithLoggerMiddleware(deps.Logger), corsMiddleWare)

	log.Printf("set routes")
	setHealthRoute(router)
	setBlogsRoute(router, deps, authMiddleWare)
	setTagsRoute(router, deps)
	setFilesRoute(router, deps, authMiddleWare)
	setAuthRoute(router, deps)
	setAdminRoute(router, deps, authMiddleWare)
	setGitHubRoute(router, deps)
	return router, nil
}

// health check
func setHealthRoute(r chi.Router) {
	r.Route("/health", func(r chi.Router) {
		hh := &handler.HealthCheckHandler{}
		r.Get("/", hh.ServeHTTP)
	})
}

// blogs
func setBlogsRoute(
	r chi.Router, deps *MuxDependencies, authMiddleWare *middleware.AuthorizationMiddleware,
) {
	r.Route("/blogs", func(r chi.Router) {
		blh := handler.NewBlogListHandler(get_blogs.NewUsecase(deps.DB, deps.BlogRepository))
		r.Get("/", blh.ServeHTTP)

		bah := handler.NewBlogAddHandler(
			create_blog.NewUsecase(deps.DB, deps.BlogRepository, deps.BlogService),
			deps.Validator)
		r.With(authMiddleWare.Middleware).Post("/", bah.ServeHTTP)

		bgh := handler.NewBlogGetHandler(
			get_blog_detail.NewUsecase(deps.DB, deps.BlogRepository), deps.JWTer)
		r.Get("/{id}", bgh.ServeHTTP)

		bdh := handler.NewBlogDeleteHandler(
			delete_blog.NewUsecase(deps.DB, deps.BlogRepository), deps.Validator)
		r.With(authMiddleWare.Middleware).Delete("/{id}", bdh.ServeHTTP)

		buh := handler.NewBlogPutHandler(
			put_blog.NewUsecase(deps.DB, deps.BlogRepository), deps.Validator)
		r.With(authMiddleWare.Middleware).Put("/{id}", buh.ServeHTTP)
	})

	r.Route("/v2/blogs", func(r chi.Router) {
		blh := handler.NewBlogGetOffsetPagingHandler(
			get_blogs_offset_paging.NewUsecase(deps.DB, deps.BlogRepositoryOffset),
		)
		r.Get("/", blh.ServeHTTP)
	})
}

// tags
func setTagsRoute(r chi.Router, deps *MuxDependencies) {
	r.Route("/tags", func(r chi.Router) {
		th := handler.NewTagListHandler(*get_tags.NewUsecase(deps.DB, deps.BlogRepository))
		r.Get("/", th.ServeHTTP)
	})
}

// files
func setFilesRoute(
	r chi.Router, deps *MuxDependencies, authMiddleWare *middleware.AuthorizationMiddleware,
) {
	r.Route("/files", func(r chi.Router) {
		gt := handler.NewGenerateThumbnailImageSignedURLHandler(
			storage_presigned_thumbnail.NewUsecase(deps.ContentsService),
			deps.Validator)
		r.With(authMiddleWare.Middleware).Post("/thumbnail/new", gt.ServeHTTP)

		gc := handler.NewGenerateContentsImageSignedURLHandler(
			storage_presigned_content.NewUsecase(deps.ContentsService),
			deps.Validator)
		r.With(authMiddleWare.Middleware).Post("/content/new", gc.ServeHTTP)
	})
}

// auth
func setAuthRoute(r chi.Router, deps *MuxDependencies) {
	r.Route("/auth", func(r chi.Router) {
		ah := handler.NewAuthLoginHandler(
			login_user.NewUsecase(deps.AuthService),
			deps.Validator,
			deps.Cookie)
		r.Post("/signin", ah.ServeHTTP)

		ash := handler.NewAuthSessionLoginHandler(login_user_session.NewUsecase(deps.AuthService))
		r.Get("/signin/me", ash.ServeHTTP)

		alh := handler.NewAuthLogoutHandler(deps.Cookie)
		r.Post("/signout", alh.ServeHTTP)
	})
}

// admin
func setAdminRoute(
	r chi.Router, deps *MuxDependencies, authMiddleWare *middleware.AuthorizationMiddleware,
) {
	r.Route("/admin", func(r chi.Router) {
		bla := handler.NewBlogListAdminHandler(get_blogs.NewUsecase(deps.DB, deps.BlogRepository))
		r.With(authMiddleWare.Middleware).Get("/blogs", bla.ServeHTTP)
	})
}

// github
func setGitHubRoute(
	r chi.Router, deps *MuxDependencies,
) {
	r.Route("/github", func(r chi.Router) {
		ghgch := handler.NewGitHubGetContributionsHandler(
			get_github_contributions.NewUsecase(deps.GitHubAPIAdapter),
		)
		r.Get("/contributions", ghgch.ServeHTTP)

		ghgchw := handler.NewGitHubGetContributionsLatestWeekHandler(
			get_github_contributions_latest_week.NewUsecase(deps.GitHubAPIAdapter, deps.Clocker),
		)
		r.Get("/contributions_latest_week", ghgchw.ServeHTTP)

	})
}
