package server

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/skiba-mateusz/RocketV2/internal/builder"
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
	builder builder.Builder
	ctx context.Context
	cancel context.CancelFunc
	cfgPath string
	reload chan struct{}
	mu 	   sync.Mutex
}

func NewDefault(logger *logger.Logger, config *config.Config, builder builder.Builder, port, cfgPath string) *DefaultServer {
	ctx, cancel := context.WithCancel(context.Background())

	srv := &DefaultServer{
		logger: logger,
		config: config,
		builder: builder,
		ctx: ctx,
		cancel: cancel,
		cfgPath: cfgPath,
		reload: make(chan struct{}),
	}

	fs := http.FileServer(http.Dir(config.BuildDir))

	mux := http.NewServeMux()
	mux.HandleFunc("/events", srv.sseHandler)
	mux.HandleFunc("/", srv.fsHandler(fs))

	srv.server =  &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}

	return srv
}

func (s *DefaultServer) Run(ctx context.Context) error {
	s.logger.Info("Server is listening on http://localhost%s", s.server.Addr)

	go s.watch(ctx)

	errChan := make(chan error, 1)
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
		close(errChan)
	}()

	select {
	case err := <-errChan:
		return fmt.Errorf("server error: %v", err)
	case <-ctx.Done():
		return s.Shutdown()
	}
}

func (s *DefaultServer) Shutdown() error {
	s.logger.Info("Shutting down server...")
	s.cancel()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown error: %v", err)
	}

	return nil
}

func (s *DefaultServer) watch(ctx context.Context) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		s.logger.Error("faild to create watcher")
		return
	}
	defer watcher.Close()

	paths := []string{
		s.config.ContentDir,
		s.config.LayoutDir,
		s.config.StaticDir,
	}

	for _, p := range paths {
		if err := s.handleWatcherPath(watcher, p); err != nil {
			s.logger.Error("failed to watch %s: %v", p, err)
		}
	}

	if s.cfgPath != "" {
		if err := watcher.Add(s.cfgPath); err != nil {
			s.logger.Error("failed to watch %s: %v", s.cfgPath, err)
		}
	}


	debounce := time.NewTimer(0)
	<-debounce.C

	var lastEvent fsnotify.Event

	for {
		select {
		case event := <-watcher.Events:
			if event.Has(fsnotify.Create) || event.Has(fsnotify.Write) || event.Has(fsnotify.Remove) {
				lastEvent = event
				debounce.Reset(time.Millisecond * 300)
			}
		case <-debounce.C:
			if filepath.Base(lastEvent.Name) == "config.yaml" {
				s.logger.Info("config changed, restart server to apply changes")
				continue
			}
			s.logger.Info("file changed (%s), rebuilding...", lastEvent.Name)
			if err := s.builder.Build(ctx); err != nil {
				s.logger.Error("rebuild failed: %v", err)
			} else {
				s.notifyReload()
			}
		case err := <-watcher.Errors:
			s.logger.Error("watcher error: %v", err)
		case <-ctx.Done():
			return
		}
	}

}

func (s *DefaultServer) handleWatcherPath(watcher *fsnotify.Watcher, root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return watcher.Add(path)
		}

		return nil
	})
}

func (s *DefaultServer) fsHandler(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasPrefix(path, "/static") || strings.HasPrefix(path, "/.well-known") {
			fs.ServeHTTP(w, r)
			return 
		}

		data, err := os.ReadFile(filepath.Join(s.config.BuildDir, path, "index.html"))
		if err != nil {
			http.NotFound(w, r)
			return 
		}
		script := `
		<script>
			window.addEventListener("load", () => {
				const es = new EventSource("/events");
				es.addEventListener("reload", () => {
						es.close();
						location.reload();
				});
				window.addEventListener("beforeunload", () => es.close());
			});
		</script>
		</body>
		`

		injected := strings.Replace(string(data), "</body>", script, 1)

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, injected)
	}
}

func (s *DefaultServer) sseHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")

	s.logger.Debug("SSE client connceted")
	defer s.logger.Debug("SSE clint disconnected")

	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()	

	for {
		s.mu.Lock()
		ch := s.reload
		s.mu.Unlock()

		select {
		case <-r.Context().Done():
			return
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			_, err := fmt.Fprintf(w, ": heartbeat\n\n")
			if err != nil {
				return
			}
			flusher.Flush()
		case <-ch:
			s.logger.Debug("event: reload")
			fmt.Fprintf(w, "event: reload\ndata: reload\n\n")
			flusher.Flush()
		}
	}
}

func (s *DefaultServer) notifyReload() {
	s.mu.Lock()
	defer s.mu.Unlock()
	close(s.reload)
	s.reload = make(chan struct{})
}