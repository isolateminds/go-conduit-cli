package cmd

import (
	"context"
	_ "embed"
	"errors"
	"io/fs"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/isolateminds/go-conduit-cli/pkg/conduit"
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
	imageTag    string
	uiImageTag  string
	detach      bool

	deploy = &cobra.Command{
		Use:              "deploy",
		Short:            "Manage a local Conduit deployment",
		Run:              runDeploy,
		PersistentPreRun: runPreRun,
	}
	rm = &cobra.Command{
		Use:   "rm",
		Short: "Remove your local Conduit deployment",
		Run:   runStop,
	}
	stop = &cobra.Command{
		Use:   "stop",
		Short: "Bring down your local Conduit deployment",
		Run:   runStop,
	}
	start = &cobra.Command{
		Use:   "start",
		Short: "Bring up your local Conduit deployment",
		Run:   runStart,
	}
	setup = &cobra.Command{
		Use:   "setup",
		Short: "Bootstrap a local Conduit deployment",
		Run:   runSetup,
	}
)

func init() {
	root.AddCommand(deploy)
	deploy.AddCommand(setup)
	deploy.AddCommand(start)
	deploy.AddCommand(stop)
	deploy.AddCommand(rm)

	//Flags
	//Deploy setup --profiles x,y,z --project-name my-proj --detach
	setup.PersistentFlags().StringSliceVar(&profiles, "profiles", []string{}, "profiles to enable")
	setup.PersistentFlags().StringVar(&projectName, "project-name", "conduit", "set the project name (defaults to conduit)")
	setup.PersistentFlags().StringVar(&imageTag, "image-tag", "latest", "set the conduit ui image tag to use (defaults to latest)")
	setup.PersistentFlags().StringVar(&uiImageTag, "ui-image-tag", "latest", "set the conduit ui image tag to use (defaults to latest)")
	setup.PersistentFlags().BoolVar(&detach, "detach", false, "run containers in the background")

	//deploy start --profiles x,y,z --detach
	start.PersistentFlags().BoolVar(&detach, "detach", false, "run containers in the background")
	start.PersistentFlags().StringSliceVar(&profiles, "profiles", []string{}, "profiles to enable")
	//deploy stop --services my-db,my-frontend
	stop.PersistentFlags().StringSliceVar(&services, "services", []string{}, "services to stop")

	//Deploy rm --services my-db,my-frontend
	rm.PersistentFlags().StringSliceVar(&services, "services", []string{}, "services to remove")
}
func runDeploy(cmd *cobra.Command, args []string) {
	if err := cmd.Help(); err != nil {
		PrintFatalError(err)
	}
}

