package cmd

import (
	"context"

	"github.com/skiba-mateusz/RocketV2/pkg/commandeer"
)

func newBuildCmd(getApp getAppFunc) *commandeer.Command {
	buildCmd := commandeer.NewCommand(
		"build",
		"Build static site",
		func(ctx context.Context, cmd *commandeer.Command, args []string) error {
			if err := assertProjectRoot(); err != nil {
				return err
			}

			app, err := getApp(cmd)
			if err != nil {
				return err
			}

			return app.Builder().Build(ctx)
		},
	)

	buildCmd.Flags.SetBool("verbose", false, "specify the verbosity")

	return buildCmd
}