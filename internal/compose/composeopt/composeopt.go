package composeopt

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/docker/compose/v2/pkg/api"
	"github.com/isolateminds/go-conduit-cli/internal/compose/types"
	"github.com/isolateminds/go-conduit-cli/internal/docker"
)

type SetComposerOptions func(opt *types.ComposerOptions) error

func Profiles(profiles ...string) SetComposerOptions {
	return func(opt *types.ComposerOptions) error {
		opt.Profiles = append(opt.Profiles, profiles...)
		return nil
	}
}

func Client(client *docker.Client) SetComposerOptions {
	return func(opt *types.ComposerOptions) error {
		opt.Client = client.Unwrap()
		return nil
	}
}

func EnvFetchUrl(url string) SetComposerOptions {
	return func(opt *types.ComposerOptions) (err error) {
		opt.Environment, err = types.LoadEnvFromURL(url)
		if err != nil {
			return fmt.Errorf("EnvFetchError: %s", err)
		}
		return nil
	}
}
func YamlFetchUrl(url string) SetComposerOptions {
	return func(opt *types.ComposerOptions) (err error) {
		opt.Yaml, err = types.LoadYamlFromURL(url)
		if err != nil {
			return fmt.Errorf("YamlFetchError: %s", err)
		}
		return nil
	}
}

/*
Provide a custom log consumer that implements:

		type LogConsumer interface {
	    	Log(containerName, message string)
	    	Err(containerName, message string)
	    	Status(container, msg string)
	    	Register(container string)
		}
*/
func CustomLogConsumer(consumer api.LogConsumer) SetComposerOptions {

	return func(opt *types.ComposerOptions) error {
		opt.LogConsumer = consumer
		return nil
	}
}

// The default docker compose logger when --detach flag is not zeroed
func DefaultComposeLogConsumer(ctx context.Context) SetComposerOptions {
	return func(opt *types.ComposerOptions) error {
		opt.LogConsumer = &logConsumer{
			ctx:        ctx,
			presenters: sync.Map{},
			width:      0,
			stdout:     os.Stdout,
			stderr:     os.Stderr,
			color:      true,
			prefix:     true,
			timestamp:  false,
		}
		return nil
	}
}
