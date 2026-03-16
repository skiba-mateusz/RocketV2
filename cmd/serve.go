package cmd

import (
	"context"

	"github.com/skiba-mateusz/RocketV2/commandeer"
	"github.com/skiba-mateusz/RocketV2/server"
)

func newServeCmd(getApp appInitFunc) *commandeer.Command {
	serveCmd := commandeer.NewCommand(
		"serve",
		"Start development server",
		func(ctx context.Context, cmd *commandeer.Command, args []string) error {
			app, err := getApp(cmd)
			if err != nil {
				return err
			}

			if err := app.builder.Build(ctx); err != nil {
				return err
			}

			srv := server.New(app.logger, app.config, cmd.Flags.GetString("port"))

			if err = srv.Run(ctx); err != nil {
				return err
			} 

			return nil
		},
	)

	serveCmd.Flags.SetString("port", "8000", "specify port to listen on")

	return serveCmd
}