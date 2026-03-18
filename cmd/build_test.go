package cmd

import (
	"fmt"
	"testing"

	"github.com/skiba-mateusz/RocketV2/internal/app"
	"github.com/skiba-mateusz/RocketV2/pkg/commandeer"
)

func TestBuildCmd_Success(t *testing.T) {
	var getApp getAppFunc
	getApp = func(cmd *commandeer.Command) (*app.App, error) {
		return &app.App{
			Bldr: &mockBuilder{},
		}, nil
	}

	buildCmd := newBuildCmd(getApp)

	err := buildCmd.ExecuteArgs(t.Context(), []string{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestBuildCmd_Error(t *testing.T) {
	var getApp getAppFunc
	getApp = func(cmd *commandeer.Command) (*app.App, error) {
		return nil, fmt.Errorf("app failed")
	}

	buildCmd := newBuildCmd(getApp)

	err := buildCmd.ExecuteArgs(t.Context(), []string{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestBuildCmd_BuilderError(t *testing.T) {
	var getApp getAppFunc
	getApp = func(cmd *commandeer.Command) (*app.App, error) {
		return &app.App{
			Bldr: &mockBuilder{err: fmt.Errorf("failed to build")},
		}, nil
	}

	buildCmd := newBuildCmd(getApp)

	err := buildCmd.ExecuteArgs(t.Context(), []string{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestBuildCmd_VerboseFlagIsSet(t *testing.T) {
	var getApp getAppFunc
	getApp = func(cmd *commandeer.Command) (*app.App, error) {
		return &app.App{
			Bldr: &mockBuilder{},
		}, nil
	}

	buildCmd := newBuildCmd(getApp)

	_ = buildCmd.ExecuteArgs(t.Context(), []string{"--verbose"})

	verboseFlag := buildCmd.Flags.GetBool("verbose")
	if verboseFlag == false {
		t.Fatal("expected verbose flag to be true, got false")
	}
}