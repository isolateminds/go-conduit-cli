package cmd

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/isolateminds/go-conduit-cli/pkg/conduit"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"k8s.io/utils/strings/slices"
)

var (
	//go:embed embed/loki.cfg.yml
	lokiCfg []byte
	//go:embed embed/prometheus.cfg.yml
	prometheusCfg []byte
)

var (
	profiles    []string
	services    []string
	projectName string
	detach      bool
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
	stop = &cobra.Command{
		Use:   "stop",
		Short: "Bring down your local Conduit deployment",
		Long:  "Bring down your local Conduit deployment",
		Run: func(cmd *cobra.Command, args []string) {
			if !IsInProjectDirectory() {
				if err := ChangeToProjectRootDir(); err != nil {
					FatalError("StartError", err)
				}
			}
			if len(services) > 0 {
				services = slices.Filter(nil, services, func(s string) bool {
					return strings.TrimSpace(s) != ""
				})
			}
			ctx := context.Background()
			con, err := conduit.NewConduitFromProject(ctx, detach, []string{})
			if err != nil {
				FatalError("StartError", err)
			}
			err = con.Stop(ctx, services)
			if err != nil {
				FatalError("StartError", err)
			}
		},
	}
	start = &cobra.Command{
		Use:   "start",
		Short: "Bring up your local Conduit deployment",
		Long:  "Bring up your local Conduit deployment",
		Run: func(cmd *cobra.Command, args []string) {
			if !IsInProjectDirectory() {
				if err := ChangeToProjectRootDir(); err != nil {
					FatalError("StartError", err)
				}
			}
			ctx := context.Background()
			con, err := conduit.NewConduitFromProject(ctx, detach, profiles)
			if err != nil {
				FatalError("StartError", err)
			}
			err = con.Json.WriteFile()
			if err != nil {
				FatalError("StartError", err)
			}
			err = con.Up(ctx)
			if err != nil {
				FatalError("StartError", err)
			}
			if detach {
				Success("Started")
				os.Exit(0)
			}

		},
	}
	setup = &cobra.Command{
		Use:   "setup",
		Short: "Bootstrap a local Conduit deployment",
		Long:  "Bootstrap a local Conduit deployment",
		Run: func(cmd *cobra.Command, args []string) {
			if IsInProjectDirectory() {
				FatalError("SetupError", errors.New("already in project directory"))
			}
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
			con, err := conduit.NewConduitBootstrapper(ctx, projectName, detach, profiles)
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
			err = con.Json.WriteFile()
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
	deploy.PersistentFlags().BoolVar(&detach, "detach", false, "run containers in the background")
	deploy.AddCommand(start)
	start.PersistentFlags().BoolVar(&detach, "detach", false, "run containers in the background")
	start.PersistentFlags().StringSliceVar(&profiles, "profiles", []string{}, "profiles to enable")
	deploy.AddCommand(stop)
	stop.PersistentFlags().StringSliceVar(&services, "services", []string{}, "services to stop")

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

// Traverse up the directory tree until conduit.json is found or we reach the root directory
// and changes to that directory
func ChangeToProjectRootDir() error {
	// Start from the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("changeToProjectRootDirError: %s", err)
	}

	for {
		conduitPath := filepath.Join(currentDir, "conduit.json")
		_, err := os.Stat(conduitPath)
		if err == nil {
			return os.Chdir(currentDir)
		}

		// Move to the parent directory
		parentDir := filepath.Dir(currentDir)
		// Check if we have reached the root directory
		if parentDir == currentDir {
			break
		}
		currentDir = parentDir
	}

	return errors.New("changeToProjectRootDirError: you are not in a project directory")
}

// Checks to see if you are in root project dir
func IsInProjectDirectory() bool {
	_, err := os.Stat("conduit.json")
	return !errors.Is(err, os.ErrNotExist)
}
