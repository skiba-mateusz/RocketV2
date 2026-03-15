package app

import (
	"github.com/skiba-mateusz/RocketV2/builder"
	"github.com/skiba-mateusz/RocketV2/config"
	"github.com/skiba-mateusz/RocketV2/logger"
	"github.com/skiba-mateusz/RocketV2/parser"
	"github.com/skiba-mateusz/RocketV2/server"
	"github.com/skiba-mateusz/RocketV2/templater"
)

type App struct {
	Logger 	*logger.Logger
	Config 	*config.Config
	Builder *builder.Builder
}

func New() (*App, error) {
	logger := logger.NewLogger()
	config, err := config.Load()
	if err != nil {
		logger.Error("failed to load config: %v", err)
		return nil, err
	}

	parser := parser.NewMarkdwonParser()
	templater, err := templater.NewGoTemplater(config)
	if err != nil {
		logger.Error("failed to load templater: %v", err)
		return nil, err
	}

	return &App{
		Logger: logger,
		Config: config,
		Builder: builder.NewBuilder(logger, config, parser, templater),
	}, nil
}

func (a *App) NewServer(port string) *server.Server {
	return server.New(a.Logger, a.Config, port)
}