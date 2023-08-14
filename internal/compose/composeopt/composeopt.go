package composeopt

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/docker/compose/v2/pkg/api"
	"github.com/isolateminds/go-conduit-cli/internal/compose/types"
	"github.com/isolateminds/go-conduit-cli/internal/docker"
	"github.com/joho/godotenv"
)

type SetComposerOptions func(opt *types.ComposerOptions) error
type TemplateFormatter interface {
	Format(in []byte) (out io.Reader, err error)
}

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
func YamlFromFile(src string) SetComposerOptions {
	return func(opt *types.ComposerOptions) (err error) {
		opt.Yaml, err = types.LoadYamlFromFile(src)
		if err != nil {
			return fmt.Errorf("YamlFromFile: %s", err)
		}
		return nil
	}
}

func EnvFromFile(src string) SetComposerOptions {
	return func(opt *types.ComposerOptions) (err error) {
		opt.Environment, err = types.LoadEnvFromFile(src)
		if err != nil {
			return fmt.Errorf("EnvFromFile: %s", err)
		}
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
func Error(message string) SetComposerOptions {
	return func(opt *types.ComposerOptions) error {
		return errors.New(message)
	}
}
func TemplateEnvFetchUrl(url string, formatter TemplateFormatter) SetComposerOptions {
	return func(opt *types.ComposerOptions) (err error) {
		if formatter == nil {
			return fmt.Errorf("TemplateEnvFetchError: Formatter must not be %v", formatter)
		}
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("TemplateEnvFetchError: %s", err)
		}
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("TemplateEnvFetchError: %s", err)
		}
		formatted, err := formatter.Format(data)
		if err != nil {
			return fmt.Errorf("TemplateEnvFetchError: %s", err)
		}
		parsed, err := godotenv.Parse(formatted)
		if err != nil {
			return fmt.Errorf("TemplateEnvFetchError: %s", err)
		}
		opt.Environment = types.NewEnvironment(parsed)
		if err != nil {
			return fmt.Errorf("TemplateEnvFetchError: %s", err)
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
