package conduit

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/isolateminds/go-conduit-cli/internal/compose/composeopt"
	"github.com/isolateminds/go-conduit-cli/internal/docker"
	"gopkg.in/yaml.v2"
)

type variableMap map[string]string
type envFormatter struct {
	formatter   composeopt.EnvFormatter
	VariableMap variableMap
}

func (tf *envFormatter) Format(in []byte) (out io.Reader, err error) {
	env := string(in)
	for k, v := range tf.VariableMap {
		env = strings.ReplaceAll(env, fmt.Sprintf("{{%s}}", k), v)
	}
	return bytes.NewReader([]byte(env)), nil
}

func newEnvFormatter(variableMap variableMap) *envFormatter {
	return &envFormatter{
		VariableMap: variableMap,
	}
}

type composeBindDbFormatter struct {
	DatabaseName string
	formatter    composeopt.YamlFormatter
}

func (f *composeBindDbFormatter) Format(in []byte) (out []byte, err error) {
	t := &DockerCompose{}
	err = yaml.Unmarshal(in, t)
	if err != nil {
		return nil, err
	}
	delete(t.Volumes, "mongo")
	delete(t.Volumes, "postgres")

	for name := range t.Services {
		switch f.DatabaseName {
		case "mongodb":
			delete(t.Services, "postgres")
			if service, ok := t.Services[name]; ok {
				service.Volumes = []string{"./database/:/data/db"}
				t.Services[name] = service
			}
		case "postgres":
			delete(t.Services, "mongodb")
			if service, ok := t.Services[name]; ok {
				service.Volumes = []string{"./database/:/var/lib/postgresql/data"}
				t.Services[name] = service
			}

		default:
			return nil, fmt.Errorf("invalid database name %s", f.DatabaseName)
		}
	}

	out, err = yaml.Marshal(t)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type bindFormatterOptions struct {
	DatabaseName string
	Client       *docker.Client
}

/*
This YamlFormatter formats the docker compose file by setting a bind mounts for the db specified
It also deletes anything database related of the opposite database so if dbName is mongodb
all postgres related keys will be deleted
*/
func newComposeBindDbFormatter(dbName string) *composeBindDbFormatter {
	return &composeBindDbFormatter{DatabaseName: dbName}
}
