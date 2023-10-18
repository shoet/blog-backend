package handlers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/shoet/blog/config"
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
	mux, err := NewMux(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create mux in NewServer(): %w", err)
	}
	srv := &http.Server{
		Handler: mux,
	}
	return &Server{srv: srv, l: l}, nil
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
