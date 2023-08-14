package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ttacon/chalk"
)

// Prints a error message and calls os.Exit(1)
func FatalError(from string, err error) {
	msg := strings.ReplaceAll(err.Error(), ":", " -> ")
	fmt.Printf("(%s): %s\n", chalk.Red.Color(from), chalk.White.Color(msg))
	os.Exit(1)
}
