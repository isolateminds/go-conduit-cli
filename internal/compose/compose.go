package compose

import (
	"context"
	"fmt"
	"sync"

	"github.com/compose-spec/compose-go/loader"
	"github.com/compose-spec/compose-go/types"
	"github.com/isolateminds/go-conduit-cli/internal/docker"
	"github.com/isolateminds/go-conduit-cli/internal/docker/imageopt"
)

type Composer struct {
	project *types.Project
	client  *docker.Client
	images  map[string]docker.Image
}

func (c *Composer) PullImages(ctx context.Context, setOFns ...imageopt.SetPullOptFn) error {
	services := c.project.AllServices()

	for _, service := range services {
		image := c.client.NewImage(service.Image)
		image.SetPullOptions(setOFns...)
		c.images[service.Image] = *image
		if err := c.client.PullImage(ctx, image); err != nil {
			return err
		}
	}
	return nil
}
func (c *Composer) PullImagesConcurrently(ctx context.Context, setOFns ...imageopt.SetPullOptFn) error {
	var wg sync.WaitGroup
	services := c.project.AllServices()
	wg.Add(len(services))
	errCh := make(chan error, len(services))

	for _, service := range services {
		image := c.client.NewImage(service.Image)
		image.SetPullOptions(setOFns...)
		c.images[service.Image] = *image
		go func(i *docker.Image) {
			defer wg.Done()
			if err := c.client.PullImage(ctx, image); err != nil {
				errCh <- err
			}
		}(image)
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			return err
		}
	}
	return nil
}

func NewComposer(client *docker.Client, projectName string, environment *Environment, yaml *Yaml) (*Composer, error) {
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

	return &Composer{
		project: project,
		client:  client,
		images:  map[string]docker.Image{},
	}, nil

}
