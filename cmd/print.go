package cmd

import (
	"fmt"
	"os"

	"github.com/ttacon/chalk"
)

// Prints a error message and then calls os.Exit(1)
func PrintFatalError(err error) {
	fmt.Println(chalk.White.Color(err.Error()))
	os.Exit(1)
}
func PrintSuccess(msg string) {
	fmt.Println(chalk.Green.Color(msg))
}
