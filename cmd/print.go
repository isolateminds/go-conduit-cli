package cmd

import (
	"fmt"
	"os"

	"github.com/ttacon/chalk"
)

// Prints a error message and calls os.Exit(1)
func FatalError(from string, err error) {
	fmt.Printf("(%s): %s\n", chalk.Red.Color(from), chalk.White.Color(err.Error()))
	os.Exit(1)
}
