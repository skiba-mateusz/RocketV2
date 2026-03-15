package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/skiba-mateusz/RocketV2/config"
	"github.com/skiba-mateusz/RocketV2/logger"
)

type Server struct {
	logger *logger.Logger
	config *config.Config
	server *http.Server
}

func New(logger *logger.Logger, config *config.Config, port string) *Server {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(config.BuildDir)))

	return &Server{
		logger: logger,
		config: config,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%s", port),
			Handler: mux,
		},
	}
}

func (s *Server) Run(ctx context.Context) error {
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

func (s *Server) Shutdown() error {
	s.logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown error: %v", err)
	}

	return nil
}