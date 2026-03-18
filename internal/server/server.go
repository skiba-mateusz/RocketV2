package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/skiba-mateusz/RocketV2/internal/config"
	"github.com/skiba-mateusz/RocketV2/pkg/logger"
)

type Server interface {
	Run(ctx context.Context) error
}

type DefaultServer struct {
	logger *logger.Logger
	config *config.Config
	server *http.Server
}

func NewDefault(logger *logger.Logger, config *config.Config, port string) *DefaultServer {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(config.BuildDir)))

	return &DefaultServer{
		logger: logger,
		config: config,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%s", port),
			Handler: mux,
		},
	}
}

func (s *DefaultServer) Run(ctx context.Context) error {
	s.logger.Info("Server is listening on http://localhost%s", s.server.Addr)

	errChan := make(chan error, 1)
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
		close(errChan)
	}()

	select {
	case err := <- errChan:
		return fmt.Errorf("server error: %v", err)
	case <-ctx.Done():
		return s.Shutdown()
	}
}

func (s *DefaultServer) Shutdown() error {
	s.logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown error: %v", err)
	}

	return nil
}