package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/skiba-mateusz/RocketV2/app"
	"github.com/skiba-mateusz/RocketV2/commandeer"
)

func newRootCmd(app *app.App) *commandeer.Command {
	rootCmd := commandeer.NewCommand(
		"RocketV2",
		"Fast SSG written in Golang",
		func(ctx context.Context, cmd *commandeer.Command, args []string) error { 
			cmd.Help()
			return nil
		},
	)

	rootCmd.Add(newBuildCmd(app))
	rootCmd.Add(newServeCmd(app))

	return rootCmd
}

func Execute() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app, err := app.New()
	if err != nil {
		return err
	}

	if err := newRootCmd(app).Execute(ctx); err != nil {
		app.Logger.Error("%v", err)
		return err
	}

	return nil
}