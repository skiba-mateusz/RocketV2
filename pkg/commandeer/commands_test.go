package commandeer

import (
	"context"
	"fmt"
	"testing"
)

func TestCommand_Creation(t *testing.T) {
	name := "test"
	desc := "example description"
	cmd := NewCommand(
		name,
		desc,
		func(ctx context.Context, cmd *Command, args []string) error {
			return nil
		},
	)

	if cmd.Name != name {
		t.Errorf("expected name to be %s, got %s", name, cmd.Name)
	}

	if cmd.Description != desc {
		t.Errorf("expected description to be %s, got %s", desc, cmd.Description)
	}

	if cmd.run == nil {
		t.Error("expected run to be set, got nil")
	}
}

func TestCommand_Add(t *testing.T) {
	root := NewCommand(
		"root",
		"root desc",
		func(ctx context.Context, cmd *Command, args []string) error {
			return nil
		},
	)

	cmds := []*Command{
		NewCommand(
			"cmd1",
			"cmd1 desc",
			func(ctx context.Context, cmd *Command, args []string) error {
				return nil
			},
		),
		NewCommand(
			"cmd2",
			"cmd2 desc",
			func(ctx context.Context, cmd *Command, args []string) error {
				return nil
			},
		),
	}

	for _, cmd := range cmds {
		if err := root.Add(cmd); err != nil {
			t.Errorf("failed to add command: %v", err)
		}
	}

	if len(root.subCommands) != len(cmds) {
		t.Errorf("expected %d subcommands, got %d", len(cmds), len(root.subCommands))
	}
}

func TestCommand_FailsIfCommandExists(t *testing.T) {
	root := NewCommand(
		"root",
		"root desc",
		func(ctx context.Context, cmd *Command, args []string) error {
			return nil
		},
	)

	sub := NewCommand(
		"sub",
		"sub desc",
		func(ctx context.Context, cmd *Command, args []string) error {
			return nil
		},
	)

	err := root.Add(sub)
	if err != nil {
		t.Fatalf("failed to add command: %v", err)
	}

	err = root.Add(sub)
	if err == nil {
		t.Fatalf("expected error when command already exists, got nil")
	}
}

func TestCommand_RunsWithNoArgs(t *testing.T) {
	called := false
	cmd := NewCommand(
		"test",
		"desc",
		func(ctx context.Context, cmd *Command, args []string) error {
			called = true
			return nil
		},
	)

	err := cmd.ExecuteArgs(t.Context(), []string{})
	if err != nil {
		t.Fatalf("failed to execute commandd: %v", err)
	}
	if !called {
		t.Fatal("expected run to be called")
	}
}

func TestCommand_RunsSubCommand(t *testing.T) {
	called := false

	root := NewCommand(
		"root",
		"root desc",
		func(ctx context.Context, cmd *Command, args []string) error {
			return nil
		},
	)
	sub := NewCommand(
		"sub",
		"sub desc",
		func(ctx context.Context, cmd *Command, args []string) error {
			called = true
			return nil
		},
	)

	err := root.Add(sub)
	if err != nil {
		t.Fatalf("failed to add command: %v", err)
	}

	err = root.ExecuteArgs(t.Context(), []string{"sub"})
	if err != nil {
		t.Fatalf("failed to execute command: %v", err)
	}
	if !called {
		t.Fatal("expected subcommand to be run")
	}
}

func TestCommand_UnknownCommand(t *testing.T) {
	root := NewCommand(
		"root",
		"root desc",
		nil,
	)

	err := root.ExecuteArgs(t.Context(), []string{"unknown"})
	if err == nil {
		t.Fatal("expected error for unknown command")
	}
}

func TestCommand_ParseFlags(t *testing.T) {
	var name string 
	var age int

	root := NewCommand(
		"root",
		"root desc",
		func(ctx context.Context, cmd *Command, args []string) error {
			name = cmd.Flags.GetString("name")
			age = cmd.Flags.GetInt("age")
			return nil
		},
	)

	expectedName := "tom"
	expectedAge := 30

	root.Flags.SetString("name", "john", "specify name")
	root.Flags.SetInt("age", 20, "specify age")

	err := root.ExecuteArgs(t.Context(), []string{"--name", expectedName, "--age", fmt.Sprint(expectedAge)})
	if err != nil {
		t.Fatalf("failed to execute command: %v", err)
	}

	if name != expectedName {
		t.Errorf("expected name to be %s, got %s", expectedName, name)
	}

	if age != expectedAge {
		t.Errorf("expected age to be %d, got %d", expectedAge, age)
	}

}