package compose

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/compose-spec/compose-go/loader"
	"github.com/compose-spec/compose-go/types"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/compose/v2/pkg/compose"
	"github.com/docker/docker/errdefs"
	"github.com/isolateminds/go-conduit-cli/internal/docker"
)

type Composer struct {
	project *types.Project
	service api.Service
}

func (c *Composer) Up(ctx context.Context) {
	err := c.service.Up(ctx, c.project, api.UpOptions{
		Create: api.CreateOptions{
			Services:             c.project.ServiceNames(),
			RemoveOrphans:        true,
			Recreate:             api.RecreateForce,
			RecreateDependencies: api.RecreateForce,
		},
		Start: api.StartOptions{
			Project:  c.project,
			Services: c.project.ServiceNames(),
			AttachTo: c.project.ServiceNames(),
		},
	})
	if err != nil {
		if errdefs.IsConflict(err) {
			err := c.service.Start(ctx, c.project.Name, api.StartOptions{
				Project:  c.project,
				AttachTo: c.project.ServiceNames(),
			})
			if err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Println("oo: ", err)
		}
	}
}

func (c *Composer) Config(ctx context.Context) ([]byte, error) {
	return c.service.Config(ctx, c.project, api.ConfigOptions{
		Format: "yaml",
	})
}

func NewComposer(projectName string, dc *docker.Client, environment *Environment, yaml *Yaml) (*Composer, error) {
	ctx := context.Background()
	configFile := types.ConfigFile{
		Content: yaml.Bytes,
	}
	configDetails := types.ConfigDetails{
		Environment: environment.Variables,
		ConfigFiles: []types.ConfigFile{configFile},
	}
	options := func(o *loader.Options) {
		o.SetProjectName(projectName, true)
	}
	project, err := loader.LoadWithContext(ctx, configDetails, options)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	client := dc.Unwrap()
	project.ApplyProfiles([]string{"mongodb"})

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
	cli, err := command.NewDockerCli(command.WithAPIClient(client))
	if err != nil {
		return nil, err
	}
	err = cli.Initialize(&flags.ClientOptions{Debug: true})
	if err != nil {
		return nil, err
	}
	service := compose.NewComposeService(cli)
	return &Composer{
		project: project,
		service: service,
	}, nil
}
