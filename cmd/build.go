package cmd

import (
	"context"

	"github.com/skiba-mateusz/RocketV2/app"
	"github.com/skiba-mateusz/RocketV2/commandeer"
)

func newBuildCmd(app *app.App) *commandeer.Command {
	buildCmd := commandeer.NewCommand(
		"build",
		"Build static site",
		func(ctx context.Context, cmd *commandeer.Command, args []string) error {
			return app.Builder.Build(ctx)
		},
	)

	buildCmd.Flags.SetInt("threads", 4, "specify threads number")
	buildCmd.Flags.SetBool("verbose", false, "specify the verbosity")

	return buildCmd
}