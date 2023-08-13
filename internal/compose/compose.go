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
	"github.com/isolateminds/go-conduit-cli/internal/compose/types"
)

type Composer struct {
	project     *ctypes.Project
	service     api.Service
	logConsumer api.LogConsumer
	options     *types.ComposerOptions
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
		return fmt.Errorf("ComposerUpError: %s", err)
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
		return nil, fmt.Errorf("NewComposerError: name cannot be empty")
	}
	options.Name = name
	for _, set := range setComposerOptions {
		if err := set(options); err != nil {
			return nil, fmt.Errorf("NewComposerError: %s", err)
		}
	}

	if options.Client == nil {
		return nil, errors.New("NewComposerError: no client provided")
	} else if options.Environment == nil {
		return nil, errors.New("NewComposerError: no environment provided")
	} else if options.Yaml == nil {
		return nil, errors.New("NewComposerError: no yaml provided")
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
			return nil, fmt.Errorf("NewComposerError: %s", err)
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
			return nil, fmt.Errorf("NewComposerError: %s", err)
		}
		err = cli.Initialize(flags.NewClientOptions())
		if err != nil {
			return nil, fmt.Errorf("NewComposerError: %s", err)
		}
		service := compose.NewComposeService(cli)
		return &Composer{
			project:     project,
			service:     service,
			options:     options,
			logConsumer: options.LogConsumer,
		}, nil
	}
}