// Sanitizes the stringslicevar flags incase there is a space eg: --services x,y, <-- trailing comma
func runPreRun(cmd *cobra.Command, args []string) {
	if len(services) > 0 {
		services = slices.Filter(nil, services, func(s string) bool {
			return strings.TrimSpace(s) != ""
		})
	}
	if len(profiles) > 0 {
		profiles = slices.Filter(nil, profiles, func(s string) bool {
			return strings.TrimSpace(s) != ""
		})
	}
}
func runRm(cmd *cobra.Command, args []string) {
	if !IsInProjectDirectory() {
		if err := ChangeToProjectRootDir(); err != nil {
			PrintFatalError(NewRemoveError(err))
		}
	}
	ctx := context.Background()
	con, err := conduit.NewConduitFromProject(ctx, detach, []string{})
	if err != nil {
		PrintFatalError(NewRemoveError(err))
	}
	err = con.Remove(ctx, services)
	if err != nil {
		PrintFatalError(NewRemoveError(err))
	}
}
func runStop(cmd *cobra.Command, args []string) {
	if !IsInProjectDirectory() {
		if err := ChangeToProjectRootDir(); err != nil {
			PrintFatalError(NewStopError(err))
		}
	}
	ctx := context.Background()
	con, err := conduit.NewConduitFromProject(ctx, detach, []string{})
	if err != nil {
		PrintFatalError(NewStopError(err))
	}
	err = con.Stop(ctx, services)
	if err != nil {
		PrintFatalError(NewStopError(err))
	}
}
func runStart(cmd *cobra.Command, args []string) {
	if !IsInProjectDirectory() {
		if err := ChangeToProjectRootDir(); err != nil {
			PrintFatalError(NewStartError(err))
		}
	}
	ctx := context.Background()
	con, err := conduit.NewConduitFromProject(ctx, detach, profiles)
	if err != nil {
		PrintFatalError(NewStartError(err))
	}
	err = con.WriteConduitJsonFile()
	if err != nil {
		PrintFatalError(NewStartError(err))
	}
	err = con.Up(ctx)
	if err != nil {
		PrintFatalError(NewStartError(err))
	}
	if detach {
		PrintSuccess("Started")
		os.Exit(0)
	}

}
func runSetup(cmd *cobra.Command, args []string) {
	if IsInProjectDirectory() {
		PrintFatalError(NewSetupError(errors.New("already in project directory")))
	}
	if err := os.Mkdir(projectName, fs.ModePerm); err != nil {
		PrintFatalError(NewSetupError(err))
	}
	if err := os.Chdir(projectName); err != nil {
		PrintFatalError(NewSetupError(err))
	}
	pDir, err := os.Getwd()
	if err != nil {
		PrintFatalError(NewSetupError(err))
	}

	//In new directory -- write config files, can can delete upon error
	deletePDir := func() { os.RemoveAll(pDir) }
	//Delete project dir uopn error unless canceled
	sigCtx, cancelSigKill := context.WithCancel(context.Background())
	handleSIGTERM(sigCtx, deletePDir)

	if err := os.WriteFile("loki.cfg.yml", lokiCfg, fs.ModePerm); err != nil {
		deletePDir()
		PrintFatalError(NewSetupError(err))
	}
	if err := os.WriteFile("prometheus.cfg.yml", prometheusCfg, fs.ModePerm); err != nil {
		deletePDir()
		PrintFatalError(NewSetupError(err))
	}
	options := &conduit.BootstrapperOptions{
		ProjectName: projectName,
		Detached:    detach,
		Profiles:    profiles,
		ImageTag:    imageTag,
		UIImageTag:  uiImageTag,
	}
	ctx := context.Background()
	con, err := conduit.NewConduitBootstrapper(ctx, options)
	if err != nil {
		deletePDir()
		PrintFatalError(NewSetupError(err))

	}
	if err := con.WriteComposeFile(); err != nil {
		deletePDir()
		PrintFatalError(NewSetupError(err))
	}
	if err := con.WriteEnvFile(); err != nil {
		deletePDir()
		PrintFatalError(NewSetupError(err))
	}
	if err := con.WriteConduitJsonFile(); err != nil {
		deletePDir()
		PrintFatalError(NewSetupError(err))
	}

	cancelSigKill()
	if err := con.Up(ctx); err != nil {
		deletePDir()
		PrintFatalError(NewSetupError(err))
	}
	//If detached print success message cause otherwise the client will be consuming docker compose logs
	if detach {
		PrintSuccess("project created")
	}
	os.Exit(0)
}

// Invokes a callback after signal termination CTRL+C
func handleSIGTERM(ctx context.Context, cb func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		select {
		case <-c:
			cb()
			return
		case <-ctx.Done():
			// Do nothing if the context is canceled before receiving the signal
			return
		}
	}()
}

// Traverse up the directory tree until conduit.json is found or we reach the root directory
// and changes to that directory
func ChangeToProjectRootDir() error {
	// Start from the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return err
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

	return errors.New("you are not in a project directory")
}

// Checks to see if you are in root project dir
func IsInProjectDirectory() bool {
	_, err := os.Stat("conduit.json")
	return !errors.Is(err, os.ErrNotExist)
}
