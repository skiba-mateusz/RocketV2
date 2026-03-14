package cmd

import (
	"fmt"

	"github.com/skiba-mateusz/RocketV2/commandeer"
)

var buildCmd = commandeer.NewCommand(
	"build",
	"Build static site",
	func(cmd *commandeer.Command, args []string) error {
		fmt.Println("build")
		fmt.Println(cmd.Flags.GetInt("threads"))
		fmt.Println(cmd.Flags.GetBool("verbose"))
		fmt.Println(args)
		return nil
	},
)

func init() {
	buildCmd.Flags.SetInt("threads", 4, "specify threads number")
	buildCmd.Flags.SetBool("verbose", false, "specify the verbosity")
	rootCmd.Add(buildCmd)
}