package cmd

import (
	"fmt"

	"github.com/skiba-mateusz/RocketV2/commandeer"
)

var rootCmd = commandeer.NewCommand(
	"RocketV2",
	"Fast SSG written in Golang",
	func(cmd *commandeer.Command, args []string) error {
		fmt.Println("root")
		fmt.Println(args)
		return nil
	},
)

func Execute() error {
	return rootCmd.Execute()
}