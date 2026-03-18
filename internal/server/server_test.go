package server

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/skiba-mateusz/RocketV2/internal/config"
	"github.com/skiba-mateusz/RocketV2/pkg/logger"
)

func setupServer(t *testing.T, port string) (*DefaultServer, string) {
	t.Helper()

	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "index.html"), []byte("<h1>test</h1>"), 0644)

	cfg := &config.Config{BuildDir: dir}
	log := logger.NewDefault(false)

	return NewDefault(log, cfg, port), dir
}

func TestServer_ServeFiles(t *testing.T) {
	srv, _ := setupServer(t, "8175")

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	go func() {
		srv.Run(ctx)
	}()
	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:8175/index.html")
	if err != nil {
		t.Fatalf("GET failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected %v, got %v", http.StatusOK, resp.StatusCode)
	}
}

func TestServer_ShutdownOnCancel(t *testing.T) {
	srv, _ := setupServer(t, "8176")

	ctx, cancel := context.WithCancel(t.Context())
	done := make(chan error, 1)

	go func() {
		done <- srv.Run(ctx)
	}()
	time.Sleep(100 * time.Millisecond)

	cancel()

	select {
	case err := <-done: 
		if err != nil {
			t.Fatalf("expected clean shutdown, got: %v", err)
		}
	case <-time.After(3*time.Second):
		t.Fatal("Server didn't shutdown within 3s")
	}
}