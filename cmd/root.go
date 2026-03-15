package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	rocketApp "github.com/skiba-mateusz/RocketV2/app"
	"github.com/skiba-mateusz/RocketV2/commandeer"
)

var app *rocketApp.App

var rootCmd = commandeer.NewCommand(
	"RocketV2",
	"Fast SSG written in Golang",
	func(ctx context.Context, cmd *commandeer.Command, args []string) error { 
		cmd.Help()
		return nil
	},
)

func Execute() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	a, err := rocketApp.New()
	if err != nil {
		return err
	}
	app = a


	if err := rootCmd.Execute(ctx); err != nil {
		app.Logger.Error("%v", err)
		return err
	}

	return nil
}