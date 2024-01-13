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
	"github.com/shoet/blog/internal/infrastracture/repository"
	"github.com/shoet/blog/internal/interfaces/cookie"
	"github.com/shoet/blog/internal/logging"
	"github.com/shoet/blog/services"
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

	db, err := infrastracture.NewDBMySQL(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create db: %w", err)
	}
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
	jwter := services.NewJWTManager(kvs, &c, []byte(cfg.JWTSecret), cfg.JWTExpiresInSec)

	blogRepo := repository.NewBlogRepository(&c)
	blogService := services.NewBlogService(db, blogRepo)

	userRepo, err := repository.NewUserRepository(&c)
	if err != nil {
		return nil, fmt.Errorf("failed to create user repository: %w", err)
	}

	authService, err := services.NewAuthService(db, userRepo, jwter)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth service: %w", err)
	}
	awsStorage, err := services.NewAWSS3StorageService(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create aws storage: %w", err)
	}

	return &MuxDependencies{
		Config:         cfg,
		DB:             db,
		BlogService:    blogService,
		AuthService:    authService,
		StorageService: awsStorage,
		JWTer:          jwter,
		Logger:         logger,
		Validator:      validator,
		Cookie:         cookie,
	}, nil
}

func (s *Server) Run(ctx context.Context) error {
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
