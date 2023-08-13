package compose

import (
	"io"
	"net/http"
	"os"
)

type Yaml struct{ Bytes []byte }

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
