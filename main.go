package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/skiba-mateusz/RocketV2/builder"
	"github.com/skiba-mateusz/RocketV2/config"
	"github.com/skiba-mateusz/RocketV2/logger"
	"github.com/skiba-mateusz/RocketV2/parser"
	"github.com/skiba-mateusz/RocketV2/templater"
)


func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger := logger.NewLogger()

	config, err := config.Load()
	if err != nil {
		handleError(err, logger)
	}

	metadataParser := parser.NewMarkdwonParser()
	goTemplater, err := templater.NewGoTemplater(config)
	if err != nil {
		handleError(err, logger)
	}

	bldr := builder.NewBuilder(logger, config, metadataParser, goTemplater)
	if err := bldr.Build(ctx); err != nil {
		handleError(err, logger)
	}
}

func handleError(err error, logger *logger.Logger) {
	if err != nil {
		logger.Error("Error: %v", err)
		os.Exit(1)
	}
}