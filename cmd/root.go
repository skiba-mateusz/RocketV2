package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/skiba-mateusz/RocketV2/internal/app"
	"github.com/skiba-mateusz/RocketV2/pkg/commandeer"
)

type getAppFunc func(cmd *commandeer.Command) (*app.App, error)

func Execute() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := newRootCmd().Execute(ctx); err != nil {
		return err
	}

	return nil
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

	var getApp func(cmd *commandeer.Command) (*app.App, error)
	getApp = func(cmd *commandeer.Command) (*app.App, error) {
		app, err := app.New(cmd.Flags.GetBool("verbose"))
		if err != nil {
			return nil, fmt.Errorf("failed to create app: %v", err)
		}
		return app, nil
	}

	rootCmd.Add(newBuildCmd(getApp))
	rootCmd.Add(newServeCmd(getApp))
	rootCmd.Add(newInitCmd())

	return rootCmd
}