package cmd

import (
	"context"
	"os"
	"path/filepath"

	"github.com/skiba-mateusz/RocketV2/pkg/commandeer"
)

func newServeCmd(getApp getAppFunc) *commandeer.Command {
	serveCmd := commandeer.NewCommand(
		"serve",
		"Start development server",
		func(ctx context.Context, cmd *commandeer.Command, args []string) error {
			if err := assertProjectRoot(); err != nil {
				return err
			}

			app, err := getApp(cmd)
			if err != nil {
				return err
			}

			if err := app.Builder().Build(ctx); err != nil {
				return err
			}

			wd, _ := os.Getwd()
			cfgPath := filepath.Join(wd, "config.yaml")

			return app.Server(cmd.Flags.GetString("port"), cfgPath).Run(ctx)
		},
	)

	serveCmd.Flags.SetString("port", "8000", "specify port to listen on")
	serveCmd.Flags.SetBool("verbose", false, "specify the verbosity")
	
	return serveCmd
}