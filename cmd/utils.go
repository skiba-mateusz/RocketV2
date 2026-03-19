package cmd

import (
	"errors"
	"fmt"
	"os"
)

func assertProjectRoot() error {
	if _ , err := os.Stat("config.yaml"); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("not a Rocket project (config.yaml not found)")
		}
		return fmt.Errorf("failed to stat config file: %v", err)
	}
	return nil
}	