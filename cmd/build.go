package cmd

import (
	"context"

	"github.com/skiba-mateusz/RocketV2/commandeer"
)

func newBuildCmd(getApp appInitFunc) *commandeer.Command {
	buildCmd := commandeer.NewCommand(
		"build",
		"Build static site",
		func(ctx context.Context, cmd *commandeer.Command, args []string) error {
			app, err := getApp(cmd)
			if err != nil {
				return err
			}

			return app.builder.Build(ctx)
		},
	)

	buildCmd.Flags.SetBool("verbose", false, "specify the verbosity")

	return buildCmd
}