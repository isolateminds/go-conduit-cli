package compose

import (
	"net/http"

	"github.com/joho/godotenv"
)

// This is an in memory environment object for use with compose files.
type Environment struct {
	// Key / Value pairs
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
func LoadEnvFromFile(src string) (*Environment, error) {
	kvPairs, err := godotenv.Read(src)
	if err != nil {
		return nil, err
	}
	return &Environment{
		Variables: kvPairs,
	}, nil
}

// This function is useful when you need to load environment variables
// from an external source, such as from a GET request response body.
func LoadEnvFromURL(url string) (env *Environment, err error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	kvPairs, err := godotenv.Parse(res.Body)
	if err != nil {
		return
	}
	return &Environment{
		Variables: kvPairs,
	}, nil
}
