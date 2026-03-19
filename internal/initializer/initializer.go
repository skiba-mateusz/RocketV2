package initializer

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed all:skeleton
var skeleton embed.FS

func Initialize(name string) error {
	wd, _ := os.Getwd()
	
	if err := os.Mkdir(name, 0755); err != nil {
		return fmt.Errorf("failed to create project dir %s: %v", name, err)
	}
	
	err := fs.WalkDir(skeleton, "skeleton", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel("skeleton", path)
		if relPath == "." {
			return nil
		}

		outPath := filepath.Join(wd, name, relPath)

		if d.IsDir() {
			if err := os.MkdirAll(outPath, 0755); err != nil {
				return fmt.Errorf("failed to create %s: %v", path, err)
			}
		} else {
			content, err := skeleton.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read embedded %s: %v", path, err)
			}
			if err = os.WriteFile(outPath, content, 0644); err != nil {
				return fmt.Errorf("failed to create file %s: %v", path, err)
			}
		}

		return  nil
	})
	if err != nil {
		return err
	}

	return nil
}