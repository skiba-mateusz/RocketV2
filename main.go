package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/skiba-mateusz/RocketV2/builder"
	"github.com/skiba-mateusz/RocketV2/cmd"
	"github.com/skiba-mateusz/RocketV2/config"
	"github.com/skiba-mateusz/RocketV2/logger"
	"github.com/skiba-mateusz/RocketV2/parser"
	"github.com/skiba-mateusz/RocketV2/server"
	"github.com/skiba-mateusz/RocketV2/templater"
)


func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
	}
	
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

	_ = server.New(logger, config, ":8000")

	// if err := srv.Run(ctx); err != nil {
	// 	handleError(err, logger)
	// }
}

func handleError(err error, logger *logger.Logger) {
	if err != nil {
		logger.Error("Error: %v", err)
		os.Exit(1)
	}
}