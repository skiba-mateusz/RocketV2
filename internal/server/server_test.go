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

type mockBuilder struct {

}

func (b *mockBuilder) Build(ctx context.Context) error {
	return nil
}

func setupServer(t *testing.T, port string) (*DefaultServer, string) {
	t.Helper()

	dir := t.TempDir()
	
	os.WriteFile(filepath.Join(dir, "index.html"), []byte("<h1>test</h1>"), 0644)
	os.WriteFile(filepath.Join(dir, "config.yaml"), nil, 0644)

	contentDir := filepath.Join(dir, "content")
	layoutDir := filepath.Join(dir, "content")
	staticDir := filepath.Join(dir, "content")

	os.MkdirAll(contentDir, 0755)
	os.MkdirAll(layoutDir, 0755)
	os.MkdirAll(staticDir, 0755)

	cfg := &config.Config{
		BuildDir: dir,
		ContentDir: contentDir,
		LayoutDir: layoutDir,
		StaticDir: staticDir,
	}
	log := logger.NewDefault(false)

	return NewDefault(log, cfg, &mockBuilder{}, port, ""), dir
}

func TestServer_ServeFiles(t *testing.T) {
	srv, _ := setupServer(t, "8175")

	ctx, cancel := context.WithCancel(t.Context())
	done := make(chan error, 1)

	go func() {
		done <- srv.Run(ctx)
	}()
	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:8175/")
	if err != nil {
		t.Fatalf("GET failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected %v, got %v", http.StatusOK, resp.StatusCode)
	}

	cancel()
	<-done
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
		time.Sleep(50 * time.Millisecond)
	case <-time.After(3*time.Second):
		t.Fatal("Server didn't shutdown within 3s")
	}
}