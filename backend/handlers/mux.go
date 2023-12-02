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
}

func NewMux(
	ctx context.Context, deps *MuxDependencies,
) (*chi.Mux, error) {
	router := chi.NewRouter()
	authMiddleWare := NewAuthorizationMiddleware(deps.JWTer)
	corsMiddleWare := NewCORSMiddleWare(deps.Config)
	router.Use(logging.WithLoggerMiddleware(*deps.Logger))

	router.Route("/", func(r chi.Router) {
		r.Use(corsMiddleWare)
		r.Route("/health", func(r chi.Router) {
			hh := &HealthCheckHandler{}
			r.Get("/", hh.ServeHTTP)
		})

		r.Route("/blogs", func(r chi.Router) {
			blh := &BlogListHandler{
				Service: deps.BlogService,
			}
			r.Get("/", blh.ServeHTTP)

			bgh := &BlogGetHandler{
				Service: deps.BlogService,
				jwter:   deps.JWTer,
			}
			r.Get("/{id}", bgh.ServeHTTP)

			// require login
			bah := &BlogAddHandler{
				Service:   deps.BlogService,
				Validator: deps.Validator,
			}
			r.With(authMiddleWare.Middleware).Post("/", bah.ServeHTTP)

			// require login
			bdh := &BlogDeleteHandler{
				Service:   deps.BlogService,
				Validator: deps.Validator,
			}
			r.With(authMiddleWare.Middleware).Delete("/", bdh.ServeHTTP)

			// require login
			buh := &BlogPutHandler{
				Service:   deps.BlogService,
				Validator: deps.Validator,
			}
			r.With(authMiddleWare.Middleware).Put("/", buh.ServeHTTP)
		})

		r.Route("/tags", func(r chi.Router) {
			th := &TagListHandler{
				Service: deps.BlogService,
			}
			r.Get("/", th.ServeHTTP)
		})

		r.Route("/files", func(r chi.Router) {
			// require login
			gt := GenerateThumbnailImageSignedURLHandler{
				StorageService: deps.StorageService,
				Validator:      deps.Validator,
			}
			r.With(authMiddleWare.Middleware).Post("/thumbnail/new", gt.ServeHTTP)

			gc := GenerateContentsImageSignedURLHandler{
				StorageService: deps.StorageService,
				Validator:      deps.Validator,
			}
			r.With(authMiddleWare.Middleware).Post("/content/new", gc.ServeHTTP)
		})

		r.Route("/auth", func(r chi.Router) {
			ah := NewAuthLoginHandler(deps.AuthService, deps.Validator, deps.Config)
			r.Post("/signin", ah.ServeHTTP)

			ash := NewAuthSessionLoginHandler(deps.AuthService, deps.Config)
			r.Get("/login/me", ash.ServeHTTP)

			alh := NewAuthLogoutHandler(deps.Config)
			r.Post("/admin/signout", alh.ServeHTTP)
		})

		r.Route("/admin", func(r chi.Router) {
			bla := &BlogListAdminHandler{
				Service: deps.BlogService,
			}
			r.With(authMiddleWare.Middleware).Get("/blogs", bla.ServeHTTP)
		})

	})
	return router, nil
}
