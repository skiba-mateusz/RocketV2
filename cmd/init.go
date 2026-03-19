package cmd

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/skiba-mateusz/RocketV2/internal/initializer"
	"github.com/skiba-mateusz/RocketV2/pkg/commandeer"
)

func newInitCmd() *commandeer.Command {
	initCmd := commandeer.NewCommand(
		"init",
		"Initialize new project",
		func(ctx context.Context, cmd *commandeer.Command, args []string) error {
			name := "my-rocket-site"
			if len(args) > 0 {
				name = args[0]
			}

			clean := filepath.Base(name)
			if clean == "." || clean == ".." {
				return fmt.Errorf("invalid project name: %s", name)
			}

			return initializer.Initialize(name)
		},
	)

	return initCmd
}