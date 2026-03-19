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
	Logger    *logger.Logger
	Config    *config.Config
	Parser    parser.Parser
	Templater templater.Templater
	Bldr   	  builder.Builder
	Srv	  	  server.Server
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
		Logger: log,
		Config: cfg,
		Templater: tmpl,
		Parser: prs,
	}, nil	
}

func (a *App) Builder() builder.Builder {
	if a.Bldr == nil {
		a.Bldr = builder.NewDefault(a.Logger, a.Config, a.Parser, a.Templater)
	}
	return a.Bldr
}

func (a *App) Server(port string) server.Server {
	if a.Srv == nil {
		a.Srv = server.NewDefault(a.Logger, a.Config, port)
	}
	return a.Srv
}