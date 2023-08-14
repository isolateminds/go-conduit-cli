package main

import (
	"fmt"
	"os"

	"github.com/isolateminds/go-conduit-cli/cmd"
)

func main() {
	cmd.Banner()
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
