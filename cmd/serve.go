package cmd

import (
	"context"

	"github.com/skiba-mateusz/RocketV2/app"
	"github.com/skiba-mateusz/RocketV2/commandeer"
)

func newServeCmd(app *app.App) *commandeer.Command {
	serveCmd := commandeer.NewCommand(
		"serve",
		"Start development server",
		func(ctx context.Context, cmd *commandeer.Command, args []string) error {
			if err := app.Builder.Build(ctx); err != nil {
				return err
			}
			
			server := app.NewServer(cmd.Flags.GetString("port"))
			return server.Run(ctx)
		},
	)

	serveCmd.Flags.SetString("port", "8000", "specify port to listen on")

	return serveCmd
}