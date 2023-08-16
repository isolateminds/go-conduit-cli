package conduit

import (
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"io/ioutil"
	"os"

	"github.com/isolateminds/go-conduit-cli/internal/compose"
	"github.com/isolateminds/go-conduit-cli/internal/compose/composeopt"
	"github.com/isolateminds/go-conduit-cli/internal/docker"
	"github.com/isolateminds/go-conduit-cli/internal/utils"
	"github.com/isolateminds/go-conduit-cli/pkg/conduit/errordefs"
)

const (
	//Hard coded for now I guess

	mongoENVTemplateURL    = "https://raw.githubusercontent.com/isolateminds/go-conduit-cli/main/templates/mongo.env"
	postgresENVTemplateURL = "https://raw.githubusercontent.com/isolateminds/go-conduit-cli/main/templates/postgres.env"

	yamlURL = "https://raw.githubusercontent.com/isolateminds/go-conduit-cli/main/content/docker-compose.yml"
)

type Conduit struct {
	composer *compose.Composer
	json     *ConduitJson
}

func (c *Conduit) Remove(ctx context.Context, services []string) error {
	return c.composer.Remove(ctx, services)
}
func (c *Conduit) Stop(ctx context.Context, services []string) error {
	return c.composer.Stop(ctx, services)
}

func (c *Conduit) Up(ctx context.Context) error {
	return c.composer.Up(ctx)
}

// For already bootstrapped projects must be in project root dir when you call this
func NewConduitFromProject(ctx context.Context, detached bool, profiles []string) (*Conduit, error) {
	//Automatically checks if connected to daemon
	client, err := docker.NewClient(ctx)
	if err != nil {
		return nil, errordefs.NewNewConduitFromProjectError(err)
	}
	//Load persisted data
	data := &ConduitJson{}
	b, err := ioutil.ReadFile("conduit.json")
	if err != nil {
		return nil, errordefs.NewNewConduitFromProjectError(err)
	}
	err = json.Unmarshal(b, data)
	if err != nil {
		return nil, errordefs.NewNewConduitFromProjectError(err)
	}

	//block databases from being added because they where added already during bootstrapping
	//block the profiles specified inside conduit.json
	blockList := append(data.Profiles, "postgres", "mongodb")
	//new profiles are ones who haven't been added before and any databases
	newProfiles := blockProfiles(profiles, blockList...)
	//the updated profiles are a mix of new and saved profiles (conduit.json)
	updatedProfiles := append(data.Profiles, newProfiles...)

	composer, err := compose.NewComposer(
		data.ProjectName,
		composeopt.Client(client),
		composeopt.EnvFromFile(".env"),
		composeopt.YamlFromFile("docker-compose.yaml"),
		withDetachedFlag(ctx, detached),
		//the profiles added here will be used with dcoker compose
		composeopt.Profiles(updatedProfiles...),
	)
	if err != nil {
		return nil, errordefs.NewNewConduitFromProjectError(err)
	}

	return &Conduit{
		composer: composer,
		json: &ConduitJson{
			ProjectName: data.ProjectName,
			Version:     data.Version,
			Database:    data.Database,
			//filter the profiles here to save the actual profiles defined in the schema
			Profiles: composer.FilterYamlProfiles(updatedProfiles),
		},
	}, nil
}

// For bootsrapping conduit projects and enabling profiles
func NewConduitBootstrapper(ctx context.Context, name string, detached bool, profiles []string) (*Conduit, error) {
	//Automatically checks if connected to daemon
	client, err := docker.NewClient(ctx)
	if err != nil {
		return nil, errordefs.NewConduitBootstrapperError(err)
	}
	db, err := getDatabaseName(profiles)
	if err != nil {
		return nil, errordefs.NewConduitBootstrapperError(err)
	}
	composer, err := compose.NewComposer(
		name,
		composeopt.Client(client),
		composeopt.YamlFetchUrl(yamlURL),
		composeopt.Profiles(profiles...),
		withDetachedFlag(ctx, detached),
		withEnvFromDatabaseProfile(ctx, db),
	)
	if err != nil {
		return nil, errordefs.NewConduitBootstrapperError(err)
	}

	return &Conduit{
		composer: composer,
		json: &ConduitJson{
			ProjectName: name,
			Database:    db,
			//filter the profiles here to save the actual profiles defined in the schema
			Profiles: composer.FilterYamlProfiles(profiles),
		},
	}, nil
}

// Writes the docker-compose.yaml to current path
func (c *Conduit) WriteComposeFile() error {
	return os.WriteFile("docker-compose.yaml", c.composer.Options.Yaml.Bytes, fs.ModePerm)
}

// Writes the .env to current path
func (c *Conduit) WriteEnvFile() error {
	return os.WriteFile(".env", c.composer.Options.Environment.Bytes, fs.ModePerm)
}

// wWrites the json file to current path
func (c *Conduit) WriteConduitJsonFile() error {
	return c.json.WriteFile()
}

// If detatched no logging will be done same as --detach or -d flag in docker compose
func withDetachedFlag(ctx context.Context, detached bool) composeopt.SetComposerOptions {
	if detached {
		return composeopt.CustomLogConsumer(nil)
	}
	return composeopt.DefaultComposeLogConsumer(ctx)
}

// Fetches either the mongodb .env  template or the postgres depending on profiles
func withEnvFromDatabaseProfile(ctx context.Context, db string) composeopt.SetComposerOptions {

	dbPass := utils.GenerateRandomString(32)
	switch db {
	case "mongodb":
		formatter := newTemplateFormatter(variableMap{
			"MasterKey":     utils.GenerateRandomString(64),
			"MongoPassword": dbPass,
		})
		return composeopt.TemplateEnvFetchUrl(mongoENVTemplateURL, formatter)
	case "postgres":
		formatter := newTemplateFormatter(variableMap{
			"MasterKey":        utils.GenerateRandomString(64),
			"PostgresPassword": dbPass,
		})
		return composeopt.TemplateEnvFetchUrl(postgresENVTemplateURL, formatter)
	default:
		return composeopt.Error("a database profile has not been given use")
	}
}

// Helper function 'blocks' profiles from being included in the result.
// It takes a list of profiles and a variable number of profiles to be blocked.
// Blocked profiles are excluded from the result list.
func blockProfiles(profiles []string, blocked ...string) []string {
	blockedSet := make(map[string]struct{})
	for _, b := range blocked {
		blockedSet[b] = struct{}{}
	}

	result := []string{}
	for _, p := range profiles {
		if _, ok := blockedSet[p]; !ok {
			result = append(result, p)
		}
	}
	return result
}

func getDatabaseName(profiles []string) (string, error) {
	var db string
	var dbs []string
	for _, profile := range profiles {
		if profile == "mongodb" {
			db = profile
			dbs = append(dbs, db)
		}
		if profile == "postgres" {
			db = profile
			dbs = append(dbs, db)
		}
	}
	if len(dbs) > 1 {
		return "", errors.New("cannot use multiple database profiles")
	}
	return db, nil
}
