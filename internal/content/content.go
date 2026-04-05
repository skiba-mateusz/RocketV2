package content

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/skiba-mateusz/RocketV2/internal/parser"
	"gopkg.in/yaml.v2"
)

func Create(path, contentDir string) (string, error) {
	cleanPath := normalizePath(path)
	metadata := populateMetadata(cleanPath)

	fm, err := writeYamlFrontmatter(metadata)
	if err != nil {
		return "", fmt.Errorf("failed to write front matter: %v", err)
	}

	outPath := filepath.Join(contentDir, cleanPath)
	if err := saveFile(outPath, fm); err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}
	
	return outPath, nil
}

func normalizePath(path string) string {
	ext := filepath.Ext(path)

	if ext == ".md" {
		return path
	}

	if ext != "" {
		path = path[:len(path)-len(ext)]
	}

	return path + ".md"
}

func populateMetadata(path string) parser.Metadata {
	base := filepath.Base(path)
	title := base[:len(base)-len(filepath.Ext(base))]
	title = strings.Join(strings.Split(title, "-"), " ")
	date := time.Now().Format(time.RFC3339)

	metadata := parser.Metadata{
		Title: title,
		Date: date,
	}

	return metadata
}

func writeYamlFrontmatter(metadata parser.Metadata) ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString("---\n")
	if err := yaml.NewEncoder(&buf).Encode(metadata); err != nil {
		return nil, fmt.Errorf("failed to encode yaml front matter: %v", err)
	}
	buf.WriteString("---\n")
	buf.WriteString(fmt.Sprintf("# %s", metadata.Title))

	return buf.Bytes(), nil
}

func saveFile(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}
	
	return nil
}