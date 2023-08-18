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
	"github.com/isolateminds/go-conduit-cli/pkg/conduit/errordefs"
	"github.com/joho/godotenv"
)

type SetComposerOptions func(opt *types.ComposerOptions) error
type EnvFormatter interface {
	Format(in []byte) (out io.Reader, err error)
}
type YamlFormatter interface {
	Format(in []byte) (out []byte, err error)
}

func WithProfiles(profiles ...string) SetComposerOptions {
	return func(opt *types.ComposerOptions) error {
		opt.Profiles = append(opt.Profiles, profiles...)
		return nil
	}
}

func WithClient(client *docker.Client) SetComposerOptions {
	return func(opt *types.ComposerOptions) error {
		opt.Client = client.Unwrap()
		return nil
	}
}
func WithYamlFromFile(src string) SetComposerOptions {
	return func(opt *types.ComposerOptions) (err error) {
		opt.Yaml, err = types.LoadYamlFromFile(src)
		if err != nil {
			return errordefs.NewYamlFileError(err)
		}
		return nil
	}
}

func WithEnvFromFile(src string) SetComposerOptions {
	return func(opt *types.ComposerOptions) (err error) {
		opt.Environment, err = types.NewEnvFromFile(src)
		if err != nil {
			return errordefs.NewEnvFileError(err)
		}
		return nil
	}
}
func WithEnvFromUrl(url string) SetComposerOptions {
	return func(opt *types.ComposerOptions) (err error) {
		opt.Environment, err = types.NewEnvFromURL(url)
		if err != nil {
			return errordefs.NewEnvFileError(err)
		}
		return nil
	}
}

/*
Returns an error for better handling when writing custom logic with SetComposerOptions

	//err message would be "MyError"
	con, err :=  compose.NewComposer(
		"project-name",
		composeopt.Error("MyError")
	)
*/
func WithError(message string) SetComposerOptions {
	return func(opt *types.ComposerOptions) error {
		return errors.New(message)
	}
}

// Fetches Yaml From url and runs the parser
func WithYamlFromUrlFormatter(url string, parser YamlFormatter) SetComposerOptions {
	return func(opt *types.ComposerOptions) error {
		if parser == nil {
			return fmt.Errorf("parser must not be %v", parser)
		}
		res, err := http.Get(url)
		if err != nil {
			return errordefs.NewYamlFileError(err)
		}
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)
		b, err := parser.Format(data)
		if err != nil {
			return errordefs.NewYamlFileError(err)
		}
		opt.Yaml = &types.Yaml{
			Bytes: b,
		}
		return nil
	}
}

// Fetches template from URL and runs the formatter
func WithEnvFromUrlFormatter(url string, formatter EnvFormatter) SetComposerOptions {
	return func(opt *types.ComposerOptions) (err error) {
		if formatter == nil {
			return fmt.Errorf("formatter must not be %v", formatter)
		}
		res, err := http.Get(url)
		if err != nil {
			return errordefs.NewEnvFileError(err)
		}
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return errordefs.NewEnvFileError(err)
		}
		formatted, err := formatter.Format(data)
		if err != nil {
			return errordefs.NewEnvFileError(err)
		}
		formattedBytes, err := io.ReadAll(formatted)
		if err != nil {
			return errordefs.NewEnvFileError(err)
		}
		parsed, err := godotenv.UnmarshalBytes(formattedBytes)
		if err != nil {
			return errordefs.NewEnvFileError(err)
		}

		opt.Environment = &types.Environment{
			Bytes:     formattedBytes,
			Variables: parsed,
		}
		if err != nil {
			return errordefs.NewEnvFileError(err)
		}
		return nil
	}
}
func WithYamlFromUrl(url string) SetComposerOptions {
	return func(opt *types.ComposerOptions) (err error) {
		opt.Yaml, err = types.LoadYamlFromURL(url)
		if err != nil {
			return errordefs.NewYamlFileError(err)
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
func WithCustomLogConsumer(consumer api.LogConsumer) SetComposerOptions {

	return func(opt *types.ComposerOptions) error {
		opt.LogConsumer = consumer
		return nil
	}
}

// The default docker compose logger when --detach flag is not zeroed
func WithDefaultComposeLogConsumer(ctx context.Context) SetComposerOptions {
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
