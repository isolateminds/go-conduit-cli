package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ttacon/chalk"
)

// Prints a error message and calls a optional cleanup function and then os.Exit(1)
func FatalError(from string, err error, cleanupFns ...func() error) {
	msg := strings.ReplaceAll(err.Error(), ":", " -> ")
	fmt.Printf("(%s): %s\n", chalk.Red.Color(from), chalk.White.Color(msg))
	for _, cleanup := range cleanupFns {
		err := cleanup()
		if err != nil {
			fmt.Printf("(%s): %s\n", chalk.Red.Color("CleanupError"), chalk.White.Color(msg))
		}
	}
	os.Exit(1)
}
func Success(msg string) {
	fmt.Println(chalk.Green.Color(msg))
}
