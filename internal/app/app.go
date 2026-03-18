package app

import (
	"fmt"

	"github.com/skiba-mateusz/RocketV2/internal/builder"
	"github.com/skiba-mateusz/RocketV2/internal/config"
	"github.com/skiba-mateusz/RocketV2/internal/parser"
	"github.com/skiba-mateusz/RocketV2/internal/server"
	"github.com/skiba-mateusz/RocketV2/internal/templater"
	"github.com/skiba-mateusz/RocketV2/pkg/logger"
)

type App struct {
	logger    *logger.Logger
	config    *config.Config
	parser    parser.Parser
	templater templater.Templater
	builder   builder.Builder
	server	  server.Server
}

func New(verbose bool) (*App, error) {
	log := logger.NewDefault(verbose)
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	tmpl, err := templater.NewGoTemplater(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create templater: %v", err)
	}

	prs := parser.NewMarkdownParser()

	return &App{
		logger: log,
		config: cfg,
		templater: tmpl,
		parser: prs,
	}, nil	
}

func (a *App) Builder() builder.Builder {
	return builder.NewDefault(a.logger, a.config, a.parser, a.templater)
}

func (a *App) Server(port string) server.Server {
	return server.NewDefault(a.logger, a.config, port)
}