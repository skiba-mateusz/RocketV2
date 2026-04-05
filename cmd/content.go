package cmd

import (
	"context"
	"fmt"

	"github.com/skiba-mateusz/RocketV2/internal/config"
	"github.com/skiba-mateusz/RocketV2/internal/content"
	"github.com/skiba-mateusz/RocketV2/pkg/commandeer"
	"github.com/skiba-mateusz/RocketV2/pkg/logger"
)

func NewContentCmd() *commandeer.Command {
	contentCmd := commandeer.NewCommand(
		"content",
		"Manage content",
		func(ctx context.Context, cmd *commandeer.Command, args []string) error {
			cmd.Help()
			return nil
		},
	)

	contentCmd.Add(NewCreateContentCmd())

	return contentCmd
}

func NewCreateContentCmd() *commandeer.Command {
	createContentCmd := commandeer.NewCommand(
		"create",
		"Create content page",
		func(ctx context.Context, cmd *commandeer.Command, args []string) error {
			if err := assertProjectRoot(); err != nil {
				return err
			}

			if len(args) == 0 {
				return fmt.Errorf("provide path, e.g. /blogs/my-first-blog")
			}

			logger := logger.NewDefault(false)
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			path := args[0]

			outPath, err := content.Create(path, cfg.ContentDir)
			if err != nil {
				return err
			}

			logger.Info("%s created successfully", outPath)
			
			return nil
		},
	)

	return createContentCmd
}