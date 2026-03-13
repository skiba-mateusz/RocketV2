package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	BaseURL      string					`yaml:"baseURL"`
	LanguageCode string					`yaml:"languageCode"`	
	Title		 string					`yaml:"title"`
	ContentDir   string					`yaml:"contentDir"`
	LayoutDir    string					`yaml:"layoutDir"`
	BuildDir     string					`yaml:"buildDir"`
	Paginate	 int					`yaml:"paginate"`
	Params       map[string]interface{} `yaml:"params"`
}

func Load() (*Config, error) {
	wd, _ := os.Getwd()
	path := filepath.Join(wd, "config.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to read file %s: %v", path, err)
	 }

	cfg := Config{
		Paginate: 	2,
		ContentDir: "content",
		LayoutDir:  "layout",
		BuildDir: 	"build",
	}
	
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal YAML file %s: %v", path, err)
	}

	return &cfg, nil
}