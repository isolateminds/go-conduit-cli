package types

import (
	"io"
	"net/http"
	"os"

	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/docker/client"
	"github.com/joho/godotenv"
)

// This is an in memory environment object for use with compose files.
type Environment struct {
	//Bytes of env file to preserve comments, whitespace and order etc.
	Bytes []byte
	// Key/Value pairs
	Variables map[string]string
}

// writes key-value pairs to the specified file destination.
func (e *Environment) WriteFile(dst string) error {
	return godotenv.Write(e.Variables, dst)
}

func NewEnvironment(kvPairs map[string]string) *Environment {
	return &Environment{
		Variables: kvPairs,
	}
}

// This function is designed for loading environment variables from a specific file source.
func NewEnvFromFile(src string) (*Environment, error) {
	b, err := os.ReadFile(src)
	if err != nil {
		return nil, err
	}
	kvPairs, err := godotenv.Read(src)
	if err != nil {
		return nil, err
	}
	return &Environment{
		Bytes:     b,
		Variables: kvPairs,
	}, nil
}

// This function is useful when you need to load environment variables
// from an external source, such as from a GET request response body.
func NewEnvFromURL(url string) (env *Environment, err error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	kvPairs, err := godotenv.Parse(res.Body)
	if err != nil {
		return
	}
	return &Environment{
		Bytes:     b,
		Variables: kvPairs,
	}, nil
}

type Yaml struct {
	Bytes []byte
}

// Implements stringer interface
func (y *Yaml) String() string {
	return string(y.Bytes)
}

// This function is designed for yaml from a specific file source.
func LoadYamlFromFile(src string) (yaml *Yaml, err error) {
	b, err := os.ReadFile(src)
	if err != nil {
		return
	}
	return &Yaml{
		Bytes: b,
	}, nil
}

// This function is useful when you need to load yaml files
// from an external source, such as from a GET request response body.
func LoadYamlFromURL(url string) (yaml *Yaml, err error) {
	res, err := http.Get(url)
	if err != nil {
		return
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	return &Yaml{
		Bytes: b,
	}, nil
}

type ComposerOptions struct {
	Name        string
	Client      client.APIClient
	Environment *Environment
	Yaml        *Yaml
	Profiles    []string
	LogConsumer api.LogConsumer
}
