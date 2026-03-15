package cmd

import (
	"context"

	"github.com/skiba-mateusz/RocketV2/commandeer"
)

var buildCmd = commandeer.NewCommand(
	"build",
	"Build static site",
	func(ctx context.Context, cmd *commandeer.Command, args []string) error {
		return app.Builder.Build(ctx)
	},
)

func init() {
	buildCmd.Flags.SetInt("threads", 4, "specify threads number")
	buildCmd.Flags.SetBool("verbose", false, "specify the verbosity")
	rootCmd.Add(buildCmd)
}