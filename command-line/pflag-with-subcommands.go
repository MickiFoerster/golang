package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

func main() {
	// Define subcommands as positional arguments
	listCmd := pflag.NewFlagSet("list", pflag.ExitOnError)
	versionCmd := pflag.NewFlagSet("version", pflag.ExitOnError)

	// Parse the command-line arguments
	pflag.Parse()

	// Check if any subcommand is provided
	if len(os.Args) < 2 {
		fmt.Println("Please provide a subcommand: list or version")
		os.Exit(1)
	}

	// Execute the appropriate subcommand logic
	switch os.Args[1] {
	case "list":
		listCmd.Parse(os.Args[2:])
		// Execute list subcommand logic here
		fmt.Println("Executing 'list' subcommand")
	case "version":
		versionCmd.Parse(os.Args[2:])
		// Execute version subcommand logic here
		fmt.Println("Executing 'version' subcommand")
	default:
		fmt.Println("Invalid subcommand. Please choose 'list' or 'version'")
		os.Exit(1)
	}
}
