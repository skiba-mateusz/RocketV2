package commandeer

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
)

type CommandFunc func(ctx context.Context, cmd *Command, args []string) error

type Command struct {
	Name        string
	Description string
	Flags 		Flags
	run         CommandFunc
	subCommands map[string]*Command
}

func NewCommand(name, description string, run CommandFunc) *Command {
	command := &Command{
		Name: name,
		Description: description,
		run: run,
	}
	command.Flags.Init(name)
	return command
}

func (c *Command) Add(cmd *Command) error {
	if c.subCommands == nil {
		c.subCommands = make(map[string]*Command)
	}

	if _, ok := c.subCommands[cmd.Name]; ok {
		return fmt.Errorf("command %s already registered", cmd.Name)
	}
	c.subCommands[cmd.Name] = cmd

	return nil
}

func (c *Command) Execute(ctx context.Context) error {
	args := os.Args[1:]
	return c.execute(ctx, args)
}

func (c *Command) Help() {	
	c.printHelp(os.Stdout)
}

func (c *Command) printHelp(w io.Writer) {
	fmt.Printf("Usage: %s COMMAND [OPTIONS]\n\n", c.Name)
	fmt.Printf("%s\n", c.Description)
	
	if len(c.subCommands) > 0 {
		fmt.Println("\nCommands:")
		for _, sub := range c.subCommands {
			fmt.Printf("  %-16s %s\n", sub.Name, sub.Description)
			if sub.Flags.isEmpty() {
				continue
			}
			sub.Flags.set.VisitAll(func(f *flag.Flag) {
				fmt.Printf("    --%-14s %s\n", f.Name, f.Usage)
			})
		}
	}

	if !c.Flags.isEmpty() {
		fmt.Println("\nFlags:")
		c.Flags.set.VisitAll(func(f *flag.Flag) {
			fmt.Printf("  --%-14s %s\n", f.Name, f.Usage)
		})
	}
}

func (c *Command) execute(ctx context.Context, args []string) error {
	if len(args) == 0 {
		if c.run != nil {
			return c.run(ctx, c, nil)
		}
		return nil
	}

	next := args[0]

	if next == "--help" || next == "-h" {
		c.Help()
		return nil
	}

	if sub, ok := c.subCommands[next]; ok {
		return sub.execute(ctx, args[1:])
	}

	if c.run != nil {
		var flagArgs, restArgs []string

		for i := 0; i < len(args); i++ {
			arg := args[i]
			if len(arg) > 0 && arg[0] == '-' {
				flagArgs = append(flagArgs, arg)
				if i+1 < len(args) && args[i+1][0] != '-' {
					i++
					flagArgs = append(flagArgs, args[i])
				}
			} else {
				restArgs = append(restArgs, arg)
			}
		}

		if err := c.Flags.set.Parse(flagArgs); err != nil {
			return fmt.Errorf("failed to parse flags: %v", err)
		}

		return c.run(ctx, c, restArgs)
	}

	if len(next) > 0 && next[0] == '-' {
		return fmt.Errorf("unknown flag %s", next)
	}

	return fmt.Errorf("unknown command %s", next)
}