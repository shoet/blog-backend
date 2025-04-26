package interfaces

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

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
	"github.com/shoet/blog/internal/logging"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	srv *http.Server
	l   net.Listener
}

func NewServer(ctx context.Context, cfg *config.Config) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.AppPort))
	if err != nil {
		return nil, fmt.Errorf("failed to create listener in NewServer(): %w", err)
	}
	log.Printf("server listening on %s", l.Addr().String())
	deps, err := BuildMuxDependencies(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to build mux dependencies in NewServer(): %w", err)
	}
	mux, err := NewMux(ctx, deps)
	if err != nil {
		return nil, fmt.Errorf("failed to create mux in NewServer(): %w", err)
	}
	srv := &http.Server{
		Handler: mux,
	}
	return &Server{srv: srv, l: l}, nil
}

func BuildMuxDependencies(ctx context.Context, cfg *config.Config) (*MuxDependencies, error) {
	logger := logging.NewLogger(os.Stdout, cfg.LogLevel)
	validator := validator.New()
	cookie := cookie.NewCookieController(cfg.Env, cfg.SiteDomain)

	log.Println("start connection DB")
	db, err := infrastracture.NewDBPostgres(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create db: %w", err)
	}
	log.Println("start connection KVS")
	kvs, err := infrastracture.NewRedisKVS(
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
	c := clocker.RealClocker{}
	jwtService := jwt_service.NewJWTService(kvs, &c, []byte(cfg.JWTSecret), cfg.JWTExpiresInSec)

	blogRepo := repository.NewBlogRepository(&c)
	blogOffsetRepo := repository.NewBlogRepositoryOffset(&c)
	blogService := blog_service.NewBlogService()

	userRepo, err := repository.NewUserRepository(&c)
	if err != nil {
		return nil, fmt.Errorf("failed to create user repository: %w", err)
	}

	commentRepo := repository.NewCommentRepository(&c)
	userProfileRepo := repository.NewUserProfileRepository()

	authService, err := auth_service.NewAuthService(db, userRepo, jwtService)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth service: %w", err)
	}

	s3Adapter, err := adapter.NewS3Adapter(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create s3 adapter: %w", err)
	}
	fileRepo := repository.NewFileRepository(cfg, s3Adapter)

	contentsService, err := contents_service.NewContentsService(s3Adapter, cfg.AWSS3ThumbnailDirectory, cfg.AWSSS3ContentImageDirectory)
	if err != nil {
		return nil, fmt.Errorf("failed to create contents service: %w", err)
	}

	gitHubAPIAdapter := adapter.NewGitHubV4APIClient(cfg.GitHubPersonalAccessToken)

	return &MuxDependencies{
		Config:                cfg,
		DB:                    db,
		BlogRepository:        blogRepo,
		BlogRepositoryOffset:  blogOffsetRepo,
		CommentRepository:     commentRepo,
		FileRepository:        fileRepo,
		UserProfileRepository: userProfileRepo,
		BlogService:           blogService,
		AuthService:           authService,
		ContentsService:       contentsService,
		JWTer:                 jwtService,
		Logger:                logger,
		Validator:             validator,
		Cookie:                cookie,
		GitHubAPIAdapter:      gitHubAPIAdapter,
		Clocker:               &c,
	}, nil
}

func (s *Server) Run(ctx context.Context) error {
	fmt.Println("server start")
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		if err := s.srv.Serve(s.l); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("failed to server in Run(): %w", err)
		}
		return nil
	})

	<-ctx.Done()

	if err := s.srv.Shutdown(context.Background()); err != nil {
		stop()
		return fmt.Errorf("failed to shutdown server in Run(): %w", err)
	}
	log.Println("server shutdowned")

	return eg.Wait()
}
