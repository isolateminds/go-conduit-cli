package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:               "goconduit",
	DisableAutoGenTag: true,

	Short: "A Golang port of ConduitPlatform/CLI with some additional features",
	Long:  "A Golang port of ConduitPlatform/CLI with some additional features",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			log.Fatal(err)
		}
	},
}

func Execute() error {
	return root.Execute()
}
