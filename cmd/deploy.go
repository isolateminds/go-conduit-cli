package cmd

import (
	"context"
	_ "embed"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/isolateminds/go-conduit-cli/pkg/conduit"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var (
	//go:embed embed/loki.cfg.yml
	lokiCfg []byte
	//go:embed embed/prometheus.cfg.yml
	prometheusCfg []byte
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
	start = &cobra.Command{
		Use: "start",
	}
	setup = &cobra.Command{
		Use:   "setup",
		Short: "Bootstrap a local Conduit deployment",
		Long:  "Bootstrap a local Conduit deployment",
		Run: func(cmd *cobra.Command, args []string) {
			err := os.Mkdir(projectName, fs.ModePerm)
			if err != nil {
				FatalError("SetupError", err)
			}
			err = os.Chdir(projectName)
			if err != nil {
				FatalError("SetupError", err)
			}
			pDir, err := os.Getwd()
			if err != nil {
				FatalError("SetupError", err)
			}
			//In new directory -- write config files, can can delete upon error
			deletePDir := func() error {
				return os.RemoveAll(pDir)
			}
			HandleSIGTERM(context.Background(), func() {
				err := deletePDir()
				if err != nil {
					FatalError("SetupError", err)
				}
			})
			err = os.WriteFile("loki.cfg.yml", lokiCfg, fs.ModePerm)
			if err != nil {
				FatalError("SetupError", err, deletePDir)
			}
			err = os.WriteFile("prometheus.cfg.yml", prometheusCfg, fs.ModePerm)
			if err != nil {
				FatalError("SetupError", err, deletePDir)
			}

			ctx := context.Background()
			con, err := conduit.NewConduitBootstrapper(ctx, projectName, detatch, profiles)
			if err != nil {
				FatalError("SetupError", err, deletePDir)
			}
			//Write .env file
			err = godotenv.Write(con.Composer.Options.Environment.Variables, ".env")
			if err != nil {
				FatalError("SetupError", err, deletePDir)
			}
			//Write compose file
			err = os.WriteFile("docker-compose.yaml", con.Composer.Options.Yaml.Bytes, fs.ModePerm)
			if err != nil {
				FatalError("SetupError", err, deletePDir)
			}
			err = conduit.WriteConduitJSON(projectName, "1.0.0", con.DatabaseName, con.Composer.FilterYamlProfiles(profiles))
			if err != nil {
				FatalError("SetupError", err, deletePDir)
			}
			err = con.Up(ctx)
			if err != nil {
				FatalError("SetupError", err, deletePDir)
			}
			Success("project created")
			os.Exit(0)

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

// Invokes a callback after signal termination CTRL+C
func HandleSIGTERM(ctx context.Context, cb func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		select {
		case <-c:
			cb()
		case <-ctx.Done():
			// Do nothing if the context is canceled before receiving the signal
		}
	}()
}
