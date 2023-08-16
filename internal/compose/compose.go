package compose

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/compose-spec/compose-go/loader"
	ctypes "github.com/compose-spec/compose-go/types"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/compose/v2/pkg/compose"
	"github.com/isolateminds/go-conduit-cli/internal/compose/composeopt"
	"github.com/isolateminds/go-conduit-cli/internal/compose/errordefs"
	"github.com/isolateminds/go-conduit-cli/internal/compose/types"
	"golang.org/x/exp/slices"
)

type Composer struct {
	project     *ctypes.Project
	service     api.Service
	logConsumer api.LogConsumer
	Options     *types.ComposerOptions
}

func (c *Composer) AllServicesNames() []string {
	result := []string{}
	for _, v := range c.project.AllServices() {
		result = append(result, v.Name)
	}
	return result
}

// Filters the underlying yaml profiles with the provided ones
// and returns the ones that only exist within the yaml - docker-compose.yaml
func (c *Composer) FilterYamlProfiles(profiles []string) []string {
	yamlProfiles := c.project.AllServices().GetProfiles()
	filtered := []string{}
	for _, profile := range profiles {
		if slices.Contains(yamlProfiles, profile) {
			filtered = append(filtered, profile)
		}
	}
	return filtered
}

func (c *Composer) Remove(ctx context.Context, services []string) error {
	err := c.checkServices(services)
	if err != nil {
		return errordefs.NewComposerRemoveError(err)
	}
	err = c.service.Remove(ctx, c.project.Name, api.RemoveOptions{
		Services: services,
		Project:  c.project,
		Stop:     true,
		Volumes:  true,
	})
	if err != nil {
		return errordefs.NewComposerRemoveError(err)
	}
	return nil
}
func (c *Composer) Stop(ctx context.Context, services []string) error {
	err := c.checkServices(services)
	if err != nil {
		return errordefs.NewComposerStopError(err)
	}
	err = c.service.Stop(ctx, c.project.Name, api.StopOptions{
		Project:  c.project,
		Services: services,
	})
	if err != nil {
		return errordefs.NewComposerStopError(err)
	}
	return nil

}
func (c *Composer) Up(ctx context.Context) error {
	err := c.service.Up(ctx, c.project, api.UpOptions{
		Create: api.CreateOptions{
			Services:             c.project.ServiceNames(),
			RemoveOrphans:        true,
			Recreate:             api.RecreateForce,
			RecreateDependencies: api.RecreateForce,
		},
		Start: api.StartOptions{
			Attach:   c.logConsumer,
			Project:  c.project,
			Services: c.project.ServiceNames(),
			AttachTo: c.project.ServiceNames(),
		},
	})
	if err != nil {
		return errordefs.NewComposerUpError(err)
	}
	return nil
}

func (c *Composer) Config(ctx context.Context) ([]byte, error) {
	return c.service.Config(ctx, c.project, api.ConfigOptions{
		Format: "yaml",
	})
}
func NewComposer(name string, setComposerOptions ...composeopt.SetComposerOptions) (*Composer, error) {
	options := &types.ComposerOptions{}
	if name == "" {
		return nil, errordefs.NewComposerError(errors.New("name must not be '' "))
	}
	options.Name = name
	for _, set := range setComposerOptions {
		if err := set(options); err != nil {
			return nil, errordefs.NewComposerError(err)
		}
	}

	if options.Client == nil {
		return nil, errordefs.NewComposerError(errors.New("no client provided"))

	} else if options.Environment == nil {
		return nil, errordefs.NewComposerError(errors.New("no environment provided"))
	} else if options.Yaml == nil {
		return nil, errordefs.NewComposerError(errors.New("no yaml provided"))
	} else {
		ctx := context.Background()
		configFile := ctypes.ConfigFile{
			Content: options.Yaml.Bytes,
		}
		configDetails := ctypes.ConfigDetails{
			Environment: options.Environment.Variables,
			ConfigFiles: []ctypes.ConfigFile{configFile},
		}
		project, err := loader.LoadWithContext(ctx, configDetails, func(o *loader.Options) {
			o.SetProjectName(options.Name, true)
		})
		if err != nil {
			return nil, errordefs.NewComposerError(err)
		}
		project.ApplyProfiles(options.Profiles)
		//Sets the proper docker compose labels this is how docker desktop
		//knows it's a compose project
		for i, s := range project.Services {
			s.CustomLabels = map[string]string{
				api.ProjectLabel:     project.Name,
				api.ServiceLabel:     s.Name,
				api.VersionLabel:     api.ComposeVersion,
				api.WorkingDirLabel:  project.WorkingDir,
				api.ConfigFilesLabel: strings.Join(project.ComposeFiles, ","),
				api.OneoffLabel:      "False",
			}
			project.Services[i] = s
		}
		cli, err := command.NewDockerCli(
			command.WithAPIClient(options.Client),
			command.WithCombinedStreams(os.Stdout),
		)
		if err != nil {
			return nil, errordefs.NewComposerError(err)
		}
		err = cli.Initialize(flags.NewClientOptions())
		if err != nil {
			return nil, errordefs.NewComposerError(err)
		}
		service := compose.NewComposeService(cli)
		return &Composer{
			project:     project,
			service:     service,
			Options:     options,
			logConsumer: options.LogConsumer,
		}, nil
	}
}

// Helper func checks to see if the services provided actually exist within the project
// such as being defined in a docker-compose.yaml
func (c *Composer) checkServices(services []string) error {
	definedServices := c.AllServicesNames()
	notDefined := []string{}
	for _, service := range services {
		if !slices.Contains(definedServices, service) {
			notDefined = append(notDefined, service)
		}
	}
	if len(notDefined) > 0 {
		return errordefs.NewComposerError(fmt.Errorf("%v are not defined services", notDefined))
	}
	return nil
}
