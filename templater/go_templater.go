package templater

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/skiba-mateusz/RocketV2/config"
)

type GoTemplater struct{
	config *config.Config
	base *template.Template
	cache map[string]*template.Template
	mu sync.RWMutex
}

func NewGoTemplater(config *config.Config) (*GoTemplater, error) {
	funcs := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
	}

	var files []string

	patterns := []string{
		filepath.Join(config.LayoutDir, "partials", "*.html"),
		filepath.Join(config.LayoutDir, "default", "baseof.html"),
	}

	for _, p := range patterns {
		matches, err := filepath.Glob(p)
		if err != nil {
			return nil, err
		}
		files = append(files, matches...)
	}

	base := template.New("base").Funcs(funcs)
	base , err := base.ParseFiles(files...)
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %v", err)
	}

	return &GoTemplater{
		config: config,
		base: base,
		cache: make(map[string]*template.Template),
	}, nil
}

func (t *GoTemplater) Render(w io.Writer, layout string, data any, templates []string) error {
	cacheKey := strings.Join(templates, "|")

	t.mu.RLock()
	tmpl, exists := t.cache[cacheKey]
	t.mu.RUnlock()

	if !exists {
		t.mu.Lock()

		if cached, ok := t.cache[cacheKey]; ok {
			tmpl = cached
		} else {
			var err error
			tmpl, err = t.base.Clone()
			if err != nil {
				t.mu.Unlock()
				return fmt.Errorf("failed to clone template: %v", err)
			}

			for _, path := range templates {
				templatePath := path
				if _, err = os.Stat(templatePath); os.IsNotExist(err) {
					templatePath = filepath.Join(t.config.LayoutDir, "default", filepath.Base(path))
				}

				tmpl, err = tmpl.ParseFiles(templatePath)
				if err != nil {
					t.mu.Unlock()
					return fmt.Errorf("failed to parse %s: %v", templatePath, err)
				}
			}
			t.cache[cacheKey] = tmpl
		}
		t.mu.Unlock()
	}

	if err := tmpl.ExecuteTemplate(w, layout, data); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	return nil
}