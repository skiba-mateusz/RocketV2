package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/skiba-mateusz/RocketV2/internal/builder"
	"github.com/skiba-mateusz/RocketV2/internal/config"
	"github.com/skiba-mateusz/RocketV2/internal/parser"
	"github.com/skiba-mateusz/RocketV2/internal/templater"
	"github.com/skiba-mateusz/RocketV2/pkg/commandeer"
	"github.com/skiba-mateusz/RocketV2/pkg/logger"
)

type app struct {
	logger 		*logger.Logger
	config 		*config.Config
	parser 	 	parser.Parser
	templater 	templater.Templater
	builder 	*builder.Builder
}

type appInitFunc func(cmd *commandeer.Command) (*app, error)

func Execute() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := newRootCmd().Execute(ctx); err != nil {
		return err
	}

	return nil
}

func newApp(cmd *commandeer.Command) (*app, error) {
	log := logger.New(cmd.Flags.GetBool("verbose"))
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	prs := parser.NewMarkdwonParser()
	tmpl, err := templater.NewGoTemplater(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create templater: %v", err)
	}

	builder := builder.New(log, cfg, prs, tmpl)

	return &app{
		logger: log,
		config: cfg,
		parser: prs,
		templater: tmpl,
		builder: builder,
	}, nil
}

func newRootCmd() *commandeer.Command {
	rootCmd := commandeer.NewCommand(
		"RocketV2",
		"Fast SSG written in Golang",
		func(ctx context.Context, cmd *commandeer.Command, args []string) error { 
			cmd.Help()
			return nil
		},
	)

	getApp := newApp

	rootCmd.Add(newBuildCmd(getApp))
	rootCmd.Add(newServeCmd(getApp))

	return rootCmd
}