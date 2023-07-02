package main

import (
	"fmt"
	"github.com/spf13/pflag"
)

func main() {
	// Define flags
	name := pflag.String("name", "", "Your name")
	age := pflag.Int("age", 0, "Your age")
	verbose := pflag.Bool("verbose", false, "Enable verbose mode")

	// Parse the command-line arguments
	pflag.Parse()

	// Access the flag values
	fmt.Println("Name:", *name)
	fmt.Println("Age:", *age)
	fmt.Println("Verbose:", *verbose)

	// Additional program logic here...
}
