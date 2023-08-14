package cmd

import (
	"context"
	"log"

	"github.com/isolateminds/go-conduit-cli/pkg/conduit"
	"github.com/spf13/cobra"
)

var (
	profiles    []string
	projectName string
	detatch     bool
	deploy      = &cobra.Command{
		Use:                   "deploy",
		DisableFlagParsing:    true,
		DisableFlagsInUseLine: true,
		DisableSuggestions:    true,
		Short:                 "Manage a local Conduit deployment",
		Long:                  "Manage a local Conduit deployment",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				log.Fatal(err)
			}
		},
	}
	setup = &cobra.Command{
		Use:   "setup",
		Short: "Bootstrap a local Conduit deployment",
		Long:  "Bootstrap a local Conduit deployment",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			c, err := conduit.NewConduit(ctx, projectName, detatch, profiles)
			if err != nil {
				FatalError("SetupError", err)
			}

			if err := c.Setup(ctx); err != nil {
				log.Fatal(err)
			}
		},
	}
)

func init() {
	root.AddCommand(deploy)
	deploy.AddCommand(setup)
	deploy.PersistentFlags().StringSliceVar(&profiles, "profiles", []string{}, "profiles to enable")
	deploy.PersistentFlags().StringVar(&projectName, "project-name", "conduit", "set the project name (defaults to conduit)")
	deploy.PersistentFlags().BoolVar(&detatch, "detach", false, "run containers in the background")
}
