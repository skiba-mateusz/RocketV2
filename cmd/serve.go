package cmd

import (
	"context"

	"github.com/skiba-mateusz/RocketV2/pkg/commandeer"
)

func newServeCmd(getApp getAppFunc) *commandeer.Command {
	serveCmd := commandeer.NewCommand(
		"serve",
		"Start development server",
		func(ctx context.Context, cmd *commandeer.Command, args []string) error {
			app, err := getApp(cmd)
			if err != nil {
				return err
			}

			if err := app.Builder().Build(ctx); err != nil {
				return err
			}

			return app.Server(cmd.Flags.GetString("port")).Run(ctx)
		},
	)

	serveCmd.Flags.SetString("port", "8000", "specify port to listen on")

	return serveCmd
}